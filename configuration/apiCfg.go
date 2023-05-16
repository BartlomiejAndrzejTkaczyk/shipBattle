package configuration

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type ApiCfg struct {
	ApiUrl        string `yaml:"apiUrl"`
	InitUrl       string `yaml:"initUrl"`
	BoardUrl      string `yaml:"boardUrl"`
	AuthTokenName string `yaml:"authTokenName"`
	FireUrl       string `yaml:"fireUrl"`
	OppoListUrl   string `yaml:"oppoListUrl"`
}

// todo check how to get relative path to file in easier way
func sourceFile() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	return filepath.Join(basePath, "apiCfg.yaml")
}

func GetApiConfig() (*ApiCfg, error) {
	f, err := os.Open(sourceFile())
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Print(err)
		}
	}(f)

	var apiCfg ApiCfg
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&apiCfg)
	if err != nil {
		return nil, err
	}

	return &apiCfg, nil
}
