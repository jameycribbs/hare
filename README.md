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
- [Example Web App](#example-web-app)

## Getting Started

### Installing

To start using Hare, install Go and run `go get`:

```sh
$ go get github.com/jameycribbs/hare
```


### Usage

#### Setting up Hare to use your JSON file(s)

Each JSON file is represented by a hare.Table.  To set things up, you need to
create a struct with an embedded pointer to a hare.Table and add a Query method
to it.

Additionally, you need to create a struct for a table's record, with
it's members cooresponding to the JSON field names, and implement 3 simple
boilerplate methods on that struct that allow it to satisfy the hare.Record
interface.

A good way to structure this is to put this boilerplate code in a "models"
package in your project.  You can find an example of this boilerplate code in the
examples/crud/models/episodes.go file.

Now you are ready to go!

Let's say you have a "data" directory with a file in it called "contacts.json".

The top-level object in Hare is a `Database`. It represents the directory on
your disk where the JSON files are located.

To open your database, you first need to new instance of a datastore.  In this
example, we are using a "disk" store:

```go
ds, err := disk.New("./data", ".json")
```
Hare also has a "ram" store for in-memory databases.

Now, you will pass the datastore to Hare's New function and it will return
a Database instance:
```go
db, err := hare.New(ds)
...
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


#### Querying a table

To query the database, you can write your query expression in pure Go and pass
it to your model's QueryContacts method as a closure.

```go
results, err := models.QueryContacts(db, func(c models.Contact) bool {
  return c.firstname == "John" && c.lastname == "Doe"
}, 0)
```


#### Associations

You can create associations (similar to "belongs_to" in Rails, but with less
features).  For example, you would create another table called "relationships" with
the fields "id" and "type" (i.e. "Spouse", "Sister", "Co-worker", etc.).  Next,
you would add a "relationship_id" field to the contacts table and you would also add
an embeded Relationship struct.  Finally, in the Contact models "AfterFind" method,
which is automatically called by Hare everytime the "Find" method is executed, you
would add code to look-up the associated relationship and populate an embedded
Relationship struct.  Take a look at the crud.go file in the "examples" directory
for an example of how this is done.

You can also mimic a "has_many" association, using a similar technique.  Take a
look at the crud example for how to do that.


#### Database Administration

There are also built-in methods you can run against the database
to create a new table or delete an existing table.


## Features

* Records for each table are stored in a newline-delimited JSON file.

* Mutexes are used for table locking.  You can have multiple readers
  or one writer for that table at one time, as long as all processes 
  share the same Database connection.

* Querying is done using Go itself.  No need to use a DSL.

* Automatic callback is run everytime a record is read, allowing you
  to do creative things like auto-populate associations, etc.
  
* The database is not read into memory, but is queried from disk, so
  no need to worry about a large dataset filling up memory.

* Two different back-end datastores to choose from:  disk or ram.

## Example Web Apps

[SnippetBox using Hare](https://www.github.com/jameycribbs/snippetbox_hare)
This is a version of the SnippetBox web application featured in Alex
Edward's outstanding book, [Let's Go](https://lets-go.alexedwards.net/),
with Hare replacing MySQL as the DBMS.  This is just a demonstration,
mainly to show how you could use Hare in a web application.

[Mule](https://www.github.com/jameycribbs/mule)
This is a simple web app that tracks expenses and possibly more.  It is
a WORK IN PROGRESS!  I plan on adding other features like a simple
calendaring function.  I have attempted to follow best practices from
the Let's Go book, but I am also structuring it to be more expandable
by separating the handlers, routes, templates, models into separate
packages rather than having everything be in the same package.
