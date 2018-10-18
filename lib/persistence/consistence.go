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

package persistence

import (
	"log"
	"reflect"

	"github.com/pkg/errors"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func (this *Persistence) GetAllDeviceInstanceUsingDeviceTypes(deviceType string) (deviceInstanceIds []string, err error) {
	deviceInstances := []model.DeviceInstance{}
	query := model.DeviceInstance{DeviceType: deviceType}

	err = this.ordf.SearchAll(&deviceInstances, query)

	if err != nil {
		return
	}
	for _, deviceInstance := range deviceInstances {
		deviceInstanceIds = append(deviceInstanceIds, deviceInstance.Id)
	}
	return
}

func (this *Persistence) ValueTypeIsConsistent(valueType model.ValueType) (err error) {
	if valueType.Id == "" {
		if valueType.Description == "" {
			return errors.New("missing description")
		}
		if valueType.Name == "" {
			return errors.New("missing name")
		}
		if valueType.BaseType == "" {
			return errors.New("missing base type")
		}
		if len(valueType.Fields) != 1 && (valueType.BaseType == model.ListBaseType || valueType.BaseType == model.MapBaseType) {
			return errors.New("Collection BaseType with more or less than one field")
		}
		if len(valueType.Fields) > 0 {
			if valueType.BaseType != model.StructBaseType && valueType.BaseType != model.IndexStructBaseType && valueType.BaseType != model.ListBaseType && valueType.BaseType != model.MapBaseType {
				return errors.New("structure subvalues with wrong base type")
			}
			for _, field := range valueType.Fields {
				err = this.ValueTypeIsConsistent(field.Type)
				if err != nil {
					return
				}
			}
		}
	} else {
		idExists, err := this.ordf.IdIsOfClass(valueType)
		if err != nil {
			log.Println(err)
			return errors.New("error on value type id check")
		}
		if !idExists {
			return errors.New("unknown valuetype id is used")
		}
	}
	return nil
}

func (this *Persistence) DeviceTypeIsConsistent(deviceType model.DeviceType) (ok bool, inconsistencies string) {
	if deviceType.Name == "" {
		return false, "missing name"
	}
	if deviceType.Description == "" {
		return false, "missing description"
	}
	if deviceType.Vendor.Name == "" && deviceType.Vendor.Id == "" {
		return false, "missing vendor"
	}
	if deviceType.DeviceClass.Name == "" && deviceType.DeviceClass.Id == "" {
		return false, "missing device class"
	}
	for _, field := range deviceType.Config {
		if field.Name == "" {
			return false, "missing config field name"
		}
	}
	for _, service := range deviceType.Services {
		ok, inconsistencies = service.IsValid()
		if !ok {
			return
		}

		serviceType := model.SmartObject{Id: service.ServiceType}
		isServiceTypeId, err := this.ordf.IdIsOfClass(serviceType)
		if err != nil {
			log.Println(err)
			return false, "error on service type id check"
		}
		if !isServiceTypeId {
			return false, "unknown service type"
		}
		for _, assignment := range service.Input {
			ok, inconsistencies = this.TypeAssignmentIsConsistent(assignment)
			if !ok {
				return false, inconsistencies
			}
		}
		for _, assignment := range service.Output {
			ok, inconsistencies = this.TypeAssignmentIsConsistent(assignment)
			if !ok {
				return false, inconsistencies
			}
		}
	}
	return true, ""
}

func (this *Persistence) TypeAssignmentIsConsistent(assignment model.TypeAssignment) (ok bool, inconsistencies string) {
	if assignment.Name == "" {
		return false, "missing name for assignment"
	}
	format := model.Format{Id: assignment.Format}
	isFormatId, err := this.ordf.IdIsOfClass(format)
	if err != nil {
		log.Println(err)
		return false, "error on format id check"
	}
	if !isFormatId {
		return false, "unknown format id"
	}

	err = this.ValueTypeIsConsistent(assignment.Type)
	if err != nil {
		return false, err.Error()
	}

	isMsgSegmentId, err := this.ordf.IdIsOfClass(assignment.MsgSegment)
	if err != nil {
		log.Println(err)
		return false, "error on msgSegment id check"
	}
	if !isMsgSegmentId {
		return false, "unknown msgSegment id"
	}

	return true, ""
}

func (this *Persistence) DeviceInstanceIsConsistent(deviceInstance model.DeviceInstance) (ok bool, inconsistencies string) {
	if deviceInstance.Url == "" {
		return false, "missing url"
	}
	if deviceInstance.Name == "" {
		return false, "missing name"
	}
	deviceType := model.DeviceType{Id: deviceInstance.DeviceType}
	deviceTypeExists, err := this.ordf.IdIsOfClass(deviceType)
	if err != nil {
		log.Println(err)
		return false, "error on deviceType id check"
	}
	if !deviceTypeExists {
		return false, "unknown deviceType id"
	}
	err = this.ordf.SelectLevel(&deviceType, -1)
	if err != nil {
		log.Println(err)
		return false, "error on deviceType check"
	}

	ok, inconsistencies = checkConfig(deviceInstance.Config, deviceType.Config)
	if !ok {
		return false, "inconsistent parameter: " + inconsistencies
	}

	for _, service := range deviceType.Services {
		endpoint := createEndpointString(service.EndpointFormat, deviceInstance.Url, service.Url, deviceInstance.Config)
		if endpoint != "" {
			collisions, err := this.EndpointCollision(endpoint)
			if err != nil {
				return false, "error: " + err.Error()
			}
			for _, collision := range collisions {
				if collision.Device != deviceInstance.Id {
					log.Println("WARNING: endpoint collision with ", collision)
					return false, "error: endpoint collision"
				}
			}
		}
	}

	return true, ""
}
func checkConfig(fields []model.ConfigField, types []model.ConfigFieldType) (bool, string) {
	fieldNames := map[string]bool{}
	fieldTypeNames := map[string]bool{}
	for _, field := range types {
		fieldTypeNames[field.Name] = true
	}
	for _, field := range fields {
		fieldNames[field.Name] = true
		if field.Value == "" {
			return false, "missing config value for '" + field.Name + "'"
		}
	}
	if !reflect.DeepEqual(fieldNames, fieldTypeNames) {
		return false, "missmatch in config field names"
	}
	return true, ""
}

func (this *Persistence) GetAllowedValues() (result model.AllowedValues) {
	serviceTypes := []model.SmartObject{}
	err := this.ordf.List(&serviceTypes, 10, 0)
	if err != nil {
		log.Println(err)
	}

	formats := []model.Format{}
	err = this.ordf.List(&formats, 10, 0)
	if err != nil {
		log.Println(err)
	}

	result = model.GetAllowedValuesBase()
	result.ServiceTypes = serviceTypes
	result.Formats = formats
	return
}
