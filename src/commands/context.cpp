#include "context.h"

#include "../git_adapter.h"
#include "../manifest.h"

#include <algorithm>
#include <fstream>
#include <iomanip>
#include <iostream>
#include <sstream>

namespace gip {
namespace commands {

namespace {

// ANSI color codes for terminal output
constexpr const char* kColorRed = "\033[31m";
constexpr const char* kColorGreen = "\033[32m";
constexpr const char* kColorYellow = "\033[33m";
constexpr const char* kColorCyan = "\033[36m";
constexpr const char* kColorMagenta = "\033[35m";
constexpr const char* kColorReset = "\033[0m";
constexpr const char* kColorBold = "\033[1m";
constexpr const char* kColorDim = "\033[2m";

void printError(const std::string& msg) {
    std::cerr << kColorRed << "[!] " << msg << kColorReset << std::endl;
}

void printSuccess(const std::string& msg) {
    std::cout << kColorGreen << "[✓] " << msg << kColorReset << std::endl;
}

void printHeader(const std::string& filePath) {
    std::cout << std::endl;
    std::cout << kColorBold << kColorCyan
              << "═══════════════════════════════════════════════════════════════" << kColorReset
              << std::endl;
    std::cout << kColorBold << "  Gip Context Report: " << kColorReset << filePath << std::endl;
    std::cout << kColorCyan << "═══════════════════════════════════════════════════════════════"
              << kColorReset << std::endl;
    std::cout << std::endl;
}

void printCommitContext(const CommitContext& ctx) {
    std::cout << kColorYellow << "┌─ Commit " << kColorBold << ctx.sha.substr(0, 7) << kColorReset;
    std::cout << kColorDim << " (" << ctx.date.substr(0, 10) << " by " << ctx.author << ")"
              << kColorReset << std::endl;
    std::cout << kColorYellow << "│" << kColorReset << std::endl;
    std::cout << kColorYellow << "│  " << kColorReset << kColorBold << ctx.message << kColorReset
              << std::endl;

    if (ctx.manifest) {
        auto manifest = Manifest::fromToon(*ctx.manifest);

        if (manifest && !manifest->entries.empty()) {
            std::cout << kColorYellow << "│" << kColorReset << std::endl;

            std::for_each(
                manifest->entries.begin(), manifest->entries.end(), [](const auto& entry) {
                    // Behavior
                    if (!entry.behavior.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorMagenta << "Intent: " << kColorReset << entry.behavior
                                  << std::endl;
                    }

                    // Rationale
                    if (!entry.rationale.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorGreen << "Rationale: " << kColorReset << entry.rationale
                                  << std::endl;
                    }

                    // Breaking
                    if (entry.breaking) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorRed << "BREAKING CHANGE" << kColorReset << std::endl;
                    }

                    // Migrations
                    if (!entry.migrations.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorRed << "Migrations: " << kColorReset;
                        for (size_t i = 0; i < entry.migrations.size(); ++i) {
                            if (i > 0) {
                                std::cout << ", ";
                            }
                            std::cout << entry.migrations[i];
                        }
                        std::cout << std::endl;
                    }

                    // Inputs
                    if (!entry.inputs.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorCyan << "Inputs: " << kColorReset;
                        for (size_t i = 0; i < entry.inputs.size(); ++i) {
                            if (i > 0) {
                                std::cout << ", ";
                            }
                            std::cout << entry.inputs[i];
                        }
                        std::cout << std::endl;
                    }

                    // Outputs
                    if (!entry.outputs.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorCyan << "Outputs: " << kColorReset << entry.outputs
                                  << std::endl;
                    }

                    // Error Model
                    if (!entry.errorModel.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorRed << "Error Model: " << kColorReset;
                        for (size_t i = 0; i < entry.errorModel.size(); ++i) {
                            if (i > 0) {
                                std::cout << ", ";
                            }
                            std::cout << entry.errorModel[i];
                        }
                        std::cout << std::endl;
                    }

                    // Preconditions
                    if (!entry.preconditions.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorCyan << "Preconditions: " << kColorReset;
                        for (size_t i = 0; i < entry.preconditions.size(); ++i) {
                            if (i > 0) {
                                std::cout << ", ";
                            }
                            std::cout << entry.preconditions[i];
                        }
                        std::cout << std::endl;
                    }

                    // Postconditions
                    if (!entry.postconditions.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorCyan << "Postconditions: " << kColorReset;
                        for (size_t i = 0; i < entry.postconditions.size(); ++i) {
                            if (i > 0) {
                                std::cout << ", ";
                            }
                            std::cout << entry.postconditions[i];
                        }
                        std::cout << std::endl;
                    }

                    // Side effects
                    if (!entry.sideEffects.empty()) {
                        std::cout << kColorYellow << "│  " << kColorReset;
                        std::cout << kColorRed << "Side Effects: " << kColorReset;
                        for (size_t i = 0; i < entry.sideEffects.size(); ++i) {
                            if (i > 0) {
                                std::cout << ", ";
                            }
                            std::cout << entry.sideEffects[i];
                        }
                        std::cout << std::endl;
                    }
                });
        }
    } else {
        std::cout << kColorYellow << "│  " << kColorReset;
        std::cout << kColorDim << "(no manifest)" << kColorReset << std::endl;
    }

