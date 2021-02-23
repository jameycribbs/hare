package hare

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/jameycribbs/hare/datastores/disk"
	"github.com/jameycribbs/hare/datastores/ram"
)

type Contact struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func (c *Contact) GetID() int {
	return c.ID
}

func (c *Contact) SetID(id int) {
	c.ID = id
}

func (c *Contact) AfterFind(db *Database) error {
	*c = Contact(*c)

	return nil
}

func runTestFns(t *testing.T, testFns []func(*Database) func(*testing.T)) {
	for i, fn := range testFns {
		tstNum := strconv.Itoa(i)

		testSetup(t)

		diskDS, err := disk.New("./testdata", ".json")
		if err != nil {
			t.Fatal(err)
		}

		diskDB, err := New(diskDS)
		if err != nil {
			t.Fatal(err)
		}
		defer diskDB.Close()

		t.Run(fmt.Sprintf("disk/%s", tstNum), fn(diskDB))

		ramDS, err := ram.New(seedData())
		if err != nil {
			t.Fatal(err)
		}

		ramDB, err := New(ramDS)
		if err != nil {
			t.Fatal(err)
		}
		defer ramDB.Close()

		t.Run(fmt.Sprintf("ram/%s", tstNum), fn(ramDB))

		testTeardown(t)
	}
}

func checkErr(t *testing.T, wantErr error, gotErr error) {
	if !errors.Is(gotErr, wantErr) {
		t.Errorf("want %v; got %v", wantErr, gotErr)
	}
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

func seedData() map[string]map[int]string {
	tblMap := make(map[string]map[int]string)
	contactsMap := make(map[int]string)

	contactsMap[1] = `{"id":1,"first_name":"John","last_name":"Doe","age":37}`
	contactsMap[2] = `{"id":2,"first_name":"Abe","last_name":"Lincoln","age":52}`
	contactsMap[3] = `{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`
	contactsMap[4] = `{"id":4,"first_name":"Helen","last_name":"Keller","age":25}`

	tblMap["contacts"] = contactsMap

	return tblMap
}
