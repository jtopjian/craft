/*
Package aptppa manages a PPA via apt-add-repository.

To see if a PPA exists:

	exists, err := aptppa.Exists(client, "chris-lea/redis-server")

To install a PPA:

	createOpts := aptppa.CreateOpts{
		Name: "chris-lea/redis-server",
	}

	err := aptppa.Create(client, createOpts)

To delete a PPA:

	err := aptppa.Delete(client, "chris-lea/redis-server")

To get a list of all installed PPAs:

	ppas, err := aptppa.List(client)
*/
package aptppa
