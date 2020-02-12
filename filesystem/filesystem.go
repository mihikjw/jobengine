package filesystem

//FileSystem holds functions for interacting with a filesystem
type FileSystem interface {
	FileExists(filepath string) (bool, error)
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
