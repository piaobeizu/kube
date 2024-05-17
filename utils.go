/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/04/01 22:22:34
 Desc     :
*/

package kube

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// type Equal struct {
// 	Key      string
// 	SortName string
// }

var sortKeys = []string{
	"EnvVar",
}

func ResourceEqual(a, b any, equals []string) bool {
	aFI := flatten(a)
	bFI := flatten(b)
	aVals := make(map[string]FlatteItem, 0)
	bVals := make(map[string]FlatteItem, 0)
	for i, equal := range equals {
		equals[i] = strings.ToLower(equal)
	}
	for _, equal := range equals {
		for k, item := range aFI {
			if regexp.MustCompile(equal).MatchString(k) {
				aVals[k] = item
			}
		}
		for k, item := range bFI {
			if regexp.MustCompile(equal).MatchString(k) {
				bVals[k] = item
			}
		}
	}
	if len(aVals) != len(bVals) {
		return false
	}
	for k, aval := range aVals {
		if _, ok := bVals[k]; !ok {
			return false
		}
		if aval.Kind != bVals[k].Kind {
			return false
		}
		if !reflect.DeepEqual(aval.Val, bVals[k].Val) {
			return false
		}
	}
	return true
}

func flatten(obj any) map[string]FlatteItem {
	result := make(map[string]FlatteItem)
	// mp := structToMap(obj)
	var f func(any, string)
	f = func(o any, prefix string) {
		var mp = make(map[string]any)
		objValue := reflect.ValueOf(o)
		if reflect.TypeOf(o).Kind() == reflect.Ptr {
			objValue = objValue.Elem()
		}
		kind := objValue.Kind()
		switch kind {
		case reflect.Slice, reflect.Array:
			if objValue.Len() == 0 {
				return
			}
			typeOfA := objValue.Index(0).Type()
			slice := reflect.MakeSlice(reflect.SliceOf(typeOfA), objValue.Len(), objValue.Len())
			reflect.Copy(slice, objValue)
			for _, key := range sortKeys {
				if typeOfA.Name() == key {
					sort.Slice(slice.Interface(), func(i, j int) bool {
						return reflectCmp(slice.Index(i).Interface(), slice.Index(j).Interface(), typeOfA.Field(0).Name)
					})
				}
			}
			for i := 0; i < slice.Len(); i++ {
				mp[fmt.Sprintf("%d", i)] = slice.Index(i).Interface()
			}
		case reflect.Struct:
			for i := 0; i < objValue.NumField(); i++ {
				field := objValue.Type().Field(i)
				// skip unexported field
				if field.PkgPath != "" {
					continue
				}
				fieldValue := objValue.Field(i).Interface()
				jsonName := field.Tag.Get("json")
				jsonName = strings.Replace(jsonName, ",inline", "", -1)
				jsonName = strings.Replace(jsonName, ",omitempty", "", -1)
				mp[jsonName] = fieldValue
			}
		case reflect.Interface:
			mp[prefix] = o
		case reflect.Func, reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.Invalid:
			return
		default:
			name := strings.Trim(prefix, ".")
			// remove the inline and omitempty
			name = strings.Replace(name, "..", ".", -1)
			result[strings.ToLower(name)] = FlatteItem{
				Name: name,
				Val:  o,
				Kind: kind.String(),
			}
		}
		for k, v := range mp {
			f(v, prefix+k+".")
		}
	}
	f(obj, "")
	return result
}

func FormatStruct(data any, indent bool) string {
	if data == nil {
		return ""
	}
	var (
		bytes []byte
		err   error
	)
	if !indent {
		bytes, err = json.Marshal(data)
	} else {
		bytes, err = json.MarshalIndent(data, "", "  ")
	}
	if err != nil {
		panic(err)
	}
	var content = string(bytes)
	content = strings.Replace(content, "\\u003c", "<", -1)
	content = strings.Replace(content, "\\u003e", ">", -1)
	content = strings.Replace(content, "\\u0026", "&", -1)
	content = strings.Replace(content, "\\\\", "", -1)
	return content
}

func reflectCmp(i, j interface{}, fieldName string) bool {
	valI := reflect.ValueOf(i).FieldByName(fieldName).Interface()
	valJ := reflect.ValueOf(j).FieldByName(fieldName).Interface()
	switch s := valI.(type) {
	case string:
		return s < valJ.(string)
	case float64:
		return s < valJ.(float64)
	case float32:
		return s < valJ.(float32)
	case int64:
		return s < valJ.(int64)
	case int32:
		return s < valJ.(int32)
	case int16:
		return s < valJ.(int16)
	case int8:
		return s < valJ.(int8)
	case uint64:
		return s < valJ.(uint64)
	case uint32:
		return s < valJ.(uint32)
	case uint16:
		return s < valJ.(uint16)
	case uint8:
		return s < valJ.(uint8)
	default:
		fmt.Println("The type is unknown")
	}
	return true
}
