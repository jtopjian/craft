package cronentry

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/jtopjian/craft/client"
	"github.com/jtopjian/craft/resources"
	"github.com/jtopjian/craft/utils"
)

const Type = "CronEntry"

// CronEntry represents an entry found in a crontab.
type CronEntry struct {
	// Name is an arbitrary name for a cron entry.
	Name string

	// Command is the command which cron will run.
	Command string

	// Minute is the minute field of the cron entry.
	Minute string

	// Hour is the hour field of the cron entry.
	Hour string

	// DayOfMonth is the day of the month field of the cron entry.
	DayOfMonth string

	// Month is the month field of the cron entry.
	Month string

	// DayOfWeek is the day of the week field of the cron entry.
	DayOfWeek string
}

// CreateOpts represents options used to create a cron entry.
type CreateOpts struct {
	// Name is an arbitrary name for the cron entry.
	Name string

	// Command is the command which cron will run.
	Command string `required:"true"`

	// Minute is the minute field of the cron entry.
	Minute string `default:"*"`

	// Hour is the hour field of the cron entry.
	Hour string `default:"*"`

	// DayOfMonth is the day of the month field of the cron entry.
	DayOfMonth string `default:"*"`

	// Month is the month field of the cron entry.
	Month string `default:"*"`

	// DayOfWeek is the day of the week field of the cron entry.
	DayOfWeek string `default:"*"`
}

// UpdateOpts represents options used to create a cron entry.
type UpdateOpts struct {
	// Command is the command which cron will run.
	Command string

	// Minute is the minute field of the cron entry.
	Minute string

	// Hour is the hour field of the cron entry.
	Hour string

	// DayOfMonth is the day of the month field of the cron entry.
	DayOfMonth string

	// Month is the month field of the cron entry.
	Month string

	// DayOfWeek is the day of the week field of the cron entry.
	DayOfWeek string
}

// Read will retrieve information about an existing cron entry.
func Read(client client.Client, user, name string) (entry CronEntry, err error) {
	client.Logger.Debugf("Reading cron entry %s for user %s", name, user)

	entries, err := cronEntryGetEntries(user)
	if err != nil {
		return
	}

	e, err := cronEntryGetEntry(entries, name)
	if err != nil {
		err = resources.NotFoundError{Type: Type, Name: name}
		return
	}

	entry.Name = name
	entry.Command = e.Command
	entry.Minute = e.Minute
	entry.Hour = e.Hour
	entry.DayOfMonth = e.DayOfMonth
	entry.Month = e.Month
	entry.DayOfWeek = e.DayOfWeek

	return
}

// Exists will retrieve information about an existing cron entry.
func Exists(client client.Client, user, name string) (exists bool, err error) {
	client.Logger.Debugf("Checking if cron entry %s exists for user %s", name, user)

	_, err = Read(client, user, name)
	if err != nil {
		if _, ok := err.(resources.NotFoundError); ok {
			err = nil
		}
		return
	}

	exists = true

	return
}

// List will retrieve all cron entries for a given user.
func List(client client.Client, user string) (entries []CronEntry, err error) {
	client.Logger.Debugf("Listing all cron entries for user %s", user)

	e, err := cronEntryGetEntries(user)
	if err != nil {
		return
	}

	for _, v := range e {
		e, err := cronEntryParseLine(v)
		if err != nil {
			err = nil
			continue
		}

		entries = append(entries, e)
	}

	return
}

