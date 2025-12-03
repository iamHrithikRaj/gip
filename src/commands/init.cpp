#include "init.h"

#include "../git_adapter.h"

#include <filesystem>
#include <fstream>
#include <iostream>

namespace fs = std::filesystem;

namespace gip {
namespace commands {

namespace {

// ANSI color codes for terminal output
constexpr const char* kColorGreen = "\033[32m";
constexpr const char* kColorCyan = "\033[36m";
constexpr const char* kColorReset = "\033[0m";

void printSuccess(const std::string& msg)
{
    std::cout << kColorGreen << "[✓] " << msg << kColorReset << std::endl;
}

void printInfo(const std::string& msg)
{
    std::cout << kColorCyan << "[i] " << msg << kColorReset << std::endl;
}

/// Template content for AI assistant instructions
const char* kCopilotInstructions = R"(# Gip - Git with Intent Protocol

This repository uses **Gip** for version control with semantic context.

## For AI Assistants

### Before Editing Files
ALWAYS run `gip context <filename>` to understand the semantic history and constraints of the file.

Example:
```bash
gip context src/tax.cpp
```

This shows you:
- **Intent**: Is this a feature, bugfix, refactor, perf, or security change?
- **Rationale**: Why was each change made?
- **Contracts**: Preconditions and postconditions
- **Side Effects**: Database writes, network calls, etc.

### Querying Context
Use filters to find specific information:
```bash
gip context src/auth.cpp --json          # Machine-readable output
gip context --all --behavior bugfix      # All bugfixes in repo
gip context --since 2025-01-01           # Recent changes
gip context --export knowledge.json      # Export for analysis
```

### When Committing
1. **DO NOT** use `git commit` directly.
2. **ALWAYS** use `gip commit -m "message"`.
3. If `gip` rejects your commit, read the error message carefully - it will provide a manifest template.
4. Append the filled manifest to your commit message.

### Commit Message Format
```
feat: your commit message here

gip:
{
  schemaVersion: "2.0",
  entries: [
    {
      file: "path/to/file",
      symbol: "function_name",
      type: "modify",
      behavior: "feature",
      rationale: "Explain WHY this change was made",
      preconditions: ["What must be true before"],
      postconditions: ["What is true after"],
      sideEffects: ["none"]
    }
  ]
}
```

### Behavior Classes
- `feature` - New functionality
- `bugfix` - Bug correction
- `refactor` - Code restructuring (no behavior change)
- `perf` - Performance optimization
- `security` - Security improvement

### Bypassing (Emergency Only)
Use `gip commit -f -m "message"` to skip manifest requirement.

### Available Commands

| Command | Description |
|---------|-------------|
| `gip context <file>` | View semantic history of a file |
| `gip context <file> --json` | Output as JSON (machine-readable) |
| `gip context --all` | Show context for all tracked files |
| `gip context --behavior <type>` | Filter by behavior class |
| `gip context --since <date>` | Filter commits since date (YYYY-MM-DD) |
| `gip context --export <file>` | Export context to JSON file |
| `gip commit -m "msg"` | Commit with manifest (required) |
| `gip commit -f -m "msg"` | Force commit without manifest |
| `gip push` | Push code AND context notes to remote |
| `gip <anything>` | Passthrough to git |
)";

auto createOrAppendCopilotInstructions(const std::string& repoRoot) -> bool
{
    fs::path githubDir = fs::path(repoRoot) / ".github";
    fs::path instructionsPath = githubDir / "copilot-instructions.md";

    // Create .github directory if it doesn't exist
    if (!fs::exists(githubDir)) {
        fs::create_directories(githubDir);
    }

    // Check if file already exists
    if (fs::exists(instructionsPath)) {
        // Read existing content
        std::ifstream inFile(instructionsPath);
        std::string content((std::istreambuf_iterator<char>(inFile)),
                            std::istreambuf_iterator<char>());
        inFile.close();

        // Check if gip instructions already present
        if (content.find("Gip - Git with Intent Protocol") != std::string::npos) {
            printInfo("Copilot instructions already contain Gip information");
            return true;
        }

        // Append to existing file
        std::ofstream outFile(instructionsPath, std::ios::app);
        outFile << "\n\n---\n\n" << kCopilotInstructions;
        outFile.close();

        printSuccess("Appended Gip instructions to existing copilot-instructions.md");
    } else {
        // Create new file
        std::ofstream outFile(instructionsPath);
        outFile << kCopilotInstructions;
        outFile.close();

        printSuccess("Created .github/copilot-instructions.md");
    }

    return true;
}

}  // anonymous namespace

auto init(const std::vector<std::string>& args) -> int
{
    GitAdapter git;

    // Check if already a git repo
    bool isRepo = git.isRepository();

    if (!isRepo) {
        // Initialize git repo first
        auto result = git.initialize();
        if (!result.success()) {
            std::cerr << "Failed to initialize git repository: " << result.stderrOutput
                      << std::endl;
            return 1;
        }
        printSuccess("Initialized git repository");
    }

    // Forward any additional args to git init
    if (!args.empty() && !isRepo) {
        std::vector<std::string> gitArgs = {"init"};
        gitArgs.insert(gitArgs.end(), args.begin(), args.end());
        git.execute(gitArgs);
    }

    // Get repo root
    std::string repoRoot = git.getRepositoryRoot();
    if (repoRoot.empty()) {
        repoRoot = ".";
    }

    // Create/append copilot instructions
    createOrAppendCopilotInstructions(repoRoot);

    printSuccess("Gip initialized successfully!");
    printInfo("AI assistants will now see instructions in .github/copilot-instructions.md");

    return 0;
}

}  // namespace commands
}  // namespace gip
