/// @file git_adapter.h
/// @brief Git operations adapter
/// @author Hrithik Raj
/// @copyright MIT License

#pragma once

#include "gip/types.h"
#include <memory>

namespace gip {

/// @brief Git adapter - wraps git operations
/// 
/// This class provides a clean interface to git operations,
/// abstracting away the underlying git implementation (shell commands or libgit2).
class GitAdapter {
public:
    /// @brief Constructor
    GitAdapter();
    
    /// @brief Destructor
    ~GitAdapter();
    
    // Non-copyable
    GitAdapter(const GitAdapter&) = delete;
    GitAdapter& operator=(const GitAdapter&) = delete;
    
    // Movable
    GitAdapter(GitAdapter&&) noexcept;
    GitAdapter& operator=(GitAdapter&&) noexcept;
    
    /// @brief Check if current directory is a git repository
    /// @return true if in a git repository
    [[nodiscard]] bool isRepository() const;
    
    /// @brief Get repository root path
    /// @return Absolute path to repository root, empty if not in a repo
    [[nodiscard]] std::string getRepositoryRoot() const;
    
    /// @brief Execute a raw git command (passthrough)
    /// @param args Command arguments (without "git" prefix)
    /// @return Result containing exit code and output
    [[nodiscard]] GitResult execute(const std::vector<std::string>& args) const;
    
    /// @brief Get list of staged files
    /// @return Vector of staged file information
    [[nodiscard]] std::vector<StagedFile> getStagedFiles() const;
    
    /// @brief Get staged diff content
    /// @return Diff content as string
    [[nodiscard]] std::string getStagedDiff() const;
    
    /// @brief Create a commit with the given message
    /// @param message Commit message
    /// @return Result of the commit operation
    [[nodiscard]] GitResult commit(const std::string& message) const;
    
    /// @brief Add a note to a commit
    /// @param commitSha Target commit SHA
    /// @param content Note content
    /// @return Result of the operation
    [[nodiscard]] GitResult addNote(const std::string& commitSha, const std::string& content) const;
    
    /// @brief Get note for a commit
    /// @param commitSha Target commit SHA
    /// @return Note content if exists
    [[nodiscard]] std::optional<std::string> getNote(const std::string& commitSha) const;
    
    /// @brief Get the current HEAD commit SHA
    /// @return Full SHA of HEAD
    [[nodiscard]] std::string getHeadSha() const;
    
    /// @brief Push to remote including notes
    /// @param remote Remote name (default: origin)
    /// @param branch Branch name (default: current branch)
    /// @return Result of the push operation
    [[nodiscard]] GitResult pushWithNotes(const std::string& remote = "origin", 
                                           const std::string& branch = "") const;
    
    /// @brief Get commits that touched a specific file
    /// @param filePath Path to the file
    /// @param limit Maximum number of commits to return
    /// @return Vector of commit contexts with manifests
    [[nodiscard]] std::vector<CommitContext> getFileHistory(const std::string& filePath, 
                                                             int limit = 20) const;
    
    /// @brief Get all tracked files in the repository
    /// @return Vector of file paths
    [[nodiscard]] std::vector<std::string> getTrackedFiles() const;
    
    /// @brief Initialize a new repository
    /// @return Result of the init operation
    [[nodiscard]] GitResult initialize() const;
    
    /// @brief Get current branch name
    /// @return Branch name
    [[nodiscard]] std::string getCurrentBranch() const;

private:
    class Impl;
    std::unique_ptr<Impl> pImpl_;
};

} // namespace gip
