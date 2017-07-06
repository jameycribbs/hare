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
	tables map[string]*table
}

// OpenDB takes a directory path pointing to one or more json files and returns
// a database connection.
func OpenDB(dbPath string) (*Database, error) {
	var err error

	db := new(Database)
	db.path = dbPath

	db.tables = make(map[string]*table)

	files, _ := ioutil.ReadDir(db.path)

	for _, file := range files {
		if !file.IsDir() {
			if file.Name() != "." && file.Name() != ".." {
				tbl := table{}

				tbl.filePtr, err = os.OpenFile(db.path+"/"+file.Name(), os.O_RDWR, 0660)
				if err != nil {
					return nil, err
				}

				tbl.initIndex()
				tbl.initLastID()

				db.tables[strings.TrimSuffix(file.Name(), tblExt)] = &tbl
			}
		}
	}

	return db, nil
}

// TableExists takes a table name and returns true if the table exits,
// false if it does not.
func (db *Database) TableExists(tblName string) bool {
	if db.tables[tblName] == nil {
		return false
	}

	return true
}

// DropTable takes a table name and deletes the associated json file.
func (db *Database) DropTable(tblName string) error {
	var err error

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

// CreateTable takes a table name and creates an associated json file.
func (db *Database) CreateTable(tblName string) (*table, error) {
	var err error

	if db.TableExists(tblName) {
		return nil, errors.New("table already exists")
	}

	tbl := table{}

	tbl.filePtr, err = os.OpenFile(db.path+"/"+tblName+tblExt, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return nil, err
	}

	tbl.initIndex()
	tbl.initLastID()

	db.tables[tblName] = &tbl

	return db.tables[tblName], nil
}

// GetTable takes a table name and returns a reference to that table.
func (db *Database) GetTable(tblName string) (*table, error) {
	if !db.TableExists(tblName) {
		return nil, errors.New("table does not exist")
	}

	return db.tables[tblName], nil
}

// Close closes all json files associated with the database.
func (db *Database) Close() {
	for _, tbl := range db.tables {
		tbl.Lock()

		if err := tbl.filePtr.Close(); err != nil {
			panic(err)
		}

		tbl.Unlock()
	}
}
