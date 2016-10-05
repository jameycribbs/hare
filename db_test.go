package hare

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

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
	var err error

	files, _ := ioutil.ReadDir("test_data")
	for _, file := range files {
		if !file.IsDir() {
			err = os.Remove("test_data/" + file.Name())
			if err != nil {
				fmt.Println("Failed to remove file", err)
			}
		}
	}

	db, err = OpenDB("test_data")
	if err != nil {
		fmt.Println("Failed to open database:", err)
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

}
