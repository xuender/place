package place

import (
	"strings"
	"time"
)

func NowFormat(format string) string {
	return TimeFormat(time.Now(), format)
}

func TimeFormat(t time.Time, format string) string {
	nf := format
	for k, v := range map[string]string{
		"yyyy": "2006",
		"mm":   "01",
		"dd":   "02",
	} {
		nf = replace(nf, k, v)
	}
	// fmt.Printf("new format: %s\n", nf)
	return t.Format(nf)
}

func replace(old, k, v string) string {
	// fmt.Printf("old: %s, k: %s, v: %s\n", old, k, v)
	ns := strings.Replace(old, k, v, -1)
	if len(k) > 1 {
		ns = replace(ns, k[1:], v[1:])
	}
	return ns
}
