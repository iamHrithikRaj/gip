package manifest

// Manifest represents a Gip change manifest for a commit
type Manifest struct {
	SchemaVersion string        `json:"schemaVersion"`          // "1.0" or "2.0"
	Commit        string        `json:"commit"`
	GlobalIntent  *GlobalIntent `json:"globalIntent,omitempty"` // v2.0: commit-level intent
	Entries       []Entry       `json:"entries"`
}

// GlobalIntent represents commit-level rationale for multi-function changes (v2.0)
type GlobalIntent struct {
	BehaviorClass []string `json:"behaviorClass"` // feature|bugfix|refactor|perf|security
	Rationale     string   `json:"rationale"`     // Why this commit exists
}

// Entry represents a single symbol/hunk modification
type Entry struct {
	Anchor               Anchor         `json:"anchor"`
	ChangeType           string         `json:"changeType"`                     // add, modify, delete, rename
	SignatureDelta       SignatureDelta `json:"signatureDelta,omitempty"`       // v2.0: optional, only when signature changes
	Contract             Contract       `json:"contract"`
	BehaviorClass        []string       `json:"behaviorClass"`                  // bugfix, feature, refactor, etc.
	SideEffects          []string       `json:"sideEffects"`                    // logs:channel, reads:file, etc.
	Compatibility        Compatibility  `json:"compatibility,omitempty"`        // v2.0: breaking change markers
	TestsTouched         []string       `json:"testsTouched,omitempty"`
	PerfBudget           *PerfBudget    `json:"perfBudget,omitempty"`
	SecurityNotes        []string       `json:"securityNotes,omitempty"`
	FeatureFlags         []string       `json:"featureFlags,omitempty"`         // v2.0: feature flag tracking
	Rationale            string         `json:"rationale"`                      // Function-specific rationale
	InheritsGlobalIntent bool           `json:"inheritsGlobalIntent,omitempty"` // v2.0: whether to use commit-level intent
}

// Anchor identifies the location of the change
type Anchor struct {
	File   string `json:"file"`
	Symbol string `json:"symbol"`
	HunkID string `json:"hunkId"`
}

// SignatureDelta captures API surface changes
type SignatureDelta struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

// Contract defines the behavioral contract
type Contract struct {
	Inputs         []string `json:"inputs,omitempty"`
	Outputs        string   `json:"outputs,omitempty"`
	Preconditions  []string `json:"preconditions"`
	Postconditions []string `json:"postconditions"`
	ErrorModel     []string `json:"errorModel"`
}

// Compatibility flags (v2.0 enhanced)
type Compatibility struct {
	Breaking      bool     `json:"breaking"`                // Is this a breaking change?
	Deprecations  []string `json:"deprecations,omitempty"`  // What's deprecated? e.g., "old_param removed in v2.0"
	Migrations    []string `json:"migrations,omitempty"`    // Migration guide, e.g., "Replace foo(x) with foo(x, strict=True)"
	
	// v1.0 compatibility (deprecated but kept for backward compat)
	BinaryBreaking     bool `json:"binaryBreaking,omitempty"`
	SourceBreaking     bool `json:"sourceBreaking,omitempty"`
	DataModelMigration bool `json:"dataModelMigration,omitempty"`
}

// PerfBudget captures performance expectations
type PerfBudget struct {
	ExpectedMaxLatencyMs int `json:"expectedMaxLatencyMs,omitempty"`
	CPUDeltaPct          int `json:"cpuDeltaPct,omitempty"`
}

// Schema version constants
const (
	SchemaVersion10 = "1.0" // Initial schema
	SchemaVersion20 = "2.0" // Added globalIntent, signatureDelta, hunk tracking
	SchemaVersionCurrent = SchemaVersion20 // Latest version
)

// Behavior class constants
const (
	BehaviorBugfix     = "bugfix"
	BehaviorFeature    = "feature"
	BehaviorRefactor   = "refactor"
	BehaviorPerf       = "perf"
	BehaviorSecurity   = "security"
	BehaviorValidation = "validation"
	BehaviorDocs       = "docs"
	BehaviorConfig     = "config"
	BehaviorMigration  = "migration"
)

// AllBehaviorClasses returns all valid behavior class options
func AllBehaviorClasses() []string {
	return []string{
		BehaviorBugfix,
		BehaviorFeature,
		BehaviorRefactor,
		BehaviorPerf,
		BehaviorSecurity,
		BehaviorValidation,
		BehaviorDocs,
		BehaviorConfig,
		BehaviorMigration,
	}
}

// Change type constants
const (
	ChangeAdd    = "add"
	ChangeModify = "modify"
	ChangeDelete = "delete"
	ChangeRename = "rename"
)
