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
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"errors"

	"github.com/cbroglie/mustache"
)

func (this *Persistence) UpdateDeviceEndpoints(device model.DeviceInstance) error {
	deviceType, err := this.GetDeviceTypeById(device.DeviceType, 3)
	if err != nil {
		return err
	}

	//delete old
	endpoints, err := this.getEndpointsByDevice(device.Id)
	if err != nil {
		return err
	}
	for _, endpoint := range endpoints {
		this.ordf.Delete(endpoint)
		if err != nil {
			return err
		}
	}

	//create new
	for _, service := range deviceType.Services {
		endpoint := model.Endpoint{
			ProtocolHandler: service.Protocol.ProtocolHandlerUrl,
			Service:         service.Id,
			Device:          device.Id,
			Endpoint:        createEndpointString(service.EndpointFormat, device.Url, service.Url, device.Config),
		}
		if endpoint.Endpoint != "" {
			tempErr := this.ordf.SetIdDeep(&endpoint)
			if tempErr != nil {
				err = tempErr
			} else {
				_, tempErr = this.ordf.Insert(endpoint)
				if tempErr != nil {
					err = tempErr
				}
			}
		}
	}
	return err
}

func (this *Persistence) DeleteEndpoints(deviceid string) (err error) {
	endpoints, err := this.getEndpointsByDevice(deviceid)
	if err != nil {
		return err
	}
	for _, endpoint := range endpoints {
		_, err = this.ordf.Delete(endpoint)
		if err != nil {
			return err
		}
	}
	return
}

func createEndpointString(format string, device string, service string, config []model.ConfigField) (result string) {
	conf := map[string]string{"device_uri": device, "service_uri": service}
	for _, field := range config {
		conf[field.Name] = field.Value
	}
	result, _ = mustache.Render(format, conf)
	return
}

func (this *Persistence) UpdateDeviceTypeEndpoints(deviceType model.DeviceType) error {
	devices, err := this.GetAllDeviceInstanceUsingDeviceTypes(deviceType.Id)
	if err != nil {
		return err
	}
	for _, deviceId := range devices {
		device, err := this.GetDeviceInstanceById(deviceId)
		if err != nil {
			return err
		}
		err = this.UpdateDeviceEndpoints(device)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Persistence) GetEndpoints(endpoint string, protocolHandler string) (result []model.Endpoint, err error) {
	err = this.ordf.SearchAll(&result, model.Endpoint{Endpoint: endpoint, ProtocolHandler: protocolHandler})
	return
}

func (this *Persistence) GetEndpointByDeviceAndService(deviceId string, serviceId string) (result model.Endpoint, err error) {
	list := []model.Endpoint{}
	err = this.ordf.Search(&list, model.Endpoint{Device: deviceId, Service: serviceId}, 1, 0)
	if err != nil {
		return result, err
	}
	if len(list) == 0 {
		return result, errors.New("no endpoint with given ids found")
	}
	return list[0], err
}

func (this *Persistence) GetEndpointsList(limit, offset int) (result []model.Endpoint, err error) {
	err = this.ordf.List(&result, limit, offset)
	return
}

func (this *Persistence) getEndpointsByDevice(deviceId string) (result []model.Endpoint, err error) {
	err = this.ordf.SearchAll(&result, model.Endpoint{Device: deviceId})
	return
}

func (this *Persistence) getEndpointsByService(serviceId string) (result []model.Endpoint, err error) {
	err = this.ordf.SearchAll(&result, model.Endpoint{Service: serviceId})
	return
}

func (this *Persistence) EndpointCollision(endpoint string) (collisions []model.Endpoint, err error) {
	err = this.ordf.SearchAll(&collisions, model.Endpoint{Endpoint: endpoint})
	return
}
