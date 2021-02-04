package strongparamstest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	. "github.com/vellotis/go-strongparams"
	"strconv"
	"testing"
)

func Test_Params(t *testing.T) {
	expectedValue := "value"
	values := mockQueryValues("root[key]=" + expectedValue)
	type RootValue struct {
		Root struct {
			Key string `params:"key,required"`
		} `params:"root,required"`
	}
	result := RootValue{}

	err := Params().Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Root.Key) {
	}
}

func Test_Params_FromQuery(t *testing.T) {
	expectedValue := "value"
	request := mockRequestWithQuery("root[key]=" + expectedValue)
	type RootValue struct {
		Root struct {
			Key string `params:"key,required"`
		} `params:"root,required"`
	}
	result := RootValue{}

	err := Params().Query(request)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Root.Key) {
	}
}

func Test_Params_FromPostForm(t *testing.T) {
	expectedValue := "value"
	request := mockRequestWithPostForm("root[key]=" + expectedValue)
	type RootValue struct {
		Root struct {
			Key string `params:"key,required"`
		} `params:"root,required"`
	}
	result := RootValue{}

	err := Params().PostForm(request)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Root.Key) {
	}
}

func Test_Require(t *testing.T) {
	expectedValue := "value"
	values := mockQueryValues("root[key]=" + expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Require("root").Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Require_FromQuery(t *testing.T) {
	expectedValue := "value"
	request := mockRequestWithQuery("root[key]=" + expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Require("root").Query(request)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Require_FromPostForm(t *testing.T) {
	expectedValue := "value"
	request := mockRequestWithPostForm("root[key]=" + expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Require("root").PostForm(request)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Require_FailsMissingKey(t *testing.T) {
	missingKey := "missingKey"
	values := mockQueryValues(missingKey + "=value")
	type KeyValue struct {
		Key string `params:"invalidKey"`
	}
	result := KeyValue{}

	err := Params().Require("root").Values(values)(&result)

	assert.EqualError(t, err, "query: missing required key: `root`")
}

func Test_RequireOne(t *testing.T) {
	expectedValue := 4
	values := mockQueryValues("key=" + fmt.Sprint(expectedValue))

	value, err := Params().RequireOne("key").Values(values)(strconv.Atoi)

	if assert.NoError(t, err) &&
		assert.Equal(t, expectedValue, value.(int)) {
	}
}

func Test_RequireOne_FromQuery(t *testing.T) {
	expectedValue := 4
	request := mockRequestWithQuery("key=" + fmt.Sprint(expectedValue))

	value, err := Params().RequireOne("key").Query(request)(strconv.Atoi)

	if assert.NoError(t, err) &&
		assert.Equal(t, expectedValue, value.(int)) {
	}
}

func Test_RequireOne_FromPostForm(t *testing.T) {
	expectedValue := 4
	request := mockRequestWithPostForm("key=" + fmt.Sprint(expectedValue))

	value, err := Params().RequireOne("key").PostForm(request)(strconv.Atoi)

	if assert.NoError(t, err) &&
		assert.Equal(t, expectedValue, value.(int)) {
	}
}

func Test_RequireOne_MissingRequiredKey(t *testing.T) {
	invalidKey := "invalidKey"
	values := mockQueryValues(invalidKey+"=someValue")

	value, err := Params().RequireOne("key").Values(values)(strconv.Atoi)

	if assert.Nil(t, value) &&
		assert.EqualError(t, err, "query: missing required key: `key`") {
	}
}

func Test_RequireOne_FailsToParse(t *testing.T) {
	failingValue := "invalid"
	values := mockQueryValues("key=" + failingValue)

	value, err := Params().RequireOne("key").Values(values)(strconv.Atoi)

	if assert.Nil(t, value) &&
		assert.EqualError(t, err, "failed to parse key: `key`, value: `invalid`: strconv.Atoi: parsing \"invalid\": invalid syntax") {
	}
}

func Test_Permit(t *testing.T) {
	expectedValue := "value"
	values := mockQueryValues("key="+expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Permit("key").Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Permit_FromQuery(t *testing.T) {
	expectedValue := "value"
	request := mockRequestWithQuery("key="+expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Permit("key").Query(request)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Permit_FromPostForm(t *testing.T) {
	expectedValue := "value"
	request := mockRequestWithPostForm("key="+expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Permit("key").PostForm(request)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Permit_NestedObject(t *testing.T) {
	someValue := "someValue"
	values := mockQueryValues("root[key1]=someValue&root[key2]=someValue")
	type RootValue struct {
		Root struct {
			Key1 string `params:"key1"`
			Key2 string `params:"key2"`
		} `params:"root"`
	}
	result := RootValue{}

	err := Params().Permit("root:{key1,key2}").Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.Equal(t, result.Root.Key1, someValue) &&
		assert.Equal(t, result.Root.Key2, someValue) {
	}
}

func Test_Require_Permit(t *testing.T) {
	expectedValue := "value"
	values := mockQueryValues("root[key]="+expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Require("root").Permit("key").Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Require_Permit_FromQuery(t *testing.T) {
	expectedValue := "value"
	request := mockRequestWithQuery("root[key]="+expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Require("root").Permit("key").Query(request)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Require_Permit_FromPostForm(t *testing.T) {
	expectedValue := "value"
	request := mockRequestWithPostForm("root[key]="+expectedValue)
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Require("root").Permit("key").PostForm(request)(&result)

	if assert.NoError(t, err) &&
		assert.NotZero(t, result) &&
		assert.Equal(t, expectedValue, result.Key) {
	}
}

func Test_Require_Permit_MissingRequiredKey(t *testing.T) {
	values := mockQueryValues("invalidKey[key]=someValue")
	type KeyValue struct {
		Key string `params:"key"`
	}
	result := KeyValue{}

	err := Params().Require("root").Permit("key").Values(values)(&result)

	assert.EqualError(t, err, "query: missing required key: `root`")
}

func Test_Require_Permit_InvalidTargetStruct(t *testing.T) {
	values := mockQueryValues("root[key1]=someValue&root[key2]=someValue")
	type InvalidStruct struct {
		SomeKey string `params:"someKey"`
	}
	result := InvalidStruct{}

	err := Params().Require("root").Permit("key1,key2").Values(values)(&result)

	if assert.Error(t, err) &&
		assert.Regexp(t, "^schema: invalid path \"key(1|2)\" \\(and 1 other error\\)$", err.Error()) {
	}
}

func Test_Require_Permit_NestedObject(t *testing.T) {
	someValue := "someValue"
	values := mockQueryValues("root[key1]=someValue&root[key2]=someValue")
	type KeyValue struct {
		Key1 string `params:"key1"`
		Key2 string `params:"key2"`
	}
	result := KeyValue{}

	err := Params().Require("root").Permit("key1,key2").Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.Equal(t, result.Key1, someValue) &&
		assert.Equal(t, result.Key2, someValue) {
	}
}

func Test_Require_Permit_AdvancedStruct(t *testing.T) {
	values := mockQueryValues("root[arr][0][key]=key1&root[arr][0][value]=value1&root[arr][1][key]=key2&root[arr][1][value]=value2")
	type Struct struct {
		Arr []struct{
			Key string `params:"key"`
			Value string  `params:"value"`
		} `params:"arr"`
	}
	result := Struct{}

	err := Params().Require("root").Permit("arr:[key,value]").Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.Equal(t, 2, len(result.Arr)) &&
		assert.Equal(t, "key1", result.Arr[0].Key) &&
		assert.Equal(t, "value1", result.Arr[0].Value) &&
		assert.Equal(t, "key2", result.Arr[1].Key) &&
		assert.Equal(t, "value2", result.Arr[1].Value) {
	}
}
