//Package hare implements a simple DBMS that stores it's data
//in newline-delimited json files.
package hare

import (
	"encoding/json"
	"errors"
	"sync"
)

var (
	ErrNoTable     = errors.New("hare: no table with that name found")
	ErrTableExists = errors.New("hare: table with that name already exists")
)

type Record interface {
	SetID(int)
	GetID() int
	AfterFind()
}

type datastorage interface {
	Close() error
	CreateTable(string) error
	DeleteRec(string, int) error
	GetLastID(string) (int, error)
	IDs(string) ([]int, error)
	InsertRec(string, int, []byte) error
	ReadRec(string, int) ([]byte, error)
	RemoveTable(string) error
	TableExists(string) bool
	TableNames() []string
	UpdateRec(string, int, []byte) error
}

type Database struct {
	store   datastorage
	locks   map[string]*sync.RWMutex
	lastIDs map[string]int
}

func New(ds datastorage) (*Database, error) {
	db := &Database{store: ds}
	db.locks = make(map[string]*sync.RWMutex)
	db.lastIDs = make(map[string]int)

	for _, tableName := range db.store.TableNames() {
		db.locks[tableName] = &sync.RWMutex{}

		lastID, err := db.store.GetLastID(tableName)
		if err != nil {
			return nil, err
		}
		db.lastIDs[tableName] = lastID
	}

	return db, nil
}

func (db *Database) Close() error {
	for _, lock := range db.locks {
		lock.Lock()
	}

	if err := db.store.Close(); err != nil {
		return err
	}

	for _, lock := range db.locks {
		lock.Unlock()
	}

	db.store = nil
	db.locks = nil
	db.lastIDs = nil

	return nil
}

func (db *Database) CreateTable(tableName string) error {
	if db.TableExists(tableName) {
		return ErrTableExists
	}

	if err := db.store.CreateTable(tableName); err != nil {
		return nil
	}

	db.locks[tableName] = &sync.RWMutex{}

	lastID, err := db.store.GetLastID(tableName)
	if err != nil {
		return err
	}
	db.lastIDs[tableName] = lastID

	return nil
}

func (db *Database) Delete(tableName string, id int) error {
	if !db.TableExists(tableName) {
		return ErrNoTable
	}

	db.locks[tableName].Lock()
	defer db.locks[tableName].Unlock()

	if err := db.store.DeleteRec(tableName, id); err != nil {
		return err
	}

	return nil
}

func (db *Database) DropTable(tableName string) error {
	if !db.TableExists(tableName) {
		return ErrNoTable
	}

	db.locks[tableName].Lock()
	defer db.locks[tableName].Unlock()

	if err := db.store.RemoveTable(tableName); err != nil {
		db.locks[tableName].Unlock()
		return err
	}

	delete(db.lastIDs, tableName)

	db.locks[tableName].Unlock()

	delete(db.locks, tableName)

	return nil
}

func (db *Database) Find(tableName string, id int, rec Record) error {
	if !db.TableExists(tableName) {
		return ErrNoTable
	}

	db.locks[tableName].RLock()
	defer db.locks[tableName].RUnlock()

	rawRec, err := db.store.ReadRec(tableName, id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawRec, rec)
	if err != nil {
		return err
	}

	rec.AfterFind()

	return nil
}

func (db *Database) IDs(tableName string) ([]int, error) {
	if !db.TableExists(tableName) {
		return nil, ErrNoTable
	}

	db.locks[tableName].Lock()
	defer db.locks[tableName].Unlock()

	ids, err := db.store.IDs(tableName)
	if err != nil {
		return nil, err
	}

	return ids, err
}

func (db *Database) Insert(tableName string, rec Record) (int, error) {
	if !db.TableExists(tableName) {
		return 0, ErrNoTable
	}

	db.locks[tableName].Lock()
	defer db.locks[tableName].Unlock()

	id := db.incrementLastID(tableName)
	rec.SetID(id)

	rawRec, err := json.Marshal(rec)
	if err != nil {
		return 0, err
	}

	if err := db.store.InsertRec(tableName, id, rawRec); err != nil {
		return 0, err
	}

	return id, nil
}

// TableExists takes a table name and returns true if the table exists,
// false if it does not.
func (db *Database) TableExists(tableName string) bool {
	return db.tableExists(tableName) && db.store.TableExists(tableName)
}

func (db *Database) Update(tableName string, rec Record) error {
	if !db.TableExists(tableName) {
		return ErrNoTable
	}

	db.locks[tableName].Lock()
	defer db.locks[tableName].Unlock()

	id := rec.GetID()

	rawRec, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	if err := db.store.UpdateRec(tableName, id, rawRec); err != nil {
		return err
	}

	return nil
}

// unexported methods

func (db *Database) incrementLastID(tableName string) int {
	lastID := db.lastIDs[tableName]

	lastID++

	db.lastIDs[tableName] = lastID

	return lastID
}

func (db *Database) tableExists(tableName string) bool {
	_, ok := db.locks[tableName]
	if !ok {
		return false
	}
	_, ok = db.lastIDs[tableName]
	if !ok {
		return false
	}

	return ok
}
