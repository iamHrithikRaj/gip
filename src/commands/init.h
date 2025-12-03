#pragma once

#include <string>
#include <vector>

namespace gip {
namespace commands {

/**
 * @brief Execute 'gip init' - initializes repo and creates AI instructions.
 *
 * Initializes a git repository (if not already initialized) and creates
 * the .github/copilot-instructions.md file to guide AI assistants
 * on proper Gip usage.
 *
 * @param args Command line arguments passed to init
 * @return Exit code (0 on success, non-zero on failure)
 */
[[nodiscard]] auto init(const std::vector<std::string>& args) -> int;

}  // namespace commands
}  // namespace gip
