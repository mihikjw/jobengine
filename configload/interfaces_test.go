package configload

import "testing"

func TestNewConfigParser1(t *testing.T) {
	result := NewConfigParser("yaml")
	if result == nil {
		t.Error("TestNewConfigParser1: Result Unexpectedly Nil")
	}
}

func TestNewConfigParser2(t *testing.T) {
	result := NewConfigParser("blobby")
	if result != nil {
		t.Error("TestNewConfigParser2: Result Unexpectedly Not Nil")
	}
}
