package place

import (
	"regexp"
)

func Split(str string) []string {
	mc := regexp.MustCompile("[\\s|,|;]+")
	ret := []string{}
	for _, s := range mc.Split(str, -1) {
		if len(s) > 0 {
			ret = append(ret, s)
		}
	}
	return ret
}
