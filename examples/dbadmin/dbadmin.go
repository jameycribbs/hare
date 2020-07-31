package main

import (
	"fmt"

	"github.com/jameycribbs/hare"
)

func main() {
	// Open the database and return a handle to it.
	db, err := hare.OpenDB("../example_data")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Here is how to create a new table in the database and get back a
	// handle to it.
	tbl, err := db.CreateTable("contacts")
	if err != nil {
		panic(err)
	}

	// If the table already exists, you can get a handle to it by
	// calling GetTable.
	tbl, err = db.GetTable("contacts")
	if err != nil {
		panic(err)
	}

	fmt.Println("Table handle:", tbl)

	// Here is how to drop a table.
	err = db.DropTable("contacts")
	if err != nil {
		panic(err)
	}
}
