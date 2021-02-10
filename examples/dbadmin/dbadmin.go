package main

import (
	"fmt"

	"github.com/jameycribbs/hare"
	"github.com/jameycribbs/hare/datastores/disk"
)

func main() {
	ds, err := disk.New("./data", ".json")
	if err != nil {
		panic(err)
	}

	db, err := hare.New(ds)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Here is how to create a new table in the database.
	err = db.CreateTable("contacts")
	if err != nil {
		panic(err)
	}

	// Here is how to check if a table exists in your database.
	if !db.TableExists("contacts") {
		fmt.Println("Table 'contacts' does not exist!")
	}

	// Here is how to drop a table.
	err = db.DropTable("contacts")
	if err != nil {
		panic(err)
	}
}
