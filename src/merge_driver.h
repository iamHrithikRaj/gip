/// @file merge_driver.h
/// @brief Merge driver for enriching Git conflict markers with manifest context
/// @author Gip Team
/// @copyright MIT License

#pragma once

#include <optional>
#include <string>
#include <vector>

namespace gip {

/// @brief Structured context extracted from a manifest for conflict enrichment
struct ConflictContext {
    std::string commitSha;                     ///< Short SHA of the commit
    std::string behaviorClass;                 ///< feature, bugfix, refactor, etc.
    std::string rationale;                     ///< Why this change was made (intent)
    std::string behaviorChanges;               ///< Description of behavioral changes
    bool breaking = false;                     ///< Is this a breaking change?
    std::vector<std::string> migrations;       ///< Migration steps if breaking
    std::vector<std::string> inputs;           ///< Input contract (args/params)
    std::string outputs;                       ///< Output contract (return type/value)
    std::string symbol;                        ///< Function/class being modified
    std::vector<std::string> preconditions;    ///< Required preconditions
    std::vector<std::string> postconditions;   ///< Expected postconditions
    std::vector<std::string> sideEffects;      ///< Known side effects
    std::vector<std::string> errorModel;       ///< Error handling approach
    std::vector<std::string> testingRequired;  ///< Required test types
    std::string testingCoverage;               ///< Coverage notes
};

/// @brief Merge driver for enriching conflict markers
///
/// This class provides the core functionality for enriching Git conflict
/// markers with structured manifest context, enabling LLMs to understand
/// the intent behind conflicting changes.
class MergeDriver {
public:
    MergeDriver() = default;
    ~MergeDriver() = default;

    // Non-copyable, movable
    MergeDriver(const MergeDriver&) = delete;
    MergeDriver& operator=(const MergeDriver&) = delete;
    MergeDriver(MergeDriver&&) noexcept = default;
    MergeDriver& operator=(MergeDriver&&) noexcept = default;

    // =========================================================================
    // Conflict Detection
    // =========================================================================

    /// @brief Check if a file has conflict markers
    /// @param filePath Path to the file to check
    /// @return true if file contains conflict markers
    [[nodiscard]] bool hasConflictMarkers(const std::string& filePath) const;

    /// @brief Get all files with conflict markers in the working directory
    /// @return Vector of file paths with conflicts
    [[nodiscard]] std::vector<std::string> getConflictedFiles() const;

    // =========================================================================
    // Conflict Enrichment
    // =========================================================================

    /// @brief Enrich conflict markers in a file with manifest context
    /// @param filePath Path to the file with conflicts
    /// @param oursSha SHA of "ours" commit (HEAD)
    /// @param theirsSha SHA of "theirs" commit (incoming)
    /// @return true if enrichment was successful
    [[nodiscard]] bool enrichConflictMarkers(const std::string& filePath,
                                             const std::string& oursSha,
                                             const std::string& theirsSha) const;

    /// @brief Enrich all conflicted files in the working directory
    /// @param oursSha SHA of "ours" commit (HEAD)
    /// @param theirsSha SHA of "theirs" commit (incoming)
    /// @return Number of files enriched
    [[nodiscard]] int enrichAllConflicts(const std::string& oursSha,
                                         const std::string& theirsSha) const;

    // =========================================================================
    // Context Extraction
    // =========================================================================

    /// @brief Extract conflict context from a manifest for a specific commit
    /// @param commitSha SHA of the commit
    /// @param filePath Path to the file (for symbol detection)
    /// @return ConflictContext if manifest exists, nullopt otherwise
    [[nodiscard]] std::optional<ConflictContext> extractContext(const std::string& commitSha,
                                                                const std::string& filePath) const;

    // =========================================================================
    // Marker Formatting
    // =========================================================================

    /// @brief Format enriched conflict marker header
    /// @param side Either "HEAD" or branch name
    /// @param description Human-readable description
    /// @param context The conflict context to embed
    /// @return Formatted marker lines
    [[nodiscard]] std::string formatEnrichedMarker(const std::string& side,
                                                   const std::string& description,
                                                   const ConflictContext& context) const;

private:
    /// @brief Parse manifest TOON content into ConflictContext
    /// @param manifest Raw manifest content
    /// @param commitSha SHA for the context
    /// @param filePath Optional file path to match specific entry
    /// @return Parsed ConflictContext
    [[nodiscard]] ConflictContext parseManifest(const std::string& manifest,
                                                const std::string& commitSha,
                                                const std::string& filePath = "") const;

    /// @brief Read file content
    /// @param filePath Path to read
    /// @return File content or empty string on error
    [[nodiscard]] std::string readFile(const std::string& filePath) const;

    /// @brief Write file content
    /// @param filePath Path to write
    /// @param content Content to write
    /// @return true if successful
    [[nodiscard]] bool writeFile(const std::string& filePath, const std::string& content) const;

    /// @brief Marker prefix for Gip context lines
    static constexpr const char* kContextPrefix = "||| Gip CONTEXT";
};

}  // namespace gip
