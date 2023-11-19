package main

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	nbastats "github.com/repository/conduit-connector-nba-stats"
)

func main() {
	sdk.Serve(nbastats.Connector)
}
