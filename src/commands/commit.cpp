#include "commit.h"

#include "../git_adapter.h"
#include "../manifest.h"

#include <algorithm>
#include <cstddef>
#include <fstream>
#include <iostream>
#include <iterator>
#include <sstream>
#include <string>
#include <utility>
#include <vector>

namespace gip::commands {

namespace {

// ANSI color codes for terminal output
// NOLINTBEGIN(readability-identifier-naming)
constexpr const char* kColorRed = "\033[31m";
constexpr const char* kColorGreen = "\033[32m";
constexpr const char* kColorYellow = "\033[33m";
constexpr const char* kColorCyan = "\033[36m";
constexpr const char* kColorReset = "\033[0m";
constexpr const char* kColorBold = "\033[1m";
// NOLINTEND(readability-identifier-naming)

void printError(const std::string& msg) {
    std::cerr << kColorRed << "[!] " << msg << kColorReset << '\n';
}

void printSuccess(const std::string& msg) {
    std::cout << kColorGreen << "[✓] " << msg << kColorReset << '\n';
}

void printInfo(const std::string& msg) {
    std::cout << kColorCyan << "[i] " << msg << kColorReset << '\n';
}

auto hasFlag(const std::vector<std::string>& args, const std::string& shortFlag,
             const std::string& longFlag) -> bool {
    return std::find(args.begin(), args.end(), shortFlag) != args.end() ||
           std::find(args.begin(), args.end(), longFlag) != args.end();
}

auto getFlagValue(const std::vector<std::string>& args, const std::string& shortFlag,
                  const std::string& longFlag) -> std::string {
    for (size_t i = 0; i < args.size(); ++i) {
        if ((args[i] == shortFlag || args[i] == longFlag) && i + 1 < args.size()) {
            return args[i + 1];
        }
        // Handle -m"message" or --message="message" format
        if (args[i].substr(0, shortFlag.size()) == shortFlag && args[i].size() > shortFlag.size()) {
            return args[i].substr(shortFlag.size());
        }
        if (args[i].substr(0, longFlag.size() + 1) == longFlag + "=") {
            return args[i].substr(longFlag.size() + 1);
        }
    }
    return "";
}

}  // anonymous namespace

// NOLINTNEXTLINE(readability-function-size)
auto commit(const std::vector<std::string>& args) -> int {
    const GitAdapter gitAdapter;

    // Check if we're in a git repo
    if (!gitAdapter.isRepository()) {
        printError("Not a git repository");
        return 1;
    }

    // Check for force flag (-f / --force)
    const bool force = hasFlag(args, "-f", "--force");

    // Get commit message
    std::string message = getFlagValue(args, "-m", "--message");
    const std::string file = getFlagValue(args, "-F", "--file");

    if (!file.empty()) {
        const std::ifstream ifs(file);
        if (ifs) {
            std::ostringstream oss;
            oss << ifs.rdbuf();
            message = oss.str();
        } else {
            printError("Could not read file: " + file);
            return 1;
        }
    }

    if (message.empty()) {
        printError("Commit message required. Use: gip commit -m \"message\" or -F <file>");
        return 1;
    }

    // Get staged files
    auto stagedFiles = gitAdapter.getStagedFiles();
    if (stagedFiles.empty()) {
        printError("No staged changes. Use 'git add' first.");
        return 1;
    }

    // Parse the commit message for manifest
    auto parseResult = ManifestParser::parse(message);

    // If force flag is set, skip manifest check
    if (force) {
        printInfo("Force mode: Skipping manifest check");

        auto result = gitAdapter.commit(parseResult.cleanMessage);
        if (!result.success()) {
            printError("Commit failed: " + result.stderrOutput);
            return 1;
        }

        printSuccess("Committed (without manifest)");
        std::cout << result.stdoutOutput << '\n';
        return 0;
    }

    // Check if manifest was provided
    if (!parseResult.hasManifest()) {
        // Generate template for error message
        std::vector<std::pair<std::string, std::string>> files;
        files.reserve(stagedFiles.size());
        std::transform(stagedFiles.begin(), stagedFiles.end(), std::back_inserter(files),
                       [](const auto& entry) { return std::make_pair(entry.path, entry.status); });

        const std::string templ = ManifestParser::generateTemplate(files);

        printError("Commit Rejected: Missing Context Manifest\n");

        std::cerr << kColorYellow << "Detected changes in:" << kColorReset << '\n';
        for (const auto& entry : stagedFiles) {
            std::cerr << "  - " << entry.path << " (" << entry.status << ")" << '\n';
        }

        std::cerr << '\n';
        std::cerr << kColorCyan
                  << "Please retry with this block appended to your commit message:" << kColorReset
                  << '\n';
        std::cerr << '\n';
        std::cerr << templ << '\n';
        std::cerr << '\n';
        std::cerr << kColorYellow << "Or use " << kColorBold << "gip commit -f" << kColorReset
                  << kColorYellow << " to force commit without manifest." << kColorReset << '\n';

        return 1;
    }

    // Validate manifest has required fields
    if (!parseResult.manifest.has_value()) {
        printError("Internal error: Manifest missing despite hasManifest() returning true");
        return 1;
    }

    const auto& manifest = *parseResult.manifest;
    bool valid = true;

    for (const auto& entry : manifest.entries) {
        if (entry.rationale.empty() || entry.rationale.find('<') != std::string::npos) {
            printError("Manifest entry for '" + entry.file + "' has incomplete rationale.");
            valid = false;
        }
        if (entry.behavior.empty() || entry.behavior.find('<') != std::string::npos) {
            printError("Manifest entry for '" + entry.file + "' has incomplete behavior.");
            valid = false;
        }
    }

    if (!valid) {
        printError("Manifest validation failed. Please fill in all <placeholder> fields.");
        return 1;
    }

    // Commit with clean message
    auto result = gitAdapter.commit(parseResult.cleanMessage);
    if (!result.success()) {
        printError("Commit failed: " + result.stderrOutput);
        return 1;
    }

    // Get the new commit SHA
    const std::string commitSha = gitAdapter.getHeadSha();

    // Store manifest in Git Notes
    const std::string toonManifest = manifest.toToon();
    auto noteResult = gitAdapter.addNote(commitSha, toonManifest);

    if (!noteResult.success()) {
        printError("Warning: Failed to store manifest in notes: " + noteResult.stderrOutput);
        // Don't fail the commit, just warn
    }

    printSuccess("Committed with manifest: " + commitSha.substr(0, 7));
    std::cout << result.stdoutOutput << '\n';

    return 0;
}

}  // namespace gip::commands
