package validate_json

// TODO remove casts with generics

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo/v2"
)

type goJson = map[string]interface{}

var (
	definitions goJson
	defErrMsg   string
	ignored     []string
)

func EquivalentToScheme(res *bytes.Buffer, scheme []byte, resSchemeName string, exclude []string) bool {
	var (
		jsonRes       goJson
		jsonResArr    []goJson
		swaggerScheme goJson
	)
	ignored = exclude

	objErr := json.Unmarshal(res.Bytes(), &jsonRes)
	if objErr != nil {
		arrErr := json.Unmarshal(res.Bytes(), &jsonResArr)
		if arrErr != nil {
			GinkgoWriter.Println("your Json-res could not be parsed. Is it of type obj or Arr?")
			return false
		}
		jsonRes = jsonResArr[0]
		GinkgoWriter.Println("Only index0 of your file will be analyzed")
	}
	json.Unmarshal(scheme, &swaggerScheme)

	schemeLocation := []string{"definitions", resSchemeName, "properties"}
	resAttributesScheme, errMsg := findObject((swaggerScheme), schemeLocation)
	if errMsg != "" {
		prtErr(errMsg)
		return false
	}

	definitions, defErrMsg = findObject((swaggerScheme), []string{"definitions"})
	if defErrMsg != "" {
		prtErr(defErrMsg)
		return false
	}

	isEqual := loopScheme(resAttributesScheme, jsonRes, resSchemeName)
	return isEqual
}

func mapReferences(schemeV interface{}, path string) interface{} {
	ref, err := findSubValue("$ref", schemeV.(goJson), path)
	if err != "" {
		return schemeV
	}
	defName := strings.Split(ref.(string), "/")[2]

	value, referenceErr := findSubValue(defName, definitions, "definitions")
	if referenceErr != "" {
		return schemeV
	}
	return value
}

func loopScheme(node goJson, jsonObj goJson, path string) bool {
	for schemeK, schemeV := range node {
		jsonValue, errMsg := findSubValue(schemeK, jsonObj, path)
		for _, ignoredField := range ignored {
			if schemeK == ignoredField {
				return true
			}
		}
		if errMsg != "" {
			GinkgoWriter.Println(errMsg)
			return false
		}

		if jsonValue == nil {
			errMsg := "The Key " + schemeK + " in path " + path + " could not be found"
			GinkgoWriter.Println(errMsg)
			return false
		}

		if !compare(schemeV, schemeK, jsonValue, path) {
			msg := path + schemeK + " does not equal it's json Value"
			GinkgoWriter.Println(msg)
			return false
		}
	}
	return true
}

func compare(schemeV interface{}, schemeK string, jsonValue interface{}, path string) bool {
	schemeV = mapReferences(schemeV, path)
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
	case "object":
		path += "/" + schemeK
		schemeObjItems, _ := findObject(schemeV.(goJson), []string{"properties"})
		return loopScheme(schemeObjItems, jsonValue.(goJson), path)

	default:
		GinkgoWriter.Println(schemeK+" has a unexpected type: "+schemeK, " found in "+path)
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

func compareArray(schemeV goJson, jsonValue []interface{}, path string) bool {
	schemeArrItems, _ := findObject(schemeV, []string{"items", "properties"})
	isEqual := loopScheme(schemeArrItems, jsonValue[0].(goJson), path)
	return isEqual
}

func checkNumberType(t string) bool {
	return (t == "float64" || t == "float32" || t == "float" || t == "integer")
}

func findObject(node goJson, path []string) (goJson, string) {
	var err bool
	for _, attribute := range path {
		node, err = findSubObject(node, attribute)
		if err {
			var logPath string
			for _, pathStr := range path {
				logPath += "/" + pathStr
			}
			errMsg := logPath + " was not found"
			return nil, errMsg
		}
	}
	return node, ""
}

func findSubObject(node goJson, attName string) (subObject goJson, err bool) {
	for k, v := range node {
		if k == attName {
			subObject = v.(goJson)
			return
		}
	}
	return nil, true
}

func prtErr(msg string) {
	GinkgoWriter.Println("SWAGVAL: " + msg)
}
