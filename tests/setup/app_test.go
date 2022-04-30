package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTest(t *testing.T) {
	InitTest()
	conf := GetTestAppConfig()
	assert.NotNil(t, conf)
	TruncateTables()
}
