/// @file merge.cpp
/// @brief Implementation of Gip merge command with enriched conflict markers
/// @author Gip Team
/// @copyright MIT License

#include "merge.h"
#include "../git_adapter.h"
#include "../merge_driver.h"

#include <iostream>

namespace gip {
namespace commands {

namespace {

constexpr const char* kColorRed = "\033[31m";
constexpr const char* kColorGreen = "\033[32m";
constexpr const char* kColorYellow = "\033[33m";
constexpr const char* kColorCyan = "\033[36m";
constexpr const char* kColorReset = "\033[0m";
constexpr const char* kColorBold = "\033[1m";

/// @brief Get the SHA of MERGE_HEAD (the branch being merged)
std::string getMergeHead(const GitAdapter& git)
{
    auto result = git.execute({"rev-parse", "--short", "MERGE_HEAD"});
    if (result.success()) {
        std::string sha = result.stdoutOutput;
        // Trim newline
        while (!sha.empty() && (sha.back() == '\n' || sha.back() == '\r')) {
            sha.pop_back();
        }
        return sha;
    }
    return "";
}

/// @brief Get the SHA of HEAD (current position)
std::string getHead(const GitAdapter& git)
{
    auto result = git.execute({"rev-parse", "--short", "HEAD"});
    if (result.success()) {
        std::string sha = result.stdoutOutput;
        while (!sha.empty() && (sha.back() == '\n' || sha.back() == '\r')) {
            sha.pop_back();
        }
        return sha;
    }
    return "";
}

/// @brief Check if we're in the middle of a merge
bool isMergeInProgress(const GitAdapter& git)
{
    auto result = git.execute({"rev-parse", "--verify", "MERGE_HEAD"});
    return result.success();
}

/// @brief Extract the branch name from merge arguments
std::string extractMergeBranch(const std::vector<std::string>& args)
{
    for (const auto& arg : args) {
        // Skip flags
        if (arg.empty() || arg[0] == '-') {
            continue;
        }
        // First non-flag argument is likely the branch
        return arg;
    }
    return "";
}

/// @brief Print help for merge conflicts
void printConflictHelp()
{
    std::cerr << "\n" << kColorCyan << "Gip Conflict Resolution:" << kColorReset << "\n";
    std::cerr << "  Conflict markers have been enriched with manifest context.\n";
    std::cerr << "  Look for " << kColorBold << "||| Gip CONTEXT" << kColorReset 
              << " lines for structured intent information.\n";
    std::cerr << "\n";
    std::cerr << kColorBold << "What the enriched markers tell you:" << kColorReset << "\n";
    std::cerr << "  • " << kColorCyan << "behaviorClass" << kColorReset << ": Is this a feature, bugfix, or refactor?\n";
    std::cerr << "  • " << kColorCyan << "rationale" << kColorReset << ": Why was this change made?\n";
    std::cerr << "  • " << kColorCyan << "preconditions/postconditions" << kColorReset << ": Expected state before/after\n";
    std::cerr << "  • " << kColorCyan << "sideEffects" << kColorReset << ": Any side effects to consider\n";
    std::cerr << "\n";
    std::cerr << "  After resolving conflicts:\n";
    std::cerr << "    git add <resolved-files>\n";
    std::cerr << "    git commit\n";
    std::cerr << "\n";
    std::cerr << "  To abort the merge:\n";
    std::cerr << "    git merge --abort\n";
}

}  // anonymous namespace

auto merge(const std::vector<std::string>& args) -> int
{
    GitAdapter git;
    
    // Check if we're in a git repository
    if (!git.isRepository()) {
        std::cerr << kColorRed << "Error: " << kColorReset 
                  << "Not a git repository" << std::endl;
        return 128;
    }
    
    // Store current HEAD before merge
    std::string headBefore = getHead(git);
    
    // Build the git merge command
    std::vector<std::string> mergeArgs = {"merge"};
    mergeArgs.insert(mergeArgs.end(), args.begin(), args.end());
    
    // Execute the merge
    auto result = git.execute(mergeArgs);
    
    // Output git's stdout
    if (!result.stdoutOutput.empty()) {
        std::cout << result.stdoutOutput;
    }
    
    // Check if we have conflicts
    if (result.exitCode != 0 && isMergeInProgress(git)) {
        // Get the commit being merged (MERGE_HEAD)
        std::string mergeHead = getMergeHead(git);
        std::string currentHead = getHead(git);
        
        if (!mergeHead.empty() && !currentHead.empty()) {
            MergeDriver driver;
            auto conflictedFiles = driver.getConflictedFiles();
            
            if (!conflictedFiles.empty()) {
                std::cerr << "\n" << kColorYellow << "Enriching conflict markers with manifest context..." 
                          << kColorReset << "\n";
                
                int enrichedCount = driver.enrichAllConflicts(currentHead, mergeHead);
                
                if (enrichedCount > 0) {
                    std::cerr << kColorGreen << "✓ " << kColorReset 
                              << "Enriched " << enrichedCount << " file(s) with Gip context\n";
                    
                    // List enriched files
                    std::cerr << "\n" << kColorCyan << "Files with enriched conflicts:" << kColorReset << "\n";
                    for (const auto& file : conflictedFiles) {
                        std::cerr << "  " << file << "\n";
                    }
                    
                    printConflictHelp();
                } else {
                    std::cerr << kColorYellow << "Note: " << kColorReset 
                              << "No manifests found for conflicting commits.\n";
                    std::cerr << "Conflict markers are standard Git format.\n";
                }
            }
        }
        
        // Still output git's stderr
        if (!result.stderrOutput.empty()) {
            std::cerr << result.stderrOutput;
        }
    } else {
        // No conflicts or other output
        if (!result.stderrOutput.empty()) {
            std::cerr << result.stderrOutput;
        }
        
        if (result.success()) {
            std::cerr << kColorGreen << "✓ " << kColorReset 
                      << "Merge completed successfully.\n";
        }
    }
    
    return result.exitCode;
}

}  // namespace commands
}  // namespace gip
