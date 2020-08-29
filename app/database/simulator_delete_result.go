package database

// SimulatorDeleteResult bridge delete result
type SimulatorDeleteResult struct {
	deletedCount int64
}

// DeletedCount get deleted count
func (result *SimulatorDeleteResult) DeletedCount() int64 {
	return result.deletedCount
}
