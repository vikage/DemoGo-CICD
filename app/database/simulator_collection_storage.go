package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-cicd/app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SimulatorCollectionStorage store collection data
type SimulatorCollectionStorage struct {
	data []map[string]interface{}
}

func checkFilters(filters []SimulatorFilter, data interface{}) (bool, []int) {
	for _, filter := range filters {
		positional := make([]int, 0)
		if filter.Match("", data, &positional) == true {
			return true, positional
		}
	}

	return false, []int{}
}

// Find find in db
func (storage *SimulatorCollectionStorage) Find(filter interface{}) (Cursor, error) {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	var matches []map[string]interface{}
	for _, document := range storage.data {
		matched, _ := checkFilters(filters, document)
		if matched {
			matches = append(matches, document)
		}
	}

	return &SimulatorCursor{data: matches, currentIndex: -1}, nil
}

// Insert insert data
func (storage *SimulatorCollectionStorage) Insert(document interface{}, opts ...*options.InsertOneOptions) (InsertOneResult, error) {
	data, err := bson.MarshalExtJSON(document, false, false)
	if err != nil {
		return nil, err
	}

	var documentMap map[string]interface{}
	err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&documentMap)
	ID := documentMap["_id"]
	if ID == nil {
		ID := utils.GenUUIDString()
		documentMap["_id"] = ID
	}

	if documentMap == nil {
		return nil, fmt.Errorf("Marshal document error: %s", err)
	}

	storage.data = append(storage.data, documentMap)

	return &SimulatorInsertOneResult{ID: ID}, nil
}

// ReplaceOne replace a document
func (storage *SimulatorCollectionStorage) ReplaceOne(filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (UpdateResult, error) {
	// TODO: Support options
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	firstMatchIndex := -1
	for index, document := range storage.data {
		matched, _ := checkFilters(filters, document)
		if matched {
			firstMatchIndex = index
			break
		}
	}

	if firstMatchIndex == -1 {
		return &SimulatorUpdateResult{matchedCount: 0, modifiedCount: 0, upsertedCount: 0, upsertedID: 0}, nil
	}

	data, err := bson.MarshalExtJSON(replacement, false, false)
	if err != nil {
		return nil, err
	}

	var documentMap map[string]interface{}
	err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&documentMap)
	if err != nil {
		return nil, err
	}

	storage.data[firstMatchIndex] = documentMap

	return &SimulatorUpdateResult{matchedCount: 1, modifiedCount: 1, upsertedCount: 0, upsertedID: 0}, nil
}

// FindOne return first match
func (storage *SimulatorCollectionStorage) FindOne(filter interface{}, opts ...*options.FindOneOptions) SingleResult {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return &SimulatorSingleResult{err: err}
	}

	firstMatchIndex := -1
	for index, document := range storage.data {
		matched, _ := checkFilters(filters, document)
		if matched {
			firstMatchIndex = index
		}
	}

	if firstMatchIndex == -1 {
		return &SimulatorSingleResult{}
	}

	return &SimulatorSingleResult{data: &(storage.data[firstMatchIndex])}
}

// DeleteOne delete first match
func (storage *SimulatorCollectionStorage) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error) {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	firstMatchIndex := -1
	for index, document := range storage.data {
		matched, _ := checkFilters(filters, document)
		if matched {
			firstMatchIndex = index
		}
	}

	if firstMatchIndex == -1 {
		return &SimulatorDeleteResult{deletedCount: 0}, nil
	}

	storage.data = append(storage.data[:firstMatchIndex], storage.data[firstMatchIndex+1:]...)
	return &SimulatorDeleteResult{deletedCount: 1}, nil
}

// DeleteMany Delete many match
func (storage *SimulatorCollectionStorage) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error) {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	newData := make([]map[string]interface{}, 0)

	for _, document := range storage.data {
		matched, _ := checkFilters(filters, document)
		if matched {
			newData = append(newData, document)
		}
	}

	deletedCount := len(storage.data) - len(newData)
	storage.data = newData

	return &SimulatorDeleteResult{deletedCount: int64(deletedCount)}, nil
}

// FindOneAndDelete find one and delete. Return match
func (storage *SimulatorCollectionStorage) FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) SingleResult {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return &SimulatorSingleResult{err: err}
	}

	firstMatchIndex := -1
	for index, document := range storage.data {
		matched, _ := checkFilters(filters, document)
		if matched {
			firstMatchIndex = index
		}
	}

	if firstMatchIndex == -1 {
		return &SimulatorSingleResult{}
	}

	matchedData := storage.data[firstMatchIndex]
	storage.data = append(storage.data[:firstMatchIndex], storage.data[firstMatchIndex+1:]...)

	return &SimulatorSingleResult{data: &matchedData}
}

