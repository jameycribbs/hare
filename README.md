![alt tag](https://raw.githubusercontent.com/jameycribbs/hare/master/hare.jpg)

Hare - A nimble little database management system written in Go
====

Hare is a pure Go database management system that stores each table as
a text file of line-delimited JSON.  It is a good fit for applications
that require a simple embedded DBMS.

## Table of Contents

- [Getting Started](#getting-started)
  - [Installing](#installing)
  - [Usage](#usage)
- [Features](#features)
- [Limitations](#limitations)

## Getting Started

### Installing

To start using Hare, install Go and run `go get`:

```sh
$ go get github.com/jameycribbs/hare/...
```

### Usage

The top-level object in Hare is a `Database`. It is represented as a directory on
your disk.

To open your database, simply use the `hare.OpenDB()` function:

```go
package main

import (
  "log"

  "github.com/jameycribbs/hare"
)

func main() {
  // OpenDB takes a directory path pointing to zero or more json files and returns
  // a database connection.
  db, err := hare.OpenDB("data")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  ...
}
```

#### Creating a table

To create a table (represented as a json file), you can use the
Database.CreateTable() function:

```go
contactsTbl, err := db.CreateTable("contacts")
if err != nil {
  log.Fatal(err)
}
```

#### Using a table

To use a table for database operations, you need to create a
structure representing the table columns, and create two methods 
on that structure:

```go
type Contact struct {
  // Required field
  ID         int    `json:"id"`
  FirstName  string `json:"firstname"`
  LastName   string `json:"lastname"`
  Phone      string `json:"phone"`
  Age        int    `json:"age"`
}

func (contact *Contact) SetID(id int) {
  contact.ID = id
}

func (contact *Contact) GetID() int {
  return contact.ID
}
```

#### Creating a record

To add a record, you can use the Table.Create() function:

```go
recID, err := contactsTbl.Create(&Contact{FirstName: "John", LastName: "Doe", Phone: "888-888-8888", Age: 21})
```


#### Finding a record

To find a record if you know the record ID, you can use the Table.Find() function:

```go
var contact Contact

err = contactsTbl.Find(recID, &contact)
if err != nil {
  log.Fatal(err)
}
```


#### Searching records

To search for a record by any field, you can use the Table.ForEachID() function
by passing it a closure that defines your query:

```go
err = contactsTbl.ForEachID(func(recID int) error {
  var contact Contact

  if err = contactsTbl.Find(recID, &contact); err != nil {
    log.Fatal(err)
  }

  if contact.FirstName == "John" && contact.LastName == "Doe" {
    fmt.Println("Contact record for John Doe:", contact)
    return hare.ForEachIDBreak{}
  }
  return nil
})
if err != nil {
  log.Fatal(err)
}
```


#### Updating a record

To add a record, you can use the Table.Update() function:

```go
var contact Contact

err = contactsTbl.Find(recID, &contact)
if err != nil {
  log.Fatal(err)
}

contact.Age = 22

if err = contactsTbl.Update(&contact); err != nil {
  log.Fatal(err)
}
```


#### Deleting a record

To delete a record, you can use the Table.Destroy() function:

```go
if err = contactsTbl.Destroy(recID); err != nil {
  log.Fatal(err)
}
```


#### Droping a table

To delete a table you can use the Database.DropTable() function:

```go
if err = db.DropTable("contacts"); err != nil {
  log.Fatal(err)
}
```


## Features

* Records for each table are stored in a newline-delimited JSON file.

* Hare uses mutexes for table locking.  You can have multiple readers
  or one writer for that table at one time, as long as all processes 
  share the same Database connection.

