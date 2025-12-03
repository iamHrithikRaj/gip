#include "git_adapter.h"

#include <array>
#include <cstdio>
#include <cstdlib>
#include <filesystem>
#include <fstream>
#include <memory>
#include <sstream>
#include <stdexcept>

#ifdef _WIN32
    #include <windows.h>
    #define popen _popen
    #define pclose _pclose
#endif

namespace gip {

namespace {

/// @brief Execute a command and capture output
/// @param cmd The command to execute
/// @return GitResult with exit code and output
GitResult executeCommand(const std::string& cmd) {
    GitResult result;
    result.exitCode = 0;

    constexpr size_t kBufferSize = 4096;
    std::array<char, kBufferSize> buffer{};
    std::string output;

    // Redirect stderr to stdout for capture
    const std::string fullCmd = cmd + " 2>&1";

    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(fullCmd.c_str(), "r"), pclose);
    if (!pipe) {
        result.exitCode = -1;
        result.stderrOutput = "Failed to execute command";
        return result;
    }

    while (fgets(buffer.data(), static_cast<int>(buffer.size()), pipe.get()) != nullptr) {
        output += buffer.data();
    }

    result.exitCode = pclose(pipe.release());

#ifdef _WIN32
    result.exitCode = result.exitCode & 0xFF;
#else
    if (WIFEXITED(result.exitCode)) {
        result.exitCode = WEXITSTATUS(result.exitCode);
    }
#endif

    if (result.exitCode == 0) {
        result.stdoutOutput = output;
    } else {
        result.stderrOutput = output;
    }

    return result;
}

/// @brief Build a command string from program and arguments
/// @param program The program to execute
/// @param args The arguments to pass
/// @return The constructed command string
std::string buildCommandString(const std::string& program, const std::vector<std::string>& args) {
    std::ostringstream cmd;
    cmd << program;

    for (const auto& arg : args) {
        cmd << " ";
        // Quote arguments with spaces, special characters, or percent signs (Windows)
        const bool needsQuoting = arg.find(' ') != std::string::npos ||
                                  arg.find('"') != std::string::npos ||
                                  arg.find('%') != std::string::npos;
        if (needsQuoting) {
            cmd << "\"";
            for (char c : arg) {
                if (c == '"') {
                    cmd << "\\";
                }
                cmd << c;
            }
            cmd << "\"";
        } else {
            cmd << arg;
        }
    }

    return cmd.str();
}

/// @brief Trim trailing newlines and carriage returns from a string
/// @param str The string to trim
/// @return The trimmed string
std::string trimNewlines(const std::string& str) {
    std::string result = str;
    while (!result.empty() && (result.back() == '\n' || result.back() == '\r')) {
        result.pop_back();
    }
    return result;
}

// Convenience aliases for shorter names
inline GitResult execCommand(const std::string& cmd) {
    return executeCommand(cmd);
}

inline std::string buildCommand(const std::string& program, const std::vector<std::string>& args) {
    return buildCommandString(program, args);
}

}  // anonymous namespace

GitAdapter::GitAdapter() = default;
GitAdapter::~GitAdapter() = default;

bool GitAdapter::isRepository() const {
    auto result = executeCommand("git rev-parse --git-dir");
    return result.success();
}

std::string GitAdapter::getRepositoryRoot() const {
    auto result = executeCommand("git rev-parse --show-toplevel");
    if (!result.success()) {
        return "";
    }
    // Trim trailing newlines
    std::string root = result.stdoutOutput;
    while (!root.empty() && (root.back() == '\n' || root.back() == '\r')) {
        root.pop_back();
    }
    return root;
}

std::string GitAdapter::getCurrentBranch() const {
    auto result = executeCommand("git rev-parse --abbrev-ref HEAD");
    if (!result.success()) {
        return "";
    }
    return trimNewlines(result.stdoutOutput);
}

