package provider

import "testing"

func TestUnsupportedProvider(t *testing.T) {
	_, e := NewProvider("will-not-exists-never-1546456", map[string]string{})
	if e.Error() != "will-not-exists-never-1546456 not supported." {
		t.Errorf("We expect an error because will-not-exists-never-1546456 is not supported")
	}
}

func TestCreateFakeProvider(t *testing.T) {
	_, e := NewProvider("fake", map[string]string{})
	if e != nil {
		t.Error(e)
	}
}
