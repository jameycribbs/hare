package hare_test

import "testing"

func TestCreateTable(t *testing.T) {
	_, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestCreateTable:", err)
	}

	if !db.TableExists("foos") {
		t.Error("TestCreateTable:", err)
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestCreateTable:", err)
	}
}
