<img src="https://raw.githubusercontent.com/jameycribbs/hare/master/hare.jpg" width="400" />

Hare - A nimble little database management system written in Go
====

Hare is a pure Go database management system that stores each table as
a text file of line-delimited JSON.  Each line of JSON represents a 
record.  It is a good fit for applications that require a simple embedded DBMS.

## Table of Contents

- [Getting Started](#getting-started)
  - [Installing](#installing)
  - [Usage](#usage)
- [Features](#features)

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
// OpenDB takes a directory path containing zero or more json files and returns
// a database connection.
db, err := hare.OpenDB("data")
if err != nil {
  ...
}
defer db.Close()

  ...
```

#### Creating a table

To create a table (represented as a json file), you can use the
Database.CreateTable() method.  This will return a handle to the
newly created table:

```go
contactsTbl, err := db.CreateTable("contacts")
```

#### Using a table

To use a table for database operations, you need to create a
structure representing the table columns, and create three
methods on that structure:

```go
type contact struct {
  // id is a required field
  id         int    `json:"id"`
  firstName  string `json:"firstname"`
  lastName   string `json:"lastname"`
  phone      string `json:"phone"`
  age        int    `json:"age"`
}

func (c *contact) SetID(id int) {
  c.id = id
}

func (c *contact) GetID() int {
  return c.id
}

func (c *contact) AfterFind() {
  *c = contact(*c)
}
```

#### Creating a record

To add a record, you can use the Table.Create() method:

```go
recID, err := contactsTbl.Create(&contact{firstName: "John", lastName: "Doe", phone: "888-888-8888", age: 21})
```


#### Finding a record

To find a record if you know the record ID, you can use the Table.Find() method:

```go
var contact contact

err = contactsTbl.Find(recID, &contact)
```


#### Querying a table

To query a table, you need to create a struct with the
table handle embedded and write one method for it that
will be the query method:

```go
type model struct {
	*hare.Table
}

func (mdl *model) query(queryFn func(rec record) bool, limit int) ([]record, error) {
	var results []record
	var err error

	for _, id := range mdl.Table.IDs() {
		r := record{}

		err = mdl.Table.Find(id, &r)
		if err != nil {
			panic(err)
		}

		if queryFn(r) {
			results = append(results, r)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}
```

Then you just need to create an instance of the model struct,
and set the embedded hare.Table to the table handle you have.
Now you are ready to start querying:

```go
mdl := model{Table: contactsTbl}

results, err := mdl.query(func(r record) bool {
  return r.firstname == "Bob" && r.lastname == "Jones"
}, 0)
```

Notice how the actual query logic is an anonymous function?
This allows you to use the full power of Go in your query
expression.


#### Updating a record

To add a record, you can use the Table.Update() method:

```go
contact.age = 22

err = contactsTbl.Update(&contact)
```


#### Deleting a record

To delete a record, you can use the Table.Destroy() method:

```go
err = contactsTbl.Destroy(3)
```


#### Droping a table

To delete a table you can use the Database.DropTable() method:

```go
err = db.DropTable("contacts")
```


## Features

* Records for each table are stored in a newline-delimited JSON file.

* Mutexes are used for table locking.  You can have multiple readers
  or one writer for that table at one time, as long as all processes 
  share the same Database connection.

* Querying is done using Go itself.  No need to use a DSL.

* The database is not read into memory, but is queried from disk, so
  no need to worry about a large dataset filling up memory.
