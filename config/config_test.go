package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func createFileForTest(t *testing.T) *os.File {
	data := []byte(`{"foo": "bar", "denny": {"sister": "jenny"}}`)
	path := filepath.Join(os.TempDir(), fmt.Sprintf("file.%d", time.Now().UnixNano()))
	fh, err := os.Create(path)
	if err != nil {
		t.Error(err)
	}
	_, err = fh.Write(data)
	if err != nil {
		t.Error(err)
	}

	return fh
}

func TestLoadWithGoodFile(t *testing.T) {
	fh := createFileForTest(t)
	path := fh.Name()
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	// Create new config
	// Load file source
	if err := New(path); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
}

func TestReadValue(t *testing.T)  {
	TestLoadWithGoodFile(t)
	ov := "bar"
	v := GetString("foo")
	if v != ov {
		t.Fatalf("Expected bar error but got %v", v)
	}
}
func TestLoadWithInvalidFile(t *testing.T) {
	fh := createFileForTest(t)
	path := fh.Name()
	defer func() {
		fh.Close()
		os.Remove(path)
	}()


	// Load file source
	err := New(path,
		"/i/do/not/exists.json")

	if err == nil {
		t.Fatal("Expected error but none !")
	}
	if !strings.Contains(fmt.Sprintf("%v", err), "/i/do/not/exists.json") {
		t.Fatalf("Expected error to contain the unexisting file but got %v", err)
	}
}


func TestWithScanToPointer(t *testing.T) {
	type Denny struct {
		Sister string
	}
	var (
		v string
		denny = &Denny{}
	)
	TestLoadWithGoodFile(t)
	ov := "bar"

	err := Scan(&v, "foo")

	if err != nil {
		t.Fatalf("Expect value but got error reading %v", err)
	}
	if v != ov {
		t.Fatalf("Expected bar error but got %v", v)
	}
	err = Scan(denny, "denny")

	if err != nil {
		t.Fatalf("Expect value but got error reading %v", err)
	}
	if denny.Sister != "jenny" {
		t.Fatalf("Expected jenny error but got %v", v)
	}
}