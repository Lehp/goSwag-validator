package validate_json

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

func EquivalentToScheme(res *bytes.Buffer, scheme []byte, resSchemeName string) bool {
	var (
		jsonRes       map[string]interface{}
		swaggerScheme map[string]interface{}
	)

	json.Unmarshal(res.Bytes(), &jsonRes)
	json.Unmarshal(scheme, &swaggerScheme)

	var swaggerSchemeVal map[string]interface{} = *&swaggerScheme
	resAttributesScheme, err := findObject((swaggerSchemeVal), []string{"definitions", resSchemeName, "properties"})

	if err != nil {
		return false
	}

	fmt.Println(jsonRes)
	fmt.Println(resAttributesScheme)

	isEqual := loopScheme(resAttributesScheme, jsonRes)

	if !isEqual {
		return isEqual
	}
	return isEqual
}

func loopScheme(node map[string]interface{}, jsonObj map[string]interface{}) (bool) {
	for schemeK, schemeV := range node {
		jsonValue, err := findSubValue(schemeK, jsonObj)
		if (err != nil) {
			fmt.Println(err.Error())
			return false
		}

		if (jsonValue == nil) {
			errMsg := "The Key " + schemeK + " is not defined"
			fmt.Println(errMsg)
			return false
		}

		if !compare(schemeV, jsonValue) {
			fmt.Println(node)
			return false
		}
	}
	return true
}

func compare(schemeV interface{}, jsonValue interface{}) bool {

	schemeDatatype, _ := findSubValue("type", schemeV.(map[string]interface{}))
	jsonType := reflect.TypeOf(jsonValue).Name()

	switch schemeDatatype {
		case "string":
			return jsonType == "string"

		case "integer":
			return checkNumberType(jsonType)

		case "number":
			return checkNumberType(jsonType)

		case "boolean":
			return jsonType == "bool"

		case "array":
			return compareArray(schemeV.(map[string]interface{}), jsonValue.([]interface{}))
			
		default:
			fmt.Println("is of a type I don't know how to handle", " found in")
			return false
	}
}

func compareArray(schemeV map[string]interface{}, jsonValue []interface{}) bool {
	schemeArrItems, _ := findObject(schemeV, []string{"items", "properties"})
	isEqual := loopScheme(schemeArrItems, jsonValue[0].(map[string]interface{}))
	return isEqual
}

func checkNumberType(t string) bool {
	return (t == "float64" || t == "float32" || t == "float" || t == "integer")
}

func findObject(node map[string]interface{}, path []string) (map[string]interface{}, error) {
	var err error
	for _, attribute := range path {
		node, err = findSubObject(node, attribute)
		if err != nil {
			return nil, err
		}
	}
	return node, nil
}

func findSubObject(node map[string]interface{}, attName string) (subObject map[string]interface{}, err error) {
	for k, v := range node {
		if k == attName {
			subObject = v.(map[string]interface{})
			return
		}
	}
	return nil, errors.New("subObject could not be found")
}

func findSubValue(key string, node map[string]interface{}) (value interface{}, err error) {
	for k, v := range node {
		if k == key {
			return v, nil
		}
	}
	return nil, errors.New("subValue could not be found" + key)
}
