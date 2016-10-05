package hare

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const dummyRune = 'X'

type table struct {
	filePtr *os.File
	rwLock  *sync.RWMutex
	lastID  int
}

type Record interface {
	SetID(int)
	GetID() int
}

type DummiesTooShortError struct {
	Msg string
}

func (e *DummiesTooShortError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

func (tbl *table) Create(rec Record) (int, error) {
	var err error
	var offset int64
	var whence int

	recID := tbl.incrementLastID()
	rec.SetID(recID)

	rawRec, err := json.Marshal(rec)
	if err != nil {
		return 0, err
	}

	// First check to see if we can fit it onto a line with a dummy record...
	offset, err = tbl.offsetToFitRec(len(rawRec))

	switch err := err.(type) {
	case nil:
		whence = 0
	case *DummiesTooShortError:
		offset = 0
		whence = 2
	default:
		return 0, err
	}

	if err = tbl.writeRec(offset, whence, rawRec); err != nil {
		return 0, err
	}

	return recID, nil
}

func (tbl *table) Destroy(recID int) error {
	var err error
	var recLength int
	var recOffset int64

	if _, recOffset, recLength, err = tbl.readRec(recID); err != nil {
		return err
	}

	if err = tbl.overwriteRec(recOffset, recLength); err != nil {
		return err
	}

	return nil
}

func (tbl *table) Find(recID int, rec Record) error {
	rawRec, _, _, err := tbl.readRec(recID)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawRec, &rec); err != nil {
		return err
	}

	return err
}

func (tbl *table) ForEachID(f func(int) error) error {
	var recMap map[string]interface{}

	r := bufio.NewReader(tbl.filePtr)

	_, err := tbl.filePtr.Seek(0, 0)
	if err != nil {
		return err
	}

	for {
		rawRec, err := r.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// If this is a record that has been deleted or is the result of an update that left extra data on the line, then skip this
		// dummy record.
		if (rawRec[0] == '\n') || (rawRec[0] == dummyRune) {
			continue
		}

		if err := json.Unmarshal(rawRec, &recMap); err != nil {
			return err
		}

		recMapID := int(recMap["id"].(float64))

		if err = f(recMapID); err != nil {
			return err
		}
	}

	return nil
}

func (tbl *table) Update(rec Record) error {
	var offset int64
	var whence int

	recID := rec.GetID()

	_, recOffset, recLength, err := tbl.readRec(recID)
	if err != nil {
		return err
	}

	rawRec, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	diff := recLength - (len(rawRec) + 1)

	if diff > 0 {
		// Changed record is smaller than record in table.

		offset = recOffset
		whence = 0

		extraData := make([]byte, diff)

		for i, _ := range extraData {
			if i == 0 {
				extraData[i] = '\n'
			} else {
				extraData[i] = 'X'
			}
		}

		rawRec = append(rawRec, extraData...)

		err = tbl.writeRec(recOffset, 0, rawRec)
		if err != nil {
			return err
		}

	} else if diff < 0 {
		// Changed record is larger than the record in table.

		// First check to see if we can fit it onto a line with a dummy record...
		offset, err = tbl.offsetToFitRec(len(rawRec))

		switch err := err.(type) {
		case nil:
			whence = 0
		case *DummiesTooShortError:
			offset = 0
			whence = 2
		default:
			return err
		}

		err = tbl.writeRec(offset, whence, rawRec)
		if err != nil {
			return err
		}

		// Turn old rec into a dummy.
		if err = tbl.overwriteRec(recOffset, recLength); err != nil {
			return err
		}
	} else {
		// Changed record is same length as record in table.

		err = tbl.writeRec(recOffset, 0, rawRec)
		if err != nil {
			return err
		}
	}

	return nil
}

//******************************************************************************
// PRIVATE METHODS
//******************************************************************************

func (tbl *table) incrementLastID() int {
	tbl.lastID += 1

	return tbl.lastID
}

func (tbl *table) initLastID() error {
	var err error
	var recMap map[string]interface{}

	r := bufio.NewReader(tbl.filePtr)

	if _, err = tbl.filePtr.Seek(0, 0); err != nil {
		return err
	}

	for {
		rawRec, err := r.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// If this is a record that has been deleted or is the result of an update that left extra data on the line, then skip this
		// dummy record.
		if (rawRec[0] == '\n') || (rawRec[0] == dummyRune) {
			continue
		}

		if err = json.Unmarshal(rawRec, &recMap); err != nil {
			return err
		}

		recMapID := int(recMap["id"].(float64))

		if recMapID > tbl.lastID {
			tbl.lastID = recMapID
		}
	}

	return nil
}

func (tbl *table) offsetToFitRec(recLengthNeeded int) (int64, error) {
	var recOffset int64
	var totalOffset int64
	var recLength int

	r := bufio.NewReader(tbl.filePtr)

	_, err := tbl.filePtr.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	for {
		rawRec, err := r.ReadBytes('\n')

		recOffset = totalOffset
		recLength = len(rawRec)
		totalOffset += int64(recLength)

		// Need to account for newline character.
		recLength -= 1

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}

		if (rawRec[0] == '\n') || (rawRec[0] == dummyRune) {
			if recLength >= recLengthNeeded {
				return recOffset, nil
			}
		}
	}

	return 0, &DummiesTooShortError{"No dummy records of sufficient length found!"}
}

func (tbl *table) overwriteRec(recOffset int64, recLength int) error {
	var err error

	// Overwrite record with XXXXXXXX...
	oldRecData := make([]byte, recLength-1)

	for i, _ := range oldRecData {
		oldRecData[i] = 'X'
	}

	err = tbl.writeRec(recOffset, 0, oldRecData)
	if err != nil {
		return err
	}

	return nil
}

func (tbl *table) readRec(recID int) ([]byte, int64, int, error) {
	var recOffset int64
	var totalOffset int64
	var recLength int
	var recMap map[string]interface{}

	r := bufio.NewReader(tbl.filePtr)

	_, err := tbl.filePtr.Seek(0, 0)
	if err != nil {
		return nil, 0, 0, err
	}

	for {
		rawRec, err := r.ReadBytes('\n')

		recOffset = totalOffset
		recLength = len(rawRec)
		totalOffset += int64(recLength)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, 0, 0, err
		}

		// If this is a record that has been deleted or is the result of an update that left extra data on the line, then skip this
		// dummy record.
		if (rawRec[0] == '\n') || (rawRec[0] == dummyRune) {
			continue
		}

		if err := json.Unmarshal(rawRec, &recMap); err != nil {
			return nil, 0, 0, err
		}

		recMapID := int(recMap["id"].(float64))

		if recMapID == recID {
			return rawRec, recOffset, recLength, err
		}
	}

	return nil, 0, 0, errors.New("Record not found!")
}

func (tbl *table) writeRec(offset int64, whence int, rec []byte) error {
	var rawRec []byte
	var err error

	w := bufio.NewWriter(tbl.filePtr)

	rawRec = append(rec, '\n')

	_, err = tbl.filePtr.Seek(offset, whence)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(rawRec)
	if err != nil {
		panic(err)
	}

	w.Flush()

	return nil
}
