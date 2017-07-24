package gitrepo

import (
	"fmt"
	"os"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "GitRepo"

// GitRepo represents a git repository on a system.
type GitRepo struct {
	// Name is the local location of the git repository.
	Name string

	// Source is the source of the git repository.
	Source string

	// Branch is the branch of the git repository.
	Branch string

	// Commit is the commit of the git repository.
	Commit string

	// Tag is the tag of the git repository.
	Tag string

	// Latest is if the git repository is at the latest update.
	Latest bool
}

// CreateOpts represents options used to create a git repo.
type CreateOpts struct {
	// Name is the destination repository path.
	Name string `required:"true"`

	// Source is the source of the git repository.
	Source string `require:"true"`

	// Owner is the user that owns the git repository.
	Owner string `default:"root"`

	// Group is the group owner of the repository.
	Group string `default:"root"`

	// Branch is the branch of the git repository.
	Branch string

	// Commit is a commit of the git repository.
	Commit string

	// Tag is a tag of the repository.
	Tag string
}

// UpdateOpts represents options used to update a git repo.
type UpdateOpts struct {
	// Owner is the user that owns the git repository.
	Owner string

	// Group is the group owner of the repository.
	Group string

	// Branch is the branch of the git repository.
	Branch string

	// Commit is a commit of the git repository.
	Commit string

	// Tag is a tag of the repository.
	Tag string

	// Latest will trigger a git pull on a branch
	Latest bool
}

func Read(client client.Client, name string) (repo GitRepo, err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Reading GitRepo %s", name)

	_, err = os.Stat(name)
	if err != nil {
		err = resources.NotFoundError{Type: Type, Name: name}
		return
	}

	_, err = os.Stat(name + "/.git/config")
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("%s is not a git repository")
		}
		return
	}

	repo.Name = name
	eo.Dir = name

	// try to determine the branch
	eo.Command = "git rev-parse --abbrev-ref HEAD"
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}
	repo.Branch = strings.TrimSpace(execResult.Stdout)

	eo.Command = "git remote update"
	execResult, err = utils.Exec(eo)
	if err != nil {
		return
	}

	eo.Command = "git status -uno"
	execResult, err = utils.Exec(eo)
	if err != nil {
		return
	}

	if strings.Contains(execResult.Stdout, "up-to-date") {
		repo.Latest = true
	}

	eo.Command = "git rev-parse HEAD"
	execResult, err = utils.Exec(eo)
	if err != nil {
		return
	}
	repo.Commit = strings.TrimSpace(execResult.Stdout)

	eo.Command = "git describe --always --tag"
	execResult, err = utils.Exec(eo)
	if err != nil {
		return
	}
	repo.Tag = strings.TrimSpace(execResult.Stdout)

	return
}

// Exists determines if a git repository exists.
func Exists(client client.Client, name string) (exists bool, err error) {
	client.Logger.Debugf("Checking if %s exists", name)

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

// List is not implemented
//func List(client client.Client) {
//	return
//}

// Create will create/clone a git repository.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Creating GitRepo")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("GitRepo Create Options: %#v", createOpts)

	eo.Command = fmt.Sprintf("git clone --quiet %s %s", createOpts.Source, createOpts.Name)
	_, err = utils.Exec(eo)
	if err != nil {
		return
	}

	eo.Dir = createOpts.Name

	if createOpts.Commit != "" {
		eo.Command = fmt.Sprintf("git checkout %s", createOpts.Commit)
		_, err = utils.Exec(eo)
		if err != nil {
			return
		}
	}

	if createOpts.Tag != "" {
		eo.Command = fmt.Sprintf("git checkout tags/%s", createOpts.Tag)
		_, err = utils.Exec(eo)
		if err != nil {
			return
		}
	}

	if createOpts.Branch != "" {
		eo.Command = fmt.Sprintf("git checkout %s", createOpts.Branch)
		_, err = utils.Exec(eo)
		if err != nil {
			return
		}
	}

	uid, err := utils.UsernameToID(createOpts.Owner)
	if err != nil {
		return
	}

	gid, err := utils.GroupToID(createOpts.Group)
	if err != nil {
		return
	}

	err = utils.ChownR(createOpts.Name, uid, gid)
	if err != nil {
		return
	}

	return
}

// Update will update an existing git repository.
func Update(client client.Client, name string, updateOpts UpdateOpts) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Updating GitRepo %s", name)

	if err = utils.BuildRequest(&updateOpts); err != nil {
		return
	}

	client.Logger.Debugf("GitRepo Update Options: %#v", updateOpts)

	eo.Dir = name

	if updateOpts.Branch != "" {
		eo.Command = fmt.Sprintf("git checkout %s", updateOpts.Branch)
		_, err = utils.Exec(eo)
		if err != nil {
			return
		}

		if updateOpts.Latest {
			eo.Command = fmt.Sprintf("git pull", updateOpts.Branch)
			_, err = utils.Exec(eo)
			if err != nil {
				return
			}
		}
	}

	if updateOpts.Commit != "" {
		eo.Command = fmt.Sprintf("git checkout %s", updateOpts.Commit)
		_, err = utils.Exec(eo)
		if err != nil {
			return
		}
	}

	if updateOpts.Tag != "" {
		eo.Command = fmt.Sprintf("git checkout tags/%s", updateOpts.Tag)
		_, err = utils.Exec(eo)
		if err != nil {
			return
		}
	}

	if updateOpts.Owner != "" && updateOpts.Group != "" {
		var uid int
		var gid int
		uid, err = utils.UsernameToID(updateOpts.Owner)
		if err != nil {
			return
		}

		gid, err = utils.GroupToID(updateOpts.Group)
		if err != nil {
			return
		}

		err = utils.ChownR(name, uid, gid)
		if err != nil {
			return
		}
	}

	return
}

// Delete will delete a git repository.
func Delete(client client.Client, name string) (err error) {
	client.Logger.Debugf("Deleting GitRepo %s", name)

	err = os.RemoveAll(name)
	if err != nil {
		return
	}

	return
}
