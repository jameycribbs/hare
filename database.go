package hare

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const tblExt = ".json"

type Database struct {
	path   string
	tables map[string]*table
}

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

func (db *Database) TableExists(tblName string) bool {
	if db.tables[tblName] == nil {
		return false
	} else {
		return true
	}
}

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

func (db *Database) CreateTable(tblName string) (*table, error) {
	var err error

	if db.TableExists(tblName) {
		return nil, errors.New("Table already exists!")
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

func (db *Database) GetTable(tblName string) (*table, error) {
	if !db.TableExists(tblName) {
		return nil, errors.New("Table does not exist!")
	}

	return db.tables[tblName], nil
}

func (db *Database) Close() {
	for _, tbl := range db.tables {
		tbl.Lock()

		if err := tbl.filePtr.Close(); err != nil {
			panic(err)
		}

		tbl.Unlock()
	}
}
