package flyapi

type STRList []string

func (list STRList) Contains(v string) bool {
	for _, element := range list {
		if element == v {
			return true
		}
	}
	return false
}
