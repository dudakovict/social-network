package tests

import (
	"fmt"
	"testing"

	"github.com/dudakovict/social-network/business/data/user/dbtest"
	"github.com/dudakovict/social-network/foundation/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}
