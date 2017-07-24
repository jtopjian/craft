package aptsource

import (
	"os"
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_aptSourceBuildFile(t *testing.T) {
	expected := "deb http://www.rabbitmq.com/debian/ testing main"

	e := entry{
		URI:          "http://www.rabbitmq.com/debian/",
		Distribution: "testing",
		Component:    "main",
	}

	actual := aptSourceBuildEntry(e, false)

	assert.Equal(t, expected, actual, "should be equal")

	expected = "deb-src http://www.rabbitmq.com/debian/ testing main"
	actual = aptSourceBuildEntry(e, true)

	assert.Equal(t, expected, actual, "should be equal")
}

func Test_aptSourceParseFile(t *testing.T) {
	expected := entry{
		URI:          "http://www.rabbitmq.com/debian/",
		Distribution: "testing",
		Component:    "main",
	}

	e := "deb http://www.rabbitmq.com/debian/ testing main"

	actual, err := aptSourceParseEntry(e)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual, "should be equal")

	expected.Source = true
	e = "deb-src http://www.rabbitmq.com/debian/ testing main"
	actual, err = aptSourceParseEntry(e)
	assert.Nil(t, err)

	assert.Equal(t, expected, actual, "should be equal")
}

func Test_AptSource_Apply(t *testing.T) {
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()
	name := "rabbitmq"

	exists, err := Exists(client, name, false)
	assert.Nil(t, err)
	assert.Equal(t, exists, false, "should be equal")

	createOpts := CreateOpts{
		Name:         "rabbitmq",
		URI:          "http://www.rabbitmq.com/debian/",
		Distribution: "testing",
		Component:    "main",
		IncludeSrc:   true,
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, name, true)
	assert.Nil(t, err)
	assert.Equal(t, exists, true, "should be equal")

	err = Delete(client, name)
	exists, err = Exists(client, name, true)
	assert.Nil(t, err)
	assert.Equal(t, exists, false, "should be equal")
}

func Test_AptSource_Test(t *testing.T) {
	t.Skip("Skipping")
	acc := os.Getenv("TEST_ACC")
	if acc == "" {
		t.Skip("TEST_ACC is not set. Skipping")
	}

	client := testhelper.TestClient()

	sources, err := List(client)
	assert.Nil(t, err)

	for _, source := range sources {
		t.Logf("%#v", source)
	}
}
