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

Each json file is represented by a hare.Table.  To set things up, you need to
create a struct with an embedded pointer to a hare.Table and add a Query method
to it.  Additionally, you need to create a struct for a table's record and
implement 3 simple boilerplate methods that allow it to satisfy the hare.Record
interface. You can find an example of the needed structs and methods in the
examples/crud/models/episodes.go file.

Once that needed structs and methods are written, you just need to create an
instance of the model struct and set it's embedded hare.Table struct pointer
to a handle you get from Hare by calling the hare.GetTable database method.
```go
var contacts models.Contacts

contacts.Table, err = db.GetTable("contacts")
if err != nil {
	panic(err)
}
```

Now you are ready to go!

#### Creating a record

To add a record, you can use the Create() method:

```go
recID, err := contacts.Create(&contact{FirstName: "John", LastName: "Doe", Phone: "888-888-8888", Age: 21})
```


#### Finding a record

To find a record if you know the record ID, you can use the Find() method:

```go
var c contact

err = contacts.Find(1, &c)
```

#### Updating a record

To update a record, you can use the Update() method:

```go
c.Age = 22

err = contacts.Update(&c)
```


#### Deleting a record

To delete a record, you can use the Destroy() method:

```go
err = contacts.Destroy(3)
```


#### Querying a table

To query the database, you can write your query in pure Go and pass it to your
model's Query method as a closure.

```go
results, err := contacts.Query(func(c models.Contact) bool {
  return c.firstname == "John" && c.lastname == "Doe"
}, 0)
```



#### Database Administration

There are also built-in methods you can run against the database
to create a new table or delete an existing table.


## Features

* Records for each table are stored in a newline-delimited JSON file.

* Mutexes are used for table locking.  You can have multiple readers
  or one writer for that table at one time, as long as all processes 
  share the same Database connection.

* Querying is done using Go itself.  No need to use a DSL.

* The database is not read into memory, but is queried from disk, so
  no need to worry about a large dataset filling up memory.
