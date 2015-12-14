# Sample CI service Written in Go

## Introduction

This is a sample CI Service implementation written in Go.  In this project, we use a simple ID mapping service as an example.


## Basic Design

### Service Package

### RESTful API server

### Registration




----

## Error Handling



### Invalid user data

If user POST data which is not valid as the input for the given service, it should return this.

* Code:
* Body:
  ```json
  {
    "code": 400
    "message": "(Human readable description of the error)"
    "error": "(Any error massage from the system.)"
  }
  ```
* Example

### Resource Not found

* Code: 404
* Body:
  ```json
  {
    "code": 404
    "message": "(Human readable description of the error)"
    "error": "(Any error massage from the system.)"
  }
  ```



### Unexpected server side error 