package aptsource

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "AptSource"

// AptSource represents an apt source entry.
type AptSource struct {
	// Name is the name of an apt source entry.
	// It is used as the name of the file which contains the entry.
	Name string

	// URI is the URI of the apt source entry.
	URI string

	// Distribution is the distribution of the apt source entry.
	Distribution string

	// Component is the component of the apt source entry.
	Component string

	// IncludeSrc denotes if a source entry is also included.
	IncludeSrc bool
}

// CreateOpts represents options used to install an apt source entry.
type CreateOpts struct {
	// Name is the name of the apt source entry.
	// It is used as the name of the file which contains the entry.
	Name string `required:"true"`

	// URI is the URI of the apt source entry.
	URI string `required:"true"`

	// Distribution is the distribution of the apt source entry.
	Distribution string `required:"true"`

	// Component is the component of the apt source entry.
	Component string

	// IncludeSrc denotes if a source entry will also be included.
	IncludeSrc bool

	// Refresh determines if an apt-get upgrade will run after the entry
	// has been created.
	Refresh bool `default:"true"`
}

// Read will retrieve information about an existing apt source entry.
func Read(client client.Client, name string) (aptSource AptSource, err error) {
	client.Logger.Debugf("Reading apt source entry %s", name)

	path := fmt.Sprintf("/etc/apt/sources.list.d/%s.list", name)
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = resources.NotFoundError{Type: Type, Name: name}
		}
		return
	}

	content, err := ioutil.ReadFile(path)
	entry, err := aptSourceParseEntry(string(content))
	if err != nil {
		return
	}

	aptSource.Name = name
	aptSource.URI = entry.URI
	aptSource.Distribution = entry.Distribution
	aptSource.Component = entry.Component

	path = fmt.Sprintf("/etc/apt/sources.list.d/%s-src.list", name)
	_, err = os.Stat(path)
	if err == nil {
		aptSource.IncludeSrc = true
	}

	return
}

// Exists will report if a given apt source entry exists on a system.
func Exists(client client.Client, name string, includeSrc bool) (exists bool, err error) {
	client.Logger.Debugf("Checking if apt source %s exists", name)

	aptSource, err := Read(client, name)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	if aptSource.Name == "" {
		return
	}

	if includeSrc {
		if !aptSource.IncludeSrc {
			return
		}
	}

	exists = true
	return
}

// List will list all apt source entries on a system.
func List(client client.Client) (aptSources []AptSource, err error) {
	client.Logger.Debugf("Listing all apt source entries")

	files, err := filepath.Glob("/etc/apt/sources.list.d/*.list")
	if err != nil {
		return
	}

	re := regexp.MustCompile("^deb")
	for _, file := range files {
		name := path.Base(file)
		name = strings.Replace(name, ".list", "", -1)

		var content []byte
		content, err = ioutil.ReadFile(file)
		if err != nil {
			return
		}

		lines := string(content)
		var aptSource AptSource
		var e entry
		// If there are multiple lines, this is assuming the source entry
		// is on the second line. Thus we'll only parse the last entry
		// in order to create one apt source object that represents the
		// whole file -- both binary and source entries.
		for _, line := range strings.Split(lines, "\n") {
			if re.MatchString(line) {
				e, err = aptSourceParseEntry(line)
				if err != nil {
					return
				}

				aptSource.Name = name
				aptSource.URI = e.URI
				aptSource.Distribution = e.Distribution
				aptSource.Component = e.Component
				aptSource.IncludeSrc = e.Source
			}
		}

		aptSources = append(aptSources, aptSource)
	}

	return

}

// Create will create an apt source entry.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Creating apt source entry")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("AptSource Create Options: %#v", createOpts)

	e := entry{
		URI:          createOpts.URI,
		Distribution: createOpts.Distribution,
		Component:    createOpts.Component,
	}

	path := fmt.Sprintf("/etc/apt/sources.list.d/%s.list", createOpts.Name)
	content := aptSourceBuildEntry(e, false)
	err = ioutil.WriteFile(path, []byte(content+"\n"), 0644)
	if err != nil {
		return
	}

	if createOpts.IncludeSrc {
		path = fmt.Sprintf("/etc/apt/sources.list.d/%s-src.list", createOpts.Name)
		content = aptSourceBuildEntry(e, true)
		err = ioutil.WriteFile(path, []byte(content+"\n"), 0644)
		if err != nil {
			return
		}
	}

	if createOpts.Refresh {
		eo.Command = fmt.Sprintf("apt-get update -qq")
		_, err = utils.Exec(eo)
		if err != nil {
			return
		}
	}

	return
}

// Update is not implemented.
func Update() (err error) {
	return
}

// Delete will delete an apt source entry.
func Delete(client client.Client, name string) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Deleting apt source entry %s", name)

	path := fmt.Sprintf("/etc/apt/sources.list.d/%s.list", name)
	err = os.Remove(path)
	if err != nil {
		return
	}

	path = fmt.Sprintf("/etc/apt/sources.list.d/%s-src.list", name)
	_, err = os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			return
		}
	}

	eo.Command = fmt.Sprintf("apt-get update -qq")
	_, err = utils.Exec(eo)
	if err != nil {
		return
	}

	return
}

// entry is an internal type that represents an apt entry.
type entry struct {
	URI          string
	Distribution string
	Component    string
	Source       bool
}

// aptSourceBuildEntry is an internal function that will build an apt source entry.
func aptSourceBuildEntry(e entry, source bool) string {
	deb := "deb"
	if source {
		deb = "deb-src"
	}

	return fmt.Sprintf("%s %s %s %s", deb, e.URI, e.Distribution, e.Component)
}

// aptSourceParseEntry is an internal function that will parse an apt source entry.
func aptSourceParseEntry(e string) (entry entry, err error) {
	v := strings.Split(e, " ")
	if len(v) != 4 {
		err = fmt.Errorf("Unable to parse %s", v)
		return
	}

	if v[0] == "deb-src" {
		entry.Source = true
	}

	entry.URI = v[1]
	entry.Distribution = v[2]
	entry.Component = v[3]

	return
}
