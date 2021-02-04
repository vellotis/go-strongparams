package strongparams

import (
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"github.com/vellotis/go-strongparams/permitter"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

type StrongParams struct {
	*strongParams
}

type strongParams struct {
	decoder     *schema.Decoder
	valueGetter func() url.Values
}

// ReturnTarget enables just features of schema.Decoder (https://github.com/gorilla/schema) without performing
// any additional checks. schema.Decoder requires a dot notation of properties eg. "root.0.key" equivalent to query
// string "root[0][key]". Given method enables using the standard query string brackets format with schema.Decoder.
// Before passing the url.Values to schema.Decoder the keys are transposed to the required dot notation.
//
// The `target` parameter shall be a pointer to the parsable struct.
type ReturnTarget func(target interface {}) error

// Params declares the *http.Request to be used for the strong parameters mechanism.
//
// **NOTE** The method panics if the passed request parameter is nil.
func Params() *StrongParams {

	return &StrongParams{
		&strongParams{
			decoder: defaultDecoder,
		},
	}
}

// StrongParams.WithDecoder instructs the strong-parameters mechanism to use explicit decoder. The new decoder will be
// used on the returned *StrongParams struct pointer not on the receiver parameter. To use new implicitly defined
// schema.Decoder, look WithDecoder method instead.
func (this *StrongParams) WithDecoder(decoder *schema.Decoder) *StrongParams {
	return &StrongParams{
		&strongParams{
			decoder: decoder,
		},
	}
}

// Query instructs the mechanism to process http.Request's url.URL property url.URL/Query() method returned url.Values.
func (this *StrongParams) Query(request *http.Request) ReturnTarget {
	return this.Values(request.URL.Query())
}

// PostForm instructs the mechanism to process http.Request's url.PostForm property's url.Values.
func (this *StrongParams) PostForm(request *http.Request) ReturnTarget {
	return this.Values(request.PostForm)
}

// Values instructs the mechanism to process url.Values from `values` parameter.
func (this *StrongParams) Values(values url.Values) ReturnTarget {
	values = cloneUrlValues(values)

	return func(target interface{}) error {
		if target == nil {
			return errors.New("`target` argument cannot be nil")
		}

		return this.decode(values, target)
	}
}

// Require instructs to ensure before processing that the processable url.Values has by `requireKey` parameter provided
// root key. The chained StrongParamsRequired.Permit rules are applied on the object found behind the parameter
// `requireKey` defined key.
//   Go: Params().Require("root").Permit("sub:{key}")
//   Whitelisted query: root[sub][key]=value
// The `requireKey` value cannot be empty. Otherwise an error is produced but not returned until executing
// ReturnTarget function.
func (this *StrongParams) Require(requireKey string) *StrongParamsRequired {
	params := StrongParamsRequired{
		&strongParamsRequired{
			strongParams: this.strongParams,
			requireKey: &requireKey,
		},
	}

	if requireKey == "" {
		params.error = errors.New("required key value cannot be empty")
	}

	return &params
}

// RequireOne instructs to ensure presence and parse a single value in url.Values. `requireKey` parameter defines the
// key to retrieve from the url.Values.
func (this *StrongParams) RequireOne(requireKey string) *StrongParamsRequireOne {
	return &StrongParamsRequireOne{
		strongParamsRequired: &strongParamsRequired{
			strongParams: this.strongParams,
			requireKey: &requireKey,
		},
	}
}

// Permit instructs to apply the rules to whitelist the keys in url.Values before decoding it to the target struct.
func (this *StrongParams) Permit(permitRule string, permitRules... string) *StrongParamsRequiredAndPermitted {
	rules, err := permitter.ParsePermitted(
		funk.Uniq(append(permitRules, permitRule)).([]string)...
	)
	return &StrongParamsRequiredAndPermitted{
		&strongParamsRequiredAndPermitted{
			strongParamsRequired: &strongParamsRequired{
				strongParams: this.strongParams,
				error: err,
			},
			permitRules: rules,
		},
	}
}

var rgxMatchStartEndBrackets = regexp.MustCompile("(?:^\\[)|(?:\\]$)")
var rgxMatchMiddleBrackets = regexp.MustCompile("(?:\\]\\[)|(?:\\[)")
func transposeToDotNotation(dotNotationQueryKey string) string {
	transposedKey := rgxMatchStartEndBrackets.ReplaceAllString(dotNotationQueryKey, "")
	transposedKey = rgxMatchMiddleBrackets.ReplaceAllString(transposedKey, ".")
	return transposedKey
}

func (this *strongParams) decode(queryValues url.Values, target interface{}) error {
	err := schema.MultiError{}
	transposedQueryValues := url.Values{}

	for key, value := range queryValues {
		if strings.ContainsRune(key, '.') {
			err[key] = errors.Errorf("`%s` contains `.` character", key)
			continue
		}

		transposedKey := transposeToDotNotation(key)
		transposedQueryValues[transposedKey] = value
	}

	if funk.NotEmpty(err) {
		return errors.WithMessage(err, "any of the query keys should not contain `.` character." +
			"The brackets query notation is transposed to a dot notation which is required by the " +
			"`github.com/gorilla/struct` decoder.")
	}

	return this.decoder.Decode(target, transposedQueryValues)
}