    std::cout << kColorYellow << "└───────────────────────────────────────────────────────────────"
              << kColorReset << std::endl;
    std::cout << std::endl;
}

/// Escape a string for JSON output
auto escapeJson(const std::string& str) -> std::string {
    std::ostringstream oss;
    for (char c : str) {
        switch (c) {
            case '"':
                oss << "\\\"";
                break;
            case '\\':
                oss << "\\\\";
                break;
            case '\n':
                oss << "\\n";
                break;
            case '\r':
                oss << "\\r";
                break;
            case '\t':
                oss << "\\t";
                break;
            default:
                oss << c;
        }
    }
    return oss.str();
}

/// Convert a vector of strings to JSON array string
auto vectorToJsonArray(const std::vector<std::string>& vec) -> std::string {
    std::ostringstream oss;
    oss << "[";
    for (size_t i = 0; i < vec.size(); ++i) {
        if (i > 0) {
            oss << ", ";
        }
        oss << "\"" << escapeJson(vec[i]) << "\"";
    }
    oss << "]";
    return oss.str();
}

/// Convert a single commit context to JSON object
auto commitToJson(const CommitContext& ctx, const std::string& indent = "    ") -> std::string {
    std::ostringstream oss;
    oss << indent << "{\n";
    oss << indent << "  \"sha\": \"" << ctx.sha << "\",\n";
    oss << indent << "  \"shortSha\": \"" << ctx.sha.substr(0, 7) << "\",\n";
    oss << indent << "  \"message\": \"" << escapeJson(ctx.message) << "\",\n";
    oss << indent << "  \"author\": \"" << escapeJson(ctx.author) << "\",\n";
    oss << indent << "  \"date\": \"" << ctx.date << "\",\n";
    oss << indent << "  \"hasManifest\": " << (ctx.manifest ? "true" : "false");

    if (ctx.manifest) {
        auto manifest = Manifest::fromToon(*ctx.manifest);
        if (manifest && !manifest->entries.empty()) {
            oss << ",\n" << indent << "  \"manifest\": {\n";
            oss << indent << "    \"schemaVersion\": \"" << manifest->schemaVersion << "\",\n";
            oss << indent << "    \"entries\": [\n";

            for (size_t i = 0; i < manifest->entries.size(); ++i) {
                const auto& entry = manifest->entries[i];
                oss << indent << "      {\n";
                oss << indent << "        \"file\": \"" << escapeJson(entry.file) << "\",\n";
                oss << indent << "        \"symbol\": \"" << escapeJson(entry.symbol) << "\",\n";
                oss << indent << "        \"type\": \"" << escapeJson(entry.type) << "\",\n";
                oss << indent << "        \"behavior\": \"" << escapeJson(entry.behavior)
                    << "\",\n";
                oss << indent << "        \"rationale\": \"" << escapeJson(entry.rationale)
                    << "\",\n";
                oss << indent << "        \"breaking\": " << (entry.breaking ? "true" : "false")
                    << ",\n";
                oss << indent << "        \"migrations\": " << vectorToJsonArray(entry.migrations)
                    << ",\n";
                oss << indent << "        \"inputs\": " << vectorToJsonArray(entry.inputs) << ",\n";
                oss << indent << "        \"outputs\": \"" << escapeJson(entry.outputs) << "\",\n";
                oss << indent << "        \"errorModel\": " << vectorToJsonArray(entry.errorModel)
                    << ",\n";
                oss << indent
                    << "        \"preconditions\": " << vectorToJsonArray(entry.preconditions)
                    << ",\n";
                oss << indent
                    << "        \"postconditions\": " << vectorToJsonArray(entry.postconditions)
                    << ",\n";
                oss << indent << "        \"sideEffects\": " << vectorToJsonArray(entry.sideEffects)
                    << "\n";
                oss << indent << "      }";
                if (i < manifest->entries.size() - 1) {
                    oss << ",";
                }
                oss << "\n";
            }

            oss << indent << "    ]\n";
            oss << indent << "  }";
        }
    }

    oss << "\n" << indent << "}";
    return oss.str();
}

/// Generate JSON output for file history
std::string historyToJson(const std::string& filePath, const std::vector<CommitContext>& history) {
    std::ostringstream oss;
    oss << "{\n";
    oss << "  \"file\": \"" << escapeJson(filePath) << "\",\n";
    oss << "  \"generatedAt\": \"" << __DATE__ << " " << __TIME__ << "\",\n";
    oss << "  \"totalCommits\": " << history.size() << ",\n";

    int withManifest = std::count_if(history.begin(), history.end(),
                                     [](const auto& ctx) { return ctx.manifest.has_value(); });

    oss << "  \"commitsWithManifest\": " << withManifest << ",\n";
    oss << "  \"commits\": [\n";

    for (size_t i = 0; i < history.size(); ++i) {
        oss << commitToJson(history[i]);
        if (i < history.size() - 1) {
            oss << ",";
        }
        oss << "\n";
    }

    oss << "  ]\n";
    oss << "}\n";

    return oss.str();
}

/// Generate JSON for all files context
std::string
allFilesToJson(const std::vector<std::pair<std::string, std::vector<CommitContext>>>& allHistory) {
    std::ostringstream oss;
    oss << "{\n";
    oss << "  \"generatedAt\": \"" << __DATE__ << " " << __TIME__ << "\",\n";
    oss << "  \"totalFiles\": " << allHistory.size() << ",\n";
    oss << "  \"files\": [\n";

    for (size_t f = 0; f < allHistory.size(); ++f) {
        const auto& [filePath, history] = allHistory[f];
        oss << "    {\n";
        oss << "      \"file\": \"" << escapeJson(filePath) << "\",\n";
        oss << "      \"totalCommits\": " << history.size() << ",\n";
        oss << "      \"commits\": [\n";

        for (size_t i = 0; i < history.size(); ++i) {
            oss << commitToJson(history[i], "        ");
            if (i < history.size() - 1) {
                oss << ",";
            }
            oss << "\n";
        }

        oss << "      ]\n";
        oss << "    }";
        if (f < allHistory.size() - 1) {
            oss << ",";
        }
        oss << "\n";
    }

    oss << "  ]\n";
    oss << "}\n";

    return oss.str();
}

/// Check if a commit matches the behavior filter
auto matchesBehaviorFilter(const CommitContext& ctx, const std::string& behaviorFilter) -> bool {
    if (behaviorFilter.empty()) {
        return true;
    }
    if (!ctx.manifest) {
        return false;
    }

    auto manifest = Manifest::fromToon(*ctx.manifest);
    if (!manifest) {
        return false;
    }

    for (const auto& entry : manifest->entries) {
        if (entry.behavior == behaviorFilter) {
            return true;
        }
    }
    return false;
}

/// Check if a commit is after the given date (YYYY-MM-DD format)
auto isAfterDate(const std::string& commitDate, const std::string& sinceDate) -> bool {
    if (sinceDate.empty()) {
        return true;
    }
    // Simple string comparison works for ISO date format (YYYY-MM-DD)
    return commitDate.substr(0, 10) >= sinceDate;
}

/// Get all tracked files in the repo
auto getAllTrackedFiles(const GitAdapter& git) -> std::vector<std::string> {
    std::vector<std::string> files;

    auto result = git.execute({"ls-files"});
    if (!result.success()) {
        return files;
    }

    std::istringstream stream(result.stdoutOutput);
    std::string line;
    while (std::getline(stream, line)) {
        if (!line.empty()) {
            files.push_back(line);
        }
    }

    return files;
}

void printUsage() {
    std::cerr << "Usage: gip context [OPTIONS] [<filename>]\n\n";
    std::cerr << "Shows the semantic history of a file with manifest context.\n\n";
    std::cerr << "OPTIONS:\n";
    std::cerr << "  --json              Output as JSON (machine-readable)\n";
    std::cerr << "  --export <file>     Export JSON to specified file\n";
    std::cerr << "  --all               Show context for all tracked files\n";
    std::cerr
        << "  --behavior <type>   Filter by behavior (feature, bugfix, refactor, perf, security)\n";
    std::cerr << "  --since <date>      Show commits since date (YYYY-MM-DD)\n";
    std::cerr << "  -h, --help          Show this help\n\n";
    std::cerr << "EXAMPLES:\n";
    std::cerr << "  gip context src/main.cpp                    Terminal output for a file\n";
    std::cerr << "  gip context src/main.cpp --json             JSON output for a file\n";
    std::cerr << "  gip context --all --export knowledge.json   Export all context to file\n";
    std::cerr << "  gip context --behavior bugfix               Show only bugfix commits\n";
    std::cerr << "  gip context --since 2025-01-01              Show commits since date\n";
}

}  // anonymous namespace

auto parseContextArgs(const std::vector<std::string>& args) -> ContextOptions {
    ContextOptions opts;

    for (size_t i = 0; i < args.size(); ++i) {
        const std::string& arg = args[i];

        if (arg == "--json") {
            opts.jsonOutput = true;
        } else if (arg == "--export" && i + 1 < args.size()) {
            opts.exportPath = args[++i];
            opts.jsonOutput = true;  // Export implies JSON
        } else if (arg == "--all") {
            opts.showAll = true;
        } else if (arg == "--behavior" && i + 1 < args.size()) {
            opts.behaviorFilter = args[++i];
        } else if (arg == "--since" && i + 1 < args.size()) {
            opts.sinceDate = args[++i];
        } else if (arg == "-h" || arg == "--help") {
            // Will be handled in main function
        } else if (!arg.empty() && arg[0] != '-') {
            opts.filePath = arg;
        }
    }

    return opts;
}

auto context(const std::vector<std::string>& args) -> int {
    // Check for help flag
    if (std::any_of(args.begin(), args.end(),
                    [](const std::string& arg) { return arg == "-h" || arg == "--help"; })) {
        printUsage();
        return 0;
    }

    ContextOptions opts = parseContextArgs(args);

    // Validate: need either a file path or --all
    if (opts.filePath.empty() && !opts.showAll) {
        printError("Usage: gip context <filename> or gip context --all");
        std::cerr << "\nRun 'gip context --help' for more options.\n";
        return 1;
    }

    GitAdapter git;

    if (!git.isRepository()) {
        printError("Not a git repository");
        return 1;
    }

    // Handle --all mode
    if (opts.showAll) {
        auto files = getAllTrackedFiles(git);

        if (files.empty()) {
            printError("No tracked files found");
            return 1;
        }

        std::vector<std::pair<std::string, std::vector<CommitContext>>> allHistory;

        for (const auto& file : files) {
            auto history = git.getFileHistory(file, 50);

            // Apply filters
            std::vector<CommitContext> filtered;
            std::copy_if(history.begin(), history.end(), std::back_inserter(filtered),
                         [&](const auto& ctx) {
                             return matchesBehaviorFilter(ctx, opts.behaviorFilter) &&
                                    isAfterDate(ctx.date, opts.sinceDate);
                         });

            if (!filtered.empty()) {
                allHistory.emplace_back(file, filtered);
            }
        }

        if (opts.jsonOutput) {
            std::string json = allFilesToJson(allHistory);

            if (!opts.exportPath.empty()) {
                std::ofstream outFile(opts.exportPath);
                if (!outFile) {
                    printError("Failed to create file: " + opts.exportPath);
                    return 1;
                }
                outFile << json;
                outFile.close();
                printSuccess("Exported context for " + std::to_string(allHistory.size()) +
                             " files to: " + opts.exportPath);
            } else {
                std::cout << json;
            }
        } else {
            // Terminal output for --all
            for (const auto& [file, history] : allHistory) {
                printHeader(file);
                for (const auto& ctx : history) {
                    printCommitContext(ctx);
                }
            }

            std::cout << kColorDim
                      << "───────────────────────────────────────────────────────────────"
                      << kColorReset << std::endl;
            std::cout << "Showing context for " << allHistory.size() << " files" << std::endl;
        }

        return 0;
    }

    // Single file mode
    auto history = git.getFileHistory(opts.filePath, 20);

    if (history.empty()) {
        printError("No commits found for: " + opts.filePath);
        std::cerr << "Make sure the file path is correct and has been committed." << std::endl;
        return 1;
    }

    // Apply filters
    std::vector<CommitContext> filtered;
    std::copy_if(history.begin(), history.end(), std::back_inserter(filtered),
                 [&](const auto& ctx) {
                     return matchesBehaviorFilter(ctx, opts.behaviorFilter) &&
                            isAfterDate(ctx.date, opts.sinceDate);
                 });

    if (filtered.empty()) {
        printError("No commits match the specified filters");
        return 1;
    }

    if (opts.jsonOutput) {
        std::string json = historyToJson(opts.filePath, filtered);

        if (!opts.exportPath.empty()) {
            std::ofstream outFile(opts.exportPath);
            if (!outFile) {
                printError("Failed to create file: " + opts.exportPath);
                return 1;
            }
            outFile << json;
            outFile.close();
            printSuccess("Exported to: " + opts.exportPath);
        } else {
            std::cout << json;
        }
    } else {
        // Terminal output
        printHeader(opts.filePath);

        int withManifest = 0;
        int withoutManifest = 0;

        for (const auto& ctx : filtered) {
            printCommitContext(ctx);

            if (ctx.manifest) {
                withManifest++;
            } else {
                withoutManifest++;
            }
        }

        // Summary
        std::cout << kColorDim << "───────────────────────────────────────────────────────────────"
                  << kColorReset << std::endl;
        std::cout << "Showing " << filtered.size() << " commits";
        std::cout << " (" << kColorGreen << withManifest << " with manifest" << kColorReset;
        std::cout << ", " << kColorDim << withoutManifest << " without" << kColorReset << ")"
                  << std::endl;

        if (!opts.behaviorFilter.empty()) {
            std::cout << "Filtered by behavior: " << kColorMagenta << opts.behaviorFilter
                      << kColorReset << std::endl;
        }
        if (!opts.sinceDate.empty()) {
            std::cout << "Filtered since: " << kColorCyan << opts.sinceDate << kColorReset
                      << std::endl;
        }
        std::cout << std::endl;
    }

    return 0;
}

}  // namespace commands
}  // namespace gip
