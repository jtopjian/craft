package utils

import (
	"fmt"
	"reflect"
	"time"
)

// BuildRequest
func BuildRequest(opts interface{}) (err error) {
	vValue := reflect.ValueOf(opts)
	if vValue.Kind() == reflect.Ptr {
		vValue = vValue.Elem()
	}

	tValue := reflect.TypeOf(opts)
	if tValue.Kind() == reflect.Ptr {
		tValue = tValue.Elem()
	}

	if vValue.Kind() == reflect.Struct {
		for i := 0; i < vValue.NumField(); i++ {
			vField := vValue.Field(i)
			tField := tValue.Field(i)

			zero := IsZero(vField)

			if requiredTag := tField.Tag.Get("required"); requiredTag == "true" {
				if zero {
					return fmt.Errorf("Missing input: %s", tField.Name)
				}
			}

			if !vField.CanSet() {
				continue
			}

			if defaultTag := tField.Tag.Get("default"); defaultTag != "" && zero {
				switch tField.Type.Name() {
				case "bool":
					switch defaultTag {
					case "true":
						vValue.Field(i).SetBool(true)
					case "false":
						vValue.Field(i).SetBool(false)
					}
				case "string":
					vValue.Field(i).SetString(defaultTag)
				}
			}
		}
	}

	return nil
}

// Shamelessly copied from github.com/gophercloud/gophercloud
func IsZero(v reflect.Value) bool {
	var t time.Time

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return false
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && IsZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(t) {
			if v.Interface().(time.Time).IsZero() {
				return true
			}
			return false
		}
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && IsZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
