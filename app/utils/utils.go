package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go-cicd/app/logger"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// GeneratePassword gen password with length
func GeneratePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return StringWithCharset(length, charset)
}

// EncryptPassword encrypt password
func EncryptPassword(password string, salt string) string {
	text := password + salt
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// StringWithCharset generate string with length by charset
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// ExistKeyInMap check exist keys exist in map
func ExistKeyInMap(m map[string]interface{}, keys []string, expectedTypes []string) bool {
	for index, key := range keys {
		value := m[key]
		expectedType := expectedTypes[index]
		realType := fmt.Sprintf("%T", value)
		if realType != expectedType {
			return false
		}
	}

	return true
}

// ExistKeyInMapString check exist key in map
func ExistKeyInMapString(m map[string]string, keys []string) bool {
	for _, key := range keys {
		value := m[key]

		if value == "" {
			return false
		}
	}

	return true
}

// GenUUIDString generate uuid
func GenUUIDString() string {
	s, _ := uuid.NewUUID()
	if len(s) == 0 {
		panic("Gen uuid fail")
	}

	return s.String()
}

// BytesFromStructValue Marshal struct value
func BytesFromStructValue(value interface{}) []byte {
	bytes, err := json.Marshal(value)
	if err != nil {
		logger.Error("Marshal error: %s value: %s", err, value)
	}

	return bytes
}

// ExpectedValueIsType check value has type in list
func ExpectedValueIsType(value interface{}, types []string) bool {
	valueType := fmt.Sprintf("%T", value)
	for _, expectedType := range types {
		if valueType == expectedType {
			return true
		}
	}

	return false
}

// GetFirstCallAgrumentsFromCalls get call args
func GetFirstCallAgrumentsFromCalls(calls []mock.Call, method string) *mock.Arguments {
	for _, call := range calls {
		if call.Method == method {
			return &call.Arguments
		}
	}

	return nil
}

// CopyMap duplicate map
func CopyMap(originalMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

// IsSameDay check 2 date is same day
func IsSameDay(date1 time.Time, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}
