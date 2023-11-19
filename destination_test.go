package nbastats_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	nbastats "github.com/repository/conduit-connector-nba-stats"
)

func TestTeardown_NoOpen(t *testing.T) {
	is := is.New(t)
	con := nbastats.NewDestination()
	err := con.Teardown(context.Background())
	is.NoErr(err)
}
