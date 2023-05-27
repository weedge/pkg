package utils

import (
	"reflect"
	"regexp"
)

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func BuildMatchRegexp(match string) (*regexp.Regexp, error) {
	var err error
	var r *regexp.Regexp

	if len(match) > 0 {
		if r, err = regexp.Compile(match); err != nil {
			return nil, err
		}
	}

	return r, nil
}
