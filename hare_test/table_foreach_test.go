package hare_test

import "testing"

func TestForEach(t *testing.T) {
	foosTbl, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestCreate:", err)
	}

	_, err = foosTbl.Create(&Foo{Bar: "test"})
	if err != nil {
		t.Error("TestCreate:", err)
	}

	var foo1 Foo

	err = foosTbl.ForEach(func(recMap map[string]interface{}) error {
		foo2 := fooFromRecMap(recMap)

		if foo2.Bar == "test" {
			foo1 = foo2
		}
		return nil
	})

	if foo1.Bar != "test" {
		t.Error("TestCreate: Expected 'test', got ", foo1.Bar)
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestCreate:", err)
	}
}
