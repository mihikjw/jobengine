package filesystem

//FileSystem holds functions for interacting with a filesystem
type FileSystem interface {
	FileExists(filepath string) (bool, error)
	DeleteFile(filepath string) error
	ReadFile(filepath string) ([]byte, error)
	WriteFile(filepath string, data []byte) error
}

//NewFileSystem acts as a constructor for filesystem interfaces
func NewFileSystem(fsType string) FileSystem {
	switch {
	case fsType == "os":
		return new(OperatingSystem)
	default:
		return nil
	}
}
