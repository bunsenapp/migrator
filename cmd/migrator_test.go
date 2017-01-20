package cmd

import (
	"fmt"
	"testing"

	"github.com/bunsenapp/migrator"
)

func TestConfigurationValidationFailureIsReturned(t *testing.T) {
	config := migrator.Configuration{}

	m := NewMigrator(config)
	if err := m.Run(); err != migrator.ErrConfigurationInvalid {
		fmt.Errorf("Error returned was not correct.")
	}
}
