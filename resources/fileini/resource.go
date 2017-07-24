package fileini

import (
	"fmt"

	"github.com/go-ini/ini"
	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "FileIni"

// FileIni represents an ini file entry on a system.
type FileIni struct {
	// FileName is the file the entry belongs to.
	FileName string

	// Key is the key of the entry.
	Key string

	// Value is the value of the entry.
	Value string

	// Section is the ini section the entry belongs to.
	Section string
}

// GetOpts represents options to read an entry in an ini file.
type GetOpts struct {
	// FileName is the name of the ini file.
	FileName string `required:"true"`

	// Key is the key of the entry.
	Key string `required:"true"`

	// Section is the section of the entry.
	Section string
}

// CreateOpts represents options to create an entry in an ini file.
type CreateOpts struct {
	// FileName is the name of the ini file.
	FileName string `required:"true"`

	// Key is the key of the entry.
	Key string `required:"true"`

	// Value is the value of the entry.
	Value string

	// Section is the section of the entry.
	Section string
}

// UpdateOpts represents options to update an entry in an ini file.
type UpdateOpts struct {
	// FileName is the name of the ini file.
	FileName string `required:"true"`

	// Key is the key of the entry.
	Key string `required:"true"`

	// Value is the value of the entry.
	Value string

	// Section is the section of the entry.
	Section string
}

// DeleteOpts represents options to delete an entry in an ini file.
type DeleteOpts struct {
	// FileName is the name of the ini file.
	FileName string `required:"true"`

	// Key is the key of the entry.
	Key string `required:"true"`

	// Section is the section of the entry.
	Section string
}

// Read will read an existing ini entry.
func Read(client client.Client, getOpts GetOpts) (entry FileIni, err error) {
	client.Logger.Debugf("Reading FileIni entry")

	if err = utils.BuildRequest(&getOpts); err != nil {
		return
	}

	client.Logger.Debugf("FileIni Read Options: %#v", getOpts)

	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, getOpts.FileName)
	if err != nil {
		return
	}

	section, err := cfg.GetSection(getOpts.Section)
	if err != nil {
		err = nil
	}

	exists := section.HasKey(getOpts.Key)
	if !exists {
		resourceTitle := fmt.Sprintf("%s/%s/%s", getOpts.FileName, getOpts.Section, getOpts.Key)
		err = resources.NotFoundError{Type: Type, Name: resourceTitle}
		return
	}

	value := section.Key(getOpts.Key).String()
	if err != nil {
		return
	}

	entry.FileName = getOpts.FileName
	entry.Section = getOpts.Section
	entry.Key = getOpts.Key
	entry.Value = value

	return
}

// Exists will determine if an ini entry exists.
func Exists(client client.Client, fileName, sectionName, keyName string) (exists bool, err error) {
	client.Logger.Debugf("Checking if entry %s in section %s of file %s exists",
		keyName, sectionName, fileName)

	getOpts := GetOpts{
		FileName: fileName,
		Section:  sectionName,
		Key:      keyName,
	}

	_, err = Read(client, getOpts)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	exists = true

	return
}

// List isn't implemented.
func List(client.Client) {
	return
}

// Create will create an entry in an ini file.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	client.Logger.Debugf("Creating FileIni entry")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, createOpts.FileName)
	if err != nil {
		return
	}

	section, err := cfg.GetSection(createOpts.Section)
	if err != nil {
		section, err = cfg.NewSection(createOpts.Section)
		if err != nil {
			return
		}
	}

	if createOpts.Value == "" {
		_, err = section.NewBooleanKey(createOpts.Key)
		if err != nil {
			return
		}
	}

	if createOpts.Value != "" {
		_, err = section.NewKey(createOpts.Key, createOpts.Value)
		if err != nil {
			return
		}
	}

	err = cfg.SaveTo(createOpts.FileName)
	if err != nil {
		return
	}

	return
}

// Update will update an existing file ini entry.
func Update(client client.Client, updateOpts UpdateOpts) (err error) {
	client.Logger.Debugf("Updating FileIni entry")

	if err = utils.BuildRequest(&updateOpts); err != nil {
		return
	}

	client.Logger.Debugf("FileIni Update Options: %#v", updateOpts)

	exists, err := Exists(client, updateOpts.FileName, updateOpts.Section, updateOpts.Key)
	if err != nil {
		return
	}

	if !exists {
		err = fmt.Errorf("FileIni entry does not exist")
		return
	}

	deleteOpts := DeleteOpts{
		FileName: updateOpts.FileName,
		Section:  updateOpts.Section,
		Key:      updateOpts.Key,
	}

	err = Delete(client, deleteOpts)
	if err != nil {
		return
	}

	createOpts := CreateOpts{
		FileName: updateOpts.FileName,
		Section:  updateOpts.Section,
		Key:      updateOpts.Key,
		Value:    updateOpts.Value,
	}

	return Create(client, createOpts)
}

func Delete(client client.Client, deleteOpts DeleteOpts) (err error) {
	client.Logger.Debugf("Deleting FileIni entry")

	if err = utils.BuildRequest(&deleteOpts); err != nil {
		return
	}

	client.Logger.Debugf("FileIni Delete Options: %#v", deleteOpts)

	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, deleteOpts.FileName)
	if err != nil {
		return
	}

	cfg.Section(deleteOpts.Section).DeleteKey(deleteOpts.Key)
	err = cfg.SaveTo(deleteOpts.FileName)
	if err != nil {
		return
	}

	return
}
