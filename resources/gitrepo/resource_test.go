package gitrepo

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_GitRepo_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()
	name := "/tmp/foo"
	source := "https://github.com/wffls/waffles"

	exists, err := Exists(client, name)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")

	createOpts := CreateOpts{
		Name:   name,
		Source: source,
		Branch: "master",
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, name)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	repo, err := Read(client, name)

	assert.Nil(t, err)
	assert.Equal(t, name, repo.Name, "should be equal")
	assert.Equal(t, "master", repo.Branch, "should be equal")

	updateOpts := UpdateOpts{
		Commit: "57290e46a2aed0",
	}

	err = Update(client, name, updateOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, name)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	repo, err = Read(client, name)
	assert.Nil(t, err)
	assert.Equal(t, name, repo.Name, "should be equal")
	assert.Equal(t, "HEAD", repo.Branch, "should be equal")
	assert.Equal(t, "57290e46a2aed06d71b32dbdd2684c9735366f8c", repo.Commit, "should be equal")

	err = Delete(client, name)

	exists, err = Exists(client, name)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")
}
