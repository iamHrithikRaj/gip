# GipStaticAnalysis.cmake
# Configure static analysis tools

function(gip_enable_clang_tidy target_name)
    find_program(CLANG_TIDY clang-tidy)
    if(CLANG_TIDY)
        message(STATUS "Found clang-tidy: ${CLANG_TIDY}")
        set_target_properties(${target_name}
            PROPERTIES
                CXX_CLANG_TIDY "${CLANG_TIDY};-checks=*,-llvmlibc-*,-fuchsia-*,-altera-*"
        )
    else()
        message(STATUS "clang-tidy not found, skipping static analysis")
    endif()
endfunction()

function(gip_enable_cppcheck target_name)
    find_program(CPPCHECK cppcheck)
    if(CPPCHECK)
        message(STATUS "Found cppcheck: ${CPPCHECK}")
        set_target_properties(${target_name}
            PROPERTIES
                CXX_CPPCHECK "${CPPCHECK};--enable=warning,performance,portability;--suppress=missingIncludeSystem"
        )
    else()
        message(STATUS "cppcheck not found, skipping")
    endif()
endfunction()

function(gip_enable_iwyu target_name)
    find_program(IWYU include-what-you-use)
    if(IWYU)
        message(STATUS "Found include-what-you-use: ${IWYU}")
        set_target_properties(${target_name}
            PROPERTIES
                CXX_INCLUDE_WHAT_YOU_USE ${IWYU}
        )
    else()
        message(STATUS "include-what-you-use not found, skipping")
    endif()
endfunction()
