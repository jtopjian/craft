package file

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "File"

// File represents a file on a system.
type File struct {
	// Name is the name of the file.
	Name string

	// Owner is the owner of the file.
	Owner string

	// Group is the group owner of the file.
	Group string

	// Mode is the permissions of the file.
	Mode string

	// MD5 is the md5sum of the file contents.
	MD5 string
}

// CreateOpts represents options used to create a file on a system.
type CreateOpts struct {
	// Name is the name of the file. This is the full path.
	Name string `required:"true"`

	// Owner is the user who owns the file.
	Owner string `default:"root"`

	// Group is the group owner of the file.
	Group string `default:"root"`

	// Mode is the permissions of the file.
	Mode string `default:"0640"`

	// Content is the file contents.
	Content string
}

// UpdateOpts represents options used to update a file on a system.
type UpdateOpts struct {
	// Owner is the user who owns the file.
	Owner string

	// Group is the group owner of the file.
	Group string

	// Mode is the permissions of the file.
	Mode string

	// Content is the file contents.
	Content string
}

// Read will read an existing file on a system.
func Read(client client.Client, fileName string) (file File, err error) {
	client.Logger.Debugf("Reading file %s", fileName)

	fi, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			err = resources.NotFoundError{Type: Type, Name: fileName}
		}
		return
	}

	mode := fi.Mode()
	if !mode.IsRegular() {
		err = fmt.Errorf("%s is not a file", fileName)
		return
	}

	uid, gid, err := utils.GetFileOwner(fileName)
	if err != nil {
		return
	}

	owner, err := utils.UIDToName(uid)
	if err != nil {
		return
	}

	group, err := utils.GIDToName(gid)
	if err != nil {
		return
	}

	var f *os.File
	f, err = os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return
	}
	fileSum := string(h.Sum(nil))

	file.Name = fileName
	file.Owner = owner
	file.Group = group
	file.Mode = mode.String()
	file.MD5 = fileSum

	return
}

// Exists will determine if a file exists on a system.
func Exists(client client.Client, fileName string) (exists bool, err error) {
	client.Logger.Debugf("Checking if file %s exists", fileName)

	file, err := Read(client, fileName)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	if file.Name != "" {
		exists = true
	}

	return
}

// Create will create a file on a system.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	client.Logger.Debugf("Creating file")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return err
	}

	client.Logger.Debugf("File Create Options: %#v", createOpts)

	mode, err := utils.StringToMode(createOpts.Mode)
	if err != nil {
		return
	}

	if createOpts.Content == "" {
		var f *os.File
		f, err = os.OpenFile(createOpts.Name, os.O_RDONLY|os.O_CREATE, mode)
		if err != nil {
			return
		}
		err = f.Close()
		if err != nil {
			return
		}
	} else {
		err = ioutil.WriteFile(createOpts.Name, []byte(createOpts.Content), mode)
		if err != nil {
			return
		}
	}

	uid, gid, err := utils.GetUIDGID(createOpts.Owner, createOpts.Group)
	if err != nil {
		return
	}

	err = os.Chown(createOpts.Name, uid, gid)
	if err != nil {
		return
	}

	return
}

// Update will update an existing file on a system.
func Update(client client.Client, fileName string, updateOpts UpdateOpts) (err error) {
	client.Logger.Debugf("Updating file %s", fileName)

	if err = utils.BuildRequest(&updateOpts); err != nil {
		return
	}

	client.Logger.Debugf("File Update Options: %#v", updateOpts)

	if updateOpts.Mode != "" {
		var mode os.FileMode
		mode, err = utils.StringToMode(updateOpts.Mode)
		if err != nil {
			return
		}

		err = os.Chmod(fileName, mode)
		if err != nil {
			return
		}
	}

	if updateOpts.Content != "" {
		var fi os.FileInfo
		fi, err = os.Stat(fileName)
		if os.IsNotExist(err) {
			err = nil
			return
		}

		mode := fi.Mode()
		err = ioutil.WriteFile(fileName, []byte(updateOpts.Content), mode)
		if err != nil {
			return
		}
	}

	if updateOpts.Owner != "" || updateOpts.Group != "" {
		var uid int
		var gid int
		uid, gid, err = utils.GetFileOwner(fileName)
		if err != nil {
			return
		}

		if updateOpts.Owner != "" {
			uid, err = utils.UsernameToID(updateOpts.Owner)
			if err != nil {
				return
			}
		}

		if updateOpts.Group != "" {
			gid, err = utils.GroupToID(updateOpts.Group)
			if err != nil {
				return
			}
		}

		err = os.Chown(fileName, uid, gid)
		if err != nil {
			return
		}
	}

	return
}

// Delete deletes a file on a system.
func Delete(client client.Client, fileName string) (err error) {
	client.Logger.Debugf("Deleting file %s", fileName)

	err = os.Remove(fileName)
	if err != nil {
		return
	}

	return
}
