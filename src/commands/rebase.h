/// @file rebase.h
/// @brief Gip rebase command with enriched conflict markers
/// @author Gip Team
/// @copyright MIT License

#pragma once

#include <string>
#include <vector>

namespace gip {
namespace commands {

/// @brief Execute rebase with enriched conflict markers
///
/// This command wraps git rebase to provide:
/// 1. Automatic note preservation (via notes.rewriteRef configuration)
/// 2. Enriched conflict markers when conflicts occur
///
/// @param args Command arguments (passed to git rebase)
/// @return Exit code (0 on success)
auto rebase(const std::vector<std::string>& args) -> int;

}  // namespace commands
}  // namespace gip
