/*
Package aptsource manages an apt source entry.

Each entry is placed into an individual file under /etc/apt/sources.list.d. The
name of the file is the "name" to reference the entry by. For example, the name
"rabbitmq" references the file /etc/apt/sources.list.d/rabbitmq.list.

To see if an entry is installed:

	exists, err := aptsource.Exists(client, "rabbitmq")

To create an entry:

	createOpts := aptsource.CreateOpts{
		Name:         "rabbitmq",
		URI:          "http://www.rabbitmq.com/debian/",
		Distribution: "testing",
		Component:    "main",
		IncludeSrc:   true,
	}

	err := aptsource.Create(client, createOpts)

To delete an entry:

	err := aptsource.Delete(client, "rabbitmq")

To get a list of all entries under /etc/apt/sources.list.d:

	entries, err := aptsource.List(client)

*/
package aptsource
