package hare_test

import "testing"

func TestCreateTable(t *testing.T) {
	var err error

	if _, err := db.CreateTable("foos"); err != nil {
		t.Error("TestCreateTable:", err)
	}

	if !db.TableExists("foos") {
		t.Error("TestCreateTable:", err)
	}

	if err = db.DropTable("foos"); err != nil {
		t.Error("TestCreateTable:", err)
	}
}
