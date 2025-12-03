# Tests CMakeLists.txt
# Unit and integration tests for Gip

# Find or fetch Catch2
include(FetchContent)

FetchContent_Declare(
    Catch2
    GIT_REPOSITORY https://github.com/catchorg/Catch2.git
    GIT_TAG        v3.4.0
)
FetchContent_MakeAvailable(Catch2)

# Include Catch2 CMake helpers
list(APPEND CMAKE_MODULE_PATH ${catch2_SOURCE_DIR}/extras)
include(CTest)
include(Catch)

# Unit tests executable
add_executable(gip_tests
    unit/test_main.cpp
    unit/test_manifest.cpp
    unit/test_context.cpp
    unit/test_merge_driver.cpp
)

target_link_libraries(gip_tests
    PRIVATE
        gip_lib
        Catch2::Catch2WithMain
)

target_include_directories(gip_tests
    PRIVATE
        ${CMAKE_SOURCE_DIR}/include
        ${CMAKE_SOURCE_DIR}/src
)

# Register tests with CTest
catch_discover_tests(gip_tests)

# E2E Integration Test (PowerShell)
find_program(POWERSHELL_EXE NAMES powershell pwsh)
if(POWERSHELL_EXE)
    add_test(
        NAME E2E_Integration_Test
        COMMAND ${POWERSHELL_EXE} -ExecutionPolicy Bypass -File ${CMAKE_CURRENT_SOURCE_DIR}/integration/test_e2e.ps1 -GipPath $<TARGET_FILE:gip>
    )
    set_tests_properties(E2E_Integration_Test PROPERTIES 
        TIMEOUT 60
        LABELS "integration"
    )
endif()

# Integration tests (if enabled)
option(GIP_BUILD_INTEGRATION_TESTS "Build integration tests" OFF)

if(GIP_BUILD_INTEGRATION_TESTS)
    add_executable(gip_integration_tests
        integration/test_git_operations.cpp
    )
    
    target_link_libraries(gip_integration_tests
        PRIVATE
            gip_lib
            Catch2::Catch2WithMain
    )
    
    catch_discover_tests(gip_integration_tests
        TEST_PREFIX "integration."
    )
endif()
