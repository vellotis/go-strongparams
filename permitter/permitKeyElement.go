package permitter

type permitKeyElement string

func (this *permitKeyElement) isPermitted(rgxResultTail [][]string) bool {
	return rgxResultTail != nil && len(rgxResultTail) == 0
}
