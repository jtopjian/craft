package fileini

import (
	"testing"

	"github.com/jtopjian/craft/testhelper"
	"github.com/stretchr/testify/assert"
)

func Test_FileIni_Apply(t *testing.T) {
	client := testhelper.TestClient()
	fileName := "test-fixtures/file.ini"
	sectionName := ""
	keyName := "boolean"

	exists, err := Exists(client, fileName, sectionName, keyName)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	sectionName = "section2"
	keyName = "debug"

	exists, err = Exists(client, fileName, sectionName, keyName)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	sectionName = "section3"
	keyName = "enabled"

	createOpts := CreateOpts{
		FileName: fileName,
		Key:      keyName,
		Value:    "false",
		Section:  sectionName,
	}

	err = Create(client, createOpts)
	assert.Nil(t, err)

	getOpts := GetOpts{
		FileName: fileName,
		Key:      keyName,
		Section:  sectionName,
	}

	_, err = Read(client, getOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName, sectionName, keyName)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	updateOpts := UpdateOpts{
		FileName: fileName,
		Key:      keyName,
		Section:  sectionName,
		Value:    "disabled",
	}

	err = Update(client, updateOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName, sectionName, keyName)
	assert.Nil(t, err)
	assert.Equal(t, true, exists, "should be equal")

	deleteOpts := DeleteOpts{
		FileName: fileName,
		Key:      keyName,
		Section:  sectionName,
	}

	err = Delete(client, deleteOpts)
	assert.Nil(t, err)

	exists, err = Exists(client, fileName, sectionName, keyName)
	assert.Nil(t, err)
	assert.Equal(t, false, exists, "should be equal")

}
