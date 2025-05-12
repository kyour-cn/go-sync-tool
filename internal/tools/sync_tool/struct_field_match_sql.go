package sync_tool

import (
	"fmt"
	"reflect"
	"strings"
)

func StructFieldMatchSQL(obj interface{}, sqlColumns []string) (string, error) {
	resultType := reflect.TypeOf(obj)
	var missingFields, redundantFields []string
	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		tag := field.Tag.Get("json")
		found := false
		for _, col := range sqlColumns {
			if col == tag {
				found = true
				break
			}
		}
		if !found {
			missingFields = append(missingFields, tag)
		}
	}

	for _, col := range sqlColumns {
		found := false
		tag := ""
		for i := 0; i < resultType.NumField(); i++ {
			field := resultType.Field(i)
			tag = field.Tag.Get("json")
			if col == tag {
				found = true
				break
			}
		}
		if !found {
			redundantFields = append(redundantFields, col)
		}
	}
	errorArr := make([]string, 0, 2)
	if len(missingFields) > 0 {
		errorArr = append(errorArr, fmt.Sprintf("\nSQL语句返回缺失字段:\n%s ", strings.Join(missingFields, ",")))
	}
	if len(redundantFields) > 0 {
		errorArr = append(errorArr, fmt.Sprintf("\nSQL语句返回多余字段:\n%s ", strings.Join(redundantFields, ",")))
	}
	if len(errorArr) > 0 {
		return strings.Join(errorArr, ""), nil
	}
	return "", nil
}
