package permittertest

import (
	"github.com/stretchr/testify/assert"
	. "github.com/vellotis/go-strongparams/permitter"
	"testing"
)

func Test_ParsePermitted_And_IsPermitted_SimpleKey(t *testing.T) {
	rule := "key"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("notPresent")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_SimpleKeyWithQuotes(t *testing.T) {
	rule := "'key'"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("notPresent")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_ObjectKey(t *testing.T) {
	rule := "key:{objKey}"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[objKey]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_ObjectKeyWithQuotes(t *testing.T) {
	rule := "key:{'objKey'}"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[objKey]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_Array(t *testing.T) {
	rule := "[]"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("[]")) &&
		assert.False(t, permittable.IsPermitted("notPresent")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_ArrayKey(t *testing.T) {
	rule := "key:[]"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_ArrayWithKeys(t *testing.T) {
	rule := "key:[objKey1,objKey2]"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[0][objKey1]")) &&
		assert.True(t, permittable.IsPermitted("key[0][objKey2]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[]")) &&
		assert.False(t, permittable.IsPermitted("key[0]")) &&
		assert.False(t, permittable.IsPermitted("key[0][]")) &&
		assert.False(t, permittable.IsPermitted("key[0][notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_NestedArrayKey(t *testing.T) {
	rule := "key:[[]]"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[0][]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[]")) &&
		assert.False(t, permittable.IsPermitted("key[0]")) &&
		assert.False(t, permittable.IsPermitted("key[0][notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_NestedArrayObjectKey(t *testing.T) {
	rule := "key:[[obj:{key}]]"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[0][0][obj][key]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[0]")) &&
		assert.False(t, permittable.IsPermitted("key[0][]")) &&
		assert.False(t, permittable.IsPermitted("key[0][obj]")) &&
		assert.False(t, permittable.IsPermitted("key[0][obj][]")) &&
		assert.False(t, permittable.IsPermitted("key[0][obj][notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_ObjectKeyWithAdditionalKeys(t *testing.T) {
	rule := "key:{add1,objKey,add2}"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[objKey]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[]")) &&
		assert.False(t, permittable.IsPermitted("key[notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_KeyInNestedObjectsOk(t *testing.T) {
	rule := "key:{x1Nested:{x2Nested:{objectKey}}}"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[x1Nested][x2Nested][objectKey]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][x2Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][x2Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][x2Nested][notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_KeyMissingInNestedObjectsOk(t *testing.T) {
	rule := "key:{x1Nested:{x2Nested:{objectKey}}}"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][x2Nested][notPresent]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][x2Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][x2Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested][x2Nested][notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_KeyInComplexStructure(t *testing.T) {
	rule := "key:{x1Nested1,x1Nested2:{x2Nested:[x4Nested:{add1,objectKey,add2}]}}"

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][objectKey]")) &&
		assert.True(t, permittable.IsPermitted("key[x1Nested1]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_KeyInComplexStructureWithSpaces(t *testing.T) {
	rule := "key : { x1Nested1 , x1Nested2 : { x2Nested : [ x4Nested : { add1 , objectKey , add2 } ] } } "

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][objectKey]")) &&
		assert.True(t, permittable.IsPermitted("key[x1Nested1]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_KeyInComplexStructureWithSpacesAndOnlyQuoted(t *testing.T) {
	rule := "'key' : { 'x1Nested1' , 'x1Nested2' : { 'x2Nested' : [ 'x4Nested' : { 'add1' , 'objectKey' , 'add2' } ] } } "

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][objectKey]")) &&
		assert.True(t, permittable.IsPermitted("key[x1Nested1]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][]")) &&
		assert.False(t, permittable.IsPermitted("key[x1Nested2][x2Nested][0][x4Nested][notPresent]")) {
	}
}

func Test_ParsePermitted_And_IsPermitted_KeyInComplexStructureWithSpacesAndOnlyQuotedIncludingSpaces(t *testing.T) {
	rule := "'key' : { 'nested1 x 1' , 'nested2 x 1' : { 'nested x 2' : [ 'nested x 4' : { 'add1' , 'object key' , 'add2' } ] } } "

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err)
		assert.True(t, permittable.IsPermitted("key[nested2 x 1][nested x 2][0][nested x 4][object key]")) &&
			assert.True(t, permittable.IsPermitted("key[nested1 x 1]")) &&
			assert.False(t, permittable.IsPermitted("key")) &&
			assert.False(t, permittable.IsPermitted("key[]")) &&
			assert.False(t, permittable.IsPermitted("key[nested2 x 1]")) &&
			assert.False(t, permittable.IsPermitted("key[nested2 x 1][]")) &&
			assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested x 2]")) &&
			assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested x 2][]")) &&
			assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested x 2][0]")) &&
			assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested x 2][0][nested x 4]")) &&
			assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested x 2][0][nested x 4][]")) &&
			assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested x 2][0][nested x 4][notPresent]")){
	}
}

func Test_ParsePermitted_And_IsPermitted_ComplexStructureWithSpacesAndOnlyQuotedIncludingSpaces(t *testing.T) {
	rule := "'key' : {" +
		" 'nested1 x 1' : {" +
			" 'nested1 x 2' : [" +
				" 'nested1 x 4' : { 'add1' , 'object1 key' , 'add2' } " +
			"] " +
		"}, " +
		" 'nested2 x 1' : {" +
			" 'nested2 x 2' : [" +
				" 'nested2 x 4' : { 'add3' , 'object2 key' , 'add4' } " +
			"] " +
		"} " +
	"} "

	permittable, err := ParsePermitted(rule)

	if assert.NoError(t, err) &&
		assert.True(t, permittable.IsPermitted("key[nested1 x 1][nested1 x 2][0][nested1 x 4][object1 key]")) &&
		assert.True(t, permittable.IsPermitted("key[nested2 x 1][nested2 x 2][0][nested2 x 4][object2 key]")) &&
		assert.False(t, permittable.IsPermitted("key")) &&
		assert.False(t, permittable.IsPermitted("key[]")) &&
		assert.False(t, permittable.IsPermitted("key[nested1 x 1]")) &&
		assert.False(t, permittable.IsPermitted("key[nested1 x 1][]")) &&
		assert.False(t, permittable.IsPermitted("key[nested1 x 1][nested1 x 2]")) &&
		assert.False(t, permittable.IsPermitted("key[nested1 x 1][nested1 x 2][]")) &&
		assert.False(t, permittable.IsPermitted("key[nested1 x 1][nested1 x 2][0]")) &&
		assert.False(t, permittable.IsPermitted("key[nested1 x 1][nested1 x 2][0][nested1 x 4]")) &&
		assert.False(t, permittable.IsPermitted("key[nested1 x 1][nested1 x 2][0][nested1 x 4][]")) &&
		assert.False(t, permittable.IsPermitted("key[nested1 x 1][nested1 x 2][0][nested1 x 4][notPresent]")) &&
		assert.False(t, permittable.IsPermitted("key[nested2 x 1]")) &&
		assert.False(t, permittable.IsPermitted("key[nested2 x 1][]")) &&
		assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested2 x 2]")) &&
		assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested2 x 2][]")) &&
		assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested2 x 2][0]")) &&
		assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested2 x 2][0][nested2 x 4]")) &&
		assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested2 x 2][0][nested2 x 4][]")) &&
		assert.False(t, permittable.IsPermitted("key[nested2 x 1][nested2 x 2][0][nested2 x 4][notPresent]")) {
	}
}
