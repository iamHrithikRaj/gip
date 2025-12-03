/// @file gip.h
/// @brief Gip - Git with Intent Protocol
/// @author Hrithik Raj
/// @copyright MIT License
///
/// Main public header for Gip library.
/// This header provides the public API for applications that want to
/// integrate Gip functionality programmatically.

#pragma once

#include "gip/git_adapter.h"
#include "gip/manifest.h"
#include "gip/types.h"
#include "gip/version.h"

/// @namespace gip
/// @brief Root namespace for Gip library
namespace gip {

/// @brief Initialize Gip library
/// @return true if initialization successful
bool initialize();

/// @brief Shutdown Gip library and cleanup resources
void shutdown();

}  // namespace gip
