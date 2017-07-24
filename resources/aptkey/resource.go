package aptkey

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
	"golang.org/x/crypto/openpgp"
)

const Type = "AptKey"

// AptKey represents a key managed by apt-key.
type AptKey struct {
	// KeyID is the short key identifier of the key.
	KeyID string

	// PublicKey is the public key.
	PublicKey string

	// Name is a name of the key maintainer.
	Name string
}

// CreateOpts represents options used to create a key via apt-key.
type CreateOpts struct {
	// KeyID is the short key identifier of the key.
	KeyID string `required:"true"`

	// KeyServer is an optional remote server to obtain the key from.
	// If KeyServer is not used, RemoteKeyFile must be used.
	KeyServer string

	// RemoteKeyFile is the URL to a public key.
	// If RemoteKeyFile is not used, KeyServer must be used.
	RemoteKeyFile string
}

// Read will read details of an existing apt-key key.
func Read(client client.Client, keyID string) (aptKey AptKey, err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Reading key %s", keyID)

	eo.Command = fmt.Sprintf("apt-key export %s", keyID)
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stdout == "" {
		err = resources.NotFoundError{Type: Type, Name: keyID}
		return
	}

	name, err := aptKeyGetName(execResult.Stdout)
	if err != nil {
		return
	}

	aptKey.KeyID = keyID
	aptKey.PublicKey = execResult.Stdout
	aptKey.Name = name

	return
}

// Exists will report if a given key exists on a system.
func Exists(client client.Client, keyID string) (exists bool, err error) {
	client.Logger.Debugf("Checking if PPA %s exists", keyID)

	aptKey, err := Read(client, keyID)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	if aptKey.KeyID != "" {
		exists = true
	}

	return
}

// List will read all apt-key managed keys installed on a system.
func List(client client.Client) (aptKeys []AptKey, err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Listing all keys via apt-key list")

	eo.Command = fmt.Sprintf("apt-key list")
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	keys := aptKeyParseList(execResult.Stdout)
	for _, keyID := range keys {
		var aptKey AptKey
		aptKey, err = Read(client, keyID)
		if err != nil {
			return
		}

		aptKeys = append(aptKeys, aptKey)
	}

	return
}

// Create will create a key via apt-key.
func Create(client client.Client, createOpts CreateOpts) (err error) {
	var eo utils.ExecOptions
	var execResult utils.ExecResult

	client.Logger.Debugf("Creating apt-key key")

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("AptKey Create Options: %#v", createOpts)

	if createOpts.RemoteKeyFile != "" {
		var key string
		key, err = aptKeyGetRemoteKeyFile(createOpts.RemoteKeyFile)
		if err != nil {
			return
		}

		var tmpfile *os.File
		tmpfile, err = ioutil.TempFile("/tmp", "apt-key")
		if err != nil {
			return
		}
		defer os.Remove(tmpfile.Name())

		if _, err = tmpfile.Write([]byte(key)); err != nil {
			return
		}

		if err = tmpfile.Close(); err != nil {
			return
		}

		eo.Command = fmt.Sprintf("apt-key add %s", tmpfile.Name())
		execResult, err = utils.Exec(eo)
		if err != nil {
			return
		}

		if execResult.Stderr != "" {
			err = fmt.Errorf("unable to add key: %s", err)
			return
		}

	}

	if createOpts.KeyServer != "" {
		eo.Command = fmt.Sprintf("apt-key adv --keyserver %s --recv-keys %s",
			createOpts.KeyServer, createOpts.KeyID)

		execResult, err = utils.Exec(eo)
		if err != nil {
			return
		}

		if execResult.Stderr != "" {
			err = fmt.Errorf("unable to add key: %s", err)
			return
		}
	}

	return
}

// Update is not implemented.
//func Update() (err error) {
//	// AptKey cannot be updated.
//	return
//}

// Delete deletes a key managed by apt-key.
func Delete(client client.Client, keyID string) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Deleting key %s", keyID)

	eo.Command = fmt.Sprintf("apt-key del %s", keyID)
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	if execResult.Stderr != "" {
		err = fmt.Errorf("unable to delete key: %s", err)
	}

	return
}

// aptKeyGetRemoteKeyFile is an internal function that will
// download a key located at a remote URL.
func aptKeyGetRemoteKeyFile(v string) (key string, err error) {
	res, err := http.Get(v)
	if err != nil {
		return
	}

	k, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	key = string(k)

	return
}

// aptKeyGetShortID is an internal function that will print the
// short key ID of a public key.
func aptKeyGetShortID(key string) (fingerprint string, err error) {
	el, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(key))
	if err != nil {
		return
	}

	if len(el) == 0 {
		err = fmt.Errorf("Error determining fingerprint of key")
		return
	}

	fingerprint = el[0].PrimaryKey.KeyIdShortString()

	return
}

// aptKeyGetName is an internal function that will get the
// maintainer name of a public key.
func aptKeyGetName(key string) (name string, err error) {
	el, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(key))
	if err != nil {
		return
	}

	if len(el) == 0 {
		err = fmt.Errorf("Error determining userid of key")
		return
	}

	identities := el[0].Identities
	for k, _ := range identities {
		if name == "" {
			name = k
		}
	}

	return
}

func aptKeyParseList(list string) (keys []string) {
	keyRe := regexp.MustCompile("^pub.+/(.+) [0-9-]+$")
	for _, line := range strings.Split(list, "\n") {
		v := keyRe.FindStringSubmatch(line)
		if v != nil {
			keys = append(keys, v[1])
		}
	}

	return
}
