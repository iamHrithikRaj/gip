#pragma once

#include <string>
#include <vector>

namespace gip {
namespace commands {

/**
 * @brief Execute 'gip commit' with manifest enforcement.
 *
 * Commits staged changes with a required semantic manifest.
 * The manifest documents the intent and contracts of the changes.
 *
 * @param args Command line arguments passed to commit
 * @return Exit code (0 on success, non-zero on failure)
 */
[[nodiscard]] auto commit(const std::vector<std::string>& args) -> int;

}  // namespace commands
}  // namespace gip
