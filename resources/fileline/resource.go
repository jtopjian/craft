package fileline

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "FileLine"

// FileLine represents a line in a file.
type FileLine struct {
	FileName string
	Line     string
}

// CreateOpts represents options used to create a line in a file.
type CreateOpts struct {
	// FileName is the name of the file.
	FileName string `required:"true"`

	// Line is the line to add to the file.
	Line string `required:"true"`

	// Match is a regular expression to match against an existing line.
	Match string
}

// GetOpts represents options to get a line in a file.
type GetOpts struct {
	// FileName is the name of the file.
	FileName string `required:"true"`

	// Line is the line to get.
	Line string

	// Match is a regular expression to match against a line.
	Match string
}

// DeleteOpts represents options to delete a line from a file.
type DeleteOpts struct {
	// FileName is the name of the file.
	FileName string `required:"true"`

	// Line is the line to delete from a file.
	Line string

	// Match is a regular expression to match against an existing line.
	Match string
}

// Read will read an existing line from a file.
func Read(client client.Client, getOpts GetOpts) (fileLine FileLine, err error) {
	client.Logger.Debug("Reading line from file")

	err = utils.BuildRequest(&getOpts)
	if err != nil {
		return
	}

	client.Logger.Debugf("FileLine Read Options: %#v", getOpts)

	lines, err := utils.FileGetLines(getOpts.FileName)
	if err != nil {
		return
	}

	var re *regexp.Regexp
	resourceTitle := fmt.Sprintf("%s/%s", getOpts.FileName, getOpts.Line)
	if getOpts.Match != "" {
		resourceTitle = fmt.Sprintf("%s/%s", getOpts.FileName, getOpts.Match)
		re = regexp.MustCompile(getOpts.Match)
	}

	for _, line := range lines {
		if getOpts.Match != "" {
			if re.MatchString(line) {
				fileLine.Line = line
			}
		} else {
			if line == getOpts.Line {
				fileLine.Line = line
			}
		}
	}

	if fileLine.Line != "" {
		fileLine.FileName = getOpts.FileName
	}

	if fileLine.Line == "" {
		err = resources.NotFoundError{Type: Type, Name: resourceTitle}
	}

	return
}

// Exists will determine if a line exists in a file.
func Exists(client client.Client, fileName, line string) (exists bool, err error) {
	client.Logger.Debugf("Checking if line %s is in file %s", line, fileName)

	getOpts := GetOpts{
		FileName: fileName,
		Line:     line,
	}

	_, err = Read(client, getOpts)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	exists = true

	return
}

// Create will create a line in a file.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	client.Logger.Debug("Adding line to file")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("FileLine Create Options: %#v", createOpts)

	var lineRe *regexp.Regexp
	if createOpts.Match != "" {
		lineRe = regexp.MustCompile(createOpts.Match)
	}

	lines, err := utils.FileGetLines(createOpts.FileName)
	if err != nil {
		return
	}

	var newLines []string
	var changed bool
	for _, line := range lines {
		if createOpts.Match != "" {
			if lineRe.MatchString(line) {
				changed = true
				line = createOpts.Line
			}
		}
		newLines = append(newLines, line)
	}

	if !changed {
		newLines = append(newLines, createOpts.Line)
	}

	newContent := strings.Join(newLines, "\n")
	err = utils.WriteFile(createOpts.FileName, newContent)
	if err != nil {
		return
	}

	return
}

// Update is not implemented.
func Update(client client.Client) {
	return
}

func Delete(client client.Client, deleteOpts DeleteOpts) (err error) {
	client.Logger.Debug("Deleting line from file")

	if err = utils.BuildRequest(&deleteOpts); err != nil {
		return
	}

	client.Logger.Debugf("FileLine Delete Options: %#v", deleteOpts)

	var lineRe *regexp.Regexp
	if deleteOpts.Match != "" {
		lineRe = regexp.MustCompile(deleteOpts.Match)
	}

	lines, err := utils.FileGetLines(deleteOpts.FileName)
	if err != nil {
		return
	}

	var newLines []string
	for _, line := range lines {
		if deleteOpts.Match != "" {
			if !lineRe.MatchString(line) {
				newLines = append(newLines, line)
			}
		} else {
			if line != deleteOpts.Line {
				newLines = append(newLines, line)
			}
		}
	}

	newContent := strings.Join(newLines, "\n")
	err = utils.WriteFile(deleteOpts.FileName, newContent)
	if err != nil {
		return
	}

	return
}
