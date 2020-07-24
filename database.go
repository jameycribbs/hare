/*
Package hare implements a simple DBMS that stores it's date
in newline-delimited json files.
*/
package hare

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const tblExt = ".json"

// Database contains attributes for the database path and all tables
// associated with the database.
type Database struct {
	path   string
	tables map[string]*Table
}

// OpenDB takes a directory path pointing to one or more json files and returns
// a pointer to a Database struct.
func OpenDB(dbPath string) (*Database, error) {
	db := new(Database)

	db.path = dbPath
	db.tables = make(map[string]*Table)

	files, err := ioutil.ReadDir(db.path)
	if err != nil {
		return nil, err
	}

	// Loop through all json files in database directory, initialize them,
	// and register them in the database.
	for _, file := range files {
		filename := file.Name()

		// If entry is sub dir, current dir, or parent dir, skip it.
		if file.IsDir() || filename == "." || filename == ".." {
			continue
		}

		if !strings.HasSuffix(filename, tblExt) {
			continue
		}

		tbl, err := openTable(db.path+"/"+filename, false)
		if err != nil {
			return nil, err
		}

		db.tables[strings.TrimSuffix(filename, tblExt)] = tbl
	}

	return db, nil
}

// Close closes all files associated with the database.
func (db *Database) Close() error {
	for _, tbl := range db.tables {
		tbl.Lock()

		if err := tbl.filePtr.Close(); err != nil {
			return err
		}

		tbl.Unlock()
	}

	db.path = ""
	db.tables = nil

	return nil
}

// CreateTable takes a table name and returns a pointer to a Table struct.
func (db *Database) CreateTable(tblName string) (*Table, error) {
	if db.TableExists(tblName) {
		return nil, errors.New("table already exists")
	}

	tbl, err := openTable(db.path+"/"+tblName+tblExt, true)
	if err != nil {
		return nil, err
	}

	db.tables[tblName] = tbl

	return db.tables[tblName], nil
}

// DropTable takes a table name and deletes the associated json file.
func (db *Database) DropTable(tblName string) error {
	tbl, err := db.GetTable(tblName)
	if err != nil {
		return err
	}

	tbl.Lock()
	defer tbl.Unlock()

	if err = tbl.filePtr.Close(); err != nil {
		return err
	}

	delete(db.tables, tblName)

	tbl = nil

	if err = os.Remove(db.path + "/" + tblName + tblExt); err != nil {
		return err
	}

	return nil
}

func (db *Database) GetTable(tblName string) (*Table, error) {
	tbl, ok := db.tables[tblName]
	if !ok {
		return nil, errors.New("table does not exist")
	}

	return tbl, nil
}

// TableExists takes a table name and returns true if the table exists,
// false if it does not.
func (db *Database) TableExists(tblName string) bool {
	_, ok := db.tables[tblName]

	return ok
}
