package validate_json

// TODO remove casts with generics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	. "github.com/onsi/ginkgo/v2"
)

type goJson = map[string]interface{}

func EquivalentToScheme(res *bytes.Buffer, scheme []byte, resSchemeName string) bool {
	var (
		jsonRes       goJson
		swaggerScheme goJson
	)

	json.Unmarshal(res.Bytes(), &jsonRes)
	json.Unmarshal(scheme, &swaggerScheme)

	schemeLocation := []string{"definitions", resSchemeName, "properties"}
	resAttributesScheme, err := findObject((swaggerScheme), schemeLocation)
	if err != nil {
		return false
	}

	isEqual := loopScheme(resAttributesScheme, jsonRes, resSchemeName)
	return isEqual
}

func loopScheme(node goJson, jsonObj goJson, path string) bool {

	for schemeK, schemeV := range node {
		jsonValue, errMsg := findSubValue(schemeK, jsonObj, path)
		if errMsg != "" {
			GinkgoWriter.Println(errMsg)
			return false
		}

		if jsonValue == nil {
			errMsg := "The Key " + schemeK + " in path " + path + " could not be found"
			fmt.Println(errMsg)
			return false
		}

		if !compare(schemeV, schemeK, jsonValue, path) {
			return false
		}
	}
	return true
}

func compare(schemeV interface{}, schemeK string, jsonValue interface{}, path string) bool {
	schemeDatatype, _ := findSubValue("type", schemeV.(goJson), path)
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
		path += "/" + schemeK
		return compareArray(schemeV.(goJson), jsonValue.([]interface{}), path)

	default:
		GinkgoWriter.Println(schemeK+" is a unexpected type ", " found in " + path)
		return false
	}
}

func findSubValue(key string, node goJson, path string) (interface{}, string) {
	for k, v := range node {
		if k == key {
			return v, ""
		}
	}
	return nil, ("json Value " + key + " could not be found in " + path)
}

// TODO path gets passed from so far up can we do anything about this?
func compareArray(schemeV goJson, jsonValue []interface{}, path string) bool {
	schemeArrItems, _ := findObject(schemeV, []string{"items", "properties"})
	isEqual := loopScheme(schemeArrItems, jsonValue[0].(goJson), path)
	return isEqual
}

func checkNumberType(t string) bool {
	return (t == "float64" || t == "float32" || t == "float" || t == "integer")
}

func findObject(node goJson, path []string) (goJson, error) {
	var err error
	for _, attribute := range path {
		node, err = findSubObject(node, attribute)
		if err != nil {
			return nil, err
		}
	}
	return node, nil
}

func findSubObject(node goJson, attName string) (subObject goJson, err error) {
	for k, v := range node {
		if k == attName {
			subObject = v.(goJson)
			return
		}
	}
	return nil, errors.New("subObject could not be found")
}
