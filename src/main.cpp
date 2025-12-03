/**
 * @file main.cpp
 * @brief Gip - Git with Intent Protocol
 *
 * A drop-in replacement for git that enforces semantic context
 * for LLM-native development workflows.
 *
 * @copyright Copyright (c) 2025
 */

#include <cstring>
#include <iostream>
#include <string>
#include <vector>

#include "commands/commit.h"
#include "commands/context.h"
#include "commands/init.h"
#include "commands/merge.h"
#include "commands/passthrough.h"
#include "commands/push.h"
#include "commands/rebase.h"

namespace {

constexpr const char* kVersion = "1.0.0";
constexpr const char* kColorCyan = "\033[36m";
constexpr const char* kColorReset = "\033[0m";
constexpr const char* kColorBold = "\033[1m";

void printVersion()
{
    std::cout << "gip version " << kVersion << std::endl;
    std::cout << "Git with Intent Protocol - LLM-native version control" << std::endl;
}

void printHelp()
{
    std::cout << kColorBold << "Gip - Git with Intent Protocol" << kColorReset << std::endl;
    std::cout << std::endl;
    std::cout << "A drop-in replacement for git that enforces semantic context" << std::endl;
    std::cout << "for LLM-native development workflows." << std::endl;
    std::cout << std::endl;
    std::cout << kColorCyan << "ENHANCED COMMANDS:" << kColorReset << std::endl;
    std::cout << "  gip init                  Initialize repo with AI instructions"
              << std::endl;
    std::cout << "  gip commit -m \"msg\"       Commit with manifest (required)" << std::endl;
    std::cout << "  gip commit -f -m \"msg\"    Force commit without manifest" << std::endl;
    std::cout << "  gip push                  Push code AND context notes to remote"
              << std::endl;
    std::cout << "  gip merge <branch>        Merge with enriched conflict markers"
              << std::endl;
    std::cout << "  gip rebase <branch>       Rebase with enriched conflict markers"
              << std::endl;
    std::cout << std::endl;
    std::cout << kColorCyan << "CONTEXT COMMANDS:" << kColorReset << std::endl;
    std::cout << "  gip context <file>              Show semantic history of a file"
              << std::endl;
    std::cout << "  gip context <file> --json       Output as JSON (machine-readable)"
              << std::endl;
    std::cout << "  gip context --all               Show context for all tracked files"
              << std::endl;
    std::cout << "  gip context --behavior <type>   Filter by behavior (feature, bugfix, etc.)"
              << std::endl;
    std::cout << "  gip context --since <date>      Filter commits since date (YYYY-MM-DD)"
              << std::endl;
    std::cout << "  gip context --export <file>     Export context to JSON file" << std::endl;
    std::cout << std::endl;
    std::cout << kColorCyan << "PASSTHROUGH:" << kColorReset << std::endl;
    std::cout << "  All other git commands are passed through directly." << std::endl;
    std::cout << "  Example: gip status, gip log, gip branch, etc." << std::endl;
    std::cout << std::endl;
    std::cout << kColorCyan << "OPTIONS:" << kColorReset << std::endl;
    std::cout << "  --version, -v         Show version" << std::endl;
    std::cout << "  --help, -h            Show this help" << std::endl;
    std::cout << std::endl;
    std::cout << "For more information: https://github.com/iamHrithikRaj/gip" << std::endl;
}

}  // anonymous namespace

auto main(int argc, char* argv[]) -> int
{
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
