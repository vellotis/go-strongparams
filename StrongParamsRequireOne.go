package strongparams

import (
	"github.com/amsokol/ignite-go-client/binary/errors"
	"net/http"
	"net/url"
)

type StrongParamsRequireOne struct {
	*strongParamsRequired
}

// ReturnOfType enables parsing and returning the single key from url.Values declared by StrongParams.RequireOne method
// or returns an error if it fails.
//
// The StringParser parameter shall be a function definition that enables parsing string value to specific required
// type.
type ReturnOfType func(StringParser) (interface {}, error)

// Query instructs the mechanism to process http.Request's url.URL property url.URL/Query() method returned url.Values.
func (this *StrongParamsRequireOne) Query(request *http.Request) ReturnOfType {
	return this.Values(request.URL.Query())
}

// PostForm instructs the mechanism to process http.Request's url.PostForm property's url.Values.
func (this *StrongParamsRequireOne) PostForm(request *http.Request) ReturnOfType {
	return this.Values(request.PostForm)
}

// Values instructs the mechanism to process url.Values from `values` parameter.
func (this *StrongParamsRequireOne) Values(values url.Values) ReturnOfType {
	values = cloneUrlValues(values)

	return func(parser StringParser) (interface{}, error) {
		assertStringParser(parser)

		if err := this.validate(values); err != nil {
			return nil, err
		}

		return callStringParser(parser, *this.requireKey, values)
	}
}

func (this *StrongParamsRequireOne) validate(values url.Values) error {
	if this.requireKey != nil {
		_, hasKey := values[*this.requireKey]
		if !hasKey {
			return errors.Errorf("query: missing required key: `%s`", *this.requireKey)
		}
	}

	return nil
}
