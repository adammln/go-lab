//validators.go
package main

import (
	"errors"
	"fmt"
	"reflect"
)

var validFields = map[string]interface{}{
	"ID":        "",
	"ParentID":  "",
	"Content":   "",
	"RankOrder": 0,
	"IsChecked": true,
	"Subtasks":  []string{},
}

func FieldValidator(fieldName string, value interface{}) error {
	if typeInterface, isExists := validFields[fieldName]; !isExists {
		return errors.New(
			fmt.Sprintf(
				`Invalid Field Name: got invalid field name "%s"`,
				fieldName,
			),
		)
	} else {
		if fieldName == "Subtasks" && value == nil {
			return nil
		}

		if reflect.TypeOf(typeInterface) != reflect.TypeOf(value) {
			return errors.New(
				fmt.Sprintf(
					`Type Mismatch: got data type from field "%s"=%s, want %s`,
					fieldName,
					reflect.TypeOf(value),
					reflect.TypeOf(typeInterface),
				),
			)
		}
		return nil
	}
}
