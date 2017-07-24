/*
Package useradd manages a system user with the useradd, usermod, and
userdel commands.

To check if a user exists:

	exists, err := useradd.Exists(client, userName)

To read a user:

	user, err := useradd.Read(client, userName)

To retrieve all users:

	users, err := useradd.List(client)

To create a user:

	createOpts := useradd.CreateOpts{
		Name: "foobar",
		UID:  "1002",
	}

	err := useradd.Create(client, createOpts)

To update a user:

	updateOpts := useradd.UpdateOpts{
		UID: "1003",
	}

	err := useradd.Update(client, updateOpts)

To delete a user:

	err : useradd.Delete(client, userName)
*/
package useradd
