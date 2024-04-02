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
				// 跳过私有属性
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
			result[strings.TrimLeft(strings.ToLower(prefix), ".")] = FlatteItem{
				Name: strings.TrimLeft(strings.ToLower(prefix), "."),
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
