Keywords: Golang, go, database, DBMS, JSON

![alt tag](https://https://github.com/jameycribbs/hare/blob/master/hare.jpg)

### Hare - A nimble little database management system for Go

Hare is a database management system that stores all records for a table as newline delimited __JSON__ strings. It can be __embedded__ into your program and is safe to use in goroutines (it uses mutexes) as long as each goroutine shares the database connection.  This makes it ok to use in a web application.

### Features

- Pure Go
- Goroutine safe (as long as each goroutine shares the database connection)
- Can utilize indexes for faster queries
- Embeddable
- Database records are stored as json strings in one file per table, making for easy external access

### How to install

~~~
go get github.com/jameycribbs/hare
~~~

### How to use

Check out examples.go in the examples directory.

### Contributions welcome!

Pull requests/forks/bug reports all welcome, and please share your thoughts, questions and feature requests in the [Issues] section or via [Email].

[Email]: mailto:jamey.cribbs@gmail.com
[Issues]: https://github.com/jameycribbs/hare/issues

