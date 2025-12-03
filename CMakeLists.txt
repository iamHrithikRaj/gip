cmake_minimum_required(VERSION 3.15)

# =============================================================================
# Project Configuration
# =============================================================================
project(gip
    VERSION 1.0.0
    DESCRIPTION "Git with Intent Protocol - LLM-native version control"
    HOMEPAGE_URL "https://github.com/iamHrithikRaj/gip"
    LANGUAGES CXX
)

# Prevent in-source builds
if(CMAKE_SOURCE_DIR STREQUAL CMAKE_BINARY_DIR)
    message(FATAL_ERROR "In-source builds are not allowed. Please use a separate build directory.")
endif()

# =============================================================================
# CMake Configuration
# =============================================================================
list(APPEND CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/cmake")

# Include custom CMake modules
include(GipVersion)
include(GipCompilerWarnings)
include(GipSanitizers)

# Set C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# Export compile commands for tooling
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Build type
if(NOT CMAKE_BUILD_TYPE AND NOT CMAKE_CONFIGURATION_TYPES)
    message(STATUS "Setting build type to 'Release' as none was specified.")
    set(CMAKE_BUILD_TYPE Release CACHE STRING "Choose the type of build." FORCE)
    set_property(CACHE CMAKE_BUILD_TYPE PROPERTY STRINGS "Debug" "Release" "MinSizeRel" "RelWithDebInfo")
endif()

# =============================================================================
# Build Options
# =============================================================================
option(GIP_BUILD_TESTS "Build unit tests" OFF)
option(GIP_BUILD_DOCS "Build documentation" OFF)
option(GIP_ENABLE_WARNINGS "Enable compiler warnings" ON)
option(GIP_ENABLE_CLANG_TIDY "Enable clang-tidy analysis" OFF)
option(GIP_ENABLE_CPPCHECK "Enable cppcheck analysis" OFF)
option(GIP_USE_LIBGIT2 "Use libgit2 instead of git CLI" ON)

# =============================================================================
# Dependencies
# =============================================================================
include(FetchContent)

# libgit2 (optional, for native git operations)
if(GIP_USE_LIBGIT2)
    FetchContent_Declare(
        libgit2
        GIT_REPOSITORY https://github.com/libgit2/libgit2.git
        GIT_TAG        v1.7.1
    )
    
    # Disable libgit2 extras we don't need
    set(BUILD_SHARED_LIBS OFF CACHE BOOL "" FORCE)
    set(BUILD_TESTS OFF CACHE BOOL "" FORCE)
    set(BUILD_CLI OFF CACHE BOOL "" FORCE)
    
    FetchContent_MakeAvailable(libgit2)
endif()

# ctoon (TOON format parser)
set(BUILD_EXAMPLES OFF CACHE BOOL "" FORCE)
set(CTOON_BUILD_TESTS OFF CACHE BOOL "" FORCE)
add_subdirectory(ctoon)

# =============================================================================
# Library Target (for testing and potential embedding)
# =============================================================================
set(GIP_LIB_SOURCES
    src/git_adapter.cpp
    src/manifest.cpp
    src/types.cpp
    src/diff_analyzer.cpp
    src/merge_driver.cpp
    src/commands/commit.cpp
    src/commands/init.cpp
    src/commands/context.cpp
    src/commands/push.cpp
    src/commands/passthrough.cpp
    src/commands/merge.cpp
    src/commands/rebase.cpp
)

add_library(gip_lib STATIC ${GIP_LIB_SOURCES})

target_include_directories(gip_lib
    PUBLIC
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:include>
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src
        ${CMAKE_CURRENT_SOURCE_DIR}/ctoon/includes
)

target_link_libraries(gip_lib
    PUBLIC
        ctoon
    PRIVATE
        $<$<BOOL:${GIP_USE_LIBGIT2}>:libgit2>
)

target_compile_definitions(gip_lib
    PRIVATE
        $<$<BOOL:${GIP_USE_LIBGIT2}>:GIP_USE_LIBGIT2>
        GIP_VERSION="${GIP_VERSION}"
        GIP_VERSION_MAJOR=${GIP_VERSION_MAJOR}
        GIP_VERSION_MINOR=${GIP_VERSION_MINOR}
        GIP_VERSION_PATCH=${GIP_VERSION_PATCH}
)

if(GIP_ENABLE_WARNINGS)
    gip_set_project_warnings(gip_lib)
endif()

gip_enable_sanitizers(gip_lib)

# =============================================================================
# Executable Target
# =============================================================================
add_executable(gitp src/main.cpp)

target_link_libraries(gitp PRIVATE gip_lib)

target_include_directories(gitp PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${CMAKE_CURRENT_SOURCE_DIR}/src
)

if(GIP_ENABLE_WARNINGS)
    gip_set_project_warnings(gitp)
endif()

gip_enable_sanitizers(gitp)

# =============================================================================
# Static Analysis
# =============================================================================
if(GIP_ENABLE_CLANG_TIDY)
    include(GipStaticAnalysis)
    gip_enable_clang_tidy(gip_lib)
    gip_enable_clang_tidy(gitp)
endif()

if(GIP_ENABLE_CPPCHECK)
    include(GipStaticAnalysis)
    gip_enable_cppcheck(gip_lib)
    gip_enable_cppcheck(gitp)
endif()

# =============================================================================
# Tests
# =============================================================================
if(GIP_BUILD_TESTS)
    enable_testing()
    add_subdirectory(tests)
endif()

# =============================================================================
# Documentation
# =============================================================================
if(GIP_BUILD_DOCS)
    find_package(Doxygen REQUIRED)
    
    set(DOXYGEN_GENERATE_HTML YES)
    set(DOXYGEN_GENERATE_MAN NO)
    set(DOXYGEN_OUTPUT_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}/docs)
    set(DOXYGEN_EXTRACT_ALL YES)
    set(DOXYGEN_EXTRACT_PRIVATE YES)
    set(DOXYGEN_USE_MDFILE_AS_MAINPAGE README.md)
    
    doxygen_add_docs(docs
        ${CMAKE_CURRENT_SOURCE_DIR}/include
        ${CMAKE_CURRENT_SOURCE_DIR}/src
        ${CMAKE_CURRENT_SOURCE_DIR}/README.md
        COMMENT "Generating API documentation"
    )
