package resources

import (
	"fmt"
)

type State string

const (
	Absent  State = "ABSENT"
	Active  State = "ACTIVE"
	Error   State = "ERROR"
	Latest  State = "LATEST"
	Present State = "PRESENT"
	Running State = "RUNNING"
	Stopped State = "STOPPED"
	Update  State = "UPDATE"
)

type Resource interface {
	Validate() error
	Exists() (bool, error)
	Create() error
	Read() error
	Update() error
	Delete() error
}

// NotFoundError is returned when a resource was not found.
type NotFoundError struct {
	Type string
	Name string
}

func (r NotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", r.Type, r.Name)
}
