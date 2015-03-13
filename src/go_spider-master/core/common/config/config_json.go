package config

import (
	"io/ioutil"
	"path/filepath"
	"github.com/bitly/simplejson"
)

func GetSpecConfig(filename string, keys ...string) (string, error) {
	jsonText, err := getJsonFile(filename)
	if err != nil {
		panic(err)
	}
	jsonBody, err := getJsonBody(jsonText)
	if err != nil {
		panic(err)
	}

	return getByKeys(jsonBody, keys)
}




// ###
func getJsonFile(filename string) ([]byte, error) {
	absFilePath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(absFilePath)
}
// ###
func getJsonBody(jsonText []byte) (*simplejson.Json, error) {
	return simplejson.NewJson(jsonText)
}
// ###
func getByKeys(jsonBody *simplejson.Json, keys []string) (string, error)  {
	for _, v := range keys {
		jsonBody = jsonBody.Get(string(v))
	}
	return jsonBody.String()
}
// ###

