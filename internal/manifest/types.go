package manifest

// Manifest represents a GIP change manifest for a commit
type Manifest struct {
	Commit  string   `json:"commit"`
	Entries []Entry  `json:"entries"`
}

// Entry represents a single symbol/hunk modification
type Entry struct {
	Anchor          Anchor          `json:"anchor"`
	ChangeType      string          `json:"changeType"`      // add, modify, delete, rename
	SignatureDelta  SignatureDelta  `json:"signatureDelta"`
	Contract        Contract        `json:"contract"`
	BehaviorClass   []string        `json:"behaviorClass"`   // bugfix, feature, refactor, etc.
	SideEffects     []string        `json:"sideEffects"`     // logs:channel, reads:file, etc.
	Compatibility   Compatibility   `json:"compatibility"`
	TestsTouched    []string        `json:"testsTouched"`
	PerfBudget      *PerfBudget     `json:"perfBudget,omitempty"`
	SecurityNotes   []string        `json:"securityNotes,omitempty"`
	FeatureFlags    []string        `json:"featureFlags,omitempty"`
	Rationale       string          `json:"rationale"`
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

// Compatibility flags
type Compatibility struct {
	BinaryBreaking      bool `json:"binaryBreaking"`
	SourceBreaking      bool `json:"sourceBreaking"`
	DataModelMigration  bool `json:"dataModelMigration"`
}

// PerfBudget captures performance expectations
type PerfBudget struct {
	ExpectedMaxLatencyMs int `json:"expectedMaxLatencyMs,omitempty"`
	CPUDeltaPct          int `json:"cpuDeltaPct,omitempty"`
}

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
