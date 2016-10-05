Keywords: Golang, go, database, DBMS, JSON

### Ivy - A simple, file-based Database Management System (DBMS) for Go

Ivy is a database management system that stores each record as a __JSON__ file. It can be __embedded__ into your program and is safe to use in goroutines (it uses mutexes) as long as each goroutine shares the database connection.  This makes it ok to use in a web application.

### Features

- Pure Go
- Goroutine safe (as long as each goroutine shares the database connection)
- Can utilize indexes for faster queries
- Embeddable
- Database records are stored as json files, making for easy external access

### How to install

~~~
go get github.com/jameycribbs/ivy
~~~

### How to use

Check out example.go in the examples directory.  For a more comprehensive example of how to use Ivy in a web application, check out [Pythia].

### Contributions welcome!

Pull requests/forks/bug reports all welcome, and please share your thoughts, questions and feature requests in the [Issues] section or via [Email].

[Email]: mailto:jamey.cribbs@gmail.com
[Issues]: https://github.com/jameycribbs/ivy/issues
[Pythia]: https://github.com/jameycribbs/pythia

