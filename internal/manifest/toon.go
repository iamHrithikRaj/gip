package manifest

import (
	"fmt"
	"regexp"
	"strings"
)

// SerializeManifest converts a Manifest to TOON format
func SerializeManifest(m *Manifest) string {
	var sb strings.Builder

	sb.WriteString("; GIP Manifest\n")
	sb.WriteString("(manifest\n")
	sb.WriteString(fmt.Sprintf("  (commit #%s)\n", m.Commit))
	sb.WriteString("  (entries\n")

	for _, entry := range m.Entries {
		sb.WriteString("    (entry\n")

		// Anchor
		sb.WriteString("      (anchor\n")
		sb.WriteString(fmt.Sprintf("        (file %s)\n", entry.Anchor.File))
		sb.WriteString(fmt.Sprintf("        (symbol %s)\n", entry.Anchor.Symbol))
		sb.WriteString(fmt.Sprintf("        (hunk %s))\n", entry.Anchor.HunkID))

		// Change type
		sb.WriteString(fmt.Sprintf("      (changeType %s)\n", entry.ChangeType))

		// Signature delta
		if entry.SignatureDelta.Before != "" || entry.SignatureDelta.After != "" {
			sb.WriteString("      (signatureDelta\n")
			sb.WriteString(fmt.Sprintf("        (before %s)\n", entry.SignatureDelta.Before))
			sb.WriteString(fmt.Sprintf("        (after %s))\n", entry.SignatureDelta.After))
		}

		// Contract
		sb.WriteString("      (contract\n")
		if len(entry.Contract.Preconditions) > 0 {
			sb.WriteString("        (preconditions\n")
			for _, pre := range entry.Contract.Preconditions {
				sb.WriteString(fmt.Sprintf("          [ \"\"\"%s\"\"\" ]\n", pre))
			}
			sb.WriteString("        )\n")
		}
		if len(entry.Contract.Postconditions) > 0 {
			sb.WriteString("        (postconditions\n")
			for _, post := range entry.Contract.Postconditions {
				sb.WriteString(fmt.Sprintf("          [ \"\"\"%s\"\"\" ]\n", post))
			}
			sb.WriteString("        )\n")
		}
		if len(entry.Contract.ErrorModel) > 0 {
			sb.WriteString("        (errorModel\n")
			for _, err := range entry.Contract.ErrorModel {
				sb.WriteString(fmt.Sprintf("          [ \"\"\"%s\"\"\" ]\n", err))
			}
			sb.WriteString("        )\n")
		}
		sb.WriteString("      )\n")

		// Behavior class
		if len(entry.BehaviorClass) > 0 {
			sb.WriteString(fmt.Sprintf("      (behaviorClass [ %s ])\n",
				strings.Join(entry.BehaviorClass, " ")))
		}

		// Side effects
		if len(entry.SideEffects) > 0 {
			sb.WriteString(fmt.Sprintf("      (sideEffects [ %s ])\n",
				strings.Join(entry.SideEffects, " ")))
		}

		// Compatibility
		sb.WriteString("      (compatibility\n")
		sb.WriteString(fmt.Sprintf("        (binaryBreaking %s)\n", boolToTOON(entry.Compatibility.BinaryBreaking)))
		sb.WriteString(fmt.Sprintf("        (sourceBreaking %s)\n", boolToTOON(entry.Compatibility.SourceBreaking)))
		sb.WriteString(fmt.Sprintf("        (dataModelMigration %s))\n", boolToTOON(entry.Compatibility.DataModelMigration)))

		// Tests touched
		if len(entry.TestsTouched) > 0 {
			sb.WriteString(fmt.Sprintf("      (testsTouched [ %s ])\n",
				strings.Join(entry.TestsTouched, " ")))
		}

		// Feature flags
		if len(entry.FeatureFlags) > 0 {
			sb.WriteString(fmt.Sprintf("      (featureFlags [ %s ])\n",
				strings.Join(entry.FeatureFlags, " ")))
		}

		// Rationale
		if entry.Rationale != "" {
			sb.WriteString(fmt.Sprintf("      (rationale \"\"\"%s\"\"\")\n", entry.Rationale))
		}

		sb.WriteString("    )\n")
	}

	sb.WriteString("  )\n")
	sb.WriteString(")\n")

	return sb.String()
}

