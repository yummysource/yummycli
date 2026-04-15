package cmdspec

// CapabilitySpec describes a top-level CLI capability command (e.g., "video").
// One JSON file per capability lives in the specs/ subdirectory. Adding a new
// JSON file automatically registers a new top-level command — no Go changes needed.
//
// This mirrors lark-cli's service-level spec: each service maps to a top-level
// command group, and each resource+method maps to an operation here.
type CapabilitySpec struct {
	// Version is the spec format version (currently "1").
	Version string `json:"version"`

	// Capability is the Cobra command name for this group, e.g. "video".
	Capability string `json:"capability"`

	// Short is the one-line description shown in --help.
	Short string `json:"short"`

	// Operations lists the subcommands available under this capability.
	// Each operation corresponds to a single Cobra leaf command.
	Operations []OperationSpec `json:"operations"`
}

// OperationSpec describes a single subcommand under a capability (e.g., "generate").
// It carries the default provider and the full flag set for that operation.
//
// This mirrors lark-cli's method-level spec: each method becomes one runnable command.
type OperationSpec struct {
	// Use is the Cobra command name, e.g. "generate".
	Use string `json:"use"`

	// Short is the one-line description shown in --help.
	Short string `json:"short"`

	// Provider is the default AI provider for this operation, e.g. "gemini".
	// Users can override via --provider flag if the flag is listed in Flags.
	Provider string `json:"provider"`

	// Flags is the ordered list of CLI flags for this operation.
	// Each flag is bound to a Cobra command by the spec builder.
	Flags []FlagSpec `json:"flags"`
}

// FlagSpec describes a single CLI flag: its name, type, default value, usage text,
// and optional constraints.
//
// This mirrors lark-cli's parameter spec. Unlike lark-cli (which uses --params JSON),
// yummycli registers each flag individually in Cobra so --help shows all options.
type FlagSpec struct {
	// Name is the flag name without the -- prefix, e.g. "aspect-ratio".
	Name string `json:"name"`

	// Type determines which Cobra flag registration method is used.
	// Valid values: "string" | "int" | "bool" | "string-array"
	Type string `json:"type"`

	// Default is the string representation of the default value.
	// For int flags, this is parsed to int at registration time.
	Default string `json:"default"`

	// Usage is the flag description shown in --help.
	Usage string `json:"usage"`

	// Required marks the flag as mandatory; the builder calls MarkFlagRequired.
	Required bool `json:"required"`

	// Enum lists the allowed values for string flags.
	// Validation happens at runtime in the operation's run function,
	// not at the Cobra layer, so --help shows the full enum list in Usage.
	Enum []string `json:"enum,omitempty"`
}
