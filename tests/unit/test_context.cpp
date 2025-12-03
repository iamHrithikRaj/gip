/// @file test_context.cpp
/// @brief Unit tests for context command
/// @author Hrithik Raj
/// @copyright MIT License

#include "gip/types.h"

#include <catch2/catch_test_macros.hpp>
#include <catch2/matchers/catch_matchers_string.hpp>

using namespace gip;
using Catch::Matchers::ContainsSubstring;

TEST_CASE("ContextOptions parsing", "[context][options]") {
    // Note: These tests would require exposing parseContextArgs or testing through the command

    SECTION("behavior class parsing") {
        REQUIRE(parseBehaviorClass("feature") == BehaviorClass::Feature);
        REQUIRE(parseBehaviorClass("bugfix") == BehaviorClass::Bugfix);
        REQUIRE(parseBehaviorClass("refactor") == BehaviorClass::Refactor);
        REQUIRE(parseBehaviorClass("perf") == BehaviorClass::Perf);
        REQUIRE(parseBehaviorClass("security") == BehaviorClass::Security);
        REQUIRE(parseBehaviorClass("unknown") == BehaviorClass::Unknown);
        REQUIRE(parseBehaviorClass("") == BehaviorClass::Unknown);
    }

    SECTION("behavior class to string") {
        REQUIRE(std::string(behaviorClassToString(BehaviorClass::Feature)) == "feature");
        REQUIRE(std::string(behaviorClassToString(BehaviorClass::Bugfix)) == "bugfix");
        REQUIRE(std::string(behaviorClassToString(BehaviorClass::Refactor)) == "refactor");
    }
}

TEST_CASE("FileStatus parsing", "[types]") {
    SECTION("parses git status codes") {
        REQUIRE(parseFileStatus('A') == FileStatus::Added);
        REQUIRE(parseFileStatus('M') == FileStatus::Modified);
        REQUIRE(parseFileStatus('D') == FileStatus::Deleted);
        REQUIRE(parseFileStatus('R') == FileStatus::Renamed);
        REQUIRE(parseFileStatus('C') == FileStatus::Copied);
        REQUIRE(parseFileStatus('?') == FileStatus::Unknown);
        REQUIRE(parseFileStatus('X') == FileStatus::Unknown);
    }

    SECTION("converts to string") {
        REQUIRE(std::string(fileStatusToString(FileStatus::Added)) == "A");
        REQUIRE(std::string(fileStatusToString(FileStatus::Modified)) == "M");
        REQUIRE(std::string(fileStatusToString(FileStatus::Deleted)) == "D");
    }
}

TEST_CASE("CommitContext", "[types]") {
    SECTION("hasManifest returns correct value") {
        CommitContext ctx;
        REQUIRE_FALSE(ctx.hasManifest());

        ctx.manifest = Manifest{};
        REQUIRE(ctx.hasManifest());
    }
}

TEST_CASE("Manifest", "[types]") {
    SECTION("empty check works") {
        Manifest m;
        REQUIRE(m.empty());
        REQUIRE(m.size() == 0);

        m.entries.push_back(ManifestEntry{});
        REQUIRE_FALSE(m.empty());
        REQUIRE(m.size() == 1);
    }
}
