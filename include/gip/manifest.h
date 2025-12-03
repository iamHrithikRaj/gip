#pragma once

#include "ctoon.h"

#include <optional>
#include <string>
#include <utility>
#include <vector>

namespace gip {

/**
 * @brief Represents a single entry in a Gip manifest.
 *
 * Each entry describes a change made to a specific file or symbol,
 * including semantic information about the nature and purpose of the change.
 */
struct ManifestEntry {
    std::string file;       ///< Path to the modified file
    std::string symbol;     ///< Primary symbol/function/class affected
    std::string type;       ///< Change type: add, modify, delete, rename
    std::string behavior;   ///< Semantic category: feature, bugfix, refactor, perf
    std::string rationale;  ///< Human-readable explanation of why

    bool breaking = false;                ///< Is this a breaking change?
    std::vector<std::string> migrations;  ///< Migration steps if breaking
    std::vector<std::string> inputs;      ///< Input contract (args/params)
    std::string outputs;                  ///< Output contract (return type/value)
    std::vector<std::string> errorModel;  ///< Error conditions/exceptions

    std::vector<std::string> preconditions;   ///< What must be true before this change
    std::vector<std::string> postconditions;  ///< What is guaranteed after this change
    std::vector<std::string> sideEffects;     ///< Other effects of this change
};

/**
 * @brief Represents a complete Gip manifest.
 *
 * A manifest contains structured semantic metadata about a commit,
 * stored in TOON format for efficient LLM consumption.
 */
struct Manifest {
    std::string schemaVersion = "2.0";   ///< Manifest schema version
    std::vector<ManifestEntry> entries;  ///< List of change entries

    /**
     * @brief Check if manifest is empty.
     * @return true if no entries
     */
    [[nodiscard]] bool empty() const noexcept {
        return entries.empty();
    }

    /**
     * @brief Get number of entries.
     * @return Entry count
     */
    [[nodiscard]] size_t size() const noexcept {
        return entries.size();
    }

    /**
     * @brief Serialize manifest to TOON format string.
     * @return TOON-encoded string representation
     */
    [[nodiscard]] auto toToon() const -> std::string;

    /**
     * @brief Convert manifest to ctoon::Value for further processing.
     * @return ctoon::Value representation of the manifest
     */
    [[nodiscard]] auto toValue() const -> ctoon::Value;

    /**
     * @brief Parse a manifest from a TOON format string.
     * @param toonStr The TOON-encoded manifest string
     * @return Parsed manifest, or nullopt if parsing fails
     */
    [[nodiscard]] static auto fromToon(const std::string& toonStr) -> std::optional<Manifest>;

    /**
     * @brief Parse a manifest from a ctoon::Value.
     * @param value The ctoon::Value to parse
     * @return Parsed manifest, or nullopt if parsing fails
     */
    [[nodiscard]] static auto fromValue(const ctoon::Value& value) -> std::optional<Manifest>;
};

/**
 * @brief Parser for extracting Gip manifests from commit messages.
 *
 * Handles the identification, extraction, and parsing of gip: blocks
 * embedded in commit messages. Supports both TOON format and simplified
 * YAML-like syntax.
 */
class ManifestParser {
public:
    /**
     * @brief Result of parsing a commit message for manifest content.
     */
    struct ParseResult {
        std::string cleanMessage;          ///< Message with manifest block stripped
        std::optional<Manifest> manifest;  ///< Parsed manifest (if found)
        std::string error;                 ///< Error message (if parsing failed)

        /**
         * @brief Check if a manifest was successfully parsed.
         * @return true if manifest is present
         */
        [[nodiscard]] auto hasManifest() const noexcept -> bool {
            return manifest.has_value();
        }

        /**
         * @brief Check if parsing encountered an error.
         * @return true if an error occurred
         */
        [[nodiscard]] auto hasError() const noexcept -> bool {
            return !error.empty();
        }

        /**
         * @brief Check if parsing was successful.
         * @return true if no error and manifest found
         */
        [[nodiscard]] auto isValid() const noexcept -> bool {
            return !hasError() && hasManifest();
        }
    };

    /**
     * @brief Parse a commit message, extracting any gip: manifest block.
     *
     * Searches for a "gip:" marker in the message and attempts to parse
     * the following content as a TOON manifest. The manifest block is
     * stripped from the returned clean message.
     *
     * @param message The full commit message to parse
     * @return ParseResult containing clean message and optional manifest
     */
    [[nodiscard]] static auto parse(const std::string& message) -> ParseResult;

    /**
     * @brief Generate a template manifest for staged files.
     *
     * Creates a pre-filled manifest template that users can edit to
     * document their changes. Includes placeholder values for semantic
     * fields.
     *
     * @param files Vector of (file_path, git_status) pairs
     * @return TOON-formatted manifest template string
     */
    [[nodiscard]] static auto
    generateTemplate(const std::vector<std::pair<std::string, std::string>>& files) -> std::string;

    /**
     * @brief Validate a manifest.
     * @param manifest The manifest to validate
     * @return Error message if invalid, empty string if valid
     */
    [[nodiscard]] static auto validate(const Manifest& manifest) -> std::string;

private:
    /**
     * @brief Locate the gip: block boundaries in a message.
     * @param message The message to search
     * @return Pair of (start, end) positions, or (npos, npos) if not found
     */
    [[nodiscard]] static auto findManifestBlock(const std::string& message)
        -> std::pair<size_t, size_t>;
};

}  // namespace gip
