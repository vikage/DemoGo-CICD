package gdi

import (
	"fmt"
	"reflect"
)

// NewContainer create a new container
func NewContainer() *Container {
	return &Container{
		storage: make(map[string]interface{}),
	}
}

// Container storage dependency infomation
type Container struct {
	storage map[string]interface{}
}

// ResolveOptions resolve options
type ResolveOptions struct {
}

// Register register a dependency
func (c *Container) Register(aType reflect.Type, fn interface{}) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("fn must be a function")
	}

	// TODO: Detect circle dependency
	signature := signatureOfType(aType)
	c.storage[signature] = fn

	return nil
}

// Resolve resolve a dependency
func (c *Container) Resolve(t reflect.Type, args ...interface{}) (interface{}, error) {
	signature := signatureOfType(t)
	fn := c.storage[signature]

	if fn == nil {
		return nil, fmt.Errorf("%s not register yet", t.Name())
	}

	// Check params
	fType := reflect.TypeOf(fn)
	numArgs := fType.NumIn()

	params := make([]reflect.Value, 0)

	for i := 0; i < numArgs; i++ {
		paramType := fType.In(i)
		param := findDependencyInArgs(paramType, args...)

		if param == nil {
			var err error
			param, err = c.Resolve(paramType)
			if err != nil {
				return nil, err
			}
		}

		params = append(params, reflect.ValueOf(param))
	}

	returned := reflect.ValueOf(fn).Call(params)
	if len(returned) == 0 {
		return nil, fmt.Errorf("Resolve %s error", t.Name())
	}

	if last := returned[len(returned)-1]; isError(last.Type()) {
		if err, _ := last.Interface().(error); err != nil {
			return nil, err
		}
	}

	return returned[0].Interface(), nil
}

func isError(t reflect.Type) bool {
	errType := reflect.TypeOf((*error)(nil)).Elem()
	return t.Implements(errType)
}

func signatureOfType(t reflect.Type) string {
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}

func findDependencyInArgs(t reflect.Type, args ...interface{}) interface{} {
	for _, arg := range args {
		if arg == nil {
			continue
		}

		argType := reflect.TypeOf(arg)

		if t.Kind() == reflect.Interface {
			if argType.Implements(t) {
				return arg
			}
		} else {
			if signatureOfType(argType) == signatureOfType(t) {
				return arg
			}
		}
	}

	return nil
}
