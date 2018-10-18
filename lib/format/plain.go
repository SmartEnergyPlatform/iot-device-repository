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
	"errors"
	"log"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func FormatToPlainText(config []model.ConfigField, value InputOutput) (string, error) {
	if len(value.Values) > 0 {
		return "", errors.New("try to interprete complex type as plain text")
	}
	return UseDeviceConfig(config, value.Value), nil
}

func ParseFromPlainText(valueType model.ValueType, value string) (result InputOutput, err error) {
	result.Type.Name = valueType.Name
	result.Type.Desc = valueType.Description
	result.Type.Id = valueType.Id
	result.Type.Base = valueType.BaseType
	result.Value = value
	if result.Type.Base == model.MapBaseType || result.Type.Base == model.StructBaseType || result.Type.Base == model.IndexStructBaseType || result.Type.Base == model.ListBaseType {
		log.Println("WARNING: using plain text format for complex type")
	}
	return
}
