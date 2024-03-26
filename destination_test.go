package nbastats_test

import (
	"context"
	"testing"

	nbastats "github.com/William-Hill/conduit-connector-nba-stats"
	"github.com/matryer/is"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := nbastats.NewDestination()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
