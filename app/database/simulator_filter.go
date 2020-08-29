package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-cicd/app/utils"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// SimulatorFilter simulator filter
type SimulatorFilter interface {
	Match(key string, data interface{}, positionals *[]int) bool
}

// {"$and": [{}, {}]}
type matchFilter struct {
	key        string
	operations []SimulatorFilter
}

type andOperation struct {
	conditionals []SimulatorFilter
}

type orOperation struct {
	conditionals []SimulatorFilter
}

type equalOperation struct {
	value interface{}
}

type neOperation struct {
	value interface{}
}

type inOperation struct {
	value []interface{}
}

func matchValueForKey(keypath string, data interface{}, expect interface{}, positionals *[]int) bool {
	if mapData, ok := data.(map[string]interface{}); ok {
		keys := strings.Split(keypath, ".")

		if len(keys) == 1 {
			return mapData[keys[0]] == expect
		}

		value := mapData[keys[0]]
		remainKeys := keys[1:]
		remainKeysString := strings.Join(remainKeys, ".")

		return matchValueForKey(remainKeysString, value, expect, positionals)
	}

	if array, ok := data.([]interface{}); ok {
		for index, item := range array {
			if matchValueForKey(keypath, item, expect, positionals) {
				*positionals = append(*positionals, index)
				return true
			}
		}
	}

	return false
}

func (filter *matchFilter) Match(key string, data interface{}, positionals *[]int) bool {
	for _, operation := range filter.operations {
		if operation.Match(filter.key, data, positionals) == false {
			return false
		}
	}

	return true
}

func (operation *equalOperation) Match(key string, data interface{}, positionals *[]int) bool {
	return matchValueForKey(key, data, operation.value, positionals)
}

func (operation *orOperation) Match(key string, data interface{}, positionals *[]int) bool {
	return false
}

func (operation *andOperation) Match(key string, data interface{}, positionals *[]int) bool {
	return false
}

func (operation *neOperation) Match(key string, data interface{}, positionals *[]int) bool {
	return !matchValueForKey(key, data, operation.value, positionals)
}

func (operation *inOperation) Match(key string, data interface{}, positionals *[]int) bool {
	return false
}

// ParseMongoFilter parse filters from mongo filter
func ParseMongoFilter(filter interface{}) ([]SimulatorFilter, error) {
	var filterMap map[string]interface{}
	data, err := bson.MarshalExtJSON(filter, false, false)
	if err != nil {
		panic(err)
	}

	json.NewDecoder(bytes.NewBuffer(data)).Decode(&filterMap)

	return ParseFilter(filterMap)
}

// ParseFilter parse filter from map
func ParseFilter(filter map[string]interface{}) ([]SimulatorFilter, error) {
	filters := make([]SimulatorFilter, 0)
	for key, value := range filter {
		var f SimulatorFilter

		if strings.Contains(key, "$") {
			// Opeation
			switch key {
			case "$eq":
				f = &equalOperation{value: value}
			case "$ne":
				f = &neOperation{value: value}
			case "$in":
				if arr, ok := value.([]interface{}); ok {
					f = &inOperation{value: arr}
				}

				if f == nil {
					return nil, fmt.Errorf("$in value must be array")
				}

			case "$and":
				andValue := value.(map[string]interface{})
				ops, err := ParseFilter(andValue)

				if err != nil {
					return nil, err
				}

				f = &andOperation{conditionals: ops}
			case "$or":
				orValue := value.(map[string]interface{})
				ops, err := ParseFilter(orValue)

				if err != nil {
					return nil, err
				}

				f = &orOperation{conditionals: ops}
			default:
				panic("Missing case")
			}
		}

		if f == nil {
			if mapFilter, ok := value.(map[string]interface{}); ok {
				ops, err := ParseFilter(mapFilter)

				if err != nil {
					return nil, err
				}

				f = &matchFilter{
					key:        key,
					operations: ops,
				}
			} else if utils.ExpectedValueIsType(value, []string{"int", "int64", "string", "float32", "float64"}) {
				eq := &equalOperation{
					value: value,
				}

				f = &matchFilter{
					key:        key,
					operations: []SimulatorFilter{eq},
				}
			}
		}

		if f != nil {
			filters = append(filters, f)
		}
	}

	return filters, nil
}
