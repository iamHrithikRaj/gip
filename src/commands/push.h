#pragma once

#include <string>
#include <vector>

namespace gip {
namespace commands {

/**
 * @brief Execute 'gip push' - pushes branch and context notes.
 *
 * Pushes the current branch to the remote, then pushes
 * the gip notes (refs/notes/gip) to share semantic context.
 *
 * @param args Command line arguments passed to push
 * @return Exit code (0 on success, non-zero on failure)
 */
[[nodiscard]] auto push(const std::vector<std::string>& args) -> int;

}  // namespace commands
}  // namespace gip
