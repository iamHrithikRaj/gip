/**
 * @file main.cpp
 * @brief Gip - Git with Intent Protocol
 *
 * A drop-in replacement for git that enforces semantic context
 * for LLM-native development workflows.
 *
 * @copyright Copyright (c) 2025
 */

#include "commands/commit.h"
#include "commands/context.h"
#include "commands/init.h"
#include "commands/merge.h"
#include "commands/passthrough.h"
#include "commands/push.h"
#include "commands/rebase.h"

#include <cstring>
#include <iostream>
#include <string>
#include <vector>

namespace {

constexpr const char* kVersion = "1.0.0";
// NOLINTBEGIN(readability-identifier-naming)
constexpr const char* kColorCyan = "\033[36m";
constexpr const char* kColorReset = "\033[0m";
constexpr const char* kColorBold = "\033[1m";
// NOLINTEND(readability-identifier-naming)

void printVersion() {
    std::cout << "gip version " << kVersion << '\n';
    std::cout << "Git with Intent Protocol - LLM-native version control" << '\n';
}

void printHelp() {
    std::cout << kColorBold << "Gip - Git with Intent Protocol" << kColorReset << '\n';
    std::cout << '\n';
    std::cout << "A drop-in replacement for git that enforces semantic context" << '\n';
    std::cout << "for LLM-native development workflows." << '\n';
    std::cout << '\n';
    std::cout << kColorCyan << "ENHANCED COMMANDS:" << kColorReset << '\n';
    std::cout << "  gip init                  Initialize repo with AI instructions" << '\n';
    std::cout << "  gip commit -m \"msg\"       Commit with manifest (required)" << '\n';
    std::cout << "  gip commit -f -m \"msg\"    Force commit without manifest" << '\n';
    std::cout << "  gip push                  Push code AND context notes to remote" << '\n';
    std::cout << "  gip merge <branch>        Merge with enriched conflict markers" << '\n';
    std::cout << "  gip rebase <branch>       Rebase with enriched conflict markers" << '\n';
    std::cout << '\n';
    std::cout << kColorCyan << "CONTEXT COMMANDS:" << kColorReset << '\n';
    std::cout << "  gip context <file>              Show semantic history of a file" << '\n';
    std::cout << "  gip context <file> --json       Output as JSON (machine-readable)" << '\n';
    std::cout << "  gip context --all               Show context for all tracked files" << '\n';
    std::cout << "  gip context --behavior <type>   Filter by behavior (feature, bugfix, etc.)"
              << '\n';
    std::cout << "  gip context --since <date>      Filter commits since date (YYYY-MM-DD)" << '\n';
    std::cout << "  gip context --export <file>     Export context to JSON file" << '\n';
    std::cout << '\n';
    std::cout << kColorCyan << "PASSTHROUGH:" << kColorReset << '\n';
    std::cout << "  All other git commands are passed through directly." << '\n';
    std::cout << "  Example: gip status, gip log, gip branch, etc." << '\n';
    std::cout << '\n';
    std::cout << kColorCyan << "OPTIONS:" << kColorReset << '\n';
    std::cout << "  --version, -v         Show version" << '\n';
    std::cout << "  --help, -h            Show this help" << '\n';
    std::cout << '\n';
    std::cout << "For more information: https://github.com/iamHrithikRaj/gip" << '\n';
}

}  // anonymous namespace

auto main(int argc, char* argv[]) -> int {
    // No arguments - show help
    if (argc < 2) {
        printHelp();
        return 0;
    }

    std::string command = argv[1];

    // Build args vector (excluding program name and command)
    std::vector<std::string> args;
    for (int i = 2; i < argc; ++i) {
        args.push_back(argv[i]);
    }

    // Version
    if (command == "--version" || command == "-v") {
        printVersion();
        return 0;
    }

    // Help
    if (command == "--help" || command == "-h" || command == "help") {
        printHelp();
        return 0;
    }

    // Enhanced commands
    if (command == "commit") {
        return gip::commands::commit(args);
    }

    if (command == "init") {
        return gip::commands::init(args);
    }

    if (command == "push") {
        return gip::commands::push(args);
    }

    if (command == "merge") {
        return gip::commands::merge(args);
    }

    if (command == "rebase") {
        return gip::commands::rebase(args);
    }

    // Gip-only commands
    if (command == "context") {
        return gip::commands::context(args);
    }

    // Passthrough to git
    std::vector<std::string> gitArgs;
    gitArgs.push_back(command);
    gitArgs.insert(gitArgs.end(), args.begin(), args.end());

    return gip::commands::passthrough(gitArgs);
}
