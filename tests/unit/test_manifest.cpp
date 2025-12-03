/// @file test_manifest.cpp
/// @brief Unit tests for manifest parsing and serialization
/// @author Hrithik Raj
/// @copyright MIT License

#include <catch2/catch_test_macros.hpp>
#include <catch2/matchers/catch_matchers_string.hpp>

#include "gip/manifest.h"
#include "gip/types.h"

using namespace gip;
using Catch::Matchers::ContainsSubstring;

TEST_CASE("ManifestParser parses valid commit messages", "[manifest][parser]") {
    SECTION("extracts manifest block from message") {
        const std::string message = R"(feat: add tax calculation

gip:
{
  schemaVersion: "2.0",
  entries: [
    {
      file: "src/tax.cpp",
      symbol: "calculate_tax",
      behavior: "feature",
      rationale: "Added 8% state tax for CA compliance"
    }
  ]
})";

        auto result = ManifestParser::parse(message);
        
        REQUIRE(result.hasManifest());
        REQUIRE_FALSE(result.hasError());
        REQUIRE(result.cleanMessage == "feat: add tax calculation");
        REQUIRE(result.manifest->entries.size() == 1);
        REQUIRE(result.manifest->entries[0].file == "src/tax.cpp");
        REQUIRE(result.manifest->entries[0].symbol == "calculate_tax");
    }
    
    SECTION("handles message without manifest") {
        const std::string message = "fix: simple typo correction";
        
        auto result = ManifestParser::parse(message);
        
        REQUIRE_FALSE(result.hasManifest());
        REQUIRE_FALSE(result.hasError());
        REQUIRE(result.cleanMessage == message);
    }
    
    SECTION("handles empty message") {
        auto result = ManifestParser::parse("");
        
        REQUIRE_FALSE(result.hasManifest());
        REQUIRE(result.cleanMessage.empty());
    }
}

TEST_CASE("ManifestParser generates valid templates", "[manifest][template]") {
    SECTION("generates template for single file") {
        std::vector<std::pair<std::string, std::string>> files = {
            {"src/main.cpp", "M"}
        };
        
        auto tmpl = ManifestParser::generateTemplate(files);
        
        REQUIRE_THAT(tmpl, ContainsSubstring("gip:"));
        REQUIRE_THAT(tmpl, ContainsSubstring("src/main.cpp"));
        REQUIRE_THAT(tmpl, ContainsSubstring("schemaVersion"));
    }
    
    SECTION("generates template for multiple files") {
        std::vector<std::pair<std::string, std::string>> files = {
            {"src/main.cpp", "M"},
            {"src/utils.cpp", "A"},
            {"src/old.cpp", "D"}
        };
        
        auto tmpl = ManifestParser::generateTemplate(files);
        
        REQUIRE_THAT(tmpl, ContainsSubstring("src/main.cpp"));
        REQUIRE_THAT(tmpl, ContainsSubstring("src/utils.cpp"));
        REQUIRE_THAT(tmpl, ContainsSubstring("src/old.cpp"));
    }
}

TEST_CASE("Manifest serialization handles TOON format", "[manifest][serializer]") {
    SECTION("round-trips manifest through TOON") {
        Manifest original;
        original.schemaVersion = "2.0";
        
        ManifestEntry entry;
        entry.file = "src/test.cpp";
        entry.symbol = "test_function";
        entry.behavior = "feature";
        entry.rationale = "Test rationale";
        entry.preconditions = {"input >= 0"};
        entry.postconditions = {"output > input"};
        
        // New fields
        entry.breaking = true;
        entry.migrations = {"Run migration script"};
        entry.inputs = {"int x", "int y"};
        entry.outputs = "int result";
        entry.errorModel = {"throws std::runtime_error"};
        
        original.entries.push_back(entry);
        
        auto toon = original.toToon();
        auto parsed = Manifest::fromToon(toon);
        
        REQUIRE(parsed.has_value());
        REQUIRE(parsed->schemaVersion == original.schemaVersion);
        REQUIRE(parsed->entries.size() == 1);
        REQUIRE(parsed->entries[0].file == entry.file);
        REQUIRE(parsed->entries[0].rationale == entry.rationale);
        
        // Verify new fields
        REQUIRE(parsed->entries[0].breaking == true);
        REQUIRE(parsed->entries[0].migrations.size() == 1);
        REQUIRE(parsed->entries[0].migrations[0] == "Run migration script");
        REQUIRE(parsed->entries[0].inputs.size() == 2);
        REQUIRE(parsed->entries[0].inputs[0] == "int x");
        REQUIRE(parsed->entries[0].outputs == "int result");
        REQUIRE(parsed->entries[0].errorModel.size() == 1);
        REQUIRE(parsed->entries[0].errorModel[0] == "throws std::runtime_error");
    }
}
