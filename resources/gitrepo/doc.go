/*
Package gitrepo manages a git repository on a system.

To check if a git repository exists:

	exists, err := gitrepo.Exists(client, path)

Note that Exists will only check if a git repository exists. It will not take
into account branch, commit, tag, etc.

To clone a git repository:

	createOpts := gitrepo.CreateOpts{
		Name:   "destination/path",
		Source: "https://github.com/foo/bar",
		Branch: "master",
	}

	err := gitrepo.Create(client, createOpts)

To update a git repository:

	updateOpts := gitrepo.UpdateOpts{
		Commit: "abcd1234",
	}

	err := gitrepo.Update(client, updateOpts)

To delete a git repository:

	err = gitrepo.Delete(client, path)
*/
package gitrepo
