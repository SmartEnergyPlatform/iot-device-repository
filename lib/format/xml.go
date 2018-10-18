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
	"encoding/xml"
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

const (
	XmlAttrFlag = "attr"
	XmlAnonym   = "anonym"
)

func FormatToXml(config []model.ConfigField, value InputOutput, addidionalInfo []model.AdditionalFormatInfo) (result string, err error) {
	xmlInfo := XmlInfo{Value: value, Config: config, AdditionalInfo: addidionalInfo}.Init()
	buffer, err := xml.MarshalIndent(xmlInfo, "", "    ")
	result = string(buffer)
	return
}

type XmlInfo struct {
	Value          InputOutput
	AdditionalInfo []model.AdditionalFormatInfo
	fieldFlags     map[string]map[string]string //map[field.id][flagName] == flag_info (e.g. attr:name=foo) || "" (e.g. attr)
	Config         []model.ConfigField
}

func (this XmlInfo) Init() XmlInfo {
	this.fieldFlags = map[string]map[string]string{}
	for _, info := range this.AdditionalInfo {
		fieldId := info.Field.Id
		this.fieldFlags[fieldId] = map[string]string{}
		for _, flag := range strings.Split(info.FormatFlag, ",") {
			flagParts := strings.Split(flag, ":")
			flagName := flagParts[0]
			this.fieldFlags[fieldId][flagName] = strings.Join(flagParts[1:], ":")
		}
	}
	return this
}

func (info XmlInfo) getParts() (attr []xml.Attr, children []XmlInfo, err error) {
	attr = []xml.Attr{}
	for _, element := range info.Value.Values {
		_, isAttr := getFieldInfo(element.FieldId, info.fieldFlags)[XmlAttrFlag]
		if isAttr {
			attr = append(attr, xml.Attr{Name: xml.Name{Local: element.Name}, Value: UseDeviceConfig(info.Config, element.Value)})
		} else {
			children = append(children, XmlInfo{AdditionalInfo: info.AdditionalInfo, Config: info.Config, fieldFlags: info.fieldFlags, Value: element})
		}
	}
	return
}

//TODO chardata
func (this XmlInfo) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	_, anonym := getFieldInfo(this.Value.FieldId, this.fieldFlags)[XmlAnonym]
	if anonym {
		e.EncodeToken(xml.CharData(this.Value.Value))
	} else {
		start.Name = xml.Name{Local: this.Value.Name, Space: ""}
		attr, childElements, err := this.getParts()
		if err != nil {
			return err
		}
		start.Attr = attr
		e.EncodeToken(start)
		for _, child := range childElements {
			e.EncodeElement(child, xml.StartElement{Name: xml.Name{Local: child.Value.Name, Space: ""}})
		}
		if this.Value.Value != "" {
			e.EncodeToken(xml.CharData(UseDeviceConfig(this.Config, this.Value.Value)))
		}
		e.EncodeToken(xml.EndElement{Name: start.Name})
	}
	return nil
}

func ParseFromXml(valueType model.ValueType, value string, infos []model.AdditionalFormatInfo) (result InputOutput, err error) {
	decoder := xml.NewDecoder(strings.NewReader(value))
	preparedIO, err := prepareInputOutput(decoder)
	if err != nil {
		return result, err
	}
	xmlInfo := XmlInfo{AdditionalInfo: infos}.Init()
	return finishIo(model.FieldType{Type: valueType, Name: preparedIO[0].Name}, preparedIO, xmlInfo.fieldFlags)
}

func prepareInputOutput(decoder *xml.Decoder) (result []InputOutput, err error) {
	var token xml.Token
	for token, err = decoder.Token(); err == nil; token, err = decoder.Token() {
		switch element := token.(type) {
		case xml.StartElement:
			child := InputOutput{Name: element.Name.Local}
			child.Values, err = prepareInputOutput(decoder)
			for _, attr := range element.Attr {
				child.Values = append(child.Values, InputOutput{Name: attr.Name.Local, Value: attr.Value})
			}
			result = append(result, child)
		case xml.CharData:
			result = append(result, InputOutput{Value: string(element)})
		case xml.EndElement:
			return
		default:
			return result, errors.New("unknown xml token type: " + reflect.TypeOf(element).Name())
		}
	}
	if err != nil && err != io.EOF {
		return
	}
	return result, nil
}

func valueTypeToMsgType(valueType model.ValueType) (result Type) {
	result.Name = valueType.Name
	result.Desc = valueType.Description
	result.Id = valueType.Id
	result.Base = valueType.BaseType
	return
}

func getFieldInfo(fieldId string, info map[string]map[string]string) (result map[string]string) {
	result, ok := info[fieldId]
	if !ok {
		result = map[string]string{}
	}
	return
}

func finishIo(field model.FieldType, prepared []InputOutput, info map[string]map[string]string) (result InputOutput, err error) {
	result.Name = field.Name
	result.FieldId = field.Id
	result.Type = valueTypeToMsgType(field.Type)

	allowedValues := model.GetAllowedValuesBase()

	thisValue, _, found := getFirstMatchingInputOutput(field.Name, prepared)
	if !found {
		return result, errors.New("cant find field :" + field.Name)
	}
	childValues := thisValue.Values

	switch {
	case allowedValues.IsPrimitive(field.Type):
		textParts := getAllMatchingInputOutput("", thisValue.Values)
		for _, part := range textParts {
			result.Value += part.Value
		}
		result.Value = strings.TrimSpace(result.Value)
	case allowedValues.IsMap(field.Type):
		//map cannot be anonym
		for _, element := range childValues {
			if element.Name != "" {
				mapField := field.Type.Fields[0]
				mapField.Name = element.Name
				subValue, err := finishIo(mapField, childValues, info)
				if err != nil {
					return result, err
				}
				result.Values = append(result.Values, subValue)
			}
		}
	default:
		for _, subField := range field.Type.Fields {
			_, anonym := getFieldInfo(subField.Id, info)[XmlAnonym]
			preparedIoSet := childValues
			if anonym {
				preparedIoSet = prepared
			}
			for _, preparedField := range getAllMatchingInputOutput(subField.Name, preparedIoSet) {
				subValue, err := finishIo(subField, []InputOutput{preparedField}, info)
				if err != nil {
					return result, err
				}
				result.Values = append(result.Values, subValue)
			}
		}
	}
	return
}

func getFirstMatchingInputOutput(name string, inputOutputs []InputOutput) (result InputOutput, index int, found bool) {
	found = false
	for index, result = range inputOutputs {
		if result.Name == name {
			found = true
			break
		}
	}
	return
}

func getAllMatchingInputOutput(name string, inputOutputs []InputOutput) (result []InputOutput) {
	for _, element := range inputOutputs {
		if element.Name == name {
			result = append(result, element)
		}
	}
	return
}
