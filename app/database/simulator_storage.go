package database

// SimulatorDatabaseStorage store simulator database
type SimulatorDatabaseStorage struct {
	data map[string]*SimulatorCollectionStorage
}

// GlobalDatabases global databases
var GlobalDatabases = make(map[string]*SimulatorDatabaseStorage)
