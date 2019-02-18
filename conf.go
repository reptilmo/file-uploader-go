// conf.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	ListenPort  string `json:"listen_port"`
	UploadDoc string `json:"upload_doc"`
	UploadPath  string `json:"upload_path"`
}

func NewConfig(path string) (*Config, error) {
	confFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer confFile.Close()

	fileData, err := ioutil.ReadAll(confFile)
	if err != nil {
		return nil, err
	}

	var c Config
	err = json.Unmarshal([]byte(fileData), &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
