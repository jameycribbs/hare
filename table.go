package hare

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"sync"
)

const dummyRune = 'X'

type table struct {
	filePtr *os.File
	sync.RWMutex
	lastID int
	index  map[int]int64
}

type Record interface {
	SetID(int)
	GetID() int
}

type DummiesTooShortError struct {
}

func (e DummiesTooShortError) Error() string {
	return ""
}

type ForEachBreak struct {
}

func (e ForEachBreak) Error() string {
	return ""
}

func (tbl *table) Create(rec Record) (int, error) {
	tbl.Lock()
	defer tbl.Unlock()

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
	case DummiesTooShortError:
		whence = 2
	default:
		return 0, err
	}

	if whence == 2 {
		offset, err = tbl.filePtr.Seek(0, 2)

		if err != nil {
			return 0, err
		}
	}

	if err = tbl.writeRec(offset, 0, rawRec); err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	tbl.index[recID] = offset

	return recID, nil
}

func (tbl *table) Destroy(recID int) error {
	var err error

	tbl.Lock()
	defer tbl.Unlock()

	rawRec, err := tbl.readRec(tbl.index[recID])
	if err != nil {
		return err
	}

	if err = tbl.overwriteRec(tbl.index[recID], len(rawRec)); err != nil {
		return err
	}

	delete(tbl.index, recID)

	return nil
}

func (tbl *table) Find(recID int, rec Record) error {
	tbl.RLock()
	defer tbl.RUnlock()

	rawRec, err := tbl.readRec(tbl.index[recID])
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawRec, &rec); err != nil {
		return err
	}

	if recID != rec.GetID() {
		return errors.New("Find Error: Record with ID of " + strconv.Itoa(recID) + " does not exist!")
	}

	return err
}

func (tbl *table) ForEach(f func(map[string]interface{}) error) error {
	var recMap map[string]interface{}

	for recID := range tbl.index {
		rawRec, err := tbl.readRec(tbl.index[recID])
		if err != nil {
			return err
		}

		if err = json.Unmarshal(rawRec, &recMap); err != nil {
			return err
		}

		err = f(recMap)

		switch err := err.(type) {
		case nil:
			continue
		case ForEachBreak:
			return nil
		default:
			return err
		}
	}

	return nil
}

func (tbl *table) Update(rec Record) error {
	tbl.Lock()
	defer tbl.Unlock()

	var offset int64
	var goToEoF bool

	recID := rec.GetID()

	oldRecOffset := tbl.index[recID]

	oldRec, err := tbl.readRec(oldRecOffset)
	if err != nil {
		return err
	}

	oldRecLen := len(oldRec)

	newRec, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	newRecLen := len(newRec)

	diff := oldRecLen - (newRecLen + 1)

	if diff > 0 {
		// Changed record is smaller than record in table.

		extraData := make([]byte, diff)

		for i, _ := range extraData {
			if i == 0 {
				extraData[i] = '\n'
			} else {
				extraData[i] = 'X'
			}
		}

		newRec = append(newRec, extraData...)

		err = tbl.writeRec(oldRecOffset, 0, newRec)
		if err != nil {
			return err
		}

	} else if diff < 0 {
		// Changed record is larger than the record in table.

		// First check to see if we can fit it onto a line with a dummy record...
		offset, err = tbl.offsetToFitRec(len(newRec))

		switch err := err.(type) {
		case nil:
		case DummiesTooShortError:
			goToEoF = true
		default:
			return err
		}

		// If we can't fit the updated record onto a line with a dummy record, then go to the End of File.
		if goToEoF {
			offset, err = tbl.filePtr.Seek(0, 2)
			if err != nil {
				return err
			}
		}

		err = tbl.writeRec(offset, 0, newRec)
		if err != nil {
			return err
		}

		// Turn old rec into a dummy.
		if err = tbl.overwriteRec(tbl.index[recID], oldRecLen); err != nil {
			return err
		}

		// Update index with new offset since record is in new position in the file.
		tbl.index[recID] = offset
	} else {
		// Changed record is same length as record in table.

		err = tbl.writeRec(tbl.index[recID], 0, newRec)
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

func (tbl *table) initIndex() error {
	var recOffset int64
	var totalOffset int64
	var recLength int
	var recMap map[string]interface{}

	tbl.index = make(map[int]int64)

	r := bufio.NewReader(tbl.filePtr)

	_, err := tbl.filePtr.Seek(0, 0)
	if err != nil {
		return err
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

		tbl.index[recMapID] = recOffset
	}

	return nil
}

func (tbl *table) initLastID() {
	tbl.lastID = 0

	for k := range tbl.index {
		if k > tbl.lastID {
			tbl.lastID = k
		}
	}
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

	return 0, DummiesTooShortError{}
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

func (tbl *table) readRec(offset int64) ([]byte, error) {
	var recMap map[string]interface{}

	r := bufio.NewReader(tbl.filePtr)

	_, err := tbl.filePtr.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	rawRec, err := r.ReadBytes('\n')

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(rawRec, &recMap); err != nil {
		return nil, err
	}

	return rawRec, err
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
