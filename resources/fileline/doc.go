/*
Package fileline manages a line in a file.

To check if a line exists:

	exists, err := fileline.Exists(client, line)

To get a line:

	getOpts := fileline.GetOpts{
		FileName: fileName,
		Line      line,
	}

	fileLine, err := fileline.Read(client, getOpts)

	getOpts = fileline.GetOpts{
		FileName: fileName,
		Match:    "^-m",
	}

	fileLine, err := fileline.Read(client, getOpts)

To create a line:

	createOpts := fileline.CreateOpts{
		FileName: fileName,
		Line:     "-m 256",
		Match:    "^-m",
	}

	err := fileline.Create(client, createOpts)

To delete a line:

	deleteOpts := fileline.DeleteOpts{
		FileName: fileName,
		Line:     line,
	}

	err = fileline.Delete(client, deleteOpts)

	deleteOpts = fileline.DeleteOpts{
		FileName: fileName,
		Match:    "^-m"
	}

	err = fileline.Delete(client, deleteOpts)
*/

package fileline
