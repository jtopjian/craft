package aptkey

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_AptKey_parseList(t *testing.T) {
	var list = `
/etc/apt/trusted.gpg
--------------------
pub   1024D/437D05B5 2004-09-12
uid                  Ubuntu Archive Automatic Signing Key <ftpmaster@ubuntu.com>
sub   2048g/79164387 2004-09-12

pub   4096R/C0B21F32 2012-05-11
uid                  Ubuntu Archive Automatic Signing Key (2012) <ftpmaster@ubuntu.com>
`

	keys := aptKeyParseList(list)
	assert.Equal(t, len(keys), 2, "should be equal")
	t.Logf("%#v", keys)
}

func Test_AptKey_List(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()

	keys, err := List(client)
	assert.Nil(t, err)

	t.Logf("%#v", keys)
}

func Test_AptKey_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()
	keyID := "6026DFCA"

	exists, err := Exists(client, keyID)
	assert.Nil(t, err)
	assert.Equal(t, exists, false, "should be equal")

	createOpts := CreateOpts{
		KeyID:         keyID,
		RemoteKeyFile: "https://www.rabbitmq.com/rabbitmq-release-signing-key.asc",
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, keyID)
	assert.Nil(t, err)
	assert.Equal(t, exists, true, "should be equal")

	err = Delete(client, keyID)
	assert.Nil(t, err)
}
