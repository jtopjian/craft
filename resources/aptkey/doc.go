/*
Package aptkey provides a way to interact with the apt-key tool.

To create a key:

	createOpts := aptkey.CreateOpts{
		KeyID:         "6026DFCA",
		RemoteKeyFile: "https://www.rabbitmq.com/rabbitmq-release-signing-key.asc",
	}

	err = aptkey.Create(client, createOpts)
	if err != nil {
		return err
	}

To check if a key exists:

	exists, err := aptkey.Exists(client, "6026DFCA")

To get information about a key:

	aptKey, err := aptkey.Read(client, "6026DFCA")

To retrieve all keys:

	aptKeys, err := aptKey.List(client)
*/
package aptkey
