/// @file merge.h
/// @brief Gip merge command with enriched conflict markers
/// @author Gip Team
/// @copyright MIT License

#pragma once

#include <string>
#include <vector>

namespace gip {
namespace commands {

/// @brief Execute merge with enriched conflict markers
///
/// This command wraps git merge to provide:
/// 1. Enriched conflict markers when conflicts occur
/// 2. Structured manifest context for LLM-assisted resolution
///
/// @param args Command arguments (passed to git merge)
/// @return Exit code (0 on success)
auto merge(const std::vector<std::string>& args) -> int;

}  // namespace commands
}  // namespace gip
