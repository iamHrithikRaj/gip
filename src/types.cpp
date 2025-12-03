#include "gip/types.h"

namespace gip {

const char* fileStatusToString(FileStatus status) noexcept {
    switch (status) {
        case FileStatus::Added: return "A";
        case FileStatus::Modified: return "M";
        case FileStatus::Deleted: return "D";
        case FileStatus::Renamed: return "R";
        case FileStatus::Copied: return "C";
        default: return "?";
    }
}

FileStatus parseFileStatus(char code) noexcept {
    switch (code) {
        case 'A': return FileStatus::Added;
        case 'M': return FileStatus::Modified;
        case 'D': return FileStatus::Deleted;
        case 'R': return FileStatus::Renamed;
        case 'C': return FileStatus::Copied;
        default: return FileStatus::Unknown;
    }
}

const char* behaviorClassToString(BehaviorClass behavior) noexcept {
    switch (behavior) {
        case BehaviorClass::Feature: return "feature";
        case BehaviorClass::Bugfix: return "bugfix";
        case BehaviorClass::Refactor: return "refactor";
        case BehaviorClass::Perf: return "perf";
        case BehaviorClass::Security: return "security";
        case BehaviorClass::Docs: return "docs";
        case BehaviorClass::Test: return "test";
        case BehaviorClass::Chore: return "chore";
        default: return "unknown";
    }
}

BehaviorClass parseBehaviorClass(const std::string& str) noexcept {
    if (str == "feature") return BehaviorClass::Feature;
    if (str == "bugfix") return BehaviorClass::Bugfix;
    if (str == "refactor") return BehaviorClass::Refactor;
    if (str == "perf") return BehaviorClass::Perf;
    if (str == "security") return BehaviorClass::Security;
    if (str == "docs") return BehaviorClass::Docs;
    if (str == "test") return BehaviorClass::Test;
    if (str == "chore") return BehaviorClass::Chore;
    return BehaviorClass::Unknown;
}

} // namespace gip
