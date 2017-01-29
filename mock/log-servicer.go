package mock

import (
	"io/ioutil"
	"log"

	"github.com/bunsenapp/migrator"
)

// MockLogServicer generates a log.Logger instance that does not output anywhere.
func MockLogServicer() migrator.LogServicer {
	logger := log.New(ioutil.Discard, "", 0)
	return logger
}
