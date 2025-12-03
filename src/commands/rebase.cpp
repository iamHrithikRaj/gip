/// @file rebase.cpp
/// @brief Implementation of Gip rebase command with enriched conflict markers
/// @author Gip Team
/// @copyright MIT License

#include "rebase.h"

#include "../git_adapter.h"
#include "../merge_driver.h"

#include <algorithm>
#include <iostream>

namespace gip::commands {

namespace {

// NOLINTBEGIN(readability-identifier-naming)
constexpr const char* kColorRed = "\033[31m";
constexpr const char* kColorGreen = "\033[32m";
constexpr const char* kColorYellow = "\033[33m";
constexpr const char* kColorCyan = "\033[36m";
constexpr const char* kColorReset = "\033[0m";
constexpr const char* kColorBold = "\033[1m";
// NOLINTEND(readability-identifier-naming)

/// @brief Get the SHA of REBASE_HEAD (the commit being applied)
std::string getRebaseHead(const GitAdapter& git) {
    auto result = git.execute({"rev-parse", "--short", "REBASE_HEAD"});
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
std::string getHead(const GitAdapter& git) {
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

/// @brief Check if we're in the middle of a rebase
bool isRebaseInProgress(const GitAdapter& git) {
    // Check for .git/rebase-merge or .git/rebase-apply directories
    auto result = git.execute({"rev-parse", "--git-path", "rebase-merge"});
    if (result.success()) {
        std::string path = result.stdoutOutput;
        while (!path.empty() && (path.back() == '\n' || path.back() == '\r')) {
            path.pop_back();
        }
        // Check if the directory exists via git command
        auto checkResult = git.execute({"rev-parse", "--verify", "REBASE_HEAD"});
        return checkResult.success();
    }
    return false;
}

/// @brief Configure git to preserve notes during rebase
void configureNotesRewrite(const GitAdapter& git) {
    // Enable note rewriting for rebase
    (void)git.execute({"config", "notes.rewrite.rebase", "true"});

    // Set the notes ref to rewrite
    (void)git.execute({"config", "notes.rewriteRef", "refs/notes/gip"});

    // Set rewrite mode to overwrite (replace old note with new)
    (void)git.execute({"config", "notes.rewriteMode", "overwrite"});
}

/// @brief Print help for rebase conflicts
void printConflictHelp() {
    std::cerr << "\n" << kColorCyan << "Gip Conflict Resolution:" << kColorReset << "\n";
    std::cerr << "  Conflict markers have been enriched with manifest context.\n";
    std::cerr << "  Look for " << kColorBold << "||| Gip CONTEXT" << kColorReset
              << " lines for structured intent information.\n";
    std::cerr << "\n";
    std::cerr << "  After resolving conflicts:\n";
    std::cerr << "    git add <resolved-files>\n";
    std::cerr << "    gip rebase --continue\n";
    std::cerr << "\n";
    std::cerr << "  To abort the rebase:\n";
    std::cerr << "    gip rebase --abort\n";
}

}  // anonymous namespace

auto rebase(const std::vector<std::string>& args) -> int {
    GitAdapter git;

    // Check if we're in a git repository
    if (!git.isRepository()) {
        std::cerr << kColorRed << "Error: " << kColorReset << "Not a git repository" << '\n';
        return 128;
    }

    // Configure notes preservation before any rebase operation
    configureNotesRewrite(git);

    // Check for --continue flag (resuming after conflict resolution)
    bool isContinue = std::any_of(args.begin(), args.end(),
                                  [](const std::string& arg) { return arg == "--continue"; });

    // Build the git rebase command
    std::vector<std::string> rebaseArgs = {"rebase"};
    rebaseArgs.insert(rebaseArgs.end(), args.begin(), args.end());

    // Store current HEAD before rebase for conflict enrichment
    // std::string headBefore = getHead(git);

    // Execute the rebase
    auto result = git.execute(rebaseArgs);

    // Output git's stdout
    if (!result.stdoutOutput.empty()) {
        std::cout << result.stdoutOutput;
    }

    // Check if we have conflicts
    if (result.exitCode != 0 && isRebaseInProgress(git)) {
        // Get the commit being rebased (the "theirs" side)
        std::string rebaseHead = getRebaseHead(git);
        std::string currentHead = getHead(git);

        if (!rebaseHead.empty() && !currentHead.empty()) {
            MergeDriver driver;
            auto conflictedFiles = driver.getConflictedFiles();

            if (!conflictedFiles.empty()) {
                std::cerr << "\n"
                          << kColorYellow << "Enriching conflict markers with manifest context..."
                          << kColorReset << "\n";

                int enrichedCount = driver.enrichAllConflicts(currentHead, rebaseHead);

                if (enrichedCount > 0) {
                    std::cerr << kColorGreen << "✓ " << kColorReset << "Enriched " << enrichedCount
                              << " file(s) with Gip context\n";

                    // List enriched files
                    std::cerr << "\n"
                              << kColorCyan << "Files with enriched conflicts:" << kColorReset
                              << "\n";
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

        if (result.success() && !isContinue) {
            std::cerr << kColorGreen << "✓ " << kColorReset
                      << "Rebase completed. Gip notes preserved.\n";
        }
    }

    return result.exitCode;
}

}  // namespace gip::commands
