/// @file merge_driver.cpp
/// @brief Implementation of merge driver for enriching Git conflict markers
/// @author Gip Team
/// @copyright MIT License

#include "merge_driver.h"

#include "git_adapter.h"
#include "manifest.h"

#include <algorithm>
#include <cstring>
#include <fstream>
#include <iostream>
#include <regex>
#include <sstream>

namespace gip {

namespace {

/// @brief Conflict marker patterns
// NOLINTBEGIN(readability-identifier-naming)
constexpr const char* kConflictStart = "<<<<<<<";
constexpr const char* kConflictMiddle = "=======";
constexpr const char* kConflictEnd = ">>>>>>>";
constexpr const char* kConflictBase = "|||||||";  // For diff3 style
// NOLINTEND(readability-identifier-naming)

/// @brief Check if a line starts with a conflict marker
bool isConflictMarker(const std::string& line, const char* marker) {
    return line.compare(0, strlen(marker), marker) == 0;
}

/// @brief Extract branch/commit info from conflict marker line
/// e.g., "<<<<<<< HEAD" -> "HEAD"
std::string extractMarkerInfo(const std::string& line) {
    auto spacePos = line.find(' ');
    if (spacePos != std::string::npos && spacePos + 1 < line.length()) {
        return line.substr(spacePos + 1);
    }
    return "";
}

/// @brief Trim whitespace from string
std::string trim(const std::string& str) {
    const auto start = str.find_first_not_of(" \t\r\n");
    if (start == std::string::npos)
        return "";
    const auto end = str.find_last_not_of(" \t\r\n");
    return str.substr(start, end - start + 1);
}

}  // anonymous namespace

// =============================================================================
// Conflict Detection
// =============================================================================

// cppcheck-suppress unusedFunction
bool MergeDriver::hasConflictMarkers(const std::string& filePath) const {
    std::string content = readFile(filePath);
    return content.find(kConflictStart) != std::string::npos;
}

std::vector<std::string> MergeDriver::getConflictedFiles() const {
    std::vector<std::string> conflictedFiles;

    GitAdapter git;
    // Use git diff --check to find files with conflicts
    auto result = git.execute({"diff", "--name-only", "--diff-filter=U"});

    if (result.success()) {
        std::istringstream stream(result.stdoutOutput);
        std::string line;
        while (std::getline(stream, line)) {
            std::string trimmed = trim(line);
            if (!trimmed.empty()) {
                conflictedFiles.push_back(trimmed);
            }
        }
    }

    return conflictedFiles;
}

// =============================================================================
// Conflict Enrichment
// =============================================================================

bool MergeDriver::enrichConflictMarkers(const std::string& filePath, const std::string& oursSha,
                                        const std::string& theirsSha) const {
    std::string content = readFile(filePath);
    if (content.empty()) {
        return false;
    }

    // Extract contexts for both sides
    auto oursContext = extractContext(oursSha, filePath);
    auto theirsContext = extractContext(theirsSha, filePath);

    // If neither side has a manifest, nothing to enrich
    if (!oursContext.has_value() && !theirsContext.has_value()) {
        return false;
    }

    // Process line by line
    std::istringstream input(content);
    std::ostringstream output;
    std::string line;

    while (std::getline(input, line)) {
        if (isConflictMarker(line, kConflictStart)) {
            // Output the original marker
            output << line << '\n';

            // Add enriched context for "ours" side
            if (oursContext.has_value()) {
                output << formatEnrichedMarker("HEAD", "Your changes", *oursContext);
            }
        } else if (isConflictMarker(line, kConflictMiddle)) {
            output << line << '\n';
        } else if (isConflictMarker(line, kConflictEnd)) {
            // Add enriched context for "theirs" side before end marker
            if (theirsContext.has_value()) {
                output << formatEnrichedMarker(extractMarkerInfo(line), "Their changes",
                                               *theirsContext);
            }
            output << line << '\n';
        } else {
            output << line << '\n';
        }
    }

    return writeFile(filePath, output.str());
}

int MergeDriver::enrichAllConflicts(const std::string& oursSha,
                                    const std::string& theirsSha) const {
    auto conflictedFiles = getConflictedFiles();
    int enrichedCount = 0;

    enrichedCount = std::count_if(conflictedFiles.begin(), conflictedFiles.end(),
                                  [this, &oursSha, &theirsSha](const std::string& file) {
                                      return enrichConflictMarkers(file, oursSha, theirsSha);
                                  });

    return enrichedCount;
}

// =============================================================================
// Context Extraction
// =============================================================================

std::optional<ConflictContext> MergeDriver::extractContext(const std::string& commitSha,
                                                           const std::string& filePath) const {
    GitAdapter git;
    auto manifest = git.getNote(commitSha);

    if (!manifest.has_value()) {
        return std::nullopt;
    }

    auto context = parseManifest(*manifest, commitSha, filePath);

    // Try to extract symbol from file path if not in manifest
    if (context.symbol.empty()) {
        // Use filename as fallback symbol
        auto lastSlash = filePath.find_last_of("/\\");
        context.symbol =
            (lastSlash != std::string::npos) ? filePath.substr(lastSlash + 1) : filePath;
    }

    return context;
}

// =============================================================================
// Marker Formatting
// =============================================================================

std::string MergeDriver::formatEnrichedMarker(const std::string& side,
                                              const std::string& description,
                                              const ConflictContext& context) const {
    std::ostringstream output;

    output << kContextPrefix << " (" << side << " - " << description << ")\n";
    output << "||| Commit: " << context.commitSha << '\n';

    if (!context.behaviorClass.empty()) {
        output << "||| behaviorClass: " << context.behaviorClass << '\n';
    }

    if (!context.behaviorChanges.empty()) {
        output << "||| behaviorChanges: " << context.behaviorChanges << '\n';
    }

    output << "||| breaking: " << (context.breaking ? "true" : "false") << '\n';

    if (context.breaking && !context.migrations.empty()) {
        for (size_t i = 0; i < context.migrations.size(); ++i) {
            output << "||| migrations[" << i << "]: " << context.migrations[i] << '\n';
        }
    }

    if (!context.inputs.empty()) {
        for (size_t i = 0; i < context.inputs.size(); ++i) {
            output << "||| inputs[" << i << "]: " << context.inputs[i] << '\n';
        }
    }

    if (!context.outputs.empty()) {
        output << "||| outputs: " << context.outputs << '\n';
    }

    // Output preconditions
    for (size_t i = 0; i < context.preconditions.size(); ++i) {
        output << "||| preconditions[" << i << "]: " << context.preconditions[i] << '\n';
    }

    // Output postconditions
    for (size_t i = 0; i < context.postconditions.size(); ++i) {
        output << "||| postconditions[" << i << "]: " << context.postconditions[i] << '\n';
    }

    if (!context.errorModel.empty()) {
        for (size_t i = 0; i < context.errorModel.size(); ++i) {
            output << "||| errorModel[" << i << "]: " << context.errorModel[i] << '\n';
        }
    }

    if (!context.rationale.empty()) {
        output << "||| rationale: " << context.rationale << '\n';
    }

    // Output side effects
    if (context.sideEffects.empty()) {
        output << "||| sideEffects: none\n";
    } else {
        for (size_t i = 0; i < context.sideEffects.size(); ++i) {
            output << "||| sideEffects[" << i << "]: " << context.sideEffects[i] << '\n';
        }
    }

    // Output testing info
    if (!context.testingRequired.empty()) {
        for (size_t i = 0; i < context.testingRequired.size(); ++i) {
            output << "||| testingRequired[" << i << "]: " << context.testingRequired[i] << '\n';
        }
    }

    if (!context.testingCoverage.empty()) {
        output << "||| testingCoverage: " << context.testingCoverage << '\n';
    }

    if (!context.symbol.empty()) {
        output << "||| symbol: " << context.symbol << '\n';
    }

    return output.str();
}

// =============================================================================
// Private Methods
// =============================================================================

ConflictContext MergeDriver::parseManifest(const std::string& manifestContent,
                                           const std::string& commitSha,
                                           const std::string& filePath) const {
    ConflictContext context;
    context.commitSha = commitSha.substr(0, 8);  // Short SHA

    // Use Manifest class to parse
    auto manifestOpt = Manifest::fromToon(manifestContent);
    if (!manifestOpt) {
        return context;
    }

    const auto& manifest = *manifestOpt;

    // Find the entry that matches the file path
    const ManifestEntry* matchedEntry = nullptr;

    if (!manifest.entries.empty()) {
        if (filePath.empty()) {
            // Default to first entry if no file path provided
            matchedEntry = &manifest.entries[0];
        } else {
            // Try to find exact match
            auto it = std::find_if(
                manifest.entries.begin(), manifest.entries.end(),
                [&filePath](const ManifestEntry& entry) { return entry.file == filePath; });
            if (it != manifest.entries.end()) {
                matchedEntry = &*it;
            }

            // If no exact match, try to find by filename
            if (!matchedEntry) {
                std::string filename = filePath;
                auto lastSlash = filename.find_last_of("/\\");
                if (lastSlash != std::string::npos) {
                    filename = filename.substr(lastSlash + 1);
                }

                for (const auto& entry : manifest.entries) {
                    std::string entryFilename = entry.file;
                    auto entryLastSlash = entryFilename.find_last_of("/\\");
                    if (entryLastSlash != std::string::npos) {
                        entryFilename = entryFilename.substr(entryLastSlash + 1);
                    }

                    if (entryFilename == filename) {
                        matchedEntry = &entry;
                        break;
                    }
                }
            }

            // Fallback to first entry if still not found
            if (!matchedEntry) {
                matchedEntry = &manifest.entries[0];
            }
        }
    }

    if (matchedEntry) {
        context.behaviorClass = matchedEntry->behavior;
        if (context.behaviorClass.empty())
            context.behaviorClass = matchedEntry->type;

        context.rationale = matchedEntry->rationale;
        context.behaviorChanges = matchedEntry->behavior;  // Use behavior as description
        context.breaking = matchedEntry->breaking;
        context.migrations = matchedEntry->migrations;
        context.inputs = matchedEntry->inputs;
        context.outputs = matchedEntry->outputs;
        context.symbol = matchedEntry->symbol;
        if (context.symbol.empty())
            context.symbol = matchedEntry->file;

        context.errorModel = matchedEntry->errorModel;
        context.preconditions = matchedEntry->preconditions;
        context.postconditions = matchedEntry->postconditions;
        context.sideEffects = matchedEntry->sideEffects;
    }

    return context;
}

std::string MergeDriver::readFile(const std::string& filePath) const {
    std::ifstream file(filePath);
    if (!file.is_open()) {
        return "";
    }

    std::ostringstream content;
    content << file.rdbuf();
    return content.str();
}

bool MergeDriver::writeFile(const std::string& filePath, const std::string& content) const {
    std::ofstream file(filePath);
    if (!file.is_open()) {
        return false;
    }

    file << content;
    return file.good();
}

}  // namespace gip
