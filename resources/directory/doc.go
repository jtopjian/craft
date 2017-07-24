/*
Package directory manages a directory on a system.

To check if a directory exists:

	exists, err := directory.Exists(client, dirName)

To create a directory:

	createOpts := directory.CreateOpts{
		Name:  "/tmp/foo",
		Owner: "root",
		Group: "root",
		Mode:  "0755",
	}

	err := directory.Create(client, CreateOpts)

To update a directory:

	updateOpts := directory.UpdateOpts{
		Mode: "0644",
	}

	err := directory.Update(client, dirName, updateOpts)

To delete a directory:

	recurse := true
	err := directory.Delete(client, dirName, recurse)
*/
package directory
