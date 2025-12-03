#include "manifest.h"

#include <algorithm>
#include <iostream>
#include <regex>
#include <sstream>

namespace gip {

// ============================================================================
// Manifest Implementation
// ============================================================================

ctoon::Value Manifest::toValue() const {
    ctoon::Object root;
    root["schemaVersion"] = ctoon::Value(ctoon::Primitive(schemaVersion));

    ctoon::Array entriesArray;
    // std::cout << "Manifest::toValue: entries.size() = " << entries.size() << std::endl;
    for (const auto& entry : entries) {
        ctoon::Object entryObj;
        entryObj["file"] = ctoon::Value(ctoon::Primitive(entry.file));
        entryObj["symbol"] = ctoon::Value(ctoon::Primitive(entry.symbol));
        entryObj["type"] = ctoon::Value(ctoon::Primitive(entry.type));
        entryObj["behavior"] = ctoon::Value(ctoon::Primitive(entry.behavior));
        entryObj["rationale"] = ctoon::Value(ctoon::Primitive(entry.rationale));
        entryObj["breaking"] = ctoon::Value(ctoon::Primitive(entry.breaking));

        // Migrations array
        ctoon::Array migArray;
        std::transform(entry.migrations.begin(), entry.migrations.end(),
                       std::back_inserter(migArray),
                       [](const std::string& mig) { return ctoon::Value(ctoon::Primitive(mig)); });
        entryObj["migrations"] = ctoon::Value(migArray);

        // Inputs array
        ctoon::Array inputsArray;
        std::transform(entry.inputs.begin(), entry.inputs.end(), std::back_inserter(inputsArray),
                       [](const std::string& inp) { return ctoon::Value(ctoon::Primitive(inp)); });
        entryObj["inputs"] = ctoon::Value(inputsArray);

        entryObj["outputs"] = ctoon::Value(ctoon::Primitive(entry.outputs));

        // Error model array
        ctoon::Array errorArray;
        std::transform(entry.errorModel.begin(), entry.errorModel.end(),
                       std::back_inserter(errorArray),
                       [](const std::string& err) { return ctoon::Value(ctoon::Primitive(err)); });
        entryObj["errorModel"] = ctoon::Value(errorArray);

        // Preconditions array
        ctoon::Array preArray;
        std::transform(entry.preconditions.begin(), entry.preconditions.end(),
                       std::back_inserter(preArray),
                       [](const std::string& pre) { return ctoon::Value(ctoon::Primitive(pre)); });
        entryObj["preconditions"] = ctoon::Value(preArray);

        // Postconditions array
        ctoon::Array postArray;
        std::transform(
            entry.postconditions.begin(), entry.postconditions.end(), std::back_inserter(postArray),
            [](const std::string& post) { return ctoon::Value(ctoon::Primitive(post)); });
        entryObj["postconditions"] = ctoon::Value(postArray);

        // Side effects array
        ctoon::Array sideArray;
        std::transform(
            entry.sideEffects.begin(), entry.sideEffects.end(), std::back_inserter(sideArray),
            [](const std::string& side) { return ctoon::Value(ctoon::Primitive(side)); });
        entryObj["sideEffects"] = ctoon::Value(sideArray);

        entriesArray.push_back(ctoon::Value(entryObj));
    }

    root["entries"] = ctoon::Value(entriesArray);
    return ctoon::Value(root);
}

std::string Manifest::toToon() const {
    // Use JSON format for now to avoid issues with nested arrays in TOON table format
    return ctoon::dumpsJson(toValue());
}

std::optional<Manifest> Manifest::fromValue(const ctoon::Value& value) {
    if (!value.isObject()) {
        return std::nullopt;
    }

    Manifest manifest;
    const auto& obj = value.asObject();

    // Schema version
    auto schemaIt = obj.find("schemaVersion");
    if (schemaIt != obj.end() && schemaIt->second.isPrimitive()) {
        manifest.schemaVersion = schemaIt->second.asPrimitive().asString();
    }

    // Entries
    auto entriesIt = obj.find("entries");
    if (entriesIt == obj.end() || !entriesIt->second.isArray()) {
        return std::nullopt;
    }

    for (const auto& entryVal : entriesIt->second.asArray()) {
        if (!entryVal.isObject()) {
            continue;
        }

        ManifestEntry entry;
        const auto& entryObj = entryVal.asObject();

        auto getStr = [&entryObj](const std::string& key) -> std::string {
            auto it = entryObj.find(key);
            if (it != entryObj.end() && it->second.isPrimitive()) {
                return it->second.asPrimitive().asString();
            }
            return "";
        };

        auto getStrArray = [&entryObj](const std::string& key) -> std::vector<std::string> {
            std::vector<std::string> result;
            auto it = entryObj.find(key);
            if (it != entryObj.end() && it->second.isArray()) {
                for (const auto& v : it->second.asArray()) {
                    if (v.isPrimitive()) {
                        result.push_back(v.asPrimitive().asString());
                    }
                }
            }
            return result;
        };

        auto getBool = [&entryObj](const std::string& key) -> bool {
            auto it = entryObj.find(key);
            if (it != entryObj.end() && it->second.isPrimitive()) {
                return it->second.asPrimitive().getBool();
            }
            return false;
        };

        entry.file = getStr("file");
        entry.symbol = getStr("symbol");
        entry.type = getStr("type");
        entry.behavior = getStr("behavior");
        entry.rationale = getStr("rationale");
        entry.breaking = getBool("breaking");
        entry.migrations = getStrArray("migrations");
        entry.inputs = getStrArray("inputs");
        entry.outputs = getStr("outputs");
        entry.errorModel = getStrArray("errorModel");
        entry.preconditions = getStrArray("preconditions");
        entry.postconditions = getStrArray("postconditions");
        entry.sideEffects = getStrArray("sideEffects");

        manifest.entries.push_back(std::move(entry));
    }

    return manifest;
}

std::optional<Manifest> Manifest::fromToon(const std::string& toonStr) {
    try {
        // Try decoding as TOON first
        auto value = ctoon::decode(toonStr);
        if (value.isObject()) {
            return fromValue(value);
        }
        // If decode returned non-object (e.g. primitive string from JSON), try JSON parser
    } catch (...) {
        // Ignore decode error and try JSON
    }

    try {
        // Fallback to JSON
        auto value = ctoon::loadsJson(toonStr);
        return fromValue(value);
    } catch (...) {
        return std::nullopt;
    }
}

// ============================================================================
// ManifestParser Implementation
// ============================================================================

std::pair<size_t, size_t> ManifestParser::findManifestBlock(const std::string& message) {
    // Look for "gip:"
    size_t start = message.find("gip:");
    if (start == std::string::npos) {
        return {std::string::npos, std::string::npos};
    }

    // Ensure it's at the start of a line or start of string
    if (start > 0 && message[start - 1] != '\n') {
        // Search again? For now, let's assume if it's not at start, it's not valid
        // But we should probably loop to find the right one.
        size_t pos = start;
        while (pos != std::string::npos) {
            pos = message.find("gip:", pos + 1);
            if (pos == std::string::npos)
                break;
            if (message[pos - 1] == '\n') {
                start = pos;
                break;
            }
        }
        if (start > 0 && message[start - 1] != '\n') {
            return {std::string::npos, std::string::npos};
        }
    }

    // Find the end of the TOON block
    size_t pos = start + 4;  // "gip:" length
    int braceDepth = 0;
    int bracketDepth = 0;
    bool inBlock = false;
    bool foundContent = false;

    while (pos < message.size()) {
        char c = message[pos];

        if (c == '{') {
            braceDepth++;
            inBlock = true;
            foundContent = true;
        } else if (c == '}') {
            braceDepth--;
        } else if (c == '[') {
            bracketDepth++;
            inBlock = true;
            foundContent = true;
        } else if (c == ']') {
            bracketDepth--;
        } else if (!isspace(c) && !inBlock && !foundContent) {
            // Found non-whitespace before any block starter?
            // Could be YAML style "key: value"
            foundContent = true;
        }

        // If we've seen content and we're back to zero depth, we're done
        if (inBlock && braceDepth == 0 && bracketDepth == 0) {
            // Find end of current line
            size_t lineEnd = message.find('\n', pos);
            if (lineEnd == std::string::npos) {
                lineEnd = message.size();
            }
            return {start, lineEnd};
        }

        pos++;
    }

    // If no braces found, look for indented block or just take the rest
    // For now, if we didn't find a closed block, let's assume it goes to the end
    // or until a line that looks like a new section (not indented)

    return {start, message.size()};
}

ManifestParser::ParseResult ManifestParser::parse(const std::string& message) {
    ParseResult result;

    auto [start, end] = findManifestBlock(message);

    if (start == std::string::npos) {
        // No manifest block found
        result.cleanMessage = message;
        return result;
    }

    // Extract the manifest content (skip "gip:" prefix)
    std::string manifestContent = message.substr(start + 4, end - start - 4);

    // Trim leading/trailing whitespace
    size_t contentStart = manifestContent.find_first_not_of(" \t\n\r");
    if (contentStart != std::string::npos) {
        manifestContent = manifestContent.substr(contentStart);
    }

    // Try to parse as TOON
    auto manifest = Manifest::fromToon(manifestContent);
    if (!manifest) {
        // Try to parse as simple YAML-like format
        // For simplicity, we'll create a basic manifest from key-value pairs
        Manifest m;
        ManifestEntry entry;

        std::istringstream stream(manifestContent);
        std::string line;

        while (std::getline(stream, line)) {
            // Trim
            size_t lineStart = line.find_first_not_of(" \t-");
            if (lineStart == std::string::npos) {
                continue;
            }
            line = line.substr(lineStart);

            size_t colonPos = line.find(':');
            if (colonPos == std::string::npos) {
                continue;
            }

            std::string key = line.substr(0, colonPos);
            std::string value = line.substr(colonPos + 1);

            // Trim value
            size_t valueStart = value.find_first_not_of(" \t\"");
            size_t valueEnd = value.find_last_not_of(" \t\"\n\r,");
            if (valueStart != std::string::npos && valueEnd != std::string::npos) {
                value = value.substr(valueStart, valueEnd - valueStart + 1);
            }

            if (key == "file") {
                entry.file = value;
            } else if (key == "symbol") {
                entry.symbol = value;
            } else if (key == "type") {
                entry.type = value;
            } else if (key == "behavior") {
                entry.behavior = value;
            } else if (key == "rationale") {
                entry.rationale = value;
            }
        }

        if (!entry.file.empty() || !entry.rationale.empty()) {
            m.entries.push_back(std::move(entry));
            result.manifest = std::move(m);
        }
    } else {
        result.manifest = std::move(manifest);
    }

    // Build clean message (strip manifest block)
    result.cleanMessage = message.substr(0, start);
    if (end < message.size()) {
        result.cleanMessage += message.substr(end);
    }

    // Trim trailing whitespace from clean message
    size_t lastNonSpace = result.cleanMessage.find_last_not_of(" \t\n\r");
    if (lastNonSpace != std::string::npos) {
        result.cleanMessage = result.cleanMessage.substr(0, lastNonSpace + 1);
    }

    return result;
}

std::string
ManifestParser::generateTemplate(const std::vector<std::pair<std::string, std::string>>& files) {
    std::ostringstream ss;

    ss << "gip:\n";
    ss << "{\n";
    ss << "  schemaVersion: \"2.0\",\n";
    ss << "  entries: [\n";

    for (size_t i = 0; i < files.size(); ++i) {
        const auto& [path, status] = files[i];

        std::string type;
        if (status == "A") {
            type = "add";
        } else if (status == "M") {
            type = "modify";
        } else if (status == "D") {
            type = "delete";
        } else if (status == "R") {
            type = "rename";
        } else {
            type = "modify";
        }

        // Extract filename as symbol placeholder
        std::string symbol = path;
        size_t lastSlash = symbol.find_last_of("/\\");
        if (lastSlash != std::string::npos) {
            symbol = symbol.substr(lastSlash + 1);
        }
        size_t dotPos = symbol.find_last_of('.');
        if (dotPos != std::string::npos) {
            symbol.resize(dotPos);
        }

        ss << "    {\n";
        ss << "      file: \"" << path << "\",\n";
        ss << "      symbol: \"" << symbol << "\",\n";
        ss << "      type: \"" << type << "\",\n";
        ss << "      behavior: \"<feature|bugfix|refactor|perf>\",\n";
        ss << "      rationale: \"<explain why this change was made>\",\n";
        ss << "      breaking: false,\n";
        ss << "      migrations: [],\n";
        ss << "      inputs: [],\n";
        ss << "      outputs: \"\",\n";
        ss << "      errorModel: [],\n";
        ss << "      preconditions: [\"<what must be true before>\"],\n";
        ss << "      postconditions: [\"<what is true after>\"]\n";
        ss << "    }";

        if (i < files.size() - 1) {
            ss << ",";
        }
        ss << "\n";
    }

    ss << "  ]\n";
    ss << "}\n";

    return ss.str();
}

std::string ManifestParser::validate(const Manifest& manifest) {
    if (manifest.entries.empty()) {
        return "Manifest has no entries";
    }

    for (size_t i = 0; i < manifest.entries.size(); ++i) {
        const auto& entry = manifest.entries[i];

        if (entry.file.empty()) {
            return "Entry " + std::to_string(i + 1) + ": file path is required";
        }

        if (entry.behavior.empty()) {
            return "Entry " + std::to_string(i + 1) + ": behavior is required";
        }

        // Validate behavior is one of the allowed values
        static const std::vector<std::string> kValidBehaviors = {
            "feature", "bugfix", "refactor", "perf", "security", "docs", "test", "chore"};

        bool validBehavior = false;
        for (const auto& valid : kValidBehaviors) {
            if (entry.behavior == valid) {
                validBehavior = true;
                break;
            }
        }

        if (!validBehavior) {
            return "Entry " + std::to_string(i + 1) + ": invalid behavior '" + entry.behavior +
                   "'. Must be one of: feature, bugfix, refactor, perf, security, docs, test, "
                   "chore";
        }

        if (entry.rationale.empty()) {
            return "Entry " + std::to_string(i + 1) + ": rationale is required";
        }
    }

    return "";  // Valid
}

}  // namespace gip
