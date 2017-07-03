package hare_test

import "testing"

func TestCreate(t *testing.T) {
	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestCreate:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("TestCreate:", err)
	}

	foo := Foo{}

	err = foosTbl.Find(id, &foo)
	if err != nil {
		t.Error("TestCreate:", err)
	}

	if foo.Bar != "test" {
		t.Error("TestCreate: Expected 'test', got ", foo.Bar)
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestCreate:", err)
	}
}
