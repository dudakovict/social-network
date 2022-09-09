package tests

import (
	"fmt"
	"testing"

	"github.com/dudakovict/social-network/business/data/user/dbtest"
	"github.com/dudakovict/social-network/foundation/docker"
)

var esc *docker.Container
var dbc *docker.Container

func TestMain(m *testing.M) {
	var err error
	esc, err = dbtest.StartGRPC()
	if err != nil {
		fmt.Println(err)
		return
	}

	dbc, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer dbtest.StopGRPC(esc)
	defer dbtest.StopDB(dbc)

	m.Run()
}
