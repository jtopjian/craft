/*
Package fileini manages an entry in an ini file.

To check if an entry exists:

	exists, err := fileini.Exists(client, fileName, sectionName, keyName)

To retrieve a single entry:

	getOpts := fileini.GetOpts{
		FileName: fileName,
		Section:  sectionName,
		Key:      keyName,
	}

	ini, err := fileini.Read(client, getOpts)

To create an entry:

	createOpts := fileini.CreateOpts{
		FileName: fileName,
		Section:  sectionName,
		Key:      keyName,
		Value:   "some value"
	}

	err := fileini.Create(client, createOpts)

To update an entry:

	updateOpts := fileini.UpdateOpts{
		FileName: fileName,
		Section:  sectionName,
		Key:      keyName,
		Value:   "some new value"
	}

	err := fileini.Update(client, updateOpts)

To delete an entry:

	deleteOpts := fileini.DeleteOpts{
		FileName: fileName,
		Section:  sectionName,
		Key:      keyName,
	}

	err := fileini.Delete(client, deleteOpts)
*/

package fileini
