package hare_test

import "testing"

func TestDestroyTable(t *testing.T) {
	_, err := db.CreateTable("foos")
	if err != nil {
		t.Error("TestDestroyTable:", err)
	}

	if !db.TableExists("foos") {
		t.Error("TestDestroyTable:", err)
	}

	err = db.DropTable("foos")
	if err != nil {
		t.Error("TestDestroyTable:", err)
	}

	if db.TableExists("foos") {
		t.Error("TestDestroyTable:", err)
	}
}