func Create(client client.Client, user string, createOpts CreateOpts) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Creating cron entry for user %s", user)

	if err = utils.BuildRequest(&createOpts); err != nil {
		return
	}

	client.Logger.Debugf("Cron Entry Create Options: %#v", createOpts)

	entries, err := cronEntryGetEntries(user)
	if err != nil {
		return
	}

	newEntry := fmt.Sprintf("%s %s %s %s %s %s # %s",
		createOpts.Minute, createOpts.Hour, createOpts.DayOfMonth, createOpts.Month,
		createOpts.DayOfWeek, createOpts.Command, createOpts.Name)
	entries = append(entries, newEntry)

	var tmpfile *os.File
	tmpfile, err = ioutil.TempFile("/tmp", "cron-entry")
	if err != nil {
		return
	}
	defer os.Remove(tmpfile.Name())

	newEntries := strings.Join(entries, "\n")
	newEntries = fmt.Sprintf("%s\n", newEntries)

	if _, err = tmpfile.Write([]byte(newEntries)); err != nil {
		return
	}

	eo.Command = fmt.Sprintf("crontab -u %s %s", user, tmpfile.Name())
	_, err = utils.Exec(eo)
	if err != nil {
		return
	}

	return
}

func Update(client client.Client, user, name string, updateOpts UpdateOpts) (err error) {
	client.Logger.Debugf("Updating cron entry %s for user %s", name, user)

	if err = utils.BuildRequest(&updateOpts); err != nil {
		return
	}

	client.Logger.Debugf("Cron Entry Update Options: %#v", updateOpts)

	err = Delete(client, user, name)
	if err != nil {
		return
	}

	createOpts := CreateOpts{
		Name:       name,
		Command:    updateOpts.Command,
		Minute:     updateOpts.Minute,
		Hour:       updateOpts.Hour,
		DayOfMonth: updateOpts.DayOfMonth,
		Month:      updateOpts.Month,
		DayOfWeek:  updateOpts.DayOfWeek,
	}

	return Create(client, user, createOpts)
}

func Delete(client client.Client, user, name string) (err error) {
	var eo utils.ExecOptions

	client.Logger.Debugf("Deleting cron entry %s for user %s", name, user)

	oldEntries, err := cronEntryGetEntries(user)
	if err != nil {
		return
	}

	var newEntries []string
	for _, entry := range oldEntries {
		if strings.Contains(entry, fmt.Sprintf("# %s", name)) {
			continue
		}
		newEntries = append(newEntries, entry)
	}

	var tmpfile *os.File
	tmpfile, err = ioutil.TempFile("/tmp", "cron-entry")
	if err != nil {
		return
	}
	defer os.Remove(tmpfile.Name())

	v := strings.Join(newEntries, "\n")
	v = fmt.Sprintf("%s\n", v)
	if _, err = tmpfile.Write([]byte(v)); err != nil {
		return
	}

	eo.Command = fmt.Sprintf("crontab -u %s %s", user, tmpfile.Name())
	_, err = utils.Exec(eo)
	if err != nil {
		return
	}

	return
}

func cronEntryGetEntries(user string) (entries []string, err error) {
	var eo utils.ExecOptions

	eo.Command = fmt.Sprintf("crontab -u %s -l", user)
	execResult, err := utils.Exec(eo)
	if err != nil {
		return
	}

	for _, v := range strings.Split(execResult.Stdout, "\n") {
		if v == "" {
			continue
		}

		entries = append(entries, v)
	}

	return
}

func cronEntryExists(entries []string, name string) (exists bool) {
	for _, line := range entries {
		if strings.Contains(line, fmt.Sprintf("# %s", name)) {
			exists = true
		}
	}

	return
}

func cronEntryGetEntry(entries []string, name string) (entry CronEntry, err error) {
	for _, line := range entries {
		if strings.Contains(line, fmt.Sprintf("# %s", name)) {
			entry, err = cronEntryParseLine(line)
			if err != nil {
				return
			}

			return
		}
	}

	err = fmt.Errorf("Entry %s not found", name)
	return
}

func cronEntryParseLine(line string) (entry CronEntry, err error) {
	re := regexp.MustCompile("^(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+(\\S+)\\s+#\\s+(.+)$")
	v := re.FindStringSubmatch(line)
	if v == nil {
		err = fmt.Errorf("Unable to parse entry [%s]", line)
		return
	}

	entry.Minute = v[1]
	entry.Hour = v[2]
	entry.DayOfMonth = v[3]
	entry.Month = v[4]
	entry.DayOfWeek = v[5]
	entry.Command = v[6]
	entry.Name = v[7]

	return
}
