package hare

import "testing"

//-----------------------------------------------------------------------------
// Tests
//-----------------------------------------------------------------------------
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
		if err.Error() != "Record not found!" {
			t.Error("TestDestroy: Expected Find error to be 'Record not found!', got ", err)
		}
	} else {
		t.Error("TestDestroy: Expected Find error, got no error.")
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestDestroy:", err)
	}
}
