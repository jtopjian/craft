package aptpkg

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "AptPkg"

// AptPkg represents a package managed by apt.
type AptPkg struct {
	// Name is the name of the package.
	Name string

	// Version is the version of the package.
	Version string

	// LatestVersion is the latest version of the package available.
	LatestVersion string
}

// CreateOpts represents options used to install a package vi apt-get.
type CreateOpts struct {
	// Name is the name of the package.
	Name string `required:"true"`

	// Version is the version of the package.
	// The following values are valid: a specific version number and "latest".
	Version string
}

// UpdateOpts represents options used to update a package vi apt-get.
type UpdateOpts struct {
	// Version is the version of the package.
	// The following values are valid: a specific version number and "latest".
	Version string
}

// Read will retrieve information about an installed apt package.
func Read(client client.Client, pkgName string) (aptPkg AptPkg, err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Reading package %s", pkgName)

	eo.Command = fmt.Sprintf("apt-cache policy %s", pkgName)
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stdout == "" {
		err = resources.NotFoundError{Type: Type, Name: pkgName}
		return
	}

	installedVersion, candidateVersion := aptPkgParseAptCache(execResult.Stdout)
	aptPkg.Name = pkgName
	aptPkg.Version = installedVersion
	aptPkg.LatestVersion = candidateVersion

	return
}

// Exists will report if a given package exists on a system.
func Exists(client client.Client, pkgName string) (exists bool, err error) {
	client.Logger.Debugf("Checking if package %s exists", pkgName)

	_, err = Read(client, pkgName)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	exists = true

	return
}

// List will retrieve all apt managed packages on a system.
func List(client client.Client) (aptPkgs []AptPkg, err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Listing all packages")

	eo.Command = "dpkg -l"
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	pkgs := aptPkgParseDpkgL(execResult.Stdout)
	for k, v := range pkgs {
		pkg := AptPkg{
			Name:    k,
			Version: v,
		}
		aptPkgs = append(aptPkgs, pkg)
	}

	return
}

// Create will install a package via apt-get.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Installing package")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("Package Create Options: %#v", createOpts)

	eo.Env = []string{
		"DEBIAN_FRONTEND=noninteractive",
		"APT_LISTBUGS_FRONTEND=none",
		"APT_LISTCHANGES_FRONTEND=none",
		"PATH=" + os.Getenv("PATH"),
	}

	var createArgs string
	if createOpts.Version != "" && createOpts.Version != "latest" {
		createArgs = fmt.Sprintf("%s=%s", createOpts.Name, createOpts.Version)
	} else {
		createArgs = createOpts.Name
	}

	eo.Command = fmt.Sprintf(
		"apt-get install -y --allow-downgrades --allow-remove-essential "+
			"--allow-change-held-packages -o DPkg::Options::=--force-confold %s",
		createArgs)

	_, err = utils.Exec(eo)
	if err != nil {
		return
	}

	return
}

// Update will update a package via apt-get.
func Update(client client.Client, pkgName string, updateOpts UpdateOpts) (err error) {
	client.Logger.Debugf("Upgrading package")

	if err = utils.BuildRequest(&updateOpts); err != nil {
		return
	}

	client.Logger.Debugf("Package Update Options: %#v", updateOpts)

	createOpts := CreateOpts{
		Name:    pkgName,
		Version: updateOpts.Version,
	}

	return Create(client, createOpts)
}

// Delete will uninstall a package via apt-get.
func Delete(client client.Client, pkgName string) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Deleting package %s", pkgName)

	eo.Env = []string{
		"DEBIAN_FRONTEND=noninteractive",
		"APT_LISTBUGS_FRONTEND=none",
		"APT_LISTCHANGES_FRONTEND=none",
		"PATH=" + os.Getenv("PATH"),
	}

	eo.Command = fmt.Sprintf("apt-get purge -q -y %s", pkgName)
	_, err = utils.Exec(eo)
	if err != nil {
		return
	}

	return
}

// apkgPkgParseAptCache is an internal function that will parse the
// output of apt-cache policy and return the version information.
func aptPkgParseAptCache(stdout string) (installed, candidate string) {
	installedRe := regexp.MustCompile("Installed: (.+)\n")
	candidateRe := regexp.MustCompile("Candidate: (.+)\n")

	if v := installedRe.FindStringSubmatch(stdout); len(v) > 1 {
		installed = v[1]
	}

	if v := candidateRe.FindStringSubmatch(stdout); len(v) > 1 {
		candidate = v[1]
	}

	return
}

// aptPkgParseDpkgL is an internal function that will parse the output
// of dpkg -l and return a list of packages and their versions.
func aptPkgParseDpkgL(stdout string) (pkgs map[string]string) {
	pkgs = make(map[string]string)
	pkgRe := regexp.MustCompile("^ii\\s+(\\S+)\\s+(\\S+)")
	for _, pkg := range strings.Split(stdout, "\n") {
		if v := pkgRe.FindStringSubmatch(pkg); v != nil {
			pkgs[v[1]] = v[2]
		}
	}

	return
}
