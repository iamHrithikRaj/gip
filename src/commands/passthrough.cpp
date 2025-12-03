#include "passthrough.h"

#include "../git_adapter.h"

#include <iostream>
#include <string>
#include <vector>

namespace gip::commands {

auto passthrough(const std::vector<std::string>& args) -> int {
    const GitAdapter git;

    auto result = git.execute(args);

    // Output whatever git outputs
    if (!result.stdoutOutput.empty()) {
        std::cout << result.stdoutOutput;
    }
    if (!result.stderrOutput.empty()) {
        std::cerr << result.stderrOutput;
    }

    return result.exitCode;
}

}  // namespace gip::commands
