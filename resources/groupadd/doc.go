/*
Package groupadd manages a system group with the groupadd, groupmod, and
groupdel commands.

To check if a group exists:

	exists, err := groupadd.Exists(client, groupName)

To read a group:

	group, err := groupadd.Read(client, groupName)

To retrieve all groups:

	groups, err := groupadd.List(client)

To create a group:

	createOpts := groupadd.CreateOpts{
		Name: "foobar",
		GID:  "1002",
	}

	err := groupadd.Create(client, createOpts)

To update a group:

	updateOpts := groupadd.UpdateOpts{
		GID: "1003",
	}

	err := groupadd.Update(client, updateOpts)

To delete a group:

	err : groupadd.Delete(client, groupName)
*/
package groupadd
