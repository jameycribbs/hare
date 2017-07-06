package hare_test

import "testing"

func TestUpdate(t *testing.T) {
	var err error

	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestUpdate:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("TestUpdate:", err)
	}

	foo := Foo{}

	if err = foosTbl.Find(id, &foo); err != nil {
		t.Error("TestUpdate:", err)
	}

	foo.Bar = "test2"

	if err = foosTbl.Update(&foo); err != nil {
		t.Error("TestUpdate:", err)
	}

	if err = foosTbl.Find(id, &foo); err != nil {
		t.Error("TestUpdate:", err)
	}

	if foo.Bar != "test2" {
		t.Error("TestUpdate: Expected 'test2', got ", foo.Bar)
	}

	if err = db.DropTable("foos"); err != nil {
		t.Error("TestUpdate:", err)
	}
}
