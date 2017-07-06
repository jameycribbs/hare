package hare_test

import (
	"strconv"
	"testing"
)

func TestDestroy(t *testing.T) {
	var err error

	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestDestroy:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("TestDestroy:", err)
	}

	if err = foosTbl.Destroy(id); err != nil {
		t.Error("TestDestroy:", err)
	}

	foo := Foo{}

	if err = foosTbl.Find(id, &foo); err != nil {
		if err.Error() != "Find Error: Record with ID of "+strconv.Itoa(id)+" does not exist!" {
			t.Error("TestDestroy: Expected Find error, got ", err)
		}
	} else {
		t.Error("TestDestroy: Expected Find error, got no error.")
	}

	if err = db.DropTable("foos"); err != nil {
		t.Error("TestDestroy:", err)
	}
}
