package disk

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
)

const dummyRune = 'X'

var (
	ErrNoRecord = errors.New("hare: no record with that id found")
)

type tableFile struct {
	ptr     *os.File
	offsets map[int]int64
}

func NewTableFile(tableName string, filePtr *os.File) (*tableFile, error) {
	var currentOffset int64
	var totalOffset int64
	var recLen int
	var recMap map[string]interface{}

	tableFile := tableFile{ptr: filePtr}
	tableFile.offsets = make(map[int]int64)

	r := bufio.NewReader(filePtr)

	_, err := filePtr.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	for {
		rec, err := r.ReadBytes('\n')

		currentOffset = totalOffset
		recLen = len(rec)
		totalOffset += int64(recLen)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// Skip dummy records.
		if (rec[0] == '\n') || (rec[0] == dummyRune) {
			continue
		}

		//Unmarshal so we can grab the record ID.
		if err := json.Unmarshal(rec, &recMap); err != nil {
			return nil, err
		}
		recMapID := int(recMap["id"].(float64))

		tableFile.offsets[recMapID] = currentOffset
	}

	return &tableFile, nil
}

func (t *tableFile) close() error {
	if err := t.ptr.Close(); err != nil {
		return err
	}

	t.offsets = nil

	return nil
}

func (t *tableFile) deleteRec(id int) error {
	offset, ok := t.offsets[id]
	if !ok {
		return ErrNoRecord
	}

	rec, err := t.readRec(id)
	if err != nil {
		return err
	}

	if err = t.overwriteRec(offset, len(rec)); err != nil {
		return err
	}

	delete(t.offsets, id)

	return nil
}

func (t *tableFile) getLastID() int {
	var lastID int

	for k := range t.offsets {
		if k > lastID {
			lastID = k
		}
	}

	return lastID
}

func (t *tableFile) ids() []int {
	ids := make([]int, len(t.offsets))

	i := 0
	for id := range t.offsets {
		ids[i] = id
		i++
	}

	return ids
}

// offsetForWritingRec takes a record length and returns the offset in the file
// where the record is to be written.  It will try to fit the record on a dummy
// line, otherwise, it will return the offset at the end of the file.
func (t *tableFile) offsetForWritingRec(recLen int) (int64, error) {
	var offset int64
	var err error

	// Can the record fit onto a line with a dummy record?
	offset, recFitErr := t.offsetToFitRec(recLen)

	switch recFitErr.(type) {
	case nil:
	case dummiesTooShortError:
		// Go to the end of the file.
		offset, err = t.ptr.Seek(0, 2)
		if err != nil {
			return 0, err
		}
	default:
		return 0, recFitErr
	}

	return offset, nil
}

// offsetToFitRec takes a record length and checks all the dummy records to see
// if any are big enough to fit the record.
func (t *tableFile) offsetToFitRec(recLenNeeded int) (int64, error) {
	var recLen int
	var offset int64
	var totalOffset int64

	r := bufio.NewReader(t.ptr)

	if _, err := t.ptr.Seek(0, 0); err != nil {
		return 0, err
	}

	for {
		rec, err := r.ReadBytes('\n')

		offset = totalOffset
		recLen = len(rec)
		totalOffset += int64(recLen)

		// Need to account for newline character.
		recLen--

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}

		// If this is a dummy record, is it big enough to
		// hold the needed record length?
		if (rec[0] == '\n') || (rec[0] == dummyRune) {
			if recLen >= recLenNeeded {
				return offset, nil
			}
		}
	}

	return 0, dummiesTooShortError{}
}

func (t *tableFile) overwriteRec(offset int64, recLen int) error {
	// Overwrite record with XXXXXXXX...
	dummyData := make([]byte, recLen-1)

	for i := range dummyData {
		dummyData[i] = 'X'
	}

	if err := t.writeRec(offset, 0, dummyData); err != nil {
		return err
	}

	return nil
}

func (t *tableFile) readRec(id int) ([]byte, error) {
	offset, ok := t.offsets[id]
	if !ok {
		return nil, ErrNoRecord
	}

	r := bufio.NewReader(t.ptr)

	if _, err := t.ptr.Seek(offset, 0); err != nil {
		return nil, err
	}

	rec, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return rec, err
}

func (t *tableFile) updateRec(id int, rec []byte) error {
	recLen := len(rec)

	oldRecOffset, ok := t.offsets[id]
	if !ok {
		return ErrNoRecord
	}

	oldRec, err := t.readRec(id)
	if err != nil {
		return err
	}

	oldRecLen := len(oldRec)

	diff := oldRecLen - (recLen + 1)

	if diff > 0 {
		// Changed record is smaller than record in table, so dummy out
		// extra space and write over old record.

		rec = append(rec, padRec(diff)...)

		if err = t.writeRec(oldRecOffset, 0, rec); err != nil {
			return err
		}

	} else if diff < 0 {
		// Changed record is larger than the record in table.

		recOffset, err := t.offsetForWritingRec(recLen)
		if err != nil {
			return err
		}

		if err = t.writeRec(recOffset, 0, rec); err != nil {
			return err
		}

		// Turn the old record into a dummy.
		if err = t.overwriteRec(oldRecOffset, oldRecLen); err != nil {
			return err
		}

		// Update the index with the new offset since the record is in a
		// new position in the file.
		t.offsets[id] = recOffset
	} else {
		// Changed record is the same length as the record in the table.
		err = t.writeRec(oldRecOffset, 0, rec)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *tableFile) writeRec(offset int64, whence int, rec []byte) error {
	var err error

	w := bufio.NewWriter(t.ptr)

	if _, err = t.ptr.Seek(offset, whence); err != nil {
		return err
	}

	if _, err = w.Write(append(rec, '\n')); err != nil {
		return err
	}

	w.Flush()

	return nil
}

func padRec(padLength int) []byte {
	extraData := make([]byte, padLength)

	extraData[0] = '\n'

	for i := 1; i < padLength; i++ {
		extraData[i] = dummyRune
	}

	return extraData
}

// dummiesTooShortError is a place to hold a custom error used
// as part of a switch.
type dummiesTooShortError struct {
}

func (e dummiesTooShortError) Error() string {
	return "all dummy records are too short"
}
