package aptppa

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "AptPPA"

// AptPPA represents a PPA installed on a system.
type AptPPA struct {
	Name string
}

// CreateOpts represents options used to install a PPA via apt-add-repository.
type CreateOpts struct {
	// Name is the name of the PPA.
	// This is what you would pass into apt-add-repistory without the "ppa:" part.
	Name string `required:"true"`

	// Refresh will trigger apt-get update after the ppa has been installed.
	Refresh bool `default:"true"`
}

// Read will return information about an installed PPA.
func Read(client client.Client, ppa string) (aptPPA AptPPA, err error) {
	client.Logger.Debugf("Reading PPA %s", ppa)

	sourceFileName, err := aptPPASourceFileName(ppa)
	if err != nil {
		return
	}

	v := "/etc/apt/sources.list.d/" + sourceFileName
	client.Logger.Debugf("PPA file: %s", v)
	_, err = os.Stat(v)
	if os.IsNotExist(err) {
		err = resources.NotFoundError{Type: Type, Name: ppa}
		return
	}

	aptPPA = AptPPA{
		Name: ppa,
	}

	return
}

// Exists determines if a PPA is installed on a system.
func Exists(client client.Client, ppa string) (exists bool, err error) {
	client.Logger.Debugf("Checking if PPA %s exists", ppa)

	aptPPA, err := Read(client, ppa)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	if aptPPA.Name != "" {
		exists = true
	}

	return
}

// List will return all PPAs installed on a system.
func List(client client.Client) (aptPPAs []AptPPA, err error) {
	client.Logger.Debugf("Listing all PPAs")

	files, err := filepath.Glob("/etc/apt/sources.list.d/*.list")
	if err != nil {
		return
	}

	lsbInfo, err := utils.GetLSBInfo()
	if err != nil {
		return
	}

	distro := fmt.Sprintf("-%s-", strings.ToLower(lsbInfo.DistributorID))
	release := fmt.Sprintf("-%s", strings.ToLower(lsbInfo.Codename))

	for _, file := range files {
		ppa := path.Base(file)
		ppa = strings.Replace(ppa, distro, "/", -1)
		ppa = strings.Replace(ppa, release, "", -1)
		ppa = strings.Replace(ppa, ".list", "", -1)

		if strings.Contains(ppa, "/") {
			aptPPA := AptPPA{
				Name: ppa,
			}

			aptPPAs = append(aptPPAs, aptPPA)
		}
	}

	return
}

// Create will install a PPA via apt-add-repository.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	var eo utils.ExecOptions
	client.Logger.Debugf("Creating PPA")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("PPA Create Options: %#v", createOpts)

	eo.Command = fmt.Sprintf("apt-add-repository -y ppa:%s", createOpts.Name)
	_, err = utils.Exec(eo)
	if err != nil {
		return
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

// Update is not implemented for aptppa.
func Update() (err error) {
	return
}

// Delete will remove a PPA from a system via apt-add-repository.
func Delete(client client.Client, ppa string) (err error) {
	var eo utils.ExecOptions
	client.Logger.Debugf("Deleting PPA %s", ppa)

	eo.Command = fmt.Sprintf("apt-add-repository -y -r ppa:%s", ppa)
	_, err = utils.Exec(eo)
	if err != nil {
		return
	}

	sourceFileName, err := aptPPASourceFileName(ppa)
	if err != nil {
		return
	}

	v := "/etc/apt/sources.list.d/" + sourceFileName
	err = os.Remove(v)
	if err != nil {
		return
	}

	eo.Command = fmt.Sprintf("apt-get update -qq")
	_, err = utils.Exec(eo)
	if err != nil {
		return
	}

	return
}

func aptPPASourceFileName(name string) (string, error) {
	lsbInfo, err := utils.GetLSBInfo()
	if err != nil {
		return "", nil
	}

	distro := fmt.Sprintf("-%s-", strings.ToLower(lsbInfo.DistributorID))
	release := strings.ToLower(lsbInfo.Codename)

	name = strings.Replace(name, "/", distro, -1)
	name = strings.Replace(name, ":", "-", -1)
	name = strings.Replace(name, ".", "_", -1)

	name = fmt.Sprintf("%s-%s.list", name, release)

	return name, nil
}
