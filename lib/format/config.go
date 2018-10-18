/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package format

import (
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"github.com/cbroglie/mustache"
)

func UseDeviceConfig(config []model.ConfigField, str string) string {
	tmpl, err := mustache.ParseString(str)
	if err != nil {
		return str
	}
	confMap := map[string]string{}
	for _, elem := range config {
		confMap[elem.Name] = elem.Value
	}
	result, err := tmpl.Render(confMap)
	if err != nil {
		return str
	}
	return result
}
