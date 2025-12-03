/// @file version.h
/// @brief Version information for Gip
/// @author Hrithik Raj
/// @copyright MIT License

#pragma once

#include <string>

namespace gip {

/// @brief Major version number
constexpr int VERSION_MAJOR = 1;

/// @brief Minor version number
constexpr int VERSION_MINOR = 0;

/// @brief Patch version number
constexpr int VERSION_PATCH = 0;

/// @brief Version string
constexpr const char* VERSION_STRING = "1.0.0";

/// @brief Full version with build info
/// @return Version string with optional build metadata
std::string getVersionFull();

/// @brief Check if version meets minimum requirement
/// @param major Minimum major version
/// @param minor Minimum minor version
/// @param patch Minimum patch version
/// @return true if current version >= specified version
bool isVersionAtLeast(int major, int minor, int patch);

}  // namespace gip
