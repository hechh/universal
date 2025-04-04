package basic

import (
	"regexp"
)

func Filter(pattern string, vals ...string) (rets []string, err error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	for _, val := range vals {
		if re.MatchString(val) {
			continue
		}
		rets = append(rets, val)
	}
	return
}

func Ifelse[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}
