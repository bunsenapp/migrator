package mock

import (
	"io/ioutil"
	"log"

	"github.com/bunsenapp/migrator"
)

// MockLogServicer generates a log.Logger instance that does not output anywhere.
func MockLogServicer() migrator.LogServicer {
	return log.New(ioutil.Discard, "", 0)
}
