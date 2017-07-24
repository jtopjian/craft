package resources

import (
	"fmt"
)

// NotFoundError is returned when a resource was not found.
type NotFoundError struct {
	Type string
	Name string
}

func (r NotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", r.Type, r.Name)
}
