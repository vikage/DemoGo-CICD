package decoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)

// SLDecoder decoder type
type SLDecoder struct {
	r io.Reader
}

// NewSLDecoder create new decoder
func NewSLDecoder(r io.Reader) *SLDecoder {
	return &SLDecoder{
		r: r,
	}
}

// Decode decode object
func (decoder *SLDecoder) Decode(v interface{}) error {
	data, _ := ioutil.ReadAll(decoder.r)
	decoder.doDecode(data, v)

	return nil
}

func (decoder *SLDecoder) doDecode(data []byte, v interface{}) {
	var jsonMap map[string]interface{}
	json.NewDecoder(bytes.NewBuffer(data)).Decode(v)
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&jsonMap)

	objectValue := reflect.ValueOf(v)
	if objectValue.Kind() == reflect.Ptr {
		objectValue = objectValue.Elem()
	}

	objectType := reflect.TypeOf(v)
	if objectType.Kind() == reflect.Ptr {
		objectType = objectType.Elem()
	}

	for i := 0; i < objectValue.NumField(); i++ {
		fieldValue := objectValue.Field(i)

		fieldType := objectType.Field(i)
		jsonTag := fieldType.Tag.Get("json")
		tagElements := strings.Split(jsonTag, ",")

		if len(tagElements) == 0 {
			continue
		}

		jsonField := tagElements[0]
		valueFromJSON := jsonMap[jsonField]
		valueFromJSONType := reflect.TypeOf(valueFromJSON)

		switch fieldValue.Interface().(type) {
		case time.Time:
			dateFormat := fieldType.Tag.Get("formatter")
			t, err := time.Parse(dateFormat, valueFromJSON.(string))
			if err == nil {
				fieldValue.Set(reflect.ValueOf(t))
			}

			continue
		}

		if valueFromJSONType.Kind() == reflect.Map && fieldValue.Kind() == reflect.Struct {
			fmt.Println(valueFromJSON)
			valueFromJSONData, err := json.Marshal(valueFromJSON)
			if err == nil {
				decoder.doDecode(valueFromJSONData, fieldValue.Addr().Interface())
			}
		}
	}
}
