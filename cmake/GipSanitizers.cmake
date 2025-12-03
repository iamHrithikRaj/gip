# GipSanitizers.cmake
# Configure sanitizers for debugging and testing

function(gip_enable_sanitizers target_name)
    if(CMAKE_CXX_COMPILER_ID STREQUAL "GNU" OR CMAKE_CXX_COMPILER_ID MATCHES ".*Clang")
        option(GIP_ENABLE_ASAN "Enable Address Sanitizer" OFF)
        option(GIP_ENABLE_UBSAN "Enable Undefined Behavior Sanitizer" OFF)
        option(GIP_ENABLE_TSAN "Enable Thread Sanitizer" OFF)
        option(GIP_ENABLE_MSAN "Enable Memory Sanitizer" OFF)

        set(SANITIZERS "")

        if(GIP_ENABLE_ASAN)
            list(APPEND SANITIZERS "address")
        endif()

        if(GIP_ENABLE_UBSAN)
            list(APPEND SANITIZERS "undefined")
        endif()

        if(GIP_ENABLE_TSAN)
            if(GIP_ENABLE_ASAN OR GIP_ENABLE_MSAN)
                message(WARNING "Thread sanitizer is not compatible with Address or Memory sanitizer")
            else()
                list(APPEND SANITIZERS "thread")
            endif()
        endif()

        if(GIP_ENABLE_MSAN)
            if(GIP_ENABLE_ASAN OR GIP_ENABLE_TSAN)
                message(WARNING "Memory sanitizer is not compatible with Address or Thread sanitizer")
            else()
                list(APPEND SANITIZERS "memory")
            endif()
        endif()

        list(JOIN SANITIZERS "," LIST_OF_SANITIZERS)

        if(LIST_OF_SANITIZERS)
            message(STATUS "Enabling sanitizers: ${LIST_OF_SANITIZERS}")
            target_compile_options(${target_name} PRIVATE -fsanitize=${LIST_OF_SANITIZERS})
            target_link_options(${target_name} PRIVATE -fsanitize=${LIST_OF_SANITIZERS})
        endif()
    elseif(MSVC)
        option(GIP_ENABLE_ASAN "Enable Address Sanitizer" OFF)
        
        if(GIP_ENABLE_ASAN)
            message(STATUS "Enabling MSVC Address Sanitizer")
            target_compile_options(${target_name} PRIVATE /fsanitize=address)
        endif()
    endif()
endfunction()
