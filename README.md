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

#### Using a table

First, you need to create a struct representing the
table's schema and then create 3 methods to satisy
the hare.Record interface:

```go
type contact struct {
  // ID is a required field
  ID         int    `json:"id"`
  FirstName  string `json:"firstname"`
  LastName   string `json:"lastname"`
  Phone      string `json:"phone"`
  Age        int    `json:"age"`
}

func (c *contact) SetID(id int) {
  c.ID = id
}

func (c *contact) GetID() int {
  return c.ID
}

func (c *contact) AfterFind() {
  *c = contact(*c)
}
```

#### Creating a record

To add a record, you can use the Table.Create() method:

```go
recID, err := contactsTbl.Create(&contact{FirstName: "John", LastName: "Doe", Phone: "888-888-8888", Age: 21})
```


#### Finding a record

To find a record if you know the record ID, you can use the Table.Find() method:

```go
var c contact

err = contactsTbl.Find(1, &c)
```

#### Updating a record

To update a record, you can use the Table.Update() method:

```go
c.Age = 22

err = contactsTbl.Update(&c)
```


#### Deleting a record

To delete a record, you can use the Table.Destroy() method:

```go
err = contactsTbl.Destroy(3)
```


#### Querying a table

To query a table, you need to create a struct with the
table handle embedded and write one method for it that
will be the query method:

```go
type contactsModel struct {
	*hare.Table
}

func (mdl *contactsModel) query(queryFn func(rec contact) bool, limit int) ([]contact, error) {
	var results []contact
	var err error

	for _, id := range mdl.Table.IDs() {
		rec := contact{}

		if err = mdl.Table.Find(id, &rec); err != nil {
			return nil, err
		}

		if queryFn(rec) {
			results = append(results, rec)
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
contactsMdl := contactsModel{Table: contactsTbl}

results, err := contactsMdl.query(func(c contact) bool {
  return c.firstname == "John" && c.lastname == "Doe"
}, 0)
```

Notice how the actual query logic is an anonymous function?
This allows you to use the full power of Go in your query
expression.



#### Creating a table

To create a new table (represented as a json file), you can use the
Database.CreateTable() method.  This will return a handle to the
newly created table:

```go
budgetTbl, err := db.CreateTable("budget")
```

#### Droping a table

To delete a table you can use the Database.DropTable() method:

```go
err = db.DropTable("budget")
```


## Features

* Records for each table are stored in a newline-delimited JSON file.

* Mutexes are used for table locking.  You can have multiple readers
  or one writer for that table at one time, as long as all processes 
  share the same Database connection.

* Querying is done using Go itself.  No need to use a DSL.

* The database is not read into memory, but is queried from disk, so
  no need to worry about a large dataset filling up memory.
