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
	"strings"
)

func ResourceEqual(a, b any, keys []string) bool {
	aFI := flatten(a)
	bFI := flatten(b)
	aVals := make(map[string]FlatteItem, 0)
	bVals := make(map[string]FlatteItem, 0)
	fmt.Printf("aFI: %v\n", Struct2Json(aFI))
	fmt.Printf("bFI: %v\n", Struct2Json(bFI))
	for i, key := range keys {
		keys[i] = strings.ToLower(key)
	}
	for _, key := range keys {
		for k, item := range aFI {
			if strings.HasPrefix(k, key) {
				aVals[k] = item
			}
		}
		for k, item := range bFI {
			if strings.HasPrefix(k, key) {
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
		if objValue.Kind() == reflect.Ptr {
			objValue = objValue.Elem()
		}
		kind := objValue.Kind()
		switch kind {
		case reflect.Map:
			mp = o.(map[string]any)
		case reflect.Slice, reflect.Array:
			for i := 0; i < objValue.Len(); i++ {
				mp[fmt.Sprintf("%d", i)] = objValue.Index(i).Interface()
			}
		case reflect.Struct:
			objType := objValue.Type()
			for i := 0; i < objValue.NumField(); i++ {
				field := objType.Field(i)
				fieldValue := objValue.Field(i).Interface()
				mp[field.Name] = fieldValue
			}
		case reflect.Interface:
			mp[prefix] = o
		case reflect.Func, reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.Invalid:
			panic("Unsupported type")
		default:
			result[strings.ToLower(prefix)] = FlatteItem{
				Name: prefix,
				Val:  o,
				Kind: reflect.TypeOf(o).String(),
			}
		}
		for k, v := range mp {
			f(v, prefix+k+".")
		}
	}
	f(obj, "")
	return result
}

func Struct2Json(data any) string {
	str, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	var content = string(str)
	content = strings.Replace(content, "\\u003c", "<", -1)
	content = strings.Replace(content, "\\u003e", ">", -1)
	content = strings.Replace(content, "\\u0026", "&", -1)
	content = strings.Replace(content, "\\\\", "", -1)
	return content
}
