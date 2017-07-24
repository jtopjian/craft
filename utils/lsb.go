package utils

import (
	"regexp"
)

type LSBInfo struct {
	DistributorID string
	Description   string
	Release       string
	Codename      string
}

func GetLSBInfo() (LSBInfo, error) {
	var eo ExecOptions

	eo.Command = "lsb_release -a"
	execResult, err := Exec(eo)
	if err != nil {
		return LSBInfo{}, err
	}

	return parseLSBInfo(execResult.Stdout), nil
}

func parseLSBInfo(stdout string) LSBInfo {
	var lsbInfo LSBInfo

	distributorRe := regexp.MustCompile("Distributor ID:\\s+(.+)\n")
	descriptionRe := regexp.MustCompile("Description:\\s+(.+)\n")
	releaseRe := regexp.MustCompile("Release:\\s+(.+)\n")
	codenameRe := regexp.MustCompile("Codename:\\s+(.+)\n")

	if v := distributorRe.FindStringSubmatch(stdout); len(v) > 1 {
		lsbInfo.DistributorID = v[1]
	}

	if v := descriptionRe.FindStringSubmatch(stdout); len(v) > 1 {
		lsbInfo.Description = v[1]
	}

	if v := releaseRe.FindStringSubmatch(stdout); len(v) > 1 {
		lsbInfo.Release = v[1]
	}

	if v := codenameRe.FindStringSubmatch(stdout); len(v) > 1 {
		lsbInfo.Codename = v[1]
	}

	return lsbInfo
}
