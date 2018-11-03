package structs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func StructFromEnv(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			value := os.Getenv(strings.ToUpper(tagValue))
			if value == "" {
				return fmt.Errorf("missing env var value: %s", strings.ToUpper(tagValue))
			}
			f := v.FieldByName(fi.Name)
			switch fi.Type.Kind() {
			case reflect.String:
				f.SetString(value)
			case reflect.Int, reflect.Int32, reflect.Int64:
				i, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				f.SetInt(int64(i))
			case reflect.Bool:
				val := false
				if value == "true" {
					val = true
				}
				f.SetBool(val)
			}
		}
	}
	return nil
}

func StructFromFile(i interface{}, envVar string) error {
	fPath := os.Getenv(envVar)
	if fPath != "" {
		raw, err := ioutil.ReadFile(fPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("%v env var is not set", envVar)
	}
}

func ToMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			fKind := v.Field(i).Type().Kind()
			key := strings.ToUpper(tagValue)
			value := v.Field(i)
			switch fKind {
			case reflect.String:
				out[key] = value.String()
			case reflect.Bool:
				out[key] = value.Bool()
			case reflect.Int, reflect.Int32, reflect.Int64:
				out[key] = value.Int()
			}
		}
	}
	return out, nil
}

func Append(obj interface{}, config map[string]string) (map[string]string, error) {
	mMap, err := ToMap(obj)
	if err != nil {
		return nil, err
	}
	for key, value := range mMap {
		config[key] = fmt.Sprintf("%v", value)
	}
	return config, nil
}
