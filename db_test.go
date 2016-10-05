package hare

import (
	"fmt"
	"os"
	"testing"
)

const dataDir = "test_data"

var db *database

type Foo struct {
	ID  int    `json:"id"`
	Bar string `json:"bar"`
}

func (foo *Foo) SetID(id int) {
	foo.ID = id
}

func (foo *Foo) GetID() int {
	return foo.ID
}

func setup() {
	_, err := os.Stat(dataDir)
	if err == nil {
		err = os.RemoveAll(dataDir)
		if err != nil {
			fmt.Println("Failed to remove test directory:", err)
			os.Exit(1)
		}
	} else if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = os.MkdirAll("test_data", 0777); err != nil {
		fmt.Println("Failed to make test directory:", err)
		os.Exit(1)
	}

	db, err = OpenDB("test_data")
	if err != nil {
		fmt.Println("Failed to open database:", err)
		os.Exit(1)
	}
}

func teardown() {
	err := os.RemoveAll(dataDir)
	if err != nil {
		fmt.Println("Failed to remove test directory:", err)
		os.Exit(1)
	}
}

//-----------------------------------------------------------------------------
// Tests
//-----------------------------------------------------------------------------
func TestCreate(t *testing.T) {
	setup()

	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("CreateTable:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("Create:", err)
	}

	foo := Foo{}

	err = foosTbl.Find(id, &foo)
	if err != nil {
		t.Error("Find failed:", err)
	}

	if foo.Bar != "test" {
		t.Error("Expected 'test', got ", foo.Bar)
	}

	teardown()
}

func TestUpdate(t *testing.T) {
	setup()

	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("CreateTable:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("Create:", err)
	}

	foo := Foo{}

	err = foosTbl.Find(id, &foo)
	if err != nil {
		t.Error("Find failed:", err)
	}

	foo.Bar = "test2"

	err = foosTbl.Update(&foo)
	if err != nil {
		t.Error("Update:", err)
	}

	err = foosTbl.Find(id, &foo)
	if err != nil {
		t.Error("Find failed:", err)
	}

	if foo.Bar != "test2" {
		t.Error("Expected 'test2', got ", foo.Bar)
	}

	teardown()
}

func TestDestroy(t *testing.T) {
	setup()

	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("CreateTable:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("Create:", err)
	}

	err = foosTbl.Destroy(id)
	if err != nil {
		t.Error("Destroy:", err)
	}

	foo := Foo{}

	err = foosTbl.Find(id, &foo)
	if err != nil {
		if err.Error() != "Record not found!" {
			t.Error("Expected Find error to be 'Record not found!', got ", err)
		}
	} else {
		t.Error("Expected Find error, got no error.")
	}

	teardown()
}
