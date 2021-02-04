package permitter

import "regexp"

type objElement map[permitKeyElement]permittable

var (
	rgxQueryPathGroup = regexp.MustCompile("(?:(?P<ArrIdx1>\\d+)|(?P<Key>[^\\[\\]]+))|(?P<Array>\\[(?:(?P<ArrIdx2>\\d+)|(?P<Object>[^\\[\\]]+))?])")
	rgxIdxArrIdx1 = rgxQueryPathGroup.SubexpIndex("ArrIdx1")
	rgxIdxKey = rgxQueryPathGroup.SubexpIndex("Key")
	rgxIdxArray = rgxQueryPathGroup.SubexpIndex("Array")
	rgxIdxArrIdx2 = rgxQueryPathGroup.SubexpIndex("ArrIdx2")
	rgxIdxObject = rgxQueryPathGroup.SubexpIndex("Object")
)

func (this *objElement) isPermitted(rgxResultTail [][]string) bool {
	if rgxResultTail != nil && len(rgxResultTail) != 0 {
		rgxResult := rgxResultTail[idxGroup]

		isRootElement := rgxResult[rgxIdxKey] != ""

		var key permitKeyElement
		if isRootElement {
			key = permitKeyElement(rgxResult[rgxIdxKey])
		} else {
			key = permitKeyElement(rgxResult[rgxIdxObject])
		}

		if maybeObj, ok := (*this)[key]; ok {
			if subElem, ok := maybeObj.(permittable); ok {
				return subElem.isPermitted(rgxResultTail[1:])
			}
		}
	}

	return false
}

func (this *objElement) IsPermitted(path string) bool {
	return isPermitted(this, path)
}

