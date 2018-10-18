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
	"encoding/json"
	"errors"
	"log"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func (this *Persistence) GetDeviceTypeList(limit int, offset int) (deviceTypes []model.DeviceType, err error) {
	err = this.ordf.List(&deviceTypes, limit, offset)
	return
}

func (this *Persistence) CheckDeviceTypeMaintenance(maintenance string) (result bool, err error) {
	searchresult := []model.DeviceType{}
	err = this.ordf.Search(&searchresult, model.DeviceType{Maintenance: []string{maintenance}}, 1, 0)
	result = len(searchresult) > 0
	return
}

func (this *Persistence) GetDeviceTypeById(id string, depth int) (deviceType model.DeviceType, err error) {
	deviceType.Id = id
	err = this.ordf.SelectLevel(&deviceType, depth)
	return
}

func (this *Persistence) GetDeepDeviceTypeById(id string) (deviceType model.DeviceType, err error) {
	deviceType.Id = id
	err = this.ordf.SelectLevel(&deviceType, -1)
	return
}

func (this *Persistence) DeviceTypeIdExists(id string) (exists bool) {
	exists, err := this.ordf.IdExists(id)
	if err != nil {
		log.Println(err)
	}
	return
}

func (this *Persistence) SetDeviceType(deviceType model.DeviceType) (err error) {
	old, err := this.GetDeepDeviceTypeById(deviceType.Id)
	if err != nil {
		return err
	}
	if old.Name == "" {
		_, err = this.ordf.Insert(deviceType)
		return
	}
	if old.Generated {
		temp := old
		temp.Name = deviceType.Name
		temp.Description = deviceType.Description
		temp.Maintenance = deviceType.Maintenance
		deviceType = temp
	}
	_, err = this.ordf.Update(old, deviceType)
	if err != nil {
		return
	}
	err = this.UpdateDeviceTypeEndpoints(deviceType)
	if err != nil {
		return
	}
	if old.ImgUrl != deviceType.ImgUrl {
		deviceInstances := []model.DeviceInstance{}
		query := model.DeviceInstance{DeviceType: deviceType.Id}
		err = this.ordf.SearchAll(&deviceInstances, query)
		if err != nil {
			return
		}
		for _, instance := range deviceInstances {
			if instance.ImgUrl == "" || instance.ImgUrl == old.ImgUrl {
				instance.ImgUrl = deviceType.ImgUrl
				err = eventsourcing.PublishDeviceInstance(instance, "")
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (this *Persistence) DeleteDeviceType(id string) (err error) {
	deviceType, err := this.GetDeepDeviceTypeById(id)
	if err != nil {
		return err
	}
	_, err = this.ordf.Delete(deviceType)
	return
}

func (this *Persistence) GetServiceById(id string) (service model.Service, err error) {
	service.Id = id
	err = this.ordf.SelectDeep(&service)
	return
}

func (this *Persistence) GetDeviceTypeIdByServiceId(serviceId string) (deviceTypeId string, err error) {
	deviceType := model.UltraShortDeviceType{Service: serviceId}
	resultList := []model.UltraShortDeviceType{}
	err = this.ordf.Search(&resultList, deviceType, 1, 0)
	if err == nil {
		if len(resultList) != 1 {
			err = errors.New("no devicetype with service " + serviceId + " found")
		} else {
			deviceTypeId = resultList[0].Id
		}
	}
	return
}

func (this *Persistence) DeviceTypeQuery(deviceType model.DeviceType) (exists bool, id string, err error) {
	if deviceType.Id != "" {
		return true, deviceType.Id, nil
	}
	deviceType.Services, exists, err = this.checkServices(deviceType.Services)
	if err != nil || !exists {
		log.Println("DEBUG: services", exists, err, deviceType)
		return exists, id, err
	}
	found := []model.DeviceType{}
	err = this.ordf.Search(&found, deviceType, 1, 0)
	if err != nil {
		return exists, id, err
	}
	exists = len(found) > 0
	if exists {
		id = found[0].Id
	}
	return
}

func (this *Persistence) checkServices(services []model.Service) (result []model.Service, exists bool, err error) {
	exists = true
	for _, service := range services {
		id := ""
		id, exists, err = this.checkService(service)
		if err != nil || !exists {
			s, _ := json.Marshal(service)
			log.Println("DEBUG: service", exists, err, string(s))
			return
		}
		result = append(result, model.Service{Id: id})
	}
	return
}

func (this *Persistence) checkService(service model.Service) (id string, exists bool, err error) {
	if service.Id != "" {
		return service.Id, true, nil
	}
	_, exists, err = this.checkTypeAssingments(service.Output)
	if err != nil || !exists {
		log.Println("DEBUG: outputs", exists, err, service)
		return id, exists, err
	}
	_, exists, err = this.checkTypeAssingments(service.Input)
	if err != nil || !exists {
		log.Println("DEBUG: inputs", exists, err, service)
		return id, exists, err
	}
	found := []model.Service{}
	err = this.ordf.Search(&found, service, 1, 0)
	if err != nil {
		return id, exists, err
	}
	exists = len(found) > 0
	if exists {
		id = found[0].Id
	} else {
		s, _ := json.Marshal(service)
		log.Println("DEBUG: service intern", exists, err, string(s))
	}
	return
}

func (this *Persistence) checkTypeAssingments(assignments []model.TypeAssignment) (result []model.TypeAssignment, exists bool, err error) {
	exists = true
	for _, assignment := range assignments {
		id := ""
		id, exists, err = this.checkTypeAssingment(assignment)
		if err != nil || !exists {
			log.Println("DEBUG: assignment", exists, err, assignment)
			return
		}
		result = append(result, model.TypeAssignment{Id: id})
	}
	return
}

func (this *Persistence) checkTypeAssingment(assignment model.TypeAssignment) (id string, exists bool, err error) {
	if assignment.Id != "" {
		return assignment.Id, true, nil
	}
	exists, id, err = this.ValueTypeQuery(assignment.Type)
	if err != nil || !exists {
		log.Println("DEBUG: valuetype", exists, err, assignment.Type)
		return
	}
	assignment.Type = model.ValueType{Id: id}
	found := []model.TypeAssignment{}
	err = this.ordf.Search(&found, assignment, 1, 0)
	if err != nil {
		return id, exists, err
	}
	exists = len(found) > 0
	if exists {
		id = found[0].Id
	}
	return
}

func (this *Persistence) QueryServiceDeviceType(service model.Service) (typeIds []string, err error) {
	deviceType := model.DeviceType{Services: []model.Service{service}}
	found := []model.DeviceType{}
	err = this.ordf.SearchAll(&found, deviceType)
	if err != nil {
		return typeIds, err
	}
	for _, dt := range found {
		typeIds = append(typeIds, dt.Id)
	}
	return
}
