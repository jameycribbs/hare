package hare

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	deleteTestFiles()
	copyTestFiles()

	os.Exit(m.Run())

	deleteTestFiles()
}

func TestOpenDB(t *testing.T) {
	db := openTestDB()
	defer db.Close()

	if reflect.TypeOf(db) != reflect.TypeOf(&Database{}) {
		t.Errorf("want %v; got %v", reflect.TypeOf(&Database{}), reflect.TypeOf(db))
	}
}

func TestCloseDB(t *testing.T) {
	db := openTestDB()
	db.Close()

	if db.TableExists("test") {
		t.Errorf("want %v; got %v", false, true)
	}
}

func TestTableExists(t *testing.T) {
	db := openTestDB()
	defer db.Close()

	if !db.TableExists("test") {
		t.Errorf("want %v; got %v", true, false)
	}

	if db.TableExists("does_not_exist") {
		t.Errorf("want %v; got %v", true, false)
	}
}

func TestGetTable(t *testing.T) {
	db := openTestDB()
	defer db.Close()

	tbl, err := db.GetTable("test")
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(tbl) != reflect.TypeOf(&Table{}) {
		t.Errorf("want %v; got %v", reflect.TypeOf(&Table{}), reflect.TypeOf(tbl))
	}
}

func TestCreateTable(t *testing.T) {
	db := openTestDB()
	defer db.Close()

	tbl, err := db.CreateTable("new_table")
	if err != nil {
		t.Fatal(err)
	}

	if reflect.TypeOf(tbl) != reflect.TypeOf(&Table{}) {
		t.Errorf("want %v; got %v", reflect.TypeOf(&Table{}), reflect.TypeOf(tbl))
	}

	if !db.TableExists("new_table") {
		t.Errorf("want %v; got %v", true, false)
	}
}

func TestDropTable(t *testing.T) {
	db := openTestDB()
	defer db.Close()

	if !db.TableExists("deletable_table") {
		t.Errorf("want %v; got %v", true, false)
	}

	err := db.DropTable("deletable_table")
	if err != nil {
		t.Fatal(err)
	}

	if db.TableExists("deletable_table") {
		t.Errorf("want %v; got %v", true, false)
	}

}

func openTestDB() *Database {
	db, err := OpenDB("test_data")
	if err != nil {
		panic(err)
	}

	return db
}

func deleteTestFiles() {
	filesToDelete := []string{
		"test_data/test.json",
		"test_data/new_table.json",
		"test_data/deletable_table.json",
	}

	for _, fileToDelete := range filesToDelete {
		if _, err := os.Stat(fileToDelete); err == nil {
			if err = os.Remove(fileToDelete); err != nil {
				panic(err)
			}
		}
	}
}

func copyTestFiles() {
	var cmds []*exec.Cmd

	cmds = append(cmds, exec.Command("cp", "test_data/test_default.json", "test_data/test.json"))
	cmds = append(cmds, exec.Command("cp", "test_data/test_default.json", "test_data/deletable_table.json"))

	for _, cmd := range cmds {
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}
