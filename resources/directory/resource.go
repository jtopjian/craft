package directory

import (
	"fmt"
	"os"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "Directory"

// Directory represents a directory on a system.
type Directory struct {
	Name  string
	Owner string
	Group string
	Mode  string
}

// CreateOpts represents options used to create a directory.
type CreateOpts struct {
	// Name is the name of the directory.
	Name string `required:"true"`

	// Owner is the user of the directory.
	Owner string `default:"root"`

	// Group is the group owner of the directory.
	Group string `default:"root"`

	// Mode is the mode/permissions of the directory.
	Mode string `default:"0755"`

	// Parents determines if the full directory path should be created.
	Parents bool
}

// UpdateOpts represents options used to update a directory.
type UpdateOpts struct {
	// Owner is the user of the directory.
	Owner string

	// Group is the group owner of the directory.
	Group string

	// Mode is the mode/permissions of the directory.
	Mode string

	// Recurse determines if changes should be applied to all files
	// and directories within the directory.
	Recurse bool
}

// Read will retrieve information about an existing directory.
func Read(client client.Client, dirName string) (dir Directory, err error) {
	client.Logger.Debugf("Reading directory %s", dirName)

	file, err := os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			err = resources.NotFoundError{Type: Type, Name: dirName}
		}
		return
	}

	mode := file.Mode()
	if !mode.IsDir() {
		err = fmt.Errorf("%s is not a directory", dirName)
		return
	}

	uid, gid, err := utils.GetFileOwner(dirName)
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

	dir.Name = dirName
	dir.Mode = mode.String()
	dir.Owner = owner
	dir.Group = group

	return
}

// Exists will determine if a directory exists on a system.
func Exists(client client.Client, dirName string) (exists bool, err error) {
	client.Logger.Debugf("Checking if directory %s exists", dirName)

	_, err = Read(client, dirName)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	exists = true

	return
}

// List is not implemented for Directory.
func List(client client.Client) {
	return
}

// Create will create a file on a system.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	client.Logger.Debugf("Creating directory")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("Directory Create Options: %#v", createOpts)

	mode, err := utils.StringToMode(createOpts.Mode)
	if err != nil {
		return
	}

	if createOpts.Parents {
		err = os.MkdirAll(createOpts.Name, mode)
		if err != nil {
			return
		}
	} else {
		err = os.Mkdir(createOpts.Name, mode)
		if err != nil {
			return
		}
	}

	uid, gid, err := utils.GetUIDGID(createOpts.Owner, createOpts.Group)
	if err != nil {
		return
	}

	if createOpts.Parents {
		err = utils.ChownR(createOpts.Name, uid, gid)
		if err != nil {
			return
		}
		err = utils.ChmodR(createOpts.Name, mode)
		if err != nil {
			return
		}
	} else {
		err = os.Chown(createOpts.Name, uid, gid)
		if err != nil {
			return
		}
	}

	return
}

// Update will update an existing directory.
func Update(client client.Client, dirName string, updateOpts UpdateOpts) (err error) {
	client.Logger.Debugf("Updating directory %s", dirName)

	if err = utils.BuildRequest(&updateOpts); err != nil {
		return
	}

	client.Logger.Debugf("Directory Update Options: %#v", updateOpts)

	if updateOpts.Mode != "" {
		var mode os.FileMode
		mode, err = utils.StringToMode(updateOpts.Mode)
		if err != nil {
			return
		}

		if updateOpts.Recurse {
			err = utils.ChmodR(dirName, mode)
			if err != nil {
				return
			}
		} else {
			err = os.Chmod(dirName, mode)
			if err != nil {
				return
			}
		}
	}

	if updateOpts.Owner != "" || updateOpts.Group != "" {
		var uid int
		var gid int

		uid, gid, err = utils.GetFileOwner(dirName)
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

		if updateOpts.Recurse {
			err = utils.ChownR(dirName, uid, gid)
			if err != nil {
				return
			}
		} else {
			err = os.Chown(dirName, uid, gid)
			if err != nil {
				return
			}
		}
	}

	return
}

// Delete will delete a file on a system.
func Delete(client client.Client, dirName string, recurse bool) (err error) {
	client.Logger.Debugf("Deleting directory %s (Recurse: %t)", dirName, recurse)

	if recurse {
		err = os.RemoveAll(dirName)
		return
	}

	err = os.Remove(dirName)

	return
}
