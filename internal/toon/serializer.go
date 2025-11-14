// Package toon handles serialization of GIP manifests to TOON (Tree Object Notation)
// format for human-readable display in Git conflict markers.
//
// TOON format uses tree-like ASCII art to represent structured data in a compact,
// readable way that fits naturally within conflict markers.
package toon

import (
	"fmt"
	"strings"

	"github.com/alpkeskin/gotoon"
	"github.com/iamHrithikRaj/gip/internal/manifest"
)

// SerializeManifest converts a Manifest to TOON format
func SerializeManifest(m *manifest.Manifest) (string, error) {
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

	return sb.String(), nil
}

// ParseManifest parses a TOON manifest string
func ParseManifest(toonStr string) (*manifest.Manifest, error) {
	// TODO: Implement TOON parsing using gotoon library
	// For now, return a placeholder
	return &manifest.Manifest{
		Commit:  "parsed",
		Entries: []manifest.Entry{},
	}, nil
}

// SerializeForConflict creates a compact TOON block for conflict injection
func SerializeForConflict(entry *manifest.Entry, commitSHA string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("||| TOON (ctx\n"))
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

// UseGoTOON demonstrates using the gotoon library for encoding/decoding
func UseGoTOON(data interface{}) (string, error) {
	// Example of using gotoon for encoding
	encoded, err := gotoon.Encode(data)
	if err != nil {
		return "", err
	}
	return encoded, nil
}
