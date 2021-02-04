package strongparamstest

import (
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/assert"
	. "github.com/vellotis/go-strongparams"
	"testing"
)

func Test_Params_ImplicitDecoder(t *testing.T)  {
	values := mockQueryValues("key=value")
	type KeyValue struct {
		Key string `params:"key"`
	}
	result:= KeyValue{}

	err := Params().Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.Equal(t, "value", result.Key) {
	}
}

func Test_WithDecoder_Params_ExplicitDecoder(t *testing.T)  {
	values := mockQueryValues("key=value")
	type KeyValue struct {
		Key string `testTag:"key"`
	}
	result:= KeyValue{}
	decoder := schema.NewDecoder()
	decoder.SetAliasTag("testTag")

	err := WithDecoderSafe(decoder).Params().Values(values)(&result)

	if assert.NoError(t, err) &&
		assert.Equal(t, "value", result.Key) {
	}
}

func Test_Params_WithDecoder_MultipleChained(t *testing.T)  {
	values := mockQueryValues("key1=value1&key2=value2")
	type KeyValue struct {
		Key string `testTag1:"key1" testTag2:"key2"`
	}
	result1:= KeyValue{}
	decoder1 := schema.NewDecoder()
	decoder1.SetAliasTag("testTag1")
	decoder1.IgnoreUnknownKeys(true)
	result2:= KeyValue{}
	decoder2 := schema.NewDecoder()
	decoder2.SetAliasTag("testTag2")
	decoder2.IgnoreUnknownKeys(true)

	params1 := Params().WithDecoder(decoder1)
	params2 := params1.WithDecoder(decoder2)

	err1 := params1.Values(values)(&result1)
	err2 := params2.Values(values)(&result2)

	if assert.NoError(t, err1) &&
		assert.NoError(t, err2) &&
		assert.Equal(t, "value1", result1.Key) &&
		assert.Equal(t, "value2", result2.Key) {
	}
}
