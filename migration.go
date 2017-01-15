package migrator

// Migration is a representation of a migration that needs to run.
type Migration struct {
	// FileName is the file name of the migration.
	FileName string

	// FileContents is the contents of the migration to run.
	FileContents string
}
