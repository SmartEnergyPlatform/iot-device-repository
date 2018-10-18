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

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func (this *Persistence) FlushDevicetypes() {
	deviceTypes := []model.DeviceType{}
	err := this.ordf.SearchAll(&deviceTypes, model.DeviceType{})
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	for _, dt := range deviceTypes {
		log.Println("flush devicetype", dt.Id)
		deepDt, err := this.GetDeepDeviceTypeById(dt.Id)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
		err = eventsourcing.PublishDeviceType(deepDt, "")
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
	}
}

func (this *Persistence) FlushGateways() {
	gateways := []model.GatewayRef{}
	err := this.ordf.SearchAll(&gateways, model.GatewayRef{})
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	for _, gw := range gateways {
		log.Println("flush gateway", gw.Id)
		gateway, err := this.GetGateway(gw.Id)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
		err = eventsourcing.PublishGateway(gateway, "")
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
	}
}

func (this *Persistence) FlushDevices() {
	log.Println("flush deviceinstances")
	instances := []model.DeviceInstance{}
	err := this.ordf.SearchAll(&instances, model.DeviceInstance{})
	if err != nil {
		log.Println("ERROR: in FlushDevices", err)
		return
	}
	for _, instance := range instances {
		log.Println("fulsh device-instance", instance.Id)
		di, err := this.GetDeviceInstanceById(instance.Id)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
		if di.ImgUrl == "" {
			dt, err := this.GetDeviceTypeById(di.DeviceType, 1)
			if err != nil {
				log.Println("ERROR:", err)
				return
			}
			di.ImgUrl = dt.ImgUrl
		}
		err = eventsourcing.PublishDeviceInstance(di, "")
		if err != nil {
			log.Println("ERROR while deleting userless instance", err)
		}
	}
}

func (this *Persistence) FlushValueTypes() {
	valuetypes := []model.ValueType{}
	err := this.ordf.SearchAll(&valuetypes, model.ValueType{})
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	for _, vt := range valuetypes {
		log.Println("flush valuetype", vt.Id)
		valuetype, err := this.GetValueTypeById(vt.Id)
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
		if valuetype.Id == "" || valuetype.Name == "" {
			log.Println("ERROR: unable to load valuetype ", vt)
		}
		err = eventsourcing.PublishValueType(valuetype, "")
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
	}
}
