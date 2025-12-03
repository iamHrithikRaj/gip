# GipVersion.cmake
# Handles version configuration for Gip

# Read version from VERSION file or use defaults
if(EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/VERSION")
    file(READ "${CMAKE_CURRENT_SOURCE_DIR}/VERSION" GIP_VERSION_STRING)
    string(STRIP "${GIP_VERSION_STRING}" GIP_VERSION_STRING)
else()
    set(GIP_VERSION_STRING "1.0.0")
endif()

# Parse version components
string(REGEX MATCH "^([0-9]+)\\.([0-9]+)\\.([0-9]+)" _ "${GIP_VERSION_STRING}")
set(GIP_VERSION_MAJOR ${CMAKE_MATCH_1})
set(GIP_VERSION_MINOR ${CMAKE_MATCH_2})
set(GIP_VERSION_PATCH ${CMAKE_MATCH_3})

# Full version with optional suffix
set(GIP_VERSION "${GIP_VERSION_MAJOR}.${GIP_VERSION_MINOR}.${GIP_VERSION_PATCH}")

# Get git commit hash for version info
find_package(Git QUIET)
if(GIT_FOUND)
    execute_process(
        COMMAND ${GIT_EXECUTABLE} rev-parse --short HEAD
        WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
        OUTPUT_VARIABLE GIP_GIT_HASH
        OUTPUT_STRIP_TRAILING_WHITESPACE
        ERROR_QUIET
    )
    execute_process(
        COMMAND ${GIT_EXECUTABLE} status --porcelain
        WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
        OUTPUT_VARIABLE GIP_GIT_STATUS
        OUTPUT_STRIP_TRAILING_WHITESPACE
        ERROR_QUIET
    )
    if(GIP_GIT_STATUS)
        set(GIP_GIT_HASH "${GIP_GIT_HASH}-dirty")
    endif()
else()
    set(GIP_GIT_HASH "unknown")
endif()

# Export version info
message(STATUS "Gip Version: ${GIP_VERSION} (${GIP_GIT_HASH})")
