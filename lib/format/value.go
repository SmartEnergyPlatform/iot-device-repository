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

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

const (
	PLAIN_ID = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#PlainText"
	JSON_ID  = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#json"
	XML_ID   = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#xml"
)

func GetFormatedValue(config []model.ConfigField, format string, value InputOutput, info []model.AdditionalFormatInfo) (result string, err error) {
	switch format {
	case JSON_ID:
		result, err = FormatToJson(config, value)
	case PLAIN_ID:
		result, err = FormatToPlainText(config, value)
	case XML_ID:
		result, err = FormatToXml(config, value, info)
	default:
		err = errors.New("unsupported format: " + format)
	}
	return
}

//content of valueType.Fields may be changed
func ParseFormat(valueType model.ValueType, format string, value string, info []model.AdditionalFormatInfo) (result InputOutput, err error) {
	switch format {
	case JSON_ID:
		result, err = ParseFromJson(valueType, value)
	case PLAIN_ID:
		result, err = ParseFromPlainText(valueType, value)
	case XML_ID:
		result, err = ParseFromXml(valueType, value, info)
	default:
		err = errors.New("unsupported format: " + format)
	}
	if err != nil {
		err = UseLiterals(&result, valueType)
	}
	return
}

func literalFieldFilter(fields []model.FieldType) (result []model.FieldType) {
	for _, field := range fields {
		vt, isLiteral := literalFilter(field.Type)
		if isLiteral {
			result = append(result, model.FieldType{Type: vt, Name: field.Name, Id: field.Id})
		}
	}
	return
}

func literalFilter(valueType model.ValueType) (result model.ValueType, isLiteral bool) {
	result.Id = valueType.Id
	result.Description = valueType.Description
	result.BaseType = valueType.BaseType
	result.Name = valueType.Name
	result.Literal = valueType.Literal
	result.Fields = literalFieldFilter(valueType.Fields)
	isLiteral = result.Literal != "" || len(result.Fields) > 0
	return
}

func UseLiterals(value *InputOutput, valueType model.ValueType) (err error) {
	filteredValueType, isLiteral := literalFilter(valueType)
	if !isLiteral {
		return
	}
	useLiteralsRecursive(value, filteredValueType)
	return
}

func useLiteralsRecursive(value *InputOutput, valueType model.ValueType) {
	if value.Value == "" && valueType.Literal != "" {
		value.Value = valueType.Literal
	}
	for _, field := range valueType.Fields {
		found := false
		for _, io := range value.Values {
			if io.FieldId == field.Id {
				found = true
				useLiteralsRecursive(&io, field.Type)
			}
		}
		if !found && !model.GetAllowedValuesBase().IsCollection(field.Type) {
			newValue := InputOutput{FieldId: field.Id, Name: field.Name, Type: Type{Name: field.Type.Name, Base: field.Type.BaseType, Id: field.Type.Id, Desc: field.Type.Description}}
			useLiteralsRecursive(&newValue, field.Type)
			value.Values = append(value.Values, newValue)
		}
	}
	return
}
