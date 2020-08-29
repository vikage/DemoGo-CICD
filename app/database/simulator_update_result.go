package database

// SimulatorUpdateResult update result bridge
type SimulatorUpdateResult struct {
	matchedCount  int64
	modifiedCount int64
	upsertedCount int64
	upsertedID    interface{}
}

// MatchedCount get matched count
func (result *SimulatorUpdateResult) MatchedCount() int64 {
	return result.matchedCount
}

// ModifiedCount get modified count
func (result *SimulatorUpdateResult) ModifiedCount() int64 {
	return result.modifiedCount
}

// UpsertedCount get upserted count
func (result *SimulatorUpdateResult) UpsertedCount() int64 {
	return result.upsertedCount
}

// UpsertedID get upserted id
func (result *SimulatorUpdateResult) UpsertedID() interface{} {
	return result.upsertedID
}
