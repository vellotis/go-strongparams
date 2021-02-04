package permitter

import (
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"regexp"
	"strings"
)

const idxGroup = 0

// Permittable exposes the common public interface that enables verifying if provided query path is permitted
type Permittable interface {
	permittable
	// IsPermitted verifies that the URL query path is permitted by the permitter rules and returns the boolean result.
	IsPermitted(path string) bool
}

type permittable interface {
	isPermitted(rgxResult [][]string) bool
}

var rgxQueryPathFull = regexp.MustCompile("^(?:[^\\[\\]]+)?(?:\\[[^\\[\\]]+\\])*(?:\\[\\])?$")
func isPermitted(perm permittable, path string) bool {
	return rgxQueryPathFull.MatchString(path) &&
		perm.isPermitted(
			rgxQueryPathGroup.FindAllStringSubmatch(path, -1),
		)
}

// ParsePermitted parses and builds query string keys whitelisting rules or returns an error. Calling
//  ParsePermitted("Literal", "Literal", "Literal")
// is equivalent to
//  ParsePermitted("[Literal, Literal, Literal]")
// .
//
// The rules can have the following elements:
//
//   • Key (KeyLiteral) defines a whitelisted key:
//
//     - string literal of ASCII numbers and letters, eg. "someKey"
//
//     - single quote (') wrapped string literal of ASCII numbers, letters and spaces, eg. "'some key'"
//
//  
//
//   • Object (ObjectLiteral) defines a whitelisted object containing any nested literals:
//
//     - "{ KeyLiteral, KeyLiteral }" eg.
//     "{ key1, key2 }"
//       matches query
//     "key1=value1&key2=value2"
//
//     - "{ KeyLiteral:Array, KeyLiteral:ObjectLiteral, KeyLiteral }" eg.
//     "{ key1:[], key2:{objKey}, key3}"
//       matches query
//     "key1[]=value&key2[objKey]=objValue&key3=keyValue"
//
//   **NOTE** rule "{}" doesn't match anything
//
//  
//
//   • Array (ArrayLiteral) defines a whitelisted array containing any nested literals:
//
//     - "[]" matches array of string values, eg.
//     "[]=value1&[]value2"
//       or
//     "[0]=value1&[1]value2"
//
//     - "[ KeyLiteral ]" eg.
//     "[ key ]"
//       matches query
//     "[0][key]=value"
//
//     - "[ KeyLiteral:Array, KeyLiteral:ObjectLiteral, KeyLiteral ]" eg.
//     "{ key1:[], key2:{objKey}, key3 }"
//       matches query
//     "[0][key1][]=value1&[0][key1][]=value2&[0][key2][objKey]=objValue&[0][key3]=keyValue"
func ParsePermitted(rules... string) (Permittable, error) {
	return buildRules(rules...)
}


// MustParsePermitted is equivalent to ParsePermitted but instead of returning an error it panics with error.
func MustParsePermitted(rules... string) Permittable {
	safe, err := ParsePermitted(rules...); if err != nil {
		panic(err)
	}
	return safe
}

func buildRules(ruleStrings ...string) (permitter Permittable, err error) {
	switch len(ruleStrings) {
	case 0:
		return nil, errors.New("")

	case 1:
		permitter, err = buildRule(ruleStrings[0])
		if err != nil {
			err = errors.Wrapf(err, "problem with rule `%s`", ruleStrings[0])
		}
		return

	default:
		arr := arrayElement{}

		for _, ruleString := range ruleStrings {
			if permitter, err = buildRule(ruleString); err != nil {
				return nil, errors.Wrapf(err, "problem with rule `%s`", ruleString)
			}

			arr = append(arr, permitter)
		}

		return &arr, nil
	}
}

var rgxValidRuleChars = regexp.MustCompile("[^\\w %,:'{}\\[\\]]")
func buildRule(ruleString string) (_ Permittable, err error) {
	builder := &permittableBuilder{}

	invalidChars := rgxValidRuleChars.FindAllString(ruleString, -1)
	if invalidChars != nil {
		invalidChars = funk.UniqString(invalidChars)
		return nil, errors.Errorf("rule `%s` contains invalid chars: %v", ruleString, invalidChars)
	}

	iterRuleString := ""
	for currRuleString := ruleString ; currRuleString != iterRuleString && currRuleString != "";  {
		iterRuleString = currRuleString
		if currRuleString, err = builder.processObjectGroups(currRuleString); err != nil { return nil, err }
		if currRuleString, err = builder.processArrayGroups(currRuleString); err != nil { return nil, err }
		currRuleString = strings.TrimSpace(currRuleString)

		elems := funk.Chain(*builder).Filter(funk.NotEmpty).Value().([]Permittable)
		count := len(elems)
		if count == 1 && rgxPermitterIdOnly.MatchString(currRuleString) {
			return elems[0], nil
		}
	}

	return nil, errors.Errorf("failed to build rule from `%s`", ruleString)
}
