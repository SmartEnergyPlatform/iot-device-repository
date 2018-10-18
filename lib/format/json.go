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
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func FormatToJson(config []model.ConfigField, value InputOutput) (result string, err error) {
	resultStruct, err := FormatToJsonStruct(config, value)
	if err != nil {
		return result, err
	}
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(false) //no encoding of < > etc
	encoder.SetIndent("", "    ")
	err = encoder.Encode(resultStruct)
	result = b.String()
	return
}

func FormatToJsonStruct(config []model.ConfigField, value InputOutput) (result interface{}, err error) {
	if !model.GetAllowedValuesBase().IsPrimitive(model.ValueType{BaseType: value.Type.Base}) {
		if model.GetAllowedValuesBase().IsSet(model.ValueType{BaseType: value.Type.Base}) {
			list := []interface{}{}
			for _, val := range value.Values {
				element, err := FormatToJsonStruct(config, val)
				if err != nil {
					return result, err
				}
				list = append(list, element)
			}
			return list, err
		} else {
			m := map[string]interface{}{}
			for _, val := range value.Values {
				m[val.Name], err = FormatToJsonStruct(config, val)
				if err != nil {
					return result, err
				}
			}
			return m, err
		}
	} else {
		effectiveValue := UseDeviceConfig(config, value.Value)
		switch value.Type.Base {
		case model.XsdBool:
			return strings.TrimSpace(effectiveValue) == "true", err
		case model.XsdInt:
			f, err := strconv.ParseFloat(effectiveValue, 64)
			return int64(f), err
		case model.XsdFloat:
			return strconv.ParseFloat(effectiveValue, 64)
		case model.XsdString:
			return effectiveValue, err
		default:
			temp, _ := json.MarshalIndent(value, "", "     ")
			log.Println(string(temp))
			return result, errors.New("for value incompatible basetype :" + value.Type.Base)
		}
	}
	return
}

func ParseFromJson(valueType model.ValueType, value string) (result InputOutput, err error) {
	var valueInterface interface{}
	err = json.Unmarshal([]byte(value), &valueInterface)
	if err != nil {
		return
	}
	result, err = ParseFromJsonInterface(valueType, valueInterface)
	if err != nil {
		log.Println(value)
		log.Println(valueInterface)
	}
	return
}

func ParseFromJsonInterface(valueType model.ValueType, valueInterface interface{}) (result InputOutput, err error) {
	result.Type.Name = valueType.Name
	result.Type.Desc = valueType.Description
	result.Type.Id = valueType.Id
	result.Type.Base = valueType.BaseType
	switch value := valueInterface.(type) {
	case map[string]interface{}:
		if result.Type.Base != model.MapBaseType && result.Type.Base != model.StructBaseType && result.Type.Base != model.IndexStructBaseType {
			log.Println("WARNING: used basetype is not consistent to map", result.Type)
		}
		for key, val := range value {
			childField, err := getChildFiled(valueType, key)
			if err != nil {
				continue
				//return result, err
			}
			child, err := ParseFromJsonInterface(childField.Type, val)
			if err != nil {
				return result, err
			}
			child.Name = key
			child.FieldId = childField.Id
			result.Values = append(result.Values, child)
		}
	case []interface{}:
		if result.Type.Base != model.ListBaseType {
			log.Println("WARNING: used basetype is not consistent to list", result.Type)
		}
		for _, val := range value {
			childField, err := getChildFiled(valueType, "")
			if err != nil {
				continue
				//return result, err
			}
			child, err := ParseFromJsonInterface(childField.Type, val)
			if err != nil {
				return result, err
			}
			child.Name = childField.Name
			child.FieldId = childField.Id
			result.Values = append(result.Values, child)
		}
	case bool:
		if result.Type.Base != model.XsdBool {
			log.Println("WARNING: used basetype is not consistent to boolean", result.Type)
		}
		if value {
			result.Value = "true"
		} else {
			result.Value = "false"
		}
	case string:
		if result.Type.Base != model.XsdString {
			log.Println("WARNING: used basetype is not consistent to string", result.Type)
		}
		result.Value = value
	case float64:
		if result.Type.Base != model.XsdInt && result.Type.Base != model.XsdFloat {
			log.Println("WARNING: used basetype is not consistent to number", result.Type)
		}
		if result.Type.Base == model.XsdInt {
			result.Value = strconv.FormatInt(int64(value), 10)
		}
		if result.Type.Base == model.XsdFloat {
			result.Value = strconv.FormatFloat(value, 'f', -1, 64)
		}
	case nil:
		return
	default:
		err = errors.New("error in ParseFromJsonInterface(): unknown interface type <<" + reflect.TypeOf(valueInterface).Name() + ">>")
	}
	return
}

func getChildFiled(valueType model.ValueType, childName string) (childField model.FieldType, err error) {
	allowed := model.GetAllowedValuesBase()
	switch {
	case allowed.IsCollection(valueType):
		return valueType.Fields[0], err
	case allowed.IsStructure(valueType):
		for _, field := range valueType.Fields {
			if field.Name == childName {
				return field, err
			}
		}
	default:
		err = errors.New("error on getChildFiled(): cant find child type")
	}
	return
}
