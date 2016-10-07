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

func TestMain(m *testing.M) {
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

	result := m.Run()

	db.Close()

	err = os.RemoveAll(dataDir)
	if err != nil {
		fmt.Println("Failed to remove test directory:", err)
		os.Exit(1)
	}

	os.Exit(result)
}

//-----------------------------------------------------------------------------
// Tests
//-----------------------------------------------------------------------------
func TestCreateTable(t *testing.T) {
	_, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestCreateTable:", err)
	}

	if !db.TableExists("foos") {
		t.Error("TestCreateTable:", err)
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestCreateTable:", err)
	}
}

func TestDestroyTable(t *testing.T) {
	_, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestDestroyTable:", err)
	}

	if !db.TableExists("foos") {
		t.Error("TestDestroyTable:", err)
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestDestroyTable:", err)
	}

	if db.TableExists("foos") {
		t.Error("TestDestroyTable:", err)
	}
}
