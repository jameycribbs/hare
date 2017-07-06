package hare_test

import "testing"

func TestCreate(t *testing.T) {
	var err error

	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestCreate:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("TestCreate:", err)
	}

	foo := Foo{}

	if err = foosTbl.Find(id, &foo); err != nil {
		t.Error("TestCreate:", err)
	}

	if foo.Bar != "test" {
		t.Error("TestCreate: Expected 'test', got ", foo.Bar)
	}

	if err = db.DropTable("foos"); err != nil {
		t.Error("TestCreate:", err)
	}
}
