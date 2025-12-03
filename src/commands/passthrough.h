#pragma once

#include <string>
#include <vector>

namespace gip {
namespace commands {

/**
 * @brief Pass through unknown commands to git.
 *
 * For commands not handled by gip (e.g., status, log, branch),
 * this forwards the command directly to git and returns its output.
 *
 * @param args Command line arguments to pass to git
 * @return Exit code from git command
 */
[[nodiscard]] auto passthrough(const std::vector<std::string>& args) -> int;

}  // namespace commands
}  // namespace gip
