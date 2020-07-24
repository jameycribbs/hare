package hare

import (
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"testing"
)

func TestAllDatabaseTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//OpenDB...

			db := openTestDB()
			defer db.Close()

			if reflect.TypeOf(db) != reflect.TypeOf(&Database{}) {
				t.Errorf("want %v; got %v", reflect.TypeOf(&Database{}), reflect.TypeOf(db))
			}
		},
		func(t *testing.T) {
			//Close...

			db := openTestDB()
			db.Close()

			if db.TableExists("test") {
				t.Errorf("want %v; got %v", false, true)
			}
		},
		func(t *testing.T) {
			//TableExists...

			db := openTestDB()
			defer db.Close()

			if !db.TableExists("test") {
				t.Errorf("want %v; got %v", true, false)
			}

			if db.TableExists("does_not_exist") {
				t.Errorf("want %v; got %v", true, false)
			}
		},
		func(t *testing.T) {
			//GetTable...

			db := openTestDB()
			defer db.Close()

			tbl, err := db.GetTable("test")
			if err != nil {
				t.Fatal(err)
			}

			if reflect.TypeOf(tbl) != reflect.TypeOf(&Table{}) {
				t.Errorf("want %v; got %v", reflect.TypeOf(&Table{}), reflect.TypeOf(tbl))
			}
		},
		func(t *testing.T) {
			//CreateTable...

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

			if _, err := os.Stat("test_data/new_table.json"); err != nil {
				t.Errorf("want %v; got %v", nil, err)
			}
		},
		func(t *testing.T) {
			//DropTable...

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

			if _, err := os.Stat("test_data/deletable.json"); err == nil {
				t.Errorf("want %v; got %v", "stat test_data/deletable_table.json: no such file or directory", nil)
			}
		},
	}

	for i, fn := range tests {
		databaseTestSetup()
		t.Run(strconv.Itoa(i), fn)
		databaseTestTeardown()
	}
}

func openTestDB() *Database {
	db, err := OpenDB("test_data")
	if err != nil {
		panic(err)
	}

	return db
}

func databaseTestSetup() {
	deleteDatabaseTestFiles()
	copyDatabaseTestFiles()
}

func databaseTestTeardown() {
	deleteDatabaseTestFiles()
}

func deleteDatabaseTestFiles() {
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

func copyDatabaseTestFiles() {
	var cmds []*exec.Cmd

	cmds = append(cmds, exec.Command("cp", "test_data/test_default.json", "test_data/test.json"))
	cmds = append(cmds, exec.Command("cp", "test_data/test_default.json", "test_data/deletable_table.json"))

	for _, cmd := range cmds {
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}
