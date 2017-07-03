package hare_test

import (
	"strconv"
	"testing"
)

func TestDestroy(t *testing.T) {
	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestDestroy:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("TestDestroy:", err)
	}

	err = foosTbl.Destroy(id)
	if err != nil {
		t.Error("TestDestroy:", err)
	}

	foo := Foo{}

	err = foosTbl.Find(id, &foo)
	if err != nil {
		if err.Error() != "Find Error: Record with ID of "+strconv.Itoa(id)+" does not exist!" {
			t.Error("TestDestroy: Expected Find error, got ", err)
		}
	} else {
		t.Error("TestDestroy: Expected Find error, got no error.")
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestDestroy:", err)
	}
}
