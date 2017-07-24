/*
Package aptpkg manages a package via apt.

To see if a package is installed:

	exists, err := aptpkg.Exists(client, "sl", "")

To install a package:

	createOpts := aptpkg.CreateOpts{
		Name: "sl",
	}

	err := aptpkg.Create(client, createOpts)
	if err != nil {
		return err
	}

To update a package:

	updateOpts := aptpkg.UpdateOpts{
		Name: "sl",
		Version: "latest",
	}

	err := aptpkg.Update(client, updateOpts)
	if err != nil {
		return err
	}

To obtain a list of all packages installed:

	pkgs, err := aptpkg.List(client)

*/
package aptpkg
