package util

import (
	"io/ioutil"

	"github.com/jtogrul/comake/types"

	"gopkg.in/yaml.v2"
)

// ReadBuildConfig parses build config from the buildfile
func ReadBuildConfig(buildfile string) (*types.Build, error) {
	configData, err := ioutil.ReadFile(buildfile)
	if err != nil {
		return nil, err
	}

	config := types.Build{}

	err = yaml.Unmarshal([]byte(configData), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
