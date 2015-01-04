//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Package handling the static configurations of the homeserver
package config

import (
	"encoding/json"
	"os"
)

var Config Configuration

type Configuration struct {
	Port                     int    `json:"port"`
	Hostname                 string `json:"hostname"`
	DefaultPowerLevel        int64  `json:"default-power-level"`
	DefaultCreatorPowerLevel int64  `json:"default-creator-power-level"`
}

func load(configfile string) Configuration {
	config := Configuration{}

	file, _ := os.Open(configfile)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		panic("could not decode or find config. Copy " + configfile + ".dist to " + configfile)
	}

	return config
}

func init() {
	Config = load("config.json")
}
