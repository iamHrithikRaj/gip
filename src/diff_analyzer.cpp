#include "diff_analyzer.h"

#include <algorithm>
#include <regex>
#include <sstream>

namespace gip {

auto DiffAnalyzer::detectLanguage(const std::string& filePath) -> std::string
{
    size_t dotPos = filePath.find_last_of('.');
    if (dotPos == std::string::npos) {
        return "unknown";
    }

    std::string ext = filePath.substr(dotPos + 1);
    std::transform(ext.begin(), ext.end(), ext.begin(), ::tolower);

    if (ext == "cpp" || ext == "cc" || ext == "cxx" || ext == "c" ||
        ext == "h" || ext == "hpp") {
        return "cpp";
    } else if (ext == "py") {
        return "python";
    } else if (ext == "js" || ext == "jsx" || ext == "ts" || ext == "tsx") {
        return "javascript";
    } else if (ext == "rs") {
        return "rust";
    } else if (ext == "go") {
        return "go";
    } else if (ext == "java") {
        return "java";
    } else if (ext == "rb") {
        return "ruby";
    } else if (ext == "cs") {
        return "csharp";
    }

    return "unknown";
}

auto DiffAnalyzer::extractSymbolName(const std::string& line,
                                      const std::string& language) -> std::string
{
    std::smatch match;

    if (language == "cpp" || language == "c") {
        // Match function definitions: returnType functionName(...)
        std::regex funcPattern(
            R"((?:[\w:]+\s+)+(\w+)\s*\([^)]*\)\s*(?:const)?\s*(?:\{|$))");
        if (std::regex_search(line, match, funcPattern)) {
            return match[1].str();
        }
        // Match class definitions
        std::regex classPattern(R"(class\s+(\w+))");
        if (std::regex_search(line, match, classPattern)) {
            return match[1].str();
        }
    } else if (language == "python") {
        // Match def functionName(...):
        std::regex funcPattern(R"(def\s+(\w+)\s*\()");
        if (std::regex_search(line, match, funcPattern)) {
            return match[1].str();
        }
        // Match class ClassName:
        std::regex classPattern(R"(class\s+(\w+))");
        if (std::regex_search(line, match, classPattern)) {
            return match[1].str();
        }
    } else if (language == "javascript") {
        // Match function declarations and arrow functions
        std::regex funcPattern(
            R"((?:function\s+(\w+)|(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s*)?\(|(\w+)\s*\([^)]*\)\s*\{))");
        if (std::regex_search(line, match, funcPattern)) {
            for (size_t i = 1; i <= 3; ++i) {
                if (match[i].matched) {
                    return match[i].str();
                }
            }
        }
        // Match class definitions
        std::regex classPattern(R"(class\s+(\w+))");
        if (std::regex_search(line, match, classPattern)) {
            return match[1].str();
        }
    } else if (language == "rust") {
        // Match fn function_name
        std::regex funcPattern(R"(fn\s+(\w+))");
        if (std::regex_search(line, match, funcPattern)) {
            return match[1].str();
        }
        // Match struct/enum/impl
        std::regex structPattern(R"((?:struct|enum|impl)\s+(\w+))");
        if (std::regex_search(line, match, structPattern)) {
            return match[1].str();
        }
    } else if (language == "go") {
        // Match func functionName
        std::regex funcPattern(R"(func\s+(?:\([^)]+\)\s*)?(\w+))");
        if (std::regex_search(line, match, funcPattern)) {
            return match[1].str();
        }
        // Match type TypeName struct
        std::regex typePattern(R"(type\s+(\w+)\s+struct)");
        if (std::regex_search(line, match, typePattern)) {
            return match[1].str();
        }
    } else if (language == "java" || language == "csharp") {
        // Match method/function definitions
        std::regex funcPattern(
            R"((?:public|private|protected|static|\s)+[\w<>\[\]]+\s+(\w+)\s*\()");
        if (std::regex_search(line, match, funcPattern)) {
            return match[1].str();
        }
        // Match class definitions
        std::regex classPattern(R"(class\s+(\w+))");
        if (std::regex_search(line, match, classPattern)) {
            return match[1].str();
        }
    }

    return "";
}

auto DiffAnalyzer::extractSymbols(const std::string& filePath,
                                   const std::string& fileDiff,
                                   const std::string& changeType)
    -> std::vector<SymbolInfo>
{
    std::vector<SymbolInfo> symbols;
    std::string language = detectLanguage(filePath);

    // Parse diff to find changed lines
    std::istringstream stream(fileDiff);
    std::string line;
    std::string currentSymbol;
    int lineNum = 0;

    while (std::getline(stream, line)) {
        // Track line numbers from @@ markers
        if (line.substr(0, 2) == "@@") {
            std::regex linePattern(R"(@@ -\d+(?:,\d+)? \+(\d+))");
            std::smatch match;
            if (std::regex_search(line, match, linePattern)) {
                lineNum = std::stoi(match[1].str());
            }
            continue;
        }

        // Check for symbol definitions in added/modified lines
        if (!line.empty() && (line[0] == '+' || line[0] == ' ')) {
            std::string content = line.substr(1);
            std::string symbol = extractSymbolName(content, language);

            if (!symbol.empty() && symbol != currentSymbol) {
                currentSymbol = symbol;

                SymbolInfo info;
                info.file = filePath;
                info.name = symbol;
                info.type = "function";  // Simplified
                info.changeType = changeType;
                info.startLine = lineNum;
                info.endLine = lineNum;

                symbols.push_back(std::move(info));
            }
        }

        if (!line.empty() && line[0] != '-') {
            lineNum++;
        }
    }

    return symbols;
}

auto DiffAnalyzer::analyze(const std::string& diff) -> std::vector<SymbolInfo>
{
    std::vector<SymbolInfo> allSymbols;

    // Split diff by file
    std::regex filePattern(R"(diff --git a/(.+?) b/(.+?)\n)");
    std::sregex_iterator it(diff.begin(), diff.end(), filePattern);
    std::sregex_iterator end;

    std::vector<std::pair<size_t, std::string>> filePositions;

    while (it != end) {
        filePositions.emplace_back(it->position(), (*it)[2].str());
        ++it;
    }

    // Extract symbols from each file's diff
    for (size_t i = 0; i < filePositions.size(); ++i) {
        size_t start = filePositions[i].first;
        size_t endPos =
            (i + 1 < filePositions.size()) ? filePositions[i + 1].first : diff.size();

        std::string fileDiff = diff.substr(start, endPos - start);
        std::string filePath = filePositions[i].second;

        // Determine change type
        std::string changeType = "modify";
        if (fileDiff.find("new file mode") != std::string::npos) {
            changeType = "add";
        } else if (fileDiff.find("deleted file mode") != std::string::npos) {
            changeType = "delete";
        }

        auto symbols = extractSymbols(filePath, fileDiff, changeType);
        allSymbols.insert(allSymbols.end(), symbols.begin(), symbols.end());
    }

    return allSymbols;
}

auto DiffAnalyzer::getChangedFiles(const std::string& diffStatus)
    -> std::vector<std::pair<std::string, std::string>>
{
    std::vector<std::pair<std::string, std::string>> files;

    std::istringstream stream(diffStatus);
    std::string line;

    while (std::getline(stream, line)) {
        if (line.empty()) {
            continue;
        }

        std::string status = line.substr(0, 1);
        size_t pathStart = line.find_first_not_of(" \t", 1);

        if (pathStart != std::string::npos) {
            std::string path = line.substr(pathStart);

            // Handle renames
            if (status == "R") {
                size_t tabPos = path.find('\t');
                if (tabPos != std::string::npos) {
                    path = path.substr(tabPos + 1);
                }
            }

            files.emplace_back(path, status);
        }
    }

    return files;
}

}  // namespace gip
