//+build integration

package database

import (
	"chatapp/pkg/util"
	"log"
	"os"
	"testing"
)

var (
	testConfig        *util.Config
	testMySQLDBSource string
)

func TestMain(m *testing.M) {
	var err error

	testConfig, err = util.ReadConfig(util.GetAbsolutePath())
	if err != nil {
		log.Fatal(err)
	}

	testMySQLDBSource = testConfig.DBConfig.MySQL.DBSource

	os.Exit(m.Run())
}