std::vector<std::string> GitAdapter::getTrackedFiles() const {
    std::vector<std::string> files;
    auto result = executeCommand("git ls-files");
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

GitResult GitAdapter::execute(const std::vector<std::string>& args) const {
    return executeCommand(buildCommandString("git", args));
}

std::vector<StagedFile> GitAdapter::getStagedFiles() const {
    std::vector<StagedFile> files;

    auto result = executeCommand("git diff --cached --name-status");
    if (!result.success()) {
        return files;
    }

    std::istringstream stream(result.stdoutOutput);
    std::string line;

    while (std::getline(stream, line)) {
        if (line.empty()) {
            continue;
        }

        StagedFile file;
        file.status = line.substr(0, 1);

        // Find the file path (after status and tab/spaces)
        const size_t pathStart = line.find_first_not_of(" \t", 1);
        if (pathStart != std::string::npos) {
            file.path = line.substr(pathStart);
            // Handle renames (format: R100\told\tnew)
            if (file.status == "R") {
                const size_t tabPos = file.path.find('\t');
                if (tabPos != std::string::npos) {
                    file.oldPath = file.path.substr(0, tabPos);
                    file.path = file.path.substr(tabPos + 1);  // Use new name
                }
            }
            files.push_back(file);
        }
    }

    return files;
}

std::string GitAdapter::getStagedDiff() const {
    auto result = executeCommand("git diff --cached");
    return result.success() ? result.stdoutOutput : "";
}

GitResult GitAdapter::commit(const std::string& message) const {
    return executeCommand(buildCommandString("git", {"commit", "-m", message}));
}

GitResult GitAdapter::addNote(const std::string& commitSha, const std::string& content) const {
    // Create a temporary file for the note content
    // This avoids issues with newlines and quoting in command line arguments
    std::string tempPath;
    try {
        auto tempDir = std::filesystem::temp_directory_path();
        // Generate a unique filename using a simple random number or timestamp would be better,
        // but commitSha is unique enough for this context.
        // Use a random suffix to avoid collisions if multiple processes work on same commit
        // (unlikely here)
        auto tempFile =
            tempDir / ("gip_note_" + commitSha + "_" + std::to_string(std::rand()) + ".txt");
        tempPath = tempFile.string();

        std::ofstream ofs(tempPath);
        if (!ofs) {
            GitResult result;
            result.exitCode = 1;
            result.stderrOutput = "Failed to create temporary file for note";
            return result;
        }
        ofs << content;
        ofs.close();
    } catch (const std::exception& e) {
        GitResult result;
        result.exitCode = 1;
        result.stderrOutput = std::string("Exception creating temp file: ") + e.what();
        return result;
    }

    auto result = executeCommand(
        buildCommandString("git", {"notes", "--ref=gip", "add", "-f", "-F", tempPath, commitSha}));

    // Cleanup
    try {
        std::filesystem::remove(tempPath);
    } catch (...) {
        // Ignore cleanup errors
    }

    return result;
}

std::optional<std::string> GitAdapter::getNote(const std::string& commitSha) const {
    auto result =
        executeCommand(buildCommandString("git", {"notes", "--ref=gip", "show", commitSha}));

    if (!result.success()) {
        return std::nullopt;
    }

    // Trim trailing newlines which git adds
    std::string content = result.stdoutOutput;
    while (!content.empty() && (content.back() == '\n' || content.back() == '\r')) {
        content.pop_back();
    }

    return content;
}

auto GitAdapter::getHeadSha() const -> std::string {
    auto result = execCommand("git rev-parse HEAD");
    if (!result.success()) {
        return "";
    }

    return trimNewlines(result.stdoutOutput);
}

auto GitAdapter::pushWithNotes(const std::string& remote, const std::string& branch) const
    -> GitResult {
    // First push the branch
    auto result = execCommand(buildCommand("git", {"push", remote, branch}));
    if (!result.success()) {
        return result;
    }

    // Then push the notes
    auto notesResult = execCommand(buildCommand("git", {"push", remote, "refs/notes/gip"}));

    // Combine results
    result.stdoutOutput += "\n" + notesResult.stdoutOutput;
    if (!notesResult.success()) {
        result.stderrOutput += "\n[gip] Warning: Failed to push notes: " + notesResult.stderrOutput;
        // Don't fail the whole push if notes fail (they might not exist yet)
    }

    return result;
}

auto GitAdapter::getFileHistory(const std::string& filePath, int limit) const
    -> std::vector<CommitContext> {
    std::vector<CommitContext> history;

    // Get commits that touched this file
    auto result = execCommand(buildCommand(
        "git", {"log", "--format=%H|%s|%an|%ai", "-n", std::to_string(limit), "--", filePath}));

    if (!result.success()) {
        return history;
    }

    std::istringstream stream(result.stdoutOutput);
    std::string line;

    while (std::getline(stream, line)) {
        if (line.empty()) {
            continue;
        }

        CommitContext ctx;

        // Parse format: SHA|message|author|date
        size_t pos1 = line.find('|');
        size_t pos2 = line.find('|', pos1 + 1);
        size_t pos3 = line.find('|', pos2 + 1);

        if (pos1 == std::string::npos || pos2 == std::string::npos || pos3 == std::string::npos) {
            continue;
        }

        ctx.sha = line.substr(0, pos1);
        ctx.shortSha = ctx.sha.substr(0, 7);
        ctx.message = line.substr(pos1 + 1, pos2 - pos1 - 1);
        ctx.author = line.substr(pos2 + 1, pos3 - pos2 - 1);
        ctx.date = line.substr(pos3 + 1);

        // Try to get the manifest note
        ctx.manifest = getNote(ctx.sha);

        history.push_back(std::move(ctx));
    }

    return history;
}

auto GitAdapter::initialize() const -> GitResult {
    return execCommand("git init");
}

}  // namespace gip
