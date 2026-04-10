// Package build exposes version metadata injected at build time via ldflags.
package build

// Version is the release version string (e.g. "1.0.0").
// Overridden at build time with:
//
//	-ldflags "-X github.com/yummysource/yummycli/internal/build.Version=<ver>"
var Version = "dev"

// Date is the build date in YYYY-MM-DD format.
// Overridden at build time with:
//
//	-ldflags "-X github.com/yummysource/yummycli/internal/build.Date=<date>"
var Date = "unknown"
