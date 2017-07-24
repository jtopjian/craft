package aptppa

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_aptPPASourceFileName(t *testing.T) {
	ppa := "chris-lea/redis-server"
	name, err := aptPPASourceFileName(ppa)
	assert.Nil(t, err)

	assert.Equal(t, name, "chris-lea-ubuntu-redis-server-xenial.list", "should be equal")
}

func Test_AptPPA_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()

	ppa := "chris-lea/redis-server"
	exists, err := Exists(client, ppa)
	assert.Nil(t, err)
	assert.Equal(t, exists, false, "should be equal")

	createOpts := CreateOpts{
		Name: ppa,
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, ppa)
	assert.Nil(t, err)
	assert.Equal(t, exists, true, "should be equal")

	err = Delete(client, ppa)
	assert.Nil(t, err)

	exists, err = Exists(client, ppa)
	assert.Nil(t, err)
	assert.Equal(t, exists, false, "should be equal")
}

func Test_AptPPA_List(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()
	ppas, err := List(client)
	assert.Nil(t, err)

	for _, ppa := range ppas {
		t.Logf("%#v", ppa)
	}
}
