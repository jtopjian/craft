/*
Package cronentry manages a cron entry for a user.

Only cron entries that are tagged with a name at the end of the entry are managed:

	0 1 * * * /path/to/command # Foobar

To see if an entry exists:

	exists, err := cronentry.Exists(client, user, name)

To create a cron entry:

	createOpts := cronentry.CreateOpts{
		Name:    "Foobar",
		Minute:  0,
		Hour:    1,
		Command: "/path/to/command",
	}

	err := cronentry.Create(client, user, createOpts)

To update a cron entry:

	updateOpts := cronentry.UpdateOpts{
		Minute:  5,
		Hour:    7,
		Command: "/path/to/command",
	}

	err := cronentry.Update(client, user, name, updateOpts)

To delete a cron entry:

	err := cronentry.Delete(client, user, name)

To get a list of all managed cron entries:

	entries, err := cronentry.List(client, user)
*/
package cronentry
