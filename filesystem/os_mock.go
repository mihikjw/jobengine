package filesystem

import "os"

//MockFileSystem is a mock class for the interface FileSystem
type MockFileSystem struct {
	FileExistsBoolResult  bool
	FileExistsErrorResult error
	DeleteFileResult      error
	ReadFileByteResult    []byte
	ReadFileErrorResult   error
	WiteFileResult        error
	GetEnvResult          string
	OpenResult            *os.File
	OpenError             error
}

//FileExists returns values from the struct
func (fs *MockFileSystem) FileExists(filepath string) (bool, error) {
	return fs.FileExistsBoolResult, fs.FileExistsErrorResult
}

//DeleteFile returns value from the struct
func (fs *MockFileSystem) DeleteFile(filepath string) error {
	return fs.DeleteFileResult
}

//ReadFile returns values from the struct
func (fs *MockFileSystem) ReadFile(filepath string) ([]byte, error) {
	return fs.ReadFileByteResult, fs.ReadFileErrorResult
}

//WriteFile returns values from the struct
func (fs *MockFileSystem) WriteFile(filepath string, data []byte) error {
	return fs.WiteFileResult
}

//GetEnv returns values from the struct
func (fs *MockFileSystem) GetEnv(name string) string {
	return fs.GetEnvResult
}

//Open returns
func (fs *MockFileSystem) Open(filepath string) (*os.File, error) {
	return fs.OpenResult, fs.OpenError
}
