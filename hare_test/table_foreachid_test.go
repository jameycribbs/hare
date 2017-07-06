package hare_test

import (
	"testing"

	"github.com/jameycribbs/hare"
)

func TestForEach(t *testing.T) {
	var err error

	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestCreate:", err)
	}

	if _, err = foosTbl.Create(&Foo{Bar: "test"}); err != nil {
		t.Error("TestCreate:", err)
	}

	var foo1 Foo

	err = foosTbl.ForEachID(func(recID int) error {
		var foo2 Foo

		if err = foosTbl.Find(recID, &foo2); err != nil {
			panic(err)
		}

		if foo2.Bar == "test" {
			foo1 = foo2
			return hare.ForEachIDBreak{}
		}
		return nil
	})

	if foo1.Bar != "test" {
		t.Error("TestCreate: Expected 'test', got ", foo1.Bar)
	}

	if err = db.DropTable("foos"); err != nil {
		t.Error("TestCreate:", err)
	}
}
