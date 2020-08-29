package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// SimulatorUpdate simulator update
type SimulatorUpdate interface {
	Update(data *map[string]interface{}, positionals *[]int) bool
}

type setOperation struct {
	keys   []string
	values []interface{}
}

type unsetOperation struct {
	keys []string
}

func updateForKey(keypath string, data interface{}, updateValue interface{}, positionals []int, isUnset bool) (bool, error) {
	if mapData, ok := data.(*map[string]interface{}); ok {
		keys := strings.Split(keypath, ".")

		if len(keys) == 1 {
			if isUnset {
				delete((*mapData), keys[0])
			} else {
				(*mapData)[keys[0]] = updateValue
			}

			return true, nil
		}

		var pointer interface{}
		if mapValue, ok := (*mapData)[keys[0]].(map[string]interface{}); ok {
			pointer = &mapValue
		}

		if arrayValue, ok := (*mapData)[keys[0]].([]interface{}); ok {
			pointer = &arrayValue
		}

		remainKeys := keys[1:]
		remainKeysString := strings.Join(remainKeys, ".")

		return updateForKey(remainKeysString, pointer, updateValue, positionals, isUnset)
	}

	if array, ok := data.(*[]interface{}); ok {
		keys := strings.Split(keypath, ".")

		if keys[0] == "$" {
			if len(positionals) == 0 {
				return false, fmt.Errorf("The positional operator did not find the match needed from the query")
			}

			position := positionals[0]
			remainKeys := keys[1:]
			remainKeysString := strings.Join(remainKeys, ".")
			item := (*array)[position]

			var pointer interface{}
			if mapValue, ok := item.(map[string]interface{}); ok {
				pointer = &mapValue
			}

			if arrayValue, ok := item.([]interface{}); ok {
				pointer = &arrayValue
			}

			return updateForKey(remainKeysString, pointer, updateValue, positionals[1:], isUnset)
		}

		for _, item := range *array {
			var pointer interface{}
			if mapValue, ok := item.(map[string]interface{}); ok {
				pointer = &mapValue
			}

			if arrayValue, ok := item.([]interface{}); ok {
				pointer = &arrayValue
			}

			_, err := updateForKey(keypath, pointer, updateValue, positionals, isUnset)
			if err != nil {
				return false, err
			}
		}

		return false, nil
	}

	return false, nil
}

func (op *setOperation) Update(data *map[string]interface{}, positionals *[]int) bool {
	matched := false
	for index, key := range op.keys {
		value := op.values[index]

		copiedPositional := *positionals
		if match, _ := updateForKey(key, data, value, copiedPositional, false); match {
			matched = match
		}
	}

	return matched
}

func (op *unsetOperation) Update(data *map[string]interface{}, positionals *[]int) bool {
	matched := false
	for _, key := range op.keys {
		copiedPositional := *positionals

		if match, _ := updateForKey(key, data, nil, copiedPositional, true); match {
			matched = true
		}
	}

	return matched
}

func newSetOperation(update interface{}) (SimulatorUpdate, error) {
	if updateMap, ok := update.(map[string]interface{}); ok {
		if len(updateMap) == 0 {
			return nil, fmt.Errorf("Update invalid")
		}

		operation := setOperation{}
		for key, value := range updateMap {
			operation.keys = append(operation.keys, key)
			operation.values = append(operation.values, value)
		}

		return &operation, nil
	}

	return nil, fmt.Errorf("Update must be a map[string]interface{}")
}

func newUnsetOperation(update interface{}) (SimulatorUpdate, error) {
	if updateMap, ok := update.(map[string]interface{}); ok {
		if len(updateMap) == 0 {
			return nil, fmt.Errorf("Update invalid")
		}

		operation := unsetOperation{}
		for key := range updateMap {
			operation.keys = append(operation.keys, key)
		}

		return &operation, nil
	}

	return nil, fmt.Errorf("Update must be a map[string]interface{}")
}

// ParseMongoUpdate parse update from mongo update struct
func ParseMongoUpdate(update interface{}) ([]SimulatorUpdate, error) {
	var updateMap map[string]interface{}
	data, err := bson.MarshalExtJSON(update, false, false)
	if err != nil {
		panic(err)
	}

	json.NewDecoder(bytes.NewBuffer(data)).Decode(&updateMap)

	return ParseUpdate(updateMap)
}

// ParseUpdate parse updates from update map
func ParseUpdate(update map[string]interface{}) ([]SimulatorUpdate, error) {
	updates := make([]SimulatorUpdate, 0)
	for key, value := range update {
		var f SimulatorUpdate

		if strings.Contains(key, "$") {
			var err error
			switch key {
			case "$set":
				f, err = newSetOperation(value)
				if err != nil {
					return nil, err
				}
			case "$unset":
				f, err = newUnsetOperation(value)
				if err != nil {
					return nil, err
				}
			default:
				panic("Missing operation")
			}
		}

		if f != nil {
			updates = append(updates, f)
		}
	}

	return updates, nil
}
