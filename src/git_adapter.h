/// @file git_adapter.h
/// @brief Git operations adapter for Gip
/// @author Hrithik Raj
/// @copyright MIT License

#pragma once

#include <optional>
#include <string>
#include <vector>

namespace gip {

/// @brief Result of a git command execution
struct GitResult {
    int exitCode = 0;
    std::string stdoutOutput;
    std::string stderrOutput;

    /// @brief Check if command succeeded
    [[nodiscard]] bool success() const noexcept {
        return exitCode == 0;
    }

    /// @brief Check if command failed
    // cppcheck-suppress unusedFunction
    [[nodiscard]] bool failed() const noexcept {
        return exitCode != 0;
    }
};

/// @brief Information about a staged file
struct StagedFile {
    std::string path;     ///< Relative path to file
    std::string status;   ///< "A" (added), "M" (modified), "D" (deleted), "R" (renamed)
    std::string oldPath;  ///< Original path (for renames)
};

/// @brief Information about a commit with its manifest
struct CommitContext {
    std::string sha;                      ///< Full commit SHA
    std::string shortSha;                 ///< Short SHA (7 chars)
    std::string message;                  ///< Commit message
    std::string author;                   ///< Author name
    std::string date;                     ///< Commit date
    std::optional<std::string> manifest;  ///< TOON-formatted manifest from notes

    /// @brief Check if commit has a manifest
    [[nodiscard]] bool hasManifest() const noexcept {
        return manifest.has_value();
    }
};

/// @brief Git adapter - wraps git operations
///
/// Provides a clean interface to git operations, abstracting away
/// the underlying implementation (shell commands or libgit2).
class GitAdapter {
public:
    GitAdapter();
    ~GitAdapter();

    // Non-copyable, movable
    GitAdapter(const GitAdapter&) = delete;
    GitAdapter& operator=(const GitAdapter&) = delete;
    GitAdapter(GitAdapter&&) noexcept = default;
    GitAdapter& operator=(GitAdapter&&) noexcept = default;

    // =========================================================================
    // Repository State
    // =========================================================================

    /// @brief Check if current directory is a git repository
    /// @return true if in a git repository
    [[nodiscard]] bool isRepository() const;

    /// @brief Get repository root path
    /// @return Absolute path to repository root, empty if not in a repo
    [[nodiscard]] std::string getRepositoryRoot() const;

    /// @brief Get the current HEAD commit SHA
    /// @return Full SHA of HEAD, empty if no commits
    [[nodiscard]] std::string getHeadSha() const;

    /// @brief Get the current branch name
    /// @return Branch name, empty if detached HEAD
    [[nodiscard]] std::string getCurrentBranch() const;

    // =========================================================================
    // Staging Area
    // =========================================================================

    /// @brief Get list of staged files
    /// @return Vector of staged file information
    [[nodiscard]] std::vector<StagedFile> getStagedFiles() const;

    /// @brief Get staged diff content
    /// @return Diff content as string
    [[nodiscard]] std::string getStagedDiff() const;

    /// @brief Get all tracked files in the repository
    /// @return Vector of file paths
    [[nodiscard]] std::vector<std::string> getTrackedFiles() const;

    // =========================================================================
    // Commit Operations
    // =========================================================================

    /// @brief Create a commit with the given message
    /// @param message Commit message
    /// @return Result of the commit operation
    [[nodiscard]] GitResult commit(const std::string& message) const;

    /// @brief Initialize a new repository
    /// @return Result of the init operation
    [[nodiscard]] GitResult initialize() const;

    // =========================================================================
    // Notes Operations
    // =========================================================================

    /// @brief Add a note to a commit
    /// @param commitSha Target commit SHA
    /// @param content Note content
    /// @return Result of the operation
    [[nodiscard]] GitResult addNote(const std::string& commitSha, const std::string& content) const;

    /// @brief Get note for a commit
    /// @param commitSha Target commit SHA
    /// @return Note content if exists
    [[nodiscard]] std::optional<std::string> getNote(const std::string& commitSha) const;

    // =========================================================================
    // Remote Operations
    // =========================================================================

    /// @brief Push to remote including notes
    /// @param remote Remote name (default: "origin")
    /// @param branch Branch name (default: current branch)
    /// @return Result of the push operation
    [[nodiscard]] GitResult pushWithNotes(const std::string& remote = "origin",
                                          const std::string& branch = "") const;

    // =========================================================================
    // History Queries
    // =========================================================================

    /// @brief Get commits that touched a specific file
    /// @param filePath Path to the file
    /// @param limit Maximum number of commits to return (default: 20)
    /// @return Vector of commit contexts with manifests
    [[nodiscard]] std::vector<CommitContext> getFileHistory(const std::string& filePath,
                                                            int limit = 20) const;

    // =========================================================================
    // Raw Execution
    // =========================================================================

    /// @brief Execute a raw git command (passthrough)
    /// @param args Command arguments (without "git" prefix)
    /// @return Result containing exit code and output
    [[nodiscard]] GitResult execute(const std::vector<std::string>& args) const;

private:
    /// @brief Notes reference name for Gip manifests
    static constexpr const char* kNotesRef = "refs/notes/gip";
};

}  // namespace gip
