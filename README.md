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
$ go get github.com/jameycribbs/hare
```


### Usage

#### Setting up Hare to use your JSON file(s)

A directory of JSON files is represented by a hare.Database. Each JSON file
needs a struct with it's members cooresponding to the JSON field names.
Additionally, you need to implement 3 simple boilerplate methods on that
struct that allow it to satisfy the hare.Record interface.

A good way to structure this is to put this boilerplate code in a "models"
package in your project.  You can find an example of this boilerplate code in the
examples/crud/models/episodes.go file.

Now you are ready to go!

Let's say you have a "data" directory with a file in it called "contacts.json".

The top-level object in Hare is a `Database`. It represents the directory on
your disk where the JSON files are located.

To open your database, you first need a new instance of a datastore.  In this
example, we are using the `Disk` datastore:

```go
ds, err := disk.New("./data", ".json")
```
Hare also has the `Ram` datastore for in-memory databases.

Now, you will pass the datastore to Hare's New function and it will return
a `Database` instance:
```go
db, err := hare.New(ds)
```


#### Creating a record

To add a record, you can use the `Insert` method:

```go
recID, err := db.Insert("contacts", &models.Contact{FirstName: "John", LastName: "Doe", Phone: "888-888-8888", Age: 21})
```


#### Finding a record

To find a record if you know the record ID, you can use the `Find` method:

```go
var c models.Contact

err = db.Find("contacts", 1, &c)
```


#### Updating a record

To update a record, you can use the `Update` method:

```go
c.Age = 22

err = db.Update("contacts", &c)
```


#### Deleting a record

To delete a record, you can use the `Delete` method:

```go
err = db.Delete("contacts", 3)
```


#### Querying

To query the database, you can write your query expression in pure Go and pass
it to your model's QueryContacts function as a closure.  You would need to create
the QueryContacts function for your model as part of setup.  You can find an
example of what this function should look like in examples/models/episodes.go.

```go
results, err := models.QueryContacts(db, func(c models.Contact) bool {
  return c.firstname == "John" && c.lastname == "Doe"
}, 0)
```


#### Associations

You can create associations (similar to "belongs_to" in Rails, but with less
features).  For example, you could create another table called "relationships" with
the fields "id" and "type" (i.e. "Spouse", "Sister", "Co-worker", etc.).  Next,
you would add a "relationship_id" field to the contacts table and you would also add
an embeded Relationship struct.  Finally, in the Contact models "AfterFind" method,
which is automatically called by Hare everytime the "Find" method is executed, you
would add code to look-up the associated relationship and populate the embedded
Relationship struct.  Take a look at the crud.go file in the "examples" directory
for an example of how this is done.

You can also mimic a "has_many" association, using a similar technique.  Take a
look at the files in the examples directory for how to do that.


#### Database Administration

There are also built-in methods you can run against the database
to create a new table or delete an existing table. Take a look at the
examples/dbadmin/dbadmin.go file for examples of how these can be used.

When Hare updates an existing record, if the changed record's length is
less than the old record's length, Hare will overwrite the old data
and pad the extra space on the line with all "X"s.

If the changed record's length is greater than the old record's length,
Hare will write the changed record at the end of the file and overwrite
the old record with all "X"s.

Similarly, when Hare deletes a record, it simply overwrites the record
with all "X"s.

Eventually, you will want to remove these obsolete records.  For an
example of how to do this, take a look at the examples/dbadmin/compact.go
file.


## Features

* Records for each table are stored in a newline-delimited JSON file.

* Mutexes are used for table locking.  You can have multiple readers
  or one writer for that table at one time, as long as all processes 
  share the same Database connection.

* Querying is done using Go itself.  No need to use a DSL.

* An AfterFind callback is run automatically, everytime a record is
  read, allowing you to do creative things like auto-populate
  associations, etc.
  
* When using the `Disk` datastore, the database is not read into
  memory, but is queried from disk, so no need to worry about a large
  dataset filling up memory.  Of course, if your database is THAT
  big, you should probably be using a real DBMS, instead of Hare!

* Two different back-end datastores to choose from:  `Disk` or `Ram`.
