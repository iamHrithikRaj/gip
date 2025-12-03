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
constexpr const char* kConflictStart = "<<<<<<<";
constexpr const char* kConflictMiddle = "=======";
constexpr const char* kConflictEnd = ">>>>>>>";
constexpr const char* kConflictBase = "|||||||";  // For diff3 style

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

/// @brief Extract value from TOON-style line: "key: value" or "key: \"value\""
std::string extractToonValue(const std::string& line) {
    auto colonPos = line.find(':');
    if (colonPos == std::string::npos)
        return "";

    std::string value = trim(line.substr(colonPos + 1));

    // Remove quotes if present
    if (value.length() >= 2 && value.front() == '"' && value.back() == '"') {
        value = value.substr(1, value.length() - 2);
    }
    return value;
}

/// @brief Extract array values from TOON content for a given key
std::vector<std::string> extractToonArray(const std::string& content, const std::string& key) {
    std::vector<std::string> result;
    std::istringstream stream(content);
    std::string line;
    bool inArray = false;

    while (std::getline(stream, line)) {
        std::string trimmedLine = trim(line);

        // Check for key start
        if (trimmedLine.find(key + ":") == 0) {
            inArray = true;
            // Check if value is on same line
            std::string value = extractToonValue(trimmedLine);
            if (!value.empty() && value != "[" && value != "{") {
                result.push_back(value);
                inArray = false;
            }
            continue;
        }

        if (inArray) {
            // Check for array item (starts with -)
            if (trimmedLine.front() == '-') {
                std::string item = trim(trimmedLine.substr(1));
                // Remove quotes if present
                if (item.length() >= 2 && item.front() == '"' && item.back() == '"') {
                    item = item.substr(1, item.length() - 2);
                }
                result.push_back(item);
            }
            // Check for end of array (next key at same indentation)
            else if (!trimmedLine.empty() && trimmedLine.find(':') != std::string::npos &&
                     line.find_first_not_of(" \t") == 0) {
                inArray = false;
            }
        }
    }

    return result;
}

/// @brief Extract single value from TOON content for a given key
/// Supports nested keys like "behavior.changes" by looking for indented content
std::string extractToonField(const std::string& content, const std::string& key) {
    // Check for nested key (e.g., "behavior.changes")
    auto dotPos = key.find('.');
    if (dotPos != std::string::npos) {
        std::string parentKey = key.substr(0, dotPos);
        std::string childKey = key.substr(dotPos + 1);

        std::istringstream stream(content);
        std::string line;
        bool inParent = false;

        while (std::getline(stream, line)) {
            std::string trimmedLine = trim(line);

            // Find parent section
            if (trimmedLine.find(parentKey + ":") == 0) {
                inParent = true;
                continue;
            }

            if (inParent) {
                // Check for child key (indented)
                if (line.find_first_not_of(" \t") > 0 && trimmedLine.find(childKey + ":") == 0) {
                    return extractToonValue(trimmedLine);
                }
                // Exit parent if we hit a non-indented line
                if (!trimmedLine.empty() && line.find_first_not_of(" \t") == 0 &&
                    trimmedLine.find(':') != std::string::npos) {
                    inParent = false;
                }
            }
        }
        return "";
    }

    // Simple key lookup
    std::istringstream stream(content);
    std::string line;

    while (std::getline(stream, line)) {
        std::string trimmedLine = trim(line);
        if (trimmedLine.find(key + ":") == 0) {
            return extractToonValue(trimmedLine);
        }
    }

    return "";
}

/// @brief Extract boolean value from TOON content
bool extractToonBool(const std::string& content, const std::string& key) {
    std::string value = extractToonField(content, key);
    return value == "true" || value == "yes" || value == "1";
}

