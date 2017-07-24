package file

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_File_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()
	fileName := "/tmp/foo.txt"

	exists, err := Exists(client, fileName)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")

	createOpts := CreateOpts{
		Name:    "/tmp/foo.txt",
		Content: "Hello, World!\n",
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	updateOpts := UpdateOpts{
		Mode:    "0777",
		Content: "Goodbye, World!\n",
	}

	err = Update(client, fileName, updateOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	err = Delete(client, fileName)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")
}
