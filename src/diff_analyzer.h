#pragma once

#include <string>
#include <utility>
#include <vector>

namespace gip {

/**
 * @brief Information about a symbol extracted from a diff.
 *
 * Represents a code symbol (function, class, method) that was
 * modified in a commit, along with metadata about the change.
 */
struct SymbolInfo {
    std::string file;        ///< Path to the source file
    std::string name;        ///< Name of the symbol
    std::string type;        ///< Type: function, class, method, variable
    std::string changeType;  ///< Change kind: add, modify, delete
    int startLine;           ///< Starting line number
    int endLine;             ///< Ending line number
};

/**
 * @brief Analyzes git diffs to extract changed symbols.
 *
 * Parses unified diff output to identify functions, classes, and
 * other code symbols that were modified. Supports multiple
 * programming languages.
 */
class DiffAnalyzer {
public:
    /**
     * @brief Analyze staged diff and extract changed symbols.
     * @param diff The unified diff output from git
     * @return Vector of symbols found in the diff
     */
    [[nodiscard]] static auto analyze(const std::string& diff) -> std::vector<SymbolInfo>;

    /**
     * @brief Get list of changed files with their status.
     *
     * Parses git diff-index or status output to extract file paths
     * and their modification status.
     *
     * @param diffStatus Output from git diff-index or similar
     * @return Vector of (filepath, status) pairs where status is A/M/D/R
     */
    [[nodiscard]] static auto getChangedFiles(const std::string& diffStatus)
        -> std::vector<std::pair<std::string, std::string>>;

private:
    /**
     * @brief Extract symbols from a single file's diff.
     * @param filePath Path to the file
     * @param fileDiff Diff content for this file
     * @param changeType Type of change (add, modify, delete)
     * @return Vector of symbols found in the file diff
     */
    [[nodiscard]] static auto extractSymbols(const std::string& filePath,
                                             const std::string& fileDiff,
                                             const std::string& changeType)
        -> std::vector<SymbolInfo>;

    /**
     * @brief Detect programming language from file extension.
     * @param filePath Path to the file
     * @return Language identifier (cpp, python, javascript, etc.)
     */
    [[nodiscard]] static auto detectLanguage(const std::string& filePath) -> std::string;

    /**
     * @brief Extract function/symbol name from a line.
     * @param line Source code line
     * @param language Language identifier
     * @return Symbol name if found, empty string otherwise
     */
    [[nodiscard]] static auto extractSymbolName(const std::string& line,
                                                const std::string& language) -> std::string;
};

}  // namespace gip
