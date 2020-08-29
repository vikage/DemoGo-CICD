package database

// SimulatorInsertOneResult insert one result bridge
type SimulatorInsertOneResult struct {
	ID interface{}
}

// InsertedID get inserted id
func (result *SimulatorInsertOneResult) InsertedID() interface{} {
	return result.ID
}
