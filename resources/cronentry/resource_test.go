package cronentry

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

var sampleCronEntry = []string{
	"*/5 4 * * * ls # Foo",
	"1 2 3 4 5 pwd # Bar",
}

func Test_cronEntryExists(t *testing.T) {
	exists := cronEntryExists(sampleCronEntry, "Foo")
	assert.Equal(t, exists, true, "should be equal")
}

func Test_CronEntry_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()

	name := "Foo"
	user := "root"

	exists, err := Exists(client, user, name)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")

	createOpts := CreateOpts{
		Name:    name,
		Minute:  "*/5",
		Hour:    "4",
		Command: "ls",
	}

	err = Create(client, user, createOpts)
	assert.Nil(t, err)
	exists, err = Exists(client, user, name)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	updateOpts := UpdateOpts{
		Minute:  "*/6",
		Hour:    "4",
		Command: "ls",
	}

	err = Update(client, user, name, updateOpts)
	assert.Nil(t, err)
	exists, err = Exists(client, user, name)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	err = Delete(client, user, name)
	exists, err = Exists(client, user, name)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")
}
