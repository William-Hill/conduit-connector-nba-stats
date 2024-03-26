package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	nbastats "github.com/William-Hill/conduit-connector-nba-stats"
)

func main() {
	sdk.Serve(nbastats.Connector)
}
