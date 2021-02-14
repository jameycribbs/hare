package ram

import "github.com/jameycribbs/hare/dberr"

// Ram is a struct that holds a map of all the
// tables in the datastore.
type Ram struct {
	tables map[string]*table
}

// New takes a map of maps with seed data
// and returns a pointer to a Ram struct.
func New(seedData map[string]map[int]string) (*Ram, error) {
	var ram Ram

	if err := ram.init(seedData); err != nil {
		return nil, err
	}

	return &ram, nil
}

// Close closes the datastore.
func (ram *Ram) Close() error {
	ram.tables = nil

	return nil
}

// CreateTable takes a table name, creates a new table
// and adds it to the map of tables in the datastore.
func (ram *Ram) CreateTable(tableName string) error {
	if ram.TableExists(tableName) {
		return dberr.ErrTableExists
	}

	ram.tables[tableName] = newTable()

	return nil
}

// DeleteRec takes a table name and a record id and deletes
// the associated record.
func (ram *Ram) DeleteRec(tableName string, id int) error {
	table, err := ram.getTable(tableName)
	if err != nil {
		return err
	}

	if err = table.deleteRec(id); err != nil {
		return err
	}

	return nil
}

// GetLastID takes a table name and returns the greatest record
// id found in the table.
func (ram *Ram) GetLastID(tableName string) (int, error) {
	table, err := ram.getTable(tableName)
	if err != nil {
		return 0, err
	}

	return table.getLastID(), nil
}

// IDs takes a table name and returns an array of all record IDs
// found in the table.
func (ram *Ram) IDs(tableName string) ([]int, error) {
	table, err := ram.getTable(tableName)
	if err != nil {
		return nil, err
	}

	return table.ids(), nil
}

// InsertRec takes a table name, a record id, and a byte array and adds
// the record to the table.
func (ram *Ram) InsertRec(tableName string, id int, rec []byte) error {
	table, err := ram.getTable(tableName)
	if err != nil {
		return err
	}

	if table.recExists(id) {
		return dberr.ErrIDExists
	}

	table.writeRec(id, rec)

	return nil
}

// ReadRec takes a table name and an id, reads the record from the
// table, and returns a populated byte array.
func (ram *Ram) ReadRec(tableName string, id int) ([]byte, error) {
	table, err := ram.getTable(tableName)
	if err != nil {
		return nil, err
	}

	rec, err := table.readRec(id)
	if err != nil {
		return nil, err
	}

	return rec, err
}

// RemoveTable takes a table name and deletes that table from the
// datastore.
func (ram *Ram) RemoveTable(tableName string) error {
	if !ram.TableExists(tableName) {
		return dberr.ErrNoTable
	}

	delete(ram.tables, tableName)

	return nil
}

// TableExists takes a table name and returns a bool indicating
// whether or not the table exists in the datastore.
func (ram *Ram) TableExists(tableName string) bool {
	_, ok := ram.tables[tableName]

	return ok
}

// TableNames returns an array of table names.
func (ram *Ram) TableNames() []string {
	var names []string

	for k := range ram.tables {
		names = append(names, k)
	}

	return names
}

// UpdateRec takes a table name, a record id, and a byte array and updates
// the table record with that id.
func (ram *Ram) UpdateRec(tableName string, id int, rec []byte) error {
	table, err := ram.getTable(tableName)
	if err != nil {
		return err
	}

	if !table.recExists(id) {
		return dberr.ErrNoRecord
	}

	table.writeRec(id, rec)

	return nil
}

//******************************************************************************
// UNEXPORTED METHODS
//******************************************************************************

func (ram *Ram) getTable(tableName string) (*table, error) {
	table, ok := ram.tables[tableName]
	if !ok {
		return nil, dberr.ErrNoTable
	}

	return table, nil
}

func (ram *Ram) getTables() ([]string, error) {
	var tableNames []string

	for name := range ram.tables {
		tableNames = append(tableNames, name)
	}

	return tableNames, nil
}

func (ram *Ram) init(seedData map[string]map[int]string) error {
	ram.tables = make(map[string]*table)

	for tableName, tableData := range seedData {
		ram.tables[tableName] = newTable()

		for id, rec := range tableData {
			if err := ram.InsertRec(tableName, id, []byte(rec)); err != nil {
				return err
			}
		}
	}

	return nil
}
