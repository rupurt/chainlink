package cmd_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func TestMain(m *testing.M) {
	cltest.SetUpDBAndRunTests(m)
}