// UpdateMany update many document
func (storage *SimulatorCollectionStorage) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error) {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	updates, err := ParseMongoUpdate(update)
	if err != nil {
		return nil, err
	}

	updateResult := SimulatorUpdateResult{}
	for _, document := range storage.data {
		matched, positionals := checkFilters(filters, document)
		if matched {
			updateResult.matchedCount++
			updateResult.modifiedCount++

			for _, updateOp := range updates {
				updateOp.Update(&document, &positionals)
			}
		}
	}

	return &updateResult, nil
}

// UpdateOne update first match document
func (storage *SimulatorCollectionStorage) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error) {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	updates, err := ParseMongoUpdate(update)
	if err != nil {
		return nil, err
	}

	for _, document := range storage.data {
		matched, positionals := checkFilters(filters, document)
		if matched {
			for _, updateOp := range updates {
				updateOp.Update(&document, &positionals)
			}

			return &SimulatorUpdateResult{matchedCount: 1, modifiedCount: 1}, nil
		}
	}

	if len(opts) > 0 {
		opt := opts[0]
		willUpsert := false

		if *opt.Upsert == true {
			willInsertDocument := make(map[string]interface{})
			if filterPrimitiveM, ok := filter.(primitive.M); ok {
				id, ok := filterPrimitiveM["_id"].(string)
				if ok == false {
					return nil, fmt.Errorf("Upsert require _id")
				}

				willInsertDocument["_id"] = id
				storage.data = append(storage.data, willInsertDocument)
				willUpsert = true
			}
		}

		if willUpsert {
			for _, document := range storage.data {
				matched, positionals := checkFilters(filters, document)
				if matched {
					for _, updateOp := range updates {
						updateOp.Update(&document, &positionals)
					}

					return &SimulatorUpdateResult{matchedCount: 1, modifiedCount: 1}, nil
				}
			}
		}
	}

	return &SimulatorUpdateResult{matchedCount: 0, modifiedCount: 0}, nil
}

// FindOneAndReplace find one and replace. Return match result
func (storage *SimulatorCollectionStorage) FindOneAndReplace(filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) SingleResult {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return &SimulatorSingleResult{err: err}
	}

	firstMatchIndex := -1
	for index, document := range storage.data {
		matched, _ := checkFilters(filters, document)
		if matched {
			firstMatchIndex = index
		}
	}

	if firstMatchIndex == -1 {
		return &SimulatorSingleResult{}
	}

	document := storage.data[firstMatchIndex]
	data, err := bson.MarshalExtJSON(replacement, false, false)
	if err != nil {
		return &SimulatorSingleResult{err: err}
	}

	var documentMap map[string]interface{}
	err = json.NewDecoder(bytes.NewBuffer(data)).Decode(&documentMap)
	if err != nil {
		return &SimulatorSingleResult{err: err}
	}

	storage.data[firstMatchIndex] = document

	return &SimulatorSingleResult{data: &document}
}

// FindOneAndUpdate find one and update. Return first match
func (storage *SimulatorCollectionStorage) FindOneAndUpdate(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) SingleResult {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return &SimulatorSingleResult{err: err}
	}

	updates, err := ParseMongoUpdate(update)
	if err != nil {
		return &SimulatorSingleResult{err: err}
	}

	firstMatchIndex := -1
	positionals := make([]int, 0)
	for index, document := range storage.data {
		matched, pos := checkFilters(filters, document)
		if matched {
			firstMatchIndex = index
			positionals = pos
		}
	}

	if firstMatchIndex == -1 {
		return &SimulatorSingleResult{}
	}

	document := storage.data[firstMatchIndex]
	originalDocument := utils.CopyMap(document)
	for _, updateOp := range updates {
		updateOp.Update(&document, &positionals)
	}

	if len(opts) == 0 {
		return &SimulatorSingleResult{data: &originalDocument}
	}

	opt := opts[0]
	if *opt.ReturnDocument == options.After {
		return &SimulatorSingleResult{data: &document}
	}

	return &SimulatorSingleResult{data: &originalDocument}
}

// CountDocuments count documents by filter
func (storage *SimulatorCollectionStorage) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	filters, err := ParseMongoFilter(filter)
	if err != nil {
		return 0, err
	}

	var matches []map[string]interface{}
	for _, document := range storage.data {
		matched, _ := checkFilters(filters, document)
		if matched {
			matches = append(matches, document)
		}
	}

	return int64(len(matches)), nil
}
