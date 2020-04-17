package services_test

import (
	"testing"

	"fmt"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func TestMain(m *testing.M) {
	fmt.Println("THIS SHOULD ONLY HAPPEN ONCE")
	cltest.SetUpDBAndRunTests(m)
}
