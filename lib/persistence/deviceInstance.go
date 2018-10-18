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
	"strings"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func (this *Persistence) GetDeviceInstanceById(id string) (deviceInstance model.DeviceInstance, err error) {
	deviceInstance.Id = id
	err = this.ordf.SelectLevel(&deviceInstance, -1)
	return
}

func (this *Persistence) DeviceInstanceIdExists(id string) (exists bool) {
	exists, err := this.ordf.IdExists(id)
	if err != nil {
		log.Println(err)
	}
	return
}

func (this *Persistence) SetDeviceInstance(deviceInstance model.DeviceInstance) (err error) {
	old, err := this.GetDeviceInstanceById(deviceInstance.Id)
	if old.Url == "" {
		log.Println("DEBUG: create DeviceInstance ", deviceInstance.Id)
		_, err = this.ordf.Insert(deviceInstance)
		if err != nil {
			return
		}
		err = this.UpdateDeviceEndpoints(deviceInstance)
		return
	}
	log.Println("DEBUG: update DeviceInstance ", old.Id)
	if err != nil {
		return err
	}
	if old.Gateway != "" && (old.Url != deviceInstance.Url || tagRemovedOrChanged(old.Tags, deviceInstance.Tags)) {
		//reset gateway hash
		gw, err := this.GetGateway(old.Gateway)
		if err != nil {
			return err
		}
		gw.Hash = ""
		err = eventsourcing.PublishGateway(gw, "")
		if err != nil {
			return err
		}
	}
	deviceInstance.Gateway = old.Gateway
	_, err = this.ordf.Update(old, deviceInstance)
	if err != nil {
		return
	}
	if old.Url != deviceInstance.Url {
		err = this.UpdateDeviceEndpoints(deviceInstance)
	}
	return
}

func tagRemovedOrChanged(oldTags []string, newTags []string) bool {
	oldTagsIndex := IndexTags(oldTags)
	newTagsIndex := IndexTags(newTags)
	for key, oldVal := range oldTagsIndex {
		newVal, ok := newTagsIndex[key]
		if !ok || newVal != oldVal {
			return true
		}
	}
	return false
}

func IndexTags(tags []string) (result map[string]string) {
	result = map[string]string{}
	for _, tag := range tags {
		parts := strings.SplitN(tag, ":", 2)
		if len(parts) != 2 {
			log.Println("ERROR: wrong tag syntax; ", tag)
			continue
		}
		result[parts[0]] = parts[1]
	}
	return result
}

func (this *Persistence) DeleteDeviceInstance(id string) (err error) {
	instance, err := this.GetDeviceInstanceById(id)
	if err != nil {
		return err
	}
	if instance.Gateway != "" {
		gw, err := this.GetGateway(instance.Gateway)
		if err != nil {
			return err
		}
		if gw.Name == "" {
			return err
		}
		gw.Devices = []model.DeviceInstance{}
		gw.Hash = ""
		err = eventsourcing.PublishGateway(gw, "")
		if gw.Name == "" {
			return err
		}
	}
	err = this.DeleteEndpoints(instance.Id)
	if err == nil {
		_, err = this.ordf.Delete(instance)
	}
	return
}

func (this *Persistence) instanceUsesProtocol(instance model.DeviceInstance, protocol string) (usesProtocol bool, err error) {
	type ProtocolProtocol struct {
		Id                 string `rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Protocol"`
		ProtocolHandlerUrl string `rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#protocol_handler_url"`
	}

	type ProtocolService struct {
		Id       string           `rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Service"`
		Protocol ProtocolProtocol `rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasProtocol"`
	}

	type ProtocolDeviceType struct {
		Id       string            `rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#DeviceType"`
		Services []ProtocolService `rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasService"`
	}

	deviceType := ProtocolDeviceType{Id: instance.DeviceType}
	err = this.ordf.SelectLevel(&deviceType, -1)

	if err != nil {
		return false, err
	}
	for _, service := range deviceType.Services {
		if service.Protocol.ProtocolHandlerUrl == protocol {
			return true, err
		}
	}
	return false, err
}

func (this *Persistence) GetDeviceServiceEntity(deviceid string) (result model.DeviceServiceEntity, err error) {
	instance, err := this.GetDeviceInstanceById(deviceid)
	if err != nil {
		return result, err
	}
	dt := model.ShortDeviceType{Id: instance.DeviceType}
	err = this.ordf.SelectLevel(&dt, -1)
	result = model.DeviceServiceEntity{Device: instance, Services: dt.Services}
	return
}
