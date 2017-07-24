package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

func ChownR(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}

func ChmodR(path string, mode os.FileMode) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chmod(name, mode)
		}
		return err
	})
}

func UsernameToID(username string) (uid int, err error) {
	v, err := user.Lookup(username)
	if err != nil {
		return uid, err
	}

	uid, err = strconv.Atoi(v.Uid)
	return
}

func UIDToName(uid int) (username string, err error) {
	uidString := strconv.Itoa(uid)
	v, err := user.LookupId(uidString)
	if err != nil {
		return
	}

	username = v.Uid
	return
}

func GroupToID(groupname string) (gid int, err error) {
	v, err := user.LookupGroup(groupname)
	if err != nil {
		err = fmt.Errorf("Unknown group: %s", groupname)
		return gid, err
	}

	gid, err = strconv.Atoi(v.Gid)
	return
}

func GIDToName(gid int) (group string, err error) {
	gidString := strconv.Itoa(gid)

	v, err := user.LookupGroupId(gidString)
	if err != nil {
		return
	}

	group = v.Gid
	return
}

func GetUIDGID(username, group string) (uid, gid int, err error) {
	uid, err = UsernameToID(username)
	if err != nil {
		return
	}

	gid, err = GroupToID(group)

	return
}

func GetFileOwner(path string) (uid, gid int, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}

	uid = int(fi.Sys().(*syscall.Stat_t).Uid)
	gid = int(fi.Sys().(*syscall.Stat_t).Gid)
	return
}

func StringToMode(v string) (mode os.FileMode, err error) {
	m, err := strconv.Atoi(v)
	if err != nil {
		return
	}

	mode = os.FileMode(m)

	return
}
