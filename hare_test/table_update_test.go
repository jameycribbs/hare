package hare_test

import "testing"

func TestUpdate(t *testing.T) {
	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestUpdate:", err)
	}

	id, err := foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("TestUpdate:", err)
	}

	foo := Foo{}

	err = foosTbl.Find(id, &foo)
	if err != nil {
		t.Error("TestUpdate:", err)
	}

	foo.Bar = "test2"

	err = foosTbl.Update(&foo)
	if err != nil {
		t.Error("TestUpdate:", err)
	}

	err = foosTbl.Find(id, &foo)
	if err != nil {
		t.Error("TestUpdate:", err)
	}

	if foo.Bar != "test2" {
		t.Error("TestUpdate: Expected 'test2', got ", foo.Bar)
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestUpdate:", err)
	}
}
