package strongparams

import (
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"github.com/vellotis/go-strongparams/permitter"
	"net/http"
	"net/url"
	"strings"
)

type StrongParamsRequired struct {
	*strongParamsRequired
}

type strongParamsRequired struct {
	*strongParams
	requireKey *string
	error      error
}

// Permit instructs to apply the rules to whitelist the keys in url.Values before decoding it to the target struct.
// The chained Permit rules are applied on the object found behind the parameter `requireKey` defined key instructed by
// StrongParams.Require method.
//   Go: Params().Require("root").Permit("sub:{key}")
//   Whitelisted query: root[sub][key]=value
// These two use cases are equivalent:
//   Permit("[key1, key2]")
//   Permit("key1", "key2")
func (this *StrongParamsRequired) Permit(permitRule  string, permitRules... string) *StrongParamsRequiredAndPermitted {
	params := StrongParamsRequiredAndPermitted{
		&strongParamsRequiredAndPermitted{
			strongParamsRequired: this.strongParamsRequired,
		},
	}

	if params.error == nil {
		params.permitRules, params.error = permitter.ParsePermitted(
			funk.Uniq(append(permitRules, permitRule)).([]string)...
		)
	}

	return &params
}

// Query instructs the mechanism to process http.Request's url.URL property url.URL/Query() method returned url.Values.
func (this *StrongParamsRequired) Query(request *http.Request) ReturnTarget {
	return this.Values(request.URL.Query())
}

// PostForm instructs the mechanism to process http.Request's url.PostForm property's url.Values.
func (this *StrongParamsRequired) PostForm(request *http.Request) ReturnTarget {
	return this.Values(request.PostForm)
}

// Values instructs the mechanism to process url.Values from `values` parameter.
func (this *StrongParamsRequired) Values(values url.Values) ReturnTarget {
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

func (this *strongParamsRequired) validateTransformAndDecode(values url.Values, target interface{}) (err error) {
	if err := this.validate(values); err != nil {
		return err
	}

	values, err = this.transform(values)
	if err != nil {
		return err
	}

	return this.decode(values, target)
}

func (this *strongParamsRequired) validate(values url.Values) error {
	if this.requireKey != nil {
		if !hasKey(values, *this.requireKey) {
			return errors.Errorf("query: missing required key: `%s`", *this.requireKey)
		}
	}

	return nil
}

func (this *strongParamsRequired) transform(values url.Values) (url.Values, error) {
	if this.requireKey != nil {
		requiredQueryValues := make(url.Values)
		for path, value := range values {
			if strings.HasPrefix(path, *this.requireKey+"[") {
				newPath := strings.Replace(path, *this.requireKey+"[", "", 1)
				newPath = strings.Replace(newPath, "]", "", 1)
				requiredQueryValues[newPath] = make([]string, len(value))
				copy(requiredQueryValues[newPath], value)
			}
		}

		return requiredQueryValues, nil
	}

	return values, nil
}
