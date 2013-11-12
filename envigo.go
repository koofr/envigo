package envigo

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Envigo(m interface{}, prefix string, getenv func(string) (string, bool)) (err error) {
	typ := reflect.TypeOf(m)
	val := reflect.ValueOf(m)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		err = fmt.Errorf("envigo: %s must be a Struct", m)
		return
	}

	if prefix != "" {
		prefix += "_"
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		v := val.Field(i)
		t := f.Type

		if !v.CanSet() {
			continue
		}

		fullName := prefix + strings.ToUpper(f.Name)

		if t.Kind() == reflect.Struct {
			if err = Envigo(v.Addr().Interface(), fullName, getenv); err != nil {
				return
			}
		} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
			if v.IsNil() {
				v.Set(reflect.New(t.Elem()))
			}

			if err = Envigo(v.Interface(), fullName, getenv); err != nil {
				return
			}
		} else {
			if s, ok := getenv(fullName); ok {
				switch t.Kind() {
				case reflect.String:
					v.SetString(s)

				case reflect.Bool:
					v.SetBool(strings.ToLower(s) == "true" || s == "1")

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					n, err := strconv.ParseInt(s, 10, 64)
					if err != nil || v.OverflowInt(n) {
						err = fmt.Errorf("envigo %s parse int error: %s", fullName, err)
						return err
					}
					v.SetInt(n)

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					n, err := strconv.ParseUint(s, 10, 64)
					if err != nil || v.OverflowUint(n) {
						err = fmt.Errorf("envigo %s parse uint error: %s", fullName, err)
						return err
					}
					v.SetUint(n)

				case reflect.Float32, reflect.Float64:
					n, err := strconv.ParseFloat(s, v.Type().Bits())
					if err != nil || v.OverflowFloat(n) {
						err = fmt.Errorf("envigo %s parse float error: %s", fullName, err)
						return err
					}
					v.SetFloat(n)
				}
			}
		}
	}

	return
}

func EnvironGetter() func(string) (string, bool) {
	items := os.Environ()

	env := make(map[string]string)

	for _, item := range items {
		parts := strings.SplitN(item, "=", 2)
		env[parts[0]] = parts[1]
	}

	return func(key string) (value string, ok bool) {
		value, ok = env[key]
		return
	}
}
