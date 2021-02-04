package strongparams

import (
	"github.com/pkg/errors"
	"net/url"
	"reflect"
)

// StringParser should be a function of:
//   func(string) {<any output type>, error}
// The interface will be checked by the runner. If it doesn't comply it will panic.
type StringParser interface{}

func assertStringParser(candidate StringParser) {
	value := reflect.ValueOf(candidate)
	if value.Kind() != reflect.Func {
		panic(errors.New("passed strongparams.StringParser is not a function of " +
			"`func(string) {<parsed out>, error}`"))
	}

	valueType := value.Type()
	if valueType.NumIn() != 1 || valueType.In(0).Kind() != reflect.String {
		panic(errors.New("strongparams.StringParser parser function must have 1 argument with `string` type"))
	}

	if valueType.NumOut() != 2 {
		panic(errors.New("strongparams.StringParser parser function must have 2 return arguments with" +
			"first as a parsed type and second as `error` type"))

	} else if !valueType.Out(1).Implements(errorInterface) {
		panic(errors.New("strongparams.StringParser parser second return argument must implement `error` type"))
	}
}

func callStringParser(parser StringParser, requireKey string, queryValues url.Values) (interface{}, error) {
	keyValue := queryValues.Get(requireKey)
	result := reflect.ValueOf(parser).Call([]reflect.Value{
		reflect.ValueOf(keyValue),
	})

	if err, ok := result[1].Interface().(error); ok {
		return nil, errors.Wrapf(err, "failed to parse key: `%s`, value: `%s`", requireKey, keyValue)
	}
	return result[0].Interface(), nil
}
