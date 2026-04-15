package cmdspec

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

// specsFS embeds all JSON files in the specs/ subdirectory at compile time.
// The embed path is relative to this file's directory, which is valid Go embed syntax.
//
// This mirrors lark-cli's embeddedMetaJSON pattern but operates on a directory
// of files rather than a single large JSON blob.
//
//go:embed specs/*.json
var specsFS embed.FS

// LoadAll parses every JSON file in the embedded specs/ directory and returns
// the resulting CapabilitySpec slice. Files that fail to parse are skipped with
// a warning rather than a fatal error so that a bad spec does not break the
// entire CLI.
//
// This mirrors lark-cli's LoadFromMeta / ListFromMetaProjects pattern: the caller
// receives a list of specs and registers commands for each one.
func LoadAll() ([]CapabilitySpec, error) {
	entries, err := fs.ReadDir(specsFS, "specs")
	if err != nil {
		return nil, fmt.Errorf("cmdspec: reading embedded specs dir: %w", err)
	}

	var specs []CapabilitySpec
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := specsFS.ReadFile("specs/" + entry.Name())
		if err != nil {
			// Non-fatal: log and continue so one bad file doesn't break all commands.
			fmt.Printf("cmdspec: warning: could not read %s: %v\n", entry.Name(), err)
			continue
		}

		var spec CapabilitySpec
		if err := json.Unmarshal(data, &spec); err != nil {
			fmt.Printf("cmdspec: warning: could not parse %s: %v\n", entry.Name(), err)
			continue
		}

		specs = append(specs, spec)
	}

	return specs, nil
}
