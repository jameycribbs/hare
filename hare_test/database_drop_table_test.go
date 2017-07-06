package hare_test

import "testing"

func TestDropTable(t *testing.T) {
	var err error

	if _, err = db.CreateTable("foos"); err != nil {
		t.Error("TestDropTable:", err)
	}

	if !db.TableExists("foos") {
		t.Error("TestDropTable:", err)
	}

	if err = db.DropTable("foos"); err != nil {
		t.Error("TestDropTable:", err)
	}

	if db.TableExists("foos") {
		t.Error("TestDropTable:", err)
	}
}
