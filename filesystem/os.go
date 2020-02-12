package filesystem

import (
	"os"
)

//OperatingSystem holds methods for interacting with the operating system filesystem
type OperatingSystem struct{}

//FileExists returns bool whether the given file exists or not
func (o *OperatingSystem) FileExists(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}
