package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseLSBInfo(t *testing.T) {
	var lsbOutput = `Distributor ID: Ubuntu
		Description:    Ubuntu 16.04.1 LTS
		Release:        16.04
		Codename:       xenial
	`

	lsbInfo := parseLSBInfo(lsbOutput)

	assert.Equal(t, lsbInfo.DistributorID, "Ubuntu", "should be equal")
	assert.Equal(t, lsbInfo.Description, "Ubuntu 16.04.1 LTS", "should be equal")
	assert.Equal(t, lsbInfo.Release, "16.04", "should be equal")
	assert.Equal(t, lsbInfo.Codename, "xenial", "should be equal")
}
