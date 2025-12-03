/// @file console.h
/// @brief Console output utilities with ANSI color support
/// @author Hrithik Raj
/// @copyright MIT License

#pragma once

#include <iostream>
#include <string>
#include <string_view>

namespace gip {
namespace console {

/// @brief ANSI color codes for terminal output
namespace colors {

/// @brief Check if terminal supports colors
[[nodiscard]] bool isColorSupported() noexcept;

/// @brief Enable or disable color output globally
void setColorEnabled(bool enabled) noexcept;

/// @brief Check if color output is enabled
[[nodiscard]] bool isColorEnabled() noexcept;

// ANSI color constants
inline constexpr std::string_view kReset = "\033[0m";
inline constexpr std::string_view kBold = "\033[1m";
inline constexpr std::string_view kDim = "\033[2m";
inline constexpr std::string_view kRed = "\033[31m";
inline constexpr std::string_view kGreen = "\033[32m";
inline constexpr std::string_view kYellow = "\033[33m";
inline constexpr std::string_view kBlue = "\033[34m";
inline constexpr std::string_view kMagenta = "\033[35m";
inline constexpr std::string_view kCyan = "\033[36m";
inline constexpr std::string_view kWhite = "\033[37m";

}  // namespace colors

/// @brief Log levels for console output
enum class LogLevel {
    kDebug,
    kInfo,
    kWarning,
    kError,
    kSuccess
};

/// @brief Print a message with the specified log level
/// @param level The log level
/// @param message The message to print
void log(LogLevel level, std::string_view message);

/// @brief Print an error message
/// @param message The error message
inline void error(std::string_view message) {
    log(LogLevel::kError, message);
}

/// @brief Print a warning message
/// @param message The warning message
inline void warning(std::string_view message) {
    log(LogLevel::kWarning, message);
}

/// @brief Print a success message
/// @param message The success message
inline void success(std::string_view message) {
    log(LogLevel::kSuccess, message);
}

/// @brief Print an info message
/// @param message The info message
inline void info(std::string_view message) {
    log(LogLevel::kInfo, message);
}

/// @brief Print a debug message (only in debug builds)
/// @param message The debug message
inline void debug([[maybe_unused]] std::string_view message) {
#ifndef NDEBUG
    log(LogLevel::kDebug, message);
#endif
}

/// @brief Print a horizontal separator line
/// @param width The width of the line
/// @param character The character to use
void printSeparator(int width = 65, char character = '─');

/// @brief Print a header with a title
/// @param title The title to display
void printHeader(std::string_view title);

}  // namespace console
}  // namespace gip
