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
	"strings"
)

func ResourceEqual(a, b any, keys []string) bool {
	aFI := flatten(a)
	bFI := flatten(b)
	aVals := make(map[string]FlatteItem, 0)
	bVals := make(map[string]FlatteItem, 0)
	for i, key := range keys {
		keys[i] = strings.ToLower(key)
	}
	for _, key := range keys {
		for k, item := range aFI {
			if regexp.MustCompile(key).MatchString(k) {
				aVals[k] = item
			}
		}
		for k, item := range bFI {
			if regexp.MustCompile(key).MatchString(k) {
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

func resourceEqual1(aFI, bFI map[string]FlatteItem, keys []string) bool {
	aVals := make(map[string]FlatteItem, 0)
	bVals := make(map[string]FlatteItem, 0)
	for i, key := range keys {
		keys[i] = strings.ToLower(key)
	}
	for _, key := range keys {
		for k, item := range aFI {
			if regexp.MustCompile(key).MatchString(k) {
				aVals[k] = item
			}
		}
		for k, item := range bFI {
			if regexp.MustCompile(key).MatchString(k) {
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
			for i := 0; i < objValue.Len(); i++ {
				mp[fmt.Sprintf("%d", i)] = objValue.Index(i).Interface()
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
		case reflect.Map, reflect.Func, reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.Invalid:
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
