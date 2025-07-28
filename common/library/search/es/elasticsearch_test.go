package es

import (
	"fmt"
	"testing"

	"github.com/star-table/startable-server/common/core/config"
)

func TestGetESClient(t *testing.T) {
	config.LoadUnitTestConfig()

	client, err := GetESClient()
	fmt.Println(err)
	fmt.Printf("%v \n", client)
}
