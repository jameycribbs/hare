package hare_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jameycribbs/hare"
)

const dataDir = "test_data"

var db *hare.Database

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

func fooFromRecMap(recMap map[string]interface{}) Foo {
	return Foo{
		ID:  int(recMap["id"].(float64)),
		Bar: recMap["bar"].(string),
	}
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

	db, err = hare.OpenDB("test_data")
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