endif()

# =============================================================================
# Installation
# =============================================================================
include(GNUInstallDirs)

install(TARGETS gitp
    RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)

install(TARGETS gip_lib
    ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
)

install(DIRECTORY include/
    DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}
)

# =============================================================================
# Package Configuration
# =============================================================================
include(CMakePackageConfigHelpers)

write_basic_package_version_file(
    "${CMAKE_CURRENT_BINARY_DIR}/gipConfigVersion.cmake"
    VERSION ${PROJECT_VERSION}
    COMPATIBILITY SameMajorVersion
)

configure_package_config_file(
    "${CMAKE_CURRENT_SOURCE_DIR}/cmake/gipConfig.cmake.in"
    "${CMAKE_CURRENT_BINARY_DIR}/gipConfig.cmake"
    INSTALL_DESTINATION ${CMAKE_INSTALL_LIBDIR}/cmake/gip
)

install(FILES
    "${CMAKE_CURRENT_BINARY_DIR}/gipConfig.cmake"
    "${CMAKE_CURRENT_BINARY_DIR}/gipConfigVersion.cmake"
    DESTINATION ${CMAKE_INSTALL_LIBDIR}/cmake/gip
)

# =============================================================================
# Summary
# =============================================================================
message(STATUS "")
message(STATUS "=== Gip Configuration Summary ===")
message(STATUS "Version:          ${GIP_VERSION} (${GIP_GIT_HASH})")
message(STATUS "Build type:       ${CMAKE_BUILD_TYPE}")
message(STATUS "C++ Standard:     ${CMAKE_CXX_STANDARD}")
message(STATUS "Compiler:         ${CMAKE_CXX_COMPILER_ID} ${CMAKE_CXX_COMPILER_VERSION}")
message(STATUS "")
message(STATUS "Options:")
message(STATUS "  Tests:          ${GIP_BUILD_TESTS}")
message(STATUS "  Documentation:  ${GIP_BUILD_DOCS}")
message(STATUS "  Warnings:       ${GIP_ENABLE_WARNINGS}")
message(STATUS "  clang-tidy:     ${GIP_ENABLE_CLANG_TIDY}")
message(STATUS "  Use libgit2:    ${GIP_USE_LIBGIT2}")
message(STATUS "")

