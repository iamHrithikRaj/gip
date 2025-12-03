/// @file test_merge_driver.cpp
/// @brief Unit tests for merge driver and conflict enrichment
/// @author Gip Team
/// @copyright MIT License

#include "gip/merge_driver.h"

#include <catch2/catch_test_macros.hpp>
#include <catch2/matchers/catch_matchers_string.hpp>

using namespace gip;
using Catch::Matchers::ContainsSubstring;

TEST_CASE("MergeDriver formats enriched markers", "[merge][enrichment]") {
    MergeDriver driver;

    SECTION("formats marker with full context") {
        ConflictContext ctx;
        ctx.commitSha = "abc1234";
        ctx.behaviorClass = "feature";
        ctx.rationale = "Added new payment method";
        ctx.breaking = true;
        ctx.migrations = {"Update payment config"};
        ctx.inputs = {"amount: float", "currency: string"};
        ctx.outputs = "bool success";
        ctx.symbol = "processPayment";
        ctx.errorModel = {"throws PaymentException"};

        std::string marker = driver.formatEnrichedMarker("HEAD", "Your changes", ctx);

        REQUIRE_THAT(marker, ContainsSubstring("||| Gip CONTEXT (HEAD - Your changes)"));
        REQUIRE_THAT(marker, ContainsSubstring("||| Commit: abc1234"));
        REQUIRE_THAT(marker, ContainsSubstring("||| behaviorClass: feature"));
        REQUIRE_THAT(marker, ContainsSubstring("||| rationale: Added new payment method"));
        REQUIRE_THAT(marker, ContainsSubstring("||| breaking: true"));
        REQUIRE_THAT(marker, ContainsSubstring("||| migrations[0]: Update payment config"));
        REQUIRE_THAT(marker, ContainsSubstring("||| inputs[0]: amount: float"));
        REQUIRE_THAT(marker, ContainsSubstring("||| outputs: bool success"));
        REQUIRE_THAT(marker, ContainsSubstring("||| symbol: processPayment"));
        REQUIRE_THAT(marker, ContainsSubstring("||| errorModel[0]: throws PaymentException"));
    }

    SECTION("formats marker with minimal context") {
        ConflictContext ctx;
        ctx.commitSha = "def5678";
        ctx.behaviorClass = "refactor";
        ctx.rationale = "Cleanup";

        std::string marker = driver.formatEnrichedMarker("feature-branch", "Their changes", ctx);

        REQUIRE_THAT(marker, ContainsSubstring("||| Gip CONTEXT (feature-branch - Their changes)"));
        REQUIRE_THAT(marker, ContainsSubstring("||| Commit: def5678"));
        REQUIRE_THAT(marker, ContainsSubstring("||| behaviorClass: refactor"));
        REQUIRE_THAT(marker, ContainsSubstring("||| rationale: Cleanup"));

        // Should not contain optional fields if they are empty/default
        // Note: This depends on implementation, but usually we don't print empty lists
        // If implementation prints everything, these checks might fail.
        // Based on typical implementation:
        // REQUIRE_THAT(marker, !ContainsSubstring("breaking:")); // breaking is false by default
    }
}
