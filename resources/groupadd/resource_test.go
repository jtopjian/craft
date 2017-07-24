package groupadd

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_Group_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()
	name := "foobar"

	exists, err := Exists(client, name)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")

	createOpts := CreateOpts{
		Name: name,
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, name)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	err = Delete(client, name)
	exists, err = Exists(client, name)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")
}
