#include "push.h"

#include "../git_adapter.h"

#include <algorithm>
#include <iostream>

namespace gip {
namespace commands {

namespace {

// ANSI color codes for terminal output
constexpr const char* kColorGreen = "\033[32m";
constexpr const char* kColorYellow = "\033[33m";
constexpr const char* kColorCyan = "\033[36m";
constexpr const char* kColorReset = "\033[0m";

void printSuccess(const std::string& msg) {
    std::cout << kColorGreen << "[✓] " << msg << kColorReset << std::endl;
}

void printInfo(const std::string& msg) {
    std::cout << kColorCyan << "[i] " << msg << kColorReset << std::endl;
}

void printWarning(const std::string& msg) {
    std::cout << kColorYellow << "[!] " << msg << kColorReset << std::endl;
}

}  // anonymous namespace

auto push(const std::vector<std::string>& args) -> int {
    GitAdapter git;

    if (!git.isRepository()) {
        std::cerr << "Not a git repository" << std::endl;
        return 1;
    }

    // Determine remote and branch
    std::string remote = "origin";
    std::string branch = "";

    // Parse args for remote and branch
    for (size_t i = 0; i < args.size(); ++i) {
        const auto& arg = args[i];
        if (arg[0] != '-') {
            if (remote == "origin") {
                remote = arg;
            } else if (branch.empty()) {
                branch = arg;
            }
        }
    }

    // If branch not specified, get current branch
    if (branch.empty()) {
        auto result = git.execute({"rev-parse", "--abbrev-ref", "HEAD"});
        if (result.success()) {
            branch = result.stdoutOutput;
            // Trim whitespace
            branch.erase(std::remove(branch.begin(), branch.end(), '\n'), branch.end());
            branch.erase(std::remove(branch.begin(), branch.end(), '\r'), branch.end());
        } else {
            branch = "main";
        }
    }

    printInfo("Pushing to " + remote + "/" + branch + "...");

    // Push branch
    std::vector<std::string> pushArgs = {"push"};
    pushArgs.insert(pushArgs.end(), args.begin(), args.end());

    // Add remote and branch if not already in args
    bool hasRemote = false;
    for (const auto& arg : args) {
        if (arg[0] != '-') {
            hasRemote = true;
            break;
        }
    }

    if (!hasRemote) {
        pushArgs.push_back(remote);
        pushArgs.push_back(branch);
    }

    auto result = git.execute(pushArgs);

    if (!result.success()) {
        std::cerr << "Push failed: " << result.stderrOutput << std::endl;
        return 1;
    }

    std::cout << result.stdoutOutput;
    if (!result.stderrOutput.empty()) {
        std::cout << result.stderrOutput;  // Git often outputs to stderr even on success
    }

    printSuccess("Pushed branch");

    // Push notes
    printInfo("Pushing context notes...");

    auto notesResult = git.execute({"push", remote, "refs/notes/gip"});

    if (notesResult.success()) {
        printSuccess("Pushed context notes");
    } else {
        // Notes might not exist yet, which is fine
        if (notesResult.stderrOutput.find("does not match any") != std::string::npos ||
            notesResult.stderrOutput.find("No refs") != std::string::npos) {
            printInfo("No context notes to push yet");
        } else {
            printWarning("Could not push notes: " + notesResult.stderrOutput);
        }
    }

    return 0;
}

}  // namespace commands
}  // namespace gip
