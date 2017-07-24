/*
Package file manages a file on a system.

To check if a file exists:

	exists, err := file.Exists(client, fileName)

To create a file:

	createOpts := file.CreateOpts{
		Name:    "/tmp/foo",
		Owner:   "root",
		Group:   "root",
		Mode:    "0755",
		Content: "Hello, World!\n",
	}

	err := file.Create(client, CreateOpts)

To update a file:

	updateOpts := file.UpdateOpts{
		Mode: "0644",
	}

	err := file.Update(client, fileName, updateOpts)

To delete a file:

	err := file.Delete(client, fileName)
*/
package file
