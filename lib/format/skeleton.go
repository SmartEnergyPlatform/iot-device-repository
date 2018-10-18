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
	"fmt"
	"log"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func GetFormatExample(db interfaces.Persistence, formatInfo model.TypeAssignment) (result string, err error) {
	formatInfo.Type, err = db.GetValueTypeById(formatInfo.Type.Id)
	if err != nil {
		return result, err
	}
	value, err := SkeletonFromAssignment(formatInfo, model.GetAllowedValuesBase())
	if err != nil {
		return result, err
	}
	return GetFormatedValue([]model.ConfigField{}, formatInfo.Format, value, formatInfo.AdditionalFormatinfo)
}

func GetBpmnSkeletonFromDeviceType(db interfaces.Persistence, deviceTypeId string, serviceId string) (result BpmnValueSkeleton, err error) {
	result.Inputs = map[string]interface{}{}
	result.Outputs = map[string]interface{}{}
	deviceType, err := db.GetDeepDeviceTypeById(deviceTypeId)
	if err != nil {
		return
	}
	var service model.Service
	for _, currentService := range deviceType.Services {
		if currentService.Id == serviceId {
			service = currentService
		}
	}
	if service.Id != serviceId {
		err = errors.New("service not found: " + serviceId)
	}

	allowedValues := model.GetAllowedValuesBase()

	for _, input := range service.Input {
		input.Type, _ = removeLiteral(input.Type)
		inputSkeleton, err := SkeletonFromAssignment(input, allowedValues)
		if err != nil {
			log.Println("ERROR in SkeletonFromAssignment()", err)
			return result, err
		}
		inputJson, err := FormatToJsonStruct([]model.ConfigField{}, inputSkeleton)
		if err != nil {
			return result, err
		}
		result.Inputs[input.Name] = inputJson
	}
	for _, output := range service.Output {
		outputSkeleton, err := SkeletonFromAssignment(output, allowedValues)
		if err != nil {
			return result, err
		}
		outputJson, err := FormatToJsonStruct([]model.ConfigField{}, outputSkeleton)
		if err != nil {
			return result, err
		}
		result.Outputs[output.Name] = outputJson
	}
	return
}

func GetBpmnSkeleton(db interfaces.Persistence, deviceInstanceId string, serviceId string) (result BpmnValueSkeleton, err error) {
	instance, err := db.GetDeviceInstanceById(deviceInstanceId)
	if err != nil {
		return
	}
	result, err = GetBpmnSkeletonFromDeviceType(db, instance.DeviceType, serviceId)
	return
}

func SkeletonFromAssignment(assignment model.TypeAssignment, allowedValues model.AllowedValues) (result InputOutput, err error) {
	result.Name = assignment.Name
	result.Type = typeFromValueType(assignment.Type)
	err = setSkeletonValueFromValueType(&result, assignment.Type, allowedValues)
	return
}

func setSkeletonValueFromValueType(skeleton *InputOutput, valueType model.ValueType, allowedValues model.AllowedValues) (err error) {
	switch {
	case allowedValues.IsPrimitive(valueType):
		switch valueType.BaseType {
		case model.XsdBool:
			skeleton.Value = "true"
		case model.XsdString:
			skeleton.Value = "STRING"
		case model.XsdFloat:
			skeleton.Value = "0.0"
		case model.XsdInt:
			skeleton.Value = "0"
		}
	case allowedValues.IsStructure(valueType):
		for _, field := range valueType.Fields {
			input := InputOutput{
				FieldId: field.Id,
				Name:    field.Name,
				Type:    typeFromValueType(field.Type),
			}
			err = setSkeletonValueFromValueType(&input, field.Type, allowedValues)
			if err != nil {
				return err
			}
			skeleton.Values = append(skeleton.Values, input)
		}
	case allowedValues.IsMap(valueType):
		if len(valueType.Fields) != 1 {
			return errors.New("Collection with more or less then one field")
		}
		subtype := valueType.Fields[0].Type
		input := InputOutput{
			FieldId: valueType.Fields[0].Id,
			Name:    "KEY",
			Type:    typeFromValueType(subtype),
		}
		err = setSkeletonValueFromValueType(&input, subtype, allowedValues)
		if err != nil {
			return err
		}
		skeleton.Values = append(skeleton.Values, input)
	case allowedValues.IsSet(valueType):
		if len(valueType.Fields) != 1 {
			return errors.New("Collection with more or less then one field")
		}
		subtype := valueType.Fields[0].Type
		input := InputOutput{
			Name:    valueType.Fields[0].Name,
			FieldId: valueType.Fields[0].Id,
			Type:    typeFromValueType(subtype),
		}
		err = setSkeletonValueFromValueType(&input, subtype, allowedValues)
		if err != nil {
			return err
		}
		skeleton.Values = append(skeleton.Values, input)
	default:
		fmt.Println("unknown base type: " + valueType.BaseType)
		return errors.New("unknown base type: " + valueType.BaseType)
	}
	return
}

func typeFromValueType(valueType model.ValueType) Type {
	return Type{Name: valueType.Name, Desc: valueType.Description, Id: valueType.Id, Base: valueType.BaseType}
}

func removeLiteralField(fields []model.FieldType) (result []model.FieldType) {
	for _, field := range fields {
		vt, isLiteral := removeLiteral(field.Type)
		if !isLiteral {
			result = append(result, model.FieldType{Type: vt, Name: field.Name, Id: field.Id})
		}
	}
	return
}

func removeLiteral(valueType model.ValueType) (result model.ValueType, isLiteral bool) {
	result.Id = valueType.Id
	result.Description = valueType.Description
	result.BaseType = valueType.BaseType
	result.Name = valueType.Name
	result.Literal = valueType.Literal
	result.Fields = removeLiteralField(valueType.Fields)
	isLiteral = valueType.Literal != "" || (model.GetAllowedValuesBase().IsStructure(valueType) && len(valueType.Fields) == 0)
	return
}