/// @brief Extract the first entry from an entries array in TOON format
/// The stored manifest has format: {entries: [{file: ..., rationale: ...}]}
/// This extracts fields from the first entry
std::string extractFromFirstEntry(const std::string& content, const std::string& key) {
    // Look for entries array and extract from first entry
    std::istringstream stream(content);
    std::string line;
    bool inEntries = false;
    bool inFirstEntry = false;
    int braceDepth = 0;

    while (std::getline(stream, line)) {
        std::string trimmedLine = trim(line);

        // Track brace depth to know when we're inside entries
        for (char c : trimmedLine) {
            if (c == '{')
                braceDepth++;
            else if (c == '}')
                braceDepth--;
        }

        // Found entries array
        if (trimmedLine.find("entries:") == 0 || trimmedLine.find("entries :") == 0) {
            inEntries = true;
            continue;
        }

        // Inside entries, look for first object
        if (inEntries && !inFirstEntry && trimmedLine.find('{') != std::string::npos) {
            inFirstEntry = true;
        }

        // Inside first entry, look for our key
        if (inFirstEntry) {
            if (trimmedLine.find(key + ":") == 0 || trimmedLine.find(key + " :") == 0) {
                return extractToonValue(trimmedLine);
            }
            // Exit if we hit closing brace of first entry
            if (trimmedLine.find('}') != std::string::npos && braceDepth <= 2) {
                break;
            }
        }
    }

    return "";
}

/// @brief Extract array from first entry in entries array
std::vector<std::string> extractArrayFromFirstEntry(const std::string& content,
                                                    const std::string& key) {
    std::vector<std::string> result;
    std::istringstream stream(content);
    std::string line;
    bool inEntries = false;
    bool inFirstEntry = false;
    bool inTargetArray = false;
    int braceDepth = 0;
    int bracketDepth = 0;

    while (std::getline(stream, line)) {
        std::string trimmedLine = trim(line);

        // Track depths
        for (char c : trimmedLine) {
            if (c == '{')
                braceDepth++;
            else if (c == '}')
                braceDepth--;
            else if (c == '[')
                bracketDepth++;
            else if (c == ']')
                bracketDepth--;
        }

        // Found entries array
        if (trimmedLine.find("entries:") == 0 || trimmedLine.find("entries :") == 0) {
            inEntries = true;
            continue;
        }

        if (inEntries && !inFirstEntry && trimmedLine.find('{') != std::string::npos) {
            inFirstEntry = true;
        }

        if (inFirstEntry) {
            // Look for target array key
            if (trimmedLine.find(key + ":") == 0 || trimmedLine.find(key + " :") == 0) {
                inTargetArray = true;
                // Check for inline array like: preconditions: ["value1", "value2"]
                auto bracketStart = trimmedLine.find('[');
                auto bracketEnd = trimmedLine.find(']');
                if (bracketStart != std::string::npos && bracketEnd != std::string::npos) {
                    // Parse inline array
                    std::string arrayContent =
                        trimmedLine.substr(bracketStart + 1, bracketEnd - bracketStart - 1);
                    std::istringstream arrayStream(arrayContent);
                    std::string item;
                    while (std::getline(arrayStream, item, ',')) {
                        item = trim(item);
                        // Remove quotes
                        if (item.length() >= 2 && item.front() == '"' && item.back() == '"') {
                            item = item.substr(1, item.length() - 2);
                        }
                        if (!item.empty()) {
                            result.push_back(item);
                        }
                    }
                    inTargetArray = false;
                }
                continue;
            }

            if (inTargetArray) {
                // Array items start with -
                if (!trimmedLine.empty() && trimmedLine[0] == '-') {
                    std::string item = trim(trimmedLine.substr(1));
                    if (item.length() >= 2 && item.front() == '"' && item.back() == '"') {
                        item = item.substr(1, item.length() - 2);
                    }
                    result.push_back(item);
                }
                // End of array
                else if (trimmedLine.find(']') != std::string::npos ||
                         (!trimmedLine.empty() && trimmedLine.find(':') != std::string::npos)) {
                    inTargetArray = false;
                }
            }

            // Exit first entry
            if (trimmedLine.find('}') != std::string::npos && braceDepth <= 2) {
                break;
            }
        }
    }

    return result;
}

}  // anonymous namespace

