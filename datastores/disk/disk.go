package disk

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ErrNoTable     = errors.New("disk: no table with that name found")
	ErrTableExists = errors.New("disk: table with that name already exists")
)

type Disk struct {
	path       string
	ext        string
	tableFiles map[string]*tableFile
}

func New(path string, ext string) (*Disk, error) {
	var dsk Disk

	dsk.path = path
	dsk.ext = ext

	if err := dsk.init(); err != nil {
		return nil, err
	}

	return &dsk, nil
}

func (dsk *Disk) Close() error {
	for _, tableFile := range dsk.tableFiles {
		if err := tableFile.close(); err != nil {
			return err
		}
	}

	dsk.path = ""
	dsk.ext = ""
	dsk.tableFiles = nil

	return nil
}

func (dsk *Disk) CreateTable(tableName string) error {
	if dsk.TableExists(tableName) {
		return ErrTableExists
	}

	filePtr, err := dsk.openFile(tableName, true)
	if err != nil {
		return err
	}

	tableFile, err := newTableFile(tableName, filePtr)
	if err != nil {
		return err
	}

	dsk.tableFiles[tableName] = tableFile

	return nil
}

func (dsk *Disk) DeleteRec(tableName string, id int) error {
	tableFile, err := dsk.getTableFile(tableName)
	if err != nil {
		return err
	}

	if err = tableFile.deleteRec(id); err != nil {
		return err
	}

	return nil
}

func (dsk *Disk) GetLastID(tableName string) (int, error) {
	tableFile, err := dsk.getTableFile(tableName)
	if err != nil {
		return 0, err
	}

	return tableFile.getLastID(), nil
}

func (dsk *Disk) IDs(tableName string) ([]int, error) {
	tableFile, err := dsk.getTableFile(tableName)
	if err != nil {
		return nil, err
	}

	return tableFile.ids(), nil
}

func (dsk *Disk) InsertRec(tableName string, id int, rec []byte) error {
	tableFile, err := dsk.getTableFile(tableName)
	if err != nil {
		return err
	}

	offset, err := tableFile.offsetForWritingRec(len(rec))
	if err != nil {
		return err
	}

	if err := tableFile.writeRec(offset, 0, rec); err != nil {
		return err
	}

	tableFile.offsets[id] = offset

	return nil
}

func (dsk *Disk) ReadRec(tableName string, id int) ([]byte, error) {
	tableFile, err := dsk.getTableFile(tableName)
	if err != nil {
		return nil, err
	}

	rec, err := tableFile.readRec(id)
	if err != nil {
		return nil, err
	}

	return rec, err
}

func (dsk *Disk) RemoveTable(tableName string) error {
	tableFile, err := dsk.getTableFile(tableName)
	if err != nil {
		return err
	}

	tableFile.close()

	if err := os.Remove(dsk.path + "/" + tableName + dsk.ext); err != nil {
		return err
	}

	delete(dsk.tableFiles, tableName)

	return nil
}

func (dsk *Disk) TableExists(tableName string) bool {
	_, ok := dsk.tableFiles[tableName]

	return ok
}

func (dsk *Disk) TableNames() []string {
	var names []string

	for k := range dsk.tableFiles {
		names = append(names, k)
	}

	return names
}

func (dsk *Disk) UpdateRec(tableName string, id int, rec []byte) error {
	tableFile, err := dsk.getTableFile(tableName)
	if err != nil {
		return err
	}

	if err = tableFile.updateRec(id, rec); err != nil {
		return err
	}

	return nil
}

//******************************************************************************
// UNEXPORTED METHODS
//******************************************************************************

func (dsk *Disk) getTableFile(tableName string) (*tableFile, error) {
	tableFile, ok := dsk.tableFiles[tableName]
	if !ok {
		return nil, ErrNoTable
	}

	return tableFile, nil
}

func (dsk *Disk) getTableNames() ([]string, error) {
	var tableNames []string

	files, err := ioutil.ReadDir(dsk.path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()

		// If entry is sub dir, current dir, or parent dir, skip it.
		if file.IsDir() || fileName == "." || fileName == ".." {
			continue
		}

		if !strings.HasSuffix(fileName, dsk.ext) {
			continue
		}

		tableNames = append(tableNames, strings.TrimSuffix(fileName, dsk.ext))
	}

	return tableNames, nil
}

func (dsk *Disk) init() error {
	dsk.tableFiles = make(map[string]*tableFile)

	tableNames, err := dsk.getTableNames()
	if err != nil {
		return err
	}

	for _, tableName := range tableNames {
		filePtr, err := dsk.openFile(tableName, false)
		if err != nil {
			return err
		}

		tableFile, err := newTableFile(tableName, filePtr)
		if err != nil {
			return err
		}

		dsk.tableFiles[tableName] = tableFile
	}

	return nil
}

func (dsk Disk) openFile(tableName string, createIfNeeded bool) (*os.File, error) {
	var osFlag int

	if createIfNeeded {
		osFlag = os.O_CREATE | os.O_RDWR
	} else {
		osFlag = os.O_RDWR
	}

	filePtr, err := os.OpenFile(dsk.path+"/"+tableName+dsk.ext, osFlag, 0660)
	if err != nil {
		return nil, err
	}

	return filePtr, nil
}

func (dsk *Disk) closeTable(tableName string) error {
	tableFile, ok := dsk.tableFiles[tableName]
	if !ok {
		return errors.New("table does not exist")
	}

	if err := tableFile.close(); err != nil {
		return err
	}

	return nil
}
