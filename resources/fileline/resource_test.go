package fileline

import (
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_FileLine_Apply(t *testing.T) {
	client := testhelper.TestClient()
	fileName := "test-fixtures/file.txt"
	line := "foo bar baz"

	exists, err := Exists(client, fileName, line)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")

	createOpts := CreateOpts{
		FileName: fileName,
		Line:     line,
		Match:    "^foo",
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName, line)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	deleteOpts := DeleteOpts{
		FileName: fileName,
		Line:     line,
	}

	err = Delete(client, deleteOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName, line)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")

	createOpts.Match = "^-m"
	createOpts.Line = "-m 512"

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName, createOpts.Line)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	createOpts.Match = "^-m"
	createOpts.Line = "-m 256"

	err = Create(client, createOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName, createOpts.Line)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")
}
