package configload

import (
	"fmt"
	"testing"

	"github.com/MichaelWittgreffe/jobengine/models"
)

type MockConfigLoad struct {
	LoadFromFileResult *models.Config
	LoadFromFileError  error
	SaveToFileError    error
}

//LoadFromFile returns mock values from struct
func (l *MockConfigLoad) LoadFromFile(version float64) (*models.Config, error) {
	return l.LoadFromFileResult, l.LoadFromFileError
}

//SaveToFile returns mock values from struct
func (l *MockConfigLoad) SaveToFile(cfg *models.Config) error {
	return l.SaveToFileError
}

//TestLoadConfig1 tests the correct execution of the function
func TestLoadConfig1(t *testing.T) {
	mLoader := &MockConfigLoad{
		LoadFromFileResult: &models.Config{},
		LoadFromFileError:  nil,
	}

	cfg, err := LoadConfig(mLoader)

	if err != nil {
		t.Errorf("TestLoadConfig1: Unexpected Error: %s", err.Error())
	}

	if cfg == nil {
		t.Error("TestLoadConfig1: Config Unexpectedly Nil")
	}
}

//TestLoadConfig2 cannot load the config from disk, loads & writes a default config instead
func TestLoadConfig2(t *testing.T) {
	mLoader := &MockConfigLoad{
		LoadFromFileResult: &models.Config{CryptoSecret: "32byteteststring_gk93;[]0=4gjday"},
		LoadFromFileError:  fmt.Errorf("Not Found"),
		SaveToFileError:    nil,
	}

	cfg, err := LoadConfig(mLoader)

	if err != nil {
		t.Errorf("TestLoadConfig2: Unexpected Error: %s", err.Error())
	}

	if cfg == nil {
		t.Error("TestLoadConfig2: Config Unexpectedly Nil")
	}
}

//TestLoadConfig3 encounters an error loading the file from disk (other than cannot be found)
func TestLoadConfig3(t *testing.T) {
	mLoader := &MockConfigLoad{
		LoadFromFileResult: nil,
		LoadFromFileError:  fmt.Errorf("Test Error"),
	}

	cfg, err := LoadConfig(mLoader)

	if err == nil {
		t.Errorf("TestLoadConfig3: Unexpectedly No Error")
	}

	if cfg != nil {
		t.Error("TestLoadConfig1: Config Unexpectedly Not Nil")
	}
}

//TestLoadConfig4 cannot save a default config to disk
func TestLoadConfig4(t *testing.T) {
	mLoader := &MockConfigLoad{
		LoadFromFileResult: &models.Config{CryptoSecret: "32byteteststring_gk93;[]0=4gjday"},
		LoadFromFileError:  fmt.Errorf("Not Found"),
		SaveToFileError:    fmt.Errorf("Test Error"),
	}

	cfg, err := LoadConfig(mLoader)

	if err == nil {
		t.Errorf("TestLoadConfig4: Unexpectedly No Error")
	}

	if cfg != nil {
		t.Error("TestLoadConfig4: Config Unexpectedly Not Nil")
	}
}

//TestLoadConfig5 has invalid arg given to function
func TestLoadConfig5(t *testing.T) {
	cfg, err := LoadConfig(nil)

	if err == nil {
		t.Errorf("TestLoadConfig5: Unexpectedly No Error")
	}

	if cfg != nil {
		t.Error("TestLoadConfig5: Config Unexpectedly Not Nil")
	}
}
