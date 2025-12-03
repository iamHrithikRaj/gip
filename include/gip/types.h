/// @file types.h
/// @brief Common types and structures for Gip
/// @author Hrithik Raj
/// @copyright MIT License

#pragma once

#include "gip/manifest.h"

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
    /// @return true if exit code is 0
    [[nodiscard]] bool success() const noexcept {
        return exitCode == 0;
    }

    /// @brief Check if command failed
    /// @return true if exit code is non-zero
    [[nodiscard]] bool failed() const noexcept {
        return exitCode != 0;
    }
};

/// @brief Status of a file in git staging area
enum class FileStatus {
    Added,     ///< New file added to staging
    Modified,  ///< Existing file modified
    Deleted,   ///< File deleted
    Renamed,   ///< File renamed
    Copied,    ///< File copied
    Unknown    ///< Unknown status
};

/// @brief Convert file status to string
/// @param status The file status
/// @return Single character status code (A, M, D, R, C, ?)
[[nodiscard]] const char* fileStatusToString(FileStatus status) noexcept;

/// @brief Parse file status from git status code
/// @param code Single character status code
/// @return Parsed FileStatus enum
[[nodiscard]] FileStatus parseFileStatus(char code) noexcept;

/// @brief Information about a staged file
struct StagedFile {
    std::string path;     ///< Relative path to file
    FileStatus status;    ///< Git status of the file
    std::string oldPath;  ///< Original path (for renames)
};

/// @brief Behavior classification for a change
enum class BehaviorClass {
    Feature,   ///< New functionality
    Bugfix,    ///< Bug correction
    Refactor,  ///< Code restructuring without behavior change
    Perf,      ///< Performance optimization
    Security,  ///< Security improvement
    Docs,      ///< Documentation only
    Test,      ///< Test addition or modification
    Chore,     ///< Maintenance tasks
    Unknown    ///< Unclassified
};

/// @brief Convert behavior class to string
/// @param behavior The behavior class
/// @return String representation
[[nodiscard]] const char* behaviorClassToString(BehaviorClass behavior) noexcept;

/// @brief Parse behavior class from string
/// @param str String representation
/// @return Parsed BehaviorClass enum
[[nodiscard]] BehaviorClass parseBehaviorClass(const std::string& str) noexcept;

/// @brief Information about a commit with its manifest
struct CommitContext {
    std::string sha;                   ///< Full commit SHA
    std::string shortSha;              ///< Short SHA (7 chars)
    std::string message;               ///< Commit message
    std::string author;                ///< Author name
    std::string email;                 ///< Author email
    std::string date;                  ///< Commit date (ISO format)
    std::optional<Manifest> manifest;  ///< Parsed manifest (if exists)

    /// @brief Check if commit has a manifest
    /// @return true if manifest exists
    [[nodiscard]] bool hasManifest() const noexcept {
        return manifest.has_value();
    }
};

}  // namespace gip
