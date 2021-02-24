package disk

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func runTestFns(t *testing.T, tests []func(t *testing.T)) {
	for i, fn := range tests {
		testSetup(t)
		t.Run(strconv.Itoa(i), fn)
		testTeardown(t)
	}
}

func newTestDisk(t *testing.T) *Disk {
	dsk, err := New("./testdata", ".json")
	if err != nil {
		t.Fatal(err)
	}

	return dsk
}

func newTestTableFile(t *testing.T) *tableFile {
	filePtr, err := os.OpenFile("./testdata/contacts.json", os.O_RDWR, 0660)
	if err != nil {
		t.Fatal(err)
	}

	tf, err := newTableFile("contacts", filePtr)
	if err != nil {
		t.Fatal(err)
	}

	return tf
}

func testSetup(t *testing.T) {
	testRemoveFiles(t)

	cmd := exec.Command("cp", "./testdata/contacts.bak", "./testdata/contacts.json")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}

func testTeardown(t *testing.T) {
	testRemoveFiles(t)
}

func testRemoveFiles(t *testing.T) {
	filesToRemove := []string{"contacts.json", "newtable.json"}

	for _, f := range filesToRemove {
		err := os.Remove("./testdata/" + f)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
}
