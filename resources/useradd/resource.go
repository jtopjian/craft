package useradd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "User"

// User represents a user on a system.
type User struct {
	// Name is the name of the user.
	Name string

	// UID is the user's UID.
	UID string

	// GID is the user's GID.
	GID string

	// Shell is the user's shell account.
	Shell string

	// HomeDir is the user's home directory.
	HomeDir string

	// Sudo is if the user has sudo access.
	Sudo bool

	// Comment is a comment of the user.
	Comment string

	// Groups are groups that the user belongs to.
	Groups []string

	// Passwd is a passwd hash of the user.
	Passwd string
}

// CreateOpts represents options used to create a user with useradd.
type CreateOpts struct {
	// Name is the name of the user.
	Name string `required:"true"`

	// UID is the uid of the user.
	UID string

	// GID is the gid of the user.
	GID string

	// Shell is the shell of the user.
	Shell string `default:"/usr/sbin/nologin"`

	// HomeDir is the user's home directory.
	HomeDir string

	// CreateHome will create the user's home directory.
	CreateHome bool

	// Sudo will give the user sudo rights.
	Sudo bool

	// System will make the account a system account.
	System bool

	// Comment is a descriptive comment of the user.
	Comment string

	// Groups are groups that the user belongs to.
	Groups []string

	// Passwd is an /etc/passwd hash of the password.
	Passwd string
}

// UpdateOpts represents options used to create a user with useradd.
type UpdateOpts struct {
	// UID is the uid of the user.
	UID string

	// GID is the gid of the user.
	GID string

	// Shell is the shell of the user.
	Shell string

	// HomeDir is the user's home directory.
	HomeDir string

	// CreateHome will create the user's home directory.
	CreateHome bool

	// Sudo will give the user sudo rights.
	Sudo bool

	// System will make the account a system account.
	System bool

	// Comment is a descriptive comment of the user.
	Comment string

	// Groups are groups that the user belongs to.
	Groups []string

	// Passwd is an /etc/passwd hash of the password.
	Passwd string
}

// Read will retrieve an existing user account.
func Read(client client.Client, name string) (user User, err error) {
	client.Logger.Debugf("Reading user %s", name)

	ge, err := getent("passwd", name)
	if len(ge) < 7 {
		err = resources.NotFoundError{Type: Type, Name: name}
		return
	}

	user.UID = ge[2]
	user.GID = ge[3]
	user.Comment = ge[4]
	user.HomeDir = ge[5]
	user.Shell = ge[6]

	ge, err = getent("shadow", name)
	if len(ge) > 0 {
		user.Passwd = ge[1]
	}

	sudoFile := fmt.Sprintf("/etc/sudoers.d/%s", name)
	_, err = os.Stat(sudoFile)
	if err == nil {
		user.Sudo = true
	}

	lines, err := utils.FileGetLines("/etc/group")
	if err != nil {
		return
	}

	groupRe := regexp.MustCompile(fmt.Sprintf("(.+):.+:.+:*%s*", name))
	for _, line := range lines {
		if v := groupRe.FindStringSubmatch(line); v != nil {
			user.Groups = append(user.Groups, v[1])
		}
	}

	return
}

