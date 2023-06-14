/*
This is a test file for gondi.go, it only tests two functions so far, so it should be fixed soon, as well as adding tests for all the other files.
*/
package gondi

import "testing"

func TestGetVersion(t *testing.T) {
	InitLibrary("")

	str := GetVersion()
	if len(str) == 0 {
		t.Error("Version string empty")
	}
	if len(str) < 7 {
		t.Error("Version string too short")
	}
}

func TestNewMetadataFrame(t *testing.T) {
	InitLibrary("")
	testString := `<data value="The quick brown fox jumps over the lazy dog" />`

	mf := NewMetadataFrame(testString)
	if mf == nil {
		t.Error("NewMetadataFrame returned nil")
		return
	}

	if mf.Data == nil {
		t.Error("NewMetadataFrame.Data is nil")
	}

	if mf.GetData() != testString {
		t.Errorf("GetData() returned %q, want %q", mf.GetData(), testString)
	}

	if mf.Length != int32(len(testString)) {
		t.Errorf("Length is %d, want %d", mf.Length, len(testString))
	}
}
