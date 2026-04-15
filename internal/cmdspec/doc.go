// Package cmdspec provides a lightweight JSON-driven capability spec system
// for dynamically registering Cobra commands.
//
// This is a minimal adaptation of lark-cli's metadata-based command approach,
// scaled to yummycli's needs. Each JSON file in the specs/ subdirectory describes
// one capability (e.g., "video") with its operations and flags. Adding a new
// capability requires only a new JSON file — no Go code changes.
//
// The mapping mirrors lark-cli's architecture:
//
//	lark-cli: services → resources → methods → parameters
//	yummycli: capabilities → operations → flags
package cmdspec