// Exists will determine if a user account exists.
func Exists(client client.Client, name string) (exists bool, err error) {
	client.Logger.Debugf("Checking if user %s exists", name)

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

// List will list all users on a system.
func List(client client.Client) (users []User, err error) {
	client.Logger.Debug("Retriving all users")

	lines, err := utils.FileGetLines("/etc/passwd")
	if err != nil {
		return
	}

	for _, line := range lines {
		var user User
		parts := strings.Split(line, ":")
		user, err = Read(client, parts[0])
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	return
}

// Create will create a user on a system.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	var eo utils.ExecOptions
	var createArgs []string

	client.Logger.Debug("Creating user")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("User Create Options: %#v", createOpts)

	if createOpts.UID != "" {
		createArgs = append(createArgs, fmt.Sprintf("-u %s", createOpts.UID))
	}

	if createOpts.GID != "" {
		createArgs = append(createArgs, fmt.Sprintf("-g %s", createOpts.GID))
	}

	if createOpts.HomeDir != "" {
		createArgs = append(createArgs, fmt.Sprintf("-d %s", createOpts.HomeDir))
	}

	if createOpts.CreateHome {
		createArgs = append(createArgs, fmt.Sprintf("-m"))
	}

	if createOpts.Shell != "" {
		createArgs = append(createArgs, fmt.Sprintf("-s %s", createOpts.Shell))
	}

	if createOpts.Passwd != "" {
		createArgs = append(createArgs, fmt.Sprintf("-p \"%s\"", createOpts.Passwd))
	}

	if createOpts.Comment != "" {
		createArgs = append(createArgs, fmt.Sprintf("-p \"%s\"", createOpts.Comment))
	}

	if len(createOpts.Groups) > 0 {
		v := strings.Join(createOpts.Groups, ",")
		createArgs = append(createArgs, fmt.Sprintf("-G %s", v))
	}

	if createOpts.System {
		createArgs = append(createArgs, fmt.Sprintf("-r"))
	}

	createArgs = append(createArgs, createOpts.Name)
	eo.Command = fmt.Sprintf("useradd %s", strings.Join(createArgs, " "))
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stderr != "" {
		err = fmt.Errorf("Error creating user %s: %s", createOpts.Name, execResult.Stderr)
		return
	}

	return
}

// Update will update an existing user on a system.
func Update(client client.Client, name string, updateOpts UpdateOpts) (err error) {
	var eo utils.ExecOptions
	var updateArgs []string

	client.Logger.Debugf("Updating user %s", name)

	if err = utils.BuildRequest(&updateOpts); err != nil {
		return
	}

	client.Logger.Debugf("User Update Options: %#v", updateOpts)

	user, err := Read(client, name)
	if err != nil {
		return
	}

	if updateOpts.UID != "" && updateOpts.UID != user.UID {
		updateArgs = append(updateArgs, fmt.Sprintf("-u %s", updateOpts.UID))
	}

	if updateOpts.GID != "" && updateOpts.GID != user.GID {
		updateArgs = append(updateArgs, fmt.Sprintf("-g %s", updateOpts.GID))
	}

	if updateOpts.Comment != "" && updateOpts.Comment != user.Comment {
		updateArgs = append(updateArgs, fmt.Sprintf("-p \"%s\"", updateOpts.Comment))
	}

	if updateOpts.HomeDir != "" && updateOpts.HomeDir != user.HomeDir {
		updateArgs = append(updateArgs, fmt.Sprintf("-d %s", updateOpts.HomeDir))
	}

	if updateOpts.Shell != "" && updateOpts.Shell != user.Shell {
		updateArgs = append(updateArgs, fmt.Sprintf("-s %s", updateOpts.Shell))
	}

	if len(updateOpts.Groups) > 0 {
		v := strings.Join(updateOpts.Groups, ",")
		updateArgs = append(updateArgs, fmt.Sprintf("-G %s", v))
	}

	updateArgs = append(updateArgs, name)

	eo.Command = fmt.Sprintf("usermod %s", strings.Join(updateArgs, " "))
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stderr != "" {
		err = fmt.Errorf("Unable to update user %s: %s", name, execResult.Stderr)
		return
	}

	return
}

// Delete will delete a user from a system.
func Delete(client client.Client, name string) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Deleting user %s", name)

	eo.Command = fmt.Sprintf("userdel %s", name)
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stderr != "" {
		err = fmt.Errorf("Unable to delete user %s: %s", name, execResult.Stderr)
		return
	}

	return
}

func getent(ent, user string) (getent []string, err error) {
	var eo utils.ExecOptions

	eo.Command = fmt.Sprintf("getent %s %s", ent, user)
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	getent = strings.Split(execResult.Stdout, ":")

	return
}
