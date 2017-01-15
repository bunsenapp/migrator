package main

import "errors"

var (
	// InvalidParameterErr is an error that is raised when a required parameter
	// is not passed in as a command line argument.
	InvalidParameterErr = errors.New("One of the required parameters was invalid. See: migrator -h")
)
