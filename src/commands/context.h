#pragma once

#include <string>
#include <vector>

namespace gip {
namespace commands {

/**
 * @brief Options for the context command.
 */
struct ContextOptions {
    std::string filePath;        ///< File to show context for (empty for --all)
    std::string exportPath;      ///< Path to export JSON (empty for terminal output)
    std::string behaviorFilter;  ///< Filter by behavior class (empty for all)
    std::string sinceDate;       ///< Filter commits since date (empty for all)
    bool jsonOutput = false;     ///< Output as JSON instead of terminal
    bool showAll = false;        ///< Show context for all files
};

/**
 * @brief Parse command line arguments into ContextOptions.
 * @param args Command line arguments to parse
 * @return Parsed options structure
 */
[[nodiscard]] auto parseContextArgs(const std::vector<std::string>& args)
    -> ContextOptions;

/**
 * @brief Execute 'gip context' with various options.
 *
 * Shows the semantic history of files with manifest context.
 * Supports filtering by behavior type, date range, and export formats.
 *
 * @param args Command line arguments passed to context
 * @return Exit code (0 on success, non-zero on failure)
 */
[[nodiscard]] auto context(const std::vector<std::string>& args) -> int;

}  // namespace commands
}  // namespace gip
