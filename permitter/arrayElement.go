package permitter

import "github.com/pkg/errors"

type arrayElement []permittable

func (this *arrayElement) isPermitted(rgxResultTail [][]string) bool {
	if len(rgxResultTail) == 0 {
		return false
	}

	isArray := rgxResultTail[0][rgxIdxArrIdx1] != "" ||
		(rgxResultTail[0][rgxIdxArray] != "" && rgxResultTail[0][rgxIdxObject] == "")
	if !isArray {
		return false
	}

	isLastElement := len(rgxResultTail) == 1
	hasNestedRule := len(*this) > 0
	if isLastElement && !hasNestedRule {
		return true
	}

	hasMoreElementsToProcess := len(rgxResultTail) > 1
	if hasMoreElementsToProcess && hasNestedRule {
		for _, subElem := range *this {
			if subElem != nil {
				if subElem.isPermitted(rgxResultTail[1:]) {
					return true
				}
			} else {
				panic(errors.New("Something went wrong"))
			}
		}
	}

	return false
}

func (this *arrayElement) IsPermitted(path string) bool {
	return isPermitted(this, path)
}