// =============================================================================
// Conflict Detection
// =============================================================================

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
            output << line << "\n";

            // Add enriched context for "ours" side
            if (oursContext.has_value()) {
                output << formatEnrichedMarker("HEAD", "Your changes", *oursContext);
            }
        } else if (isConflictMarker(line, kConflictMiddle)) {
            output << line << "\n";
        } else if (isConflictMarker(line, kConflictEnd)) {
            // Add enriched context for "theirs" side before end marker
            if (theirsContext.has_value()) {
                output << formatEnrichedMarker(extractMarkerInfo(line), "Their changes",
                                               *theirsContext);
            }
            output << line << "\n";
        } else {
            output << line << "\n";
        }
    }

    return writeFile(filePath, output.str());
}

int MergeDriver::enrichAllConflicts(const std::string& oursSha,
                                    const std::string& theirsSha) const {
    auto conflictedFiles = getConflictedFiles();
    int enrichedCount = 0;

    for (const auto& file : conflictedFiles) {
        if (enrichConflictMarkers(file, oursSha, theirsSha)) {
            ++enrichedCount;
        }
    }

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
    output << "||| Commit: " << context.commitSha << "\n";

    if (!context.behaviorClass.empty()) {
        output << "||| behaviorClass: " << context.behaviorClass << "\n";
    }

    if (!context.behaviorChanges.empty()) {
        output << "||| behaviorChanges: " << context.behaviorChanges << "\n";
    }

    output << "||| breaking: " << (context.breaking ? "true" : "false") << "\n";

    if (context.breaking && !context.migrations.empty()) {
        for (size_t i = 0; i < context.migrations.size(); ++i) {
            output << "||| migrations[" << i << "]: " << context.migrations[i] << "\n";
        }
    }

    if (!context.inputs.empty()) {
        for (size_t i = 0; i < context.inputs.size(); ++i) {
            output << "||| inputs[" << i << "]: " << context.inputs[i] << "\n";
        }
    }

    if (!context.outputs.empty()) {
        output << "||| outputs: " << context.outputs << "\n";
    }

    // Output preconditions
    for (size_t i = 0; i < context.preconditions.size(); ++i) {
        output << "||| preconditions[" << i << "]: " << context.preconditions[i] << "\n";
    }

    // Output postconditions
    for (size_t i = 0; i < context.postconditions.size(); ++i) {
        output << "||| postconditions[" << i << "]: " << context.postconditions[i] << "\n";
    }

    if (!context.errorModel.empty()) {
        for (size_t i = 0; i < context.errorModel.size(); ++i) {
            output << "||| errorModel[" << i << "]: " << context.errorModel[i] << "\n";
        }
    }

    if (!context.rationale.empty()) {
        output << "||| rationale: " << context.rationale << "\n";
    }

    // Output side effects
    if (context.sideEffects.empty()) {
        output << "||| sideEffects: none\n";
    } else {
        for (size_t i = 0; i < context.sideEffects.size(); ++i) {
            output << "||| sideEffects[" << i << "]: " << context.sideEffects[i] << "\n";
        }
    }

    // Output testing info
    if (!context.testingRequired.empty()) {
        for (size_t i = 0; i < context.testingRequired.size(); ++i) {
            output << "||| testingRequired[" << i << "]: " << context.testingRequired[i] << "\n";
        }
    }

    if (!context.testingCoverage.empty()) {
        output << "||| testingCoverage: " << context.testingCoverage << "\n";
    }

    if (!context.symbol.empty()) {
        output << "||| symbol: " << context.symbol << "\n";
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
            for (const auto& entry : manifest.entries) {
                if (entry.file == filePath) {
                    matchedEntry = &entry;
                    break;
                }
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
