package helpers

import "reflect"

func IsStructEmpty(s interface{}) bool {
	v := reflect.ValueOf(s)

	// Pastikan input adalah struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}

	// Cek setiap field apakah memiliki nilai default
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			return false
		}
	}
	return true
}
