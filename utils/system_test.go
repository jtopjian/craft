package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UsernameToID(t *testing.T) {
	uid, err := UsernameToID("root")
	assert.Nil(t, err)
	assert.Equal(t, uid, 0, "should be equal")
}

func Test_GroupToID(t *testing.T) {
	gid, err := GroupToID("root")
	assert.Nil(t, err)
	assert.Equal(t, gid, 0, "should be equal")
}