// ParseManifest parses a TOON manifest string
func ParseManifest(toonStr string) (*Manifest, error) {
	// Use the gotoon library to parse
	// For now, implement a simple regex-based parser to get it working
	m := &Manifest{
		Entries: []Entry{},
	}

	// Extract commit SHA
	commitRe := regexp.MustCompile(`\(commit\s+"([^"]+)"`)
	if matches := commitRe.FindStringSubmatch(toonStr); len(matches) > 1 {
		m.Commit = matches[1]
	}

	// For each entry block, extract the data
	// This is a simplified parser - in production, use proper TOON parsing
	entryRe := regexp.MustCompile(`(?s)\(entry\s+(.*?)\n\s*\)\s*\]`)
	entryMatches := entryRe.FindAllStringSubmatch(toonStr, -1)

	for _, entryMatch := range entryMatches {
		if len(entryMatch) < 2 {
			continue
		}
		entryBlock := entryMatch[1]

		entry := Entry{
			Anchor: Anchor{},
			Contract: Contract{
				Preconditions:  []string{},
				Postconditions: []string{},
				ErrorModel:     []string{},
			},
			BehaviorClass: []string{},
			SideEffects:   []string{},
		}

		// Extract anchor fields
		if matches := regexp.MustCompile(`\(file\s+"([^"]+)"`).FindStringSubmatch(entryBlock); len(matches) > 1 {
			entry.Anchor.File = matches[1]
		}
		if matches := regexp.MustCompile(`\(symbol\s+"([^"]+)"`).FindStringSubmatch(entryBlock); len(matches) > 1 {
			entry.Anchor.Symbol = matches[1]
		}

		// Extract contract fields
		if matches := regexp.MustCompile(`\(precondition\s+"([^"]+)"`).FindStringSubmatch(entryBlock); len(matches) > 1 {
			entry.Contract.Preconditions = []string{matches[1]}
		}
		if matches := regexp.MustCompile(`\(postcondition\s+"([^"]+)"`).FindStringSubmatch(entryBlock); len(matches) > 1 {
			entry.Contract.Postconditions = []string{matches[1]}
		}
		if matches := regexp.MustCompile(`\(error-model\s+"([^"]+)"`).FindStringSubmatch(entryBlock); len(matches) > 1 {
			entry.Contract.ErrorModel = []string{matches[1]}
		}

		// Extract rationale
		if matches := regexp.MustCompile(`\(rationale\s+"([^"]+)"`).FindStringSubmatch(entryBlock); len(matches) > 1 {
			entry.Rationale = matches[1]
		}

		m.Entries = append(m.Entries, entry)
	}

	return m, nil
}

// SerializeForConflict creates a compact TOON block for conflict injection
func SerializeForConflict(entry *Entry, commitSHA string) string {
	var sb strings.Builder

	sb.WriteString("||| TOON (ctx\n")
	sb.WriteString(fmt.Sprintf("|||   (commit #%s)\n", commitSHA))
	sb.WriteString(fmt.Sprintf("|||   (anchor (file %s) (symbol %s))\n",
		entry.Anchor.File, entry.Anchor.Symbol))

	// Compact contract
	if len(entry.Contract.Preconditions) > 0 || len(entry.Contract.Postconditions) > 0 {
		sb.WriteString("|||   (contract\n")
		if len(entry.Contract.Preconditions) > 0 {
			sb.WriteString(fmt.Sprintf("|||     (preconditions [ \"\"\"%s\"\"\" ])\n",
				strings.Join(entry.Contract.Preconditions, "; ")))
		}
		if len(entry.Contract.Postconditions) > 0 {
			sb.WriteString(fmt.Sprintf("|||     (postconditions [ \"\"\"%s\"\"\" ])\n",
				strings.Join(entry.Contract.Postconditions, "; ")))
		}
		if len(entry.Contract.ErrorModel) > 0 {
			sb.WriteString(fmt.Sprintf("|||     (errorModel [ \"\"\"%s\"\"\" ])\n",
				strings.Join(entry.Contract.ErrorModel, "; ")))
		}
		sb.WriteString("|||   )\n")
	}

	if len(entry.BehaviorClass) > 0 {
		sb.WriteString(fmt.Sprintf("|||   (behaviorClass [ %s ])\n",
			strings.Join(entry.BehaviorClass, " ")))
	}

	if entry.Rationale != "" {
		sb.WriteString(fmt.Sprintf("|||   (rationale \"\"\"%s\"\"\")\n", entry.Rationale))
	}

	sb.WriteString("||| )\n")

	return sb.String()
}

func boolToTOON(b bool) string {
	if b {
		return "+"
	}
	return "-"
}
