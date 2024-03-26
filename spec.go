package nbastats

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// version is set during the build process with ldflags (see Makefile).
// Default version matches default from runtime/debug.
var version = "v0.1.0"

// Specification returns the connector's specification.
func Specification() sdk.Specification {
	return sdk.Specification{
		Name:        "nba-stats",
		Summary:     "<describe your connector>",
		Description: "<describe your connector in detail>",
		Version:     version,
		Author:      "<your name>",
	}
}
