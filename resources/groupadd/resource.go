package groupadd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "Group"

// Group represents a group managed by groupadd
type Group struct {
	// Name is the name of the group
	Name string

	// GID is the group id of the group
	GID string
}

// CreateOpts represents options used to create a group with groupadd.
type CreateOpts struct {
	// Name is the name of the group.
	Name string `required:"true"`

	// GID is the group id
	GID string
}

// UpdateOpts represents options used to update a group with groupmod.
type UpdateOpts struct {
	// GID is the group id
	GID string
}

// Read will read an existing group.
func Read(client client.Client, name string) (group Group, err error) {
	client.Logger.Debugf("Reading Group %s", name)

	gid, err := utils.GroupToID(name)
	if err != nil {
		err = resources.NotFoundError{Type: Type, Name: name}
		return
	}

	group.Name = name
	group.GID = strconv.Itoa(gid)

	return
}

// Exists will determine if a group exists.
func Exists(client client.Client, name string) (exists bool, err error) {
	client.Logger.Debugf("Checking if group %s exists", name)

	_, err = Read(client, name)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	exists = true
	return
}

// List will list all groups on a system.
func List(client client.Client) (groups []Group, err error) {
	client.Logger.Debug("Retrieving all groups")

	lines, err := utils.FileGetLines("/etc/group")
	if err != nil {
		return
	}

	for _, line := range lines {
		var group Group
		parts := strings.Split(line, ":")
		group.Name = parts[0]
		group.GID = parts[2]

		groups = append(groups, group)
	}

	return
}

// Create will create a group via groupadd.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	var eo utils.ExecOptions
	var createArgs []string

	client.Logger.Debug("Creating group")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("Group Create Options: %#v", createOpts)

	if createOpts.GID != "" {
		createArgs = append(createArgs, fmt.Sprintf("-g %s", createOpts.GID))
	}

	createArgs = append(createArgs, createOpts.Name)
	eo.Command = fmt.Sprintf("groupadd %s", strings.Join(createArgs, " "))
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stderr != "" {
		err = fmt.Errorf("Error adding group: %s", execResult.Stderr)
		return
	}

	return
}

// Update will update a group via groupmod.
func Update(client client.Client, name string, updateOpts UpdateOpts) (err error) {
	var eo utils.ExecOptions
	var updateArgs []string

	client.Logger.Debugf("Updating Group %s", name)

	if err = utils.BuildRequest(&updateOpts); err != nil {
		return
	}

	client.Logger.Debugf("Group Update Options: %#v", updateOpts)

	if updateOpts.GID != "" {
		updateArgs = append(updateArgs, fmt.Sprintf("-g %s", updateOpts.GID))
	}

	eo.Command = fmt.Sprintf("groupmod %s %s", strings.Join(updateArgs, " "), name)
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stderr != "" {
		err = fmt.Errorf("Error updating group: %s", execResult.Stderr)
		return
	}

	return
}

// Delete will delete a group via groupdel.
func Delete(client client.Client, name string) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Deleting Group %s", name)

	eo.Command = fmt.Sprintf("groupdel %s", name)
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stderr != "" {
		err = fmt.Errorf("Error updating group: %s", execResult.Stderr)
		return
	}

	return
}
