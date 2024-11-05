package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

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

func PrintStructValues(s interface{}) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		log.Printf("%s: %v\n", fieldType.Name, field.Interface())
	}
}

func ReadHTMLFile(fileName string) (string, error) {
	path := fmt.Sprintf("frontend/%s", fileName)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
