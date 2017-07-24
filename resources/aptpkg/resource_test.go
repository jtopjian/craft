package aptpkg

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_aptPkgParseAptCache(t *testing.T) {
	var stdout = `
		sl:
			Installed: (none)
			Candidate: 3.03-17build1
			Version table:
				 3.03-17build1 500
						500 http://nova.clouds.archive.ubuntu.com/ubuntu xenial/universe amd64 Packages
	`

	installed, candidate := aptPkgParseAptCache(stdout)

	assert.Equal(t, installed, "(none)", "should be equal")
	assert.Equal(t, candidate, "3.03-17build1", "should be equal")
}

func Test_AptPkg_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()
	pkgName := "sl"

	createOpts := CreateOpts{
		Name: "sl",
	}

	err := Create(client, createOpts)
	assert.Nil(t, err)

	err = Delete(client, pkgName)
	assert.Nil(t, err)
}

func Test_AptPkg_List(t *testing.T) {
	t.Skip("Skipping")

	client := testhelper.TestClient()
	pkgs, err := List(client)
	assert.Nil(t, err)

	for _, pkg := range pkgs {
		t.Logf("%#v", pkg)
	}
}
