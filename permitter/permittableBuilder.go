package permitter

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

type permittableBuilder []Permittable

var permitterIdPattern = "@(?P<PermitterID>\\d+)@"
var rgxPermitterIdOnly = regexp.MustCompile("^"+permitterIdPattern+"$")
var keyPattern = "'[%\\w ]+'|[%\\w]+"
var keysPattern = "^\\s*(?:"+permitterIdPattern+"|(?P<Key>"+keyPattern+"))\\s*(?:,|$)"
var rgxKeys = regexp.MustCompile(keysPattern)
var rgxKeysInObj = regexp.MustCompile("" +
	"\\{(?P<Object>[@%'\\w ,:]+)\\}" +
	"|" +
	"^(?P<Keys>[@%'\\w ,\\:]+)$")
var rgxMerge = regexp.MustCompile("(?P<Key>"+keyPattern+")"+"\\s*:\\s*"+permitterIdPattern+"\\s*(?:,|$)")
func (this *permittableBuilder) processObjectGroups(ruleString string) (string, error) {
	rule := ruleString

	for _, keysResult := range rgxKeysInObj.FindAllStringSubmatch(rule, -1) {
		keys := keysResult[rgxKeysInObj.SubexpIndex("Keys")] +
			keysResult[rgxKeysInObj.SubexpIndex("Object")]
		obj := objElement{}

		for _, mergeResult := range rgxMerge.FindAllStringSubmatch(keys, -1) {
			permitterId := mergeResult[rgxMerge.SubexpIndex("PermitterID")]
			key := mergeResult[rgxMerge.SubexpIndex("Key")]
			keyElement := permitKeyElement(noQuotes(key))

			if idx, err := strconv.Atoi(permitterId); err != nil {
				return "", errors.New("Silent PANIC situation! Submit a bug report. Permitter index is not a numeric " +
					"value near: " + ruleString)

			} else if !(0 <= idx || idx < len(*this)) {
				return "", errors.New("Silent PANIC situation! Submit a bug report. Permitter index is greater than " +
					"the amount of permitter references near: " + ruleString)

			} else {
				obj[keyElement] = (*this)[idx]
				(*this)[idx] = nil
			}

			keys = removeProcessedRule(keys, mergeResult[idxGroup])
		}

		for toRemove := ""; strings.TrimSpace(keys) != ""; keys = strings.TrimSpace(removeProcessedRule(keys, toRemove)) {
			keyResult := rgxKeys.FindStringSubmatch(keys); if keyResult == nil {
				return "", errors.New("invalid key declaration near: " + keys)
			}
			permitterId := keyResult[rgxKeys.SubexpIndex("PermitterID")]
			key := keyResult[rgxKeys.SubexpIndex("Key")]
			keyElement := permitKeyElement(noQuotes(key))

			if permitterId != "" {
				if idx, err := strconv.Atoi(permitterId); err != nil {
					return "", errors.New("Silent PANIC situation! Submit a bug report. Permitter index is not a numeric " +
						"value near: " + ruleString)

				} else if !(0 <= idx || idx < len(*this)) {
					return "", errors.New("Silent PANIC situation! Submit a bug report. Permitter index is greater than " +
						"the amount of permitter references near: " + ruleString)

				} else if _, ok := (*this)[idx].(*objElement); ok {
					obj[keyElement] = (*this)[idx]
					(*this)[idx] = nil
					toRemove = keyResult[0]
				}

			} else {
				obj[keyElement] = &keyElement
				toRemove = keyResult[0]
			}
		}

		*this = append(*this, &obj)
		idx := len(*this)-1
		rule = strings.ReplaceAll(rule, keysResult[idxGroup], fmt.Sprintf("@%d@", idx))
	}

	return rule, nil
}

var rgxKeysInArr = regexp.MustCompile("\\[(?P<Array>[@%'\\w ,:]*)\\]")
func (this *permittableBuilder) processArrayGroups(ruleString string) (transformedRule string, err error) {
	rule := ruleString

	for _, arrayResult := range rgxKeysInArr.FindAllStringSubmatch(rule, -1) {
		arrayData := arrayResult[rgxKeysInArr.SubexpIndex("Array")]
		arr := arrayElement{}

		for _, mergeResult := range rgxMerge.FindAllStringSubmatch(arrayData, -1) {
			if mergedData, err := this.processObjectGroups(mergeResult[0]); err != nil {
				return "", err
			} else {
				arrayData = strings.Replace(arrayData, mergeResult[0], mergedData, 1)
			}
		}

		for toRemove := ""; strings.TrimSpace(arrayData) != ""; arrayData = strings.TrimSpace(removeProcessedRule(arrayData, toRemove)) {
			keyResult := rgxKeys.FindStringSubmatch(arrayData); if keyResult == nil {
				return "", errors.New("invalid array rule declaration near: " + ruleString)
			}
			permitterId := keyResult[rgxKeys.SubexpIndex("PermitterID")]
			key := keyResult[rgxKeys.SubexpIndex("Key")]

			if permitterId != "" {
				if idx, err := strconv.Atoi(permitterId); err != nil {
					return "", errors.Errorf("Any of the rule keys in `%s` should not contain format `@\\d+@`. " +
						"Illegal value near: %s", arrayData, ruleString)

				} else if !(0 <= idx || idx < len(*this)) {
					return "", errors.New("Silent PANIC situation! Possible bug. Permitter index is not a numeric " +
						"value near: " + ruleString)

				} else {
					arr = append(arr, (*this)[idx])
					(*this)[idx] = nil
				}

			} else {
				keyElement := permitKeyElement(noQuotes(key))
				singleKeyObjElement := objElement{
					keyElement: &keyElement,
				}
				arr = append(arr, &singleKeyObjElement)
			}

			toRemove = keyResult[0]
		}

		*this = append(*this, &arr)
		idx := len(*this)-1
		rule = strings.ReplaceAll(rule, arrayResult[idxGroup], fmt.Sprintf("@%d@", idx))
	}

	return rule, nil
}
