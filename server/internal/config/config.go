package config

import (
	"io/ioutil"

	"github.com/MasterDimmy/jsonc"
)

type configItemProject struct {
	Name     string   // "TestProject"
	Key      string   // "i can upload and check for updates with this key"
	Branches []string // "main" , "test", "dev" ...
}

type configBody struct {
	Listen   string              // ip:port to listen to
	Storage  string              //path to file storage
	Projects []configItemProject //An array that defines descriptions of the projects hosted
}

var Config configBody //allow multithreaded read

func Read() error {
	buf, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}
	err = jsonc.Unmarshal(buf, &Config)
	return err
}
