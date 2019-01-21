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

package interfaces

import (
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

type Persistence interface {
	SetId(element interface{}) error

	//Border between repositories

	GetAllDeviceInstanceUsingDeviceTypes(deviceType string) (deviceInstanceIds []string, err error)
	DeviceInstanceIsConsistent(deviceInstance model.DeviceInstance) (ok bool, inconsistencies string)
	DeviceTypeIsConsistent(deviceType model.DeviceType) (ok bool, inconsistencies string)
	GetAllowedValues() model.AllowedValues

	//DeviceType Methods

	GetServiceById(id string) (model.Service, error)

	GetDeviceTypeIdByServiceId(serviceId string) (deviceTypeId string, err error)

	CheckDeviceTypeMaintenance(maintenance string) (result bool, err error)
	GetDeepDeviceTypeById(id string) (model.DeviceType, error)
	GetDeviceTypeById(id string, depth int) (model.DeviceType, error)
	DeviceTypeIdExists(string) bool
	SetDeviceType(dt model.DeviceType) error
	DeleteDeviceType(string) error
	DeviceTypeQuery(deviceType model.DeviceType) (exists bool, id string, err error)
	QueryServiceDeviceType(service model.Service) (typeIds []string, err error)

	//DeviceInstance Methods

	GetDeviceInstanceById(id string) (model.DeviceInstance, error)
	DeviceInstanceIdExists(string) bool
	SetDeviceInstance(dt model.DeviceInstance) error
	DeleteDeviceInstance(id string) error
	GetDeviceServiceEntity(deviceid string) (result model.DeviceServiceEntity, err error)

	//Search
	SearchValueType(query string, limit int, offset int) (valueTypes []model.ValueType, err error)
	SearchProtocol(query string, limit int, offset int) (protocols []model.Protocol, err error)
	SearchText(resultList interface{}, regexSearch interface{}, limit int, offset int) (err error)

	//ValueType Methods
	GetValueTypeList(limit int, offset int) ([]model.ValueType, error)
	GetValueTypeById(id string) (model.ValueType, error)
	ValueTypeQuery(valueType model.ValueType) (exists bool, id string, err error)
	CreateValueType(model.ValueType) (err error)
	DeleteValueType(id string) (err error)
	CheckValueTypeDelete(id string) (err error)
	ValueTypeIsConsistent(valueType model.ValueType) (err error)
	ValueTypeIdExists(id string) (exists bool, err error)

	//other structs
	CreateVendor(model.Vendor) (id string, err error)
	CreateProtocol(model.Protocol) (id string, err error)
	CreateDeviceClass(model.DeviceClass) (id string, err error)
	ListDeviceClass(limit int, offset int) ([]model.DeviceClass, error)
	ListVendor(limit int, offset int) ([]model.Vendor, error)
	DeleteVendor(id string) error
	DeleteDeviceClass(id string) error

	//gateway
	GetGateway(id string) (model.Gateway, error)
	DeleteGateway(id string) error
	GetGatewayName(id string) (name string, err error)
	GetGatewayNameByDevice(id string) (name string, err error)
	SetGatway(id string, name string, hash string, devices []string) (err error)
	ProvideGateway(id string, owner string) (gateway model.Gateway, isNew bool, err error)
	GatewayCheckCommit(id string, ref model.GatewayRef) (err error)
	CheckClearGateway(id string) error

	//endpoint
	DeleteEndpoints(deviceid string) (err error)
	UpdateDeviceEndpoints(device model.DeviceInstance) error     //delete, update, insert
	UpdateDeviceTypeEndpoints(deviceType model.DeviceType) error //delete, update, insert
	GetEndpoints(endpoint string, protocolHandler string) (result []model.Endpoint, err error)
	GetEndpointByDeviceAndService(deviceId string, serviceId string) (result model.Endpoint, err error)
	GetEndpointsList(limit, offset int) (result []model.Endpoint, err error)

	GetProtocolByUri(uri string) (result model.Protocol, err error)

	FlushDevices()
	FlushGateways()
	FlushDevicetypes()
	FlushValueTypes()
}
