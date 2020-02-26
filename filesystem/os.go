package filesystem

import (
	"fmt"
	"io/ioutil"
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

//DeleteFile is a small wrapper for os.Remove, errors are type *PathError
func (o *OperatingSystem) DeleteFile(filepath string) error {
	return os.Remove(filepath)
}

//ReadFile loads a file into memory and returns as a byte slice
func (o *OperatingSystem) ReadFile(filepath string) ([]byte, error) {
	rawFileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return rawFileContent, nil
}

//WriteFile writes a file onto disk
func (o *OperatingSystem) WriteFile(filepath string, data []byte) error {
	if len(filepath) <= 0 || len(data) <= 0 {
		return fmt.Errorf("Invalid Args")
	}

	return ioutil.WriteFile(filepath, data, 0644)
}

//GetEnv returns the requested environment variable or nothing if it was not found
func (o *OperatingSystem) GetEnv(name string) string {
	return os.Getenv(name)
}

//Open returns a handle to a file
func (o *OperatingSystem) Open(filepath string) (*os.File, error) {
	return os.Open(filepath)
}
