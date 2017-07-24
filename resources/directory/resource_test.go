package directory

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_Directory_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()
	dirName := "/tmp/foo"

	exists, err := Exists(client, dirName)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")

	createOpts := CreateOpts{
		Name: dirName,
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)
	exists, err = Exists(client, dirName)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	updateOpts := UpdateOpts{
		Mode: "0644",
	}

	err = Update(client, dirName, updateOpts)
	assert.Nil(t, err)
	exists, err = Exists(client, dirName)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	err = Delete(client, dirName, false)
	assert.Nil(t, err)
	exists, err = Exists(client, dirName)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")
}
