package strongparams

import (
	"github.com/pkg/errors"
	"github.com/vellotis/go-strongparams/permitter"
	"net/http"
	"net/url"
)

type StrongParamsRequiredAndPermitted struct {
	*strongParamsRequiredAndPermitted
}

type strongParamsRequiredAndPermitted struct {
	*strongParamsRequired
	permitRules permitter.Permittable
}

// Query instructs the mechanism to process http.Request's url.URL property url.URL/Query() method returned url.Values.
func (this *StrongParamsRequiredAndPermitted) Query(request *http.Request) ReturnTarget {
	return this.Values(request.URL.Query())
}

// PostForm instructs the mechanism to process http.Request's url.PostForm property's url.Values.
func (this *StrongParamsRequiredAndPermitted) PostForm(request *http.Request) ReturnTarget {
	return this.Values(request.PostForm)
}

// Values instructs the mechanism to process url.Values from `values` parameter.
func (this *StrongParamsRequiredAndPermitted) Values(values url.Values) ReturnTarget {
	values = cloneUrlValues(values)

	if this.error == nil {
		this.error = this.validate(values)
	}

	return func(target interface{}) error {
		if this.error != nil {
			return this.error
		} else if target == nil {
			return errors.New("`target` argument cannot be nil")
		}

		return this.validateTransformAndDecode(values, target)
	}
}

func (this *StrongParamsRequiredAndPermitted) validateTransformAndDecode(values url.Values, target interface{}) (err error) {
	if err := this.validate(values); err != nil {
		return err
	}

	values, err = this.transform(values)
	if err != nil {
		return err
	}

	return this.decode(values, target)
}

func (this *strongParamsRequiredAndPermitted) validate(values url.Values) error {
	if err := this.strongParamsRequired.validate(values); err != nil {
		return err
	}

	return nil
}

func (this *strongParamsRequiredAndPermitted) transform(queryValues url.Values) (_ url.Values, err error) {
	queryValues, err = this.strongParamsRequired.transform(queryValues)
	if err != nil {
		return nil, err
	}

	for queryKeyPath := range queryValues {
		if !this.permitRules.IsPermitted(queryKeyPath) {
			queryValues.Del(queryKeyPath)
		}
	}

	return queryValues, nil
}
