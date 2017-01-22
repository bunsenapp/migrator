package cmd

import (
	"testing"

	"github.com/bunsenapp/migrator"
	"github.com/bunsenapp/migrator/mock"
)

func TestConfigurationValidationFailureIsReturned(t *testing.T) {
	config := migrator.Configuration{}

	m, err := NewMigrator(config, nil)
	if err != nil {
		t.Errorf("Error occurred whilst creating migrator.")
	}

	if err = m.Run(); err == nil || err != migrator.ErrConfigurationInvalid {
		t.Errorf("error returned was not correct.")
	}
}

func TestNotInitialisedDatabaseServicerResultsInError(t *testing.T) {
	config := mock.ValidConfiguration()

	m := Migrator{
		Configuration: config,
	}
	if err := m.Run(); err == nil || err != migrator.ErrDbServicerNotInitialised {
		t.Errorf("error returned was not correct.")
	}
}
