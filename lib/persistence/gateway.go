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
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

const GATEWAY_NONE = ""

func (this *Persistence) GetGatewayNameByDevice(id string) (name string, err error) {
	ref := model.DeviceToGateway{Id: id}
	err = this.ordf.SelectLevel(&ref, -1)
	name = ref.Gateway.Name
	return
}

func (this *Persistence) GetGateway(id string) (gateway model.Gateway, err error) {
	gateway.Id = id
	err = this.ordf.SelectLevel(&gateway, -1)
	return
}

func (this *Persistence) DeleteGateway(id string) error {
	gateway, err := this.GetGateway(id)
	if err != nil {
		return err
	}
	if gateway.Name == "" {
		return nil
	}
	_, err = this.ordf.Delete(gateway)
	return err
}

func (this *Persistence) getGatewayFlat(id string) (gateway model.GatewayFlat, err error) {
	gateway.Id = id
	err = this.ordf.SelectLevel(&gateway, -1)
	return
}

func (this *Persistence) changeDeviceGateway(deviceId string, gatewayId string) (err error) {
	exists, err := this.ordf.IdExists(deviceId)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	current := model.DeviceGatewayRelation{Id: deviceId}
	err = this.ordf.Select(&current)
	if err != nil {
		return err
	}
	new := model.DeviceGatewayRelation{Id: deviceId, Gateway: gatewayId}
	_, err = this.ordf.Update(current, new)
	return err
}

func (this *Persistence) changeDeviceListGateway(deviceIds []string, gatewayId string) (err error) {
	for _, id := range deviceIds {
		tempErr := this.changeDeviceGateway(id, gatewayId)
		if tempErr != nil {
			log.Println("ERROR: changeDeviceListGateway(): ", tempErr)
			err = tempErr
		}
	}
	return err
}

func (this *Persistence) CheckClearGateway(id string) error {
	gw, err := this.GetGateway(id)
	if err != nil {
		return err
	}
	if gw.Name == "" {
		return errors.New("cannot clear not existing gateway")
	}
	return err
}

func (this *Persistence) GetGatewayName(id string) (name string, err error) {
	current := model.GatewayName{Id: id}
	err = this.ordf.Select(&current)
	return current.Name, err
}

func gatewayDeviceDiff(old []string, new []string) (add []string, remove []string) {
	compare := func(X, Y []string) []string {
		m := make(map[string]int)
		for _, y := range Y {
			m[y]++
		}
		var ret []string
		for _, x := range X {
			if m[x] > 0 {
				m[x]--
				continue
			}
			ret = append(ret, x)
		}
		return ret
	}
	return compare(new, old), compare(old, new)
}

func (this *Persistence) GatewayCheckCommit(id string, ref model.GatewayRef) (err error) {
	gateway, err := this.GetGateway(id)
	if err != nil {
		return err
	}
	if gateway.Name == "" {
		return errors.New("gateway (" + id + ") does not exist")
	}
	mayOwn, err := this.gatewayMayOwnDevices(id, ref.Devices)
	if err != nil {
		return err
	}
	if !mayOwn {
		log.Println("ERROR: gateway ("+id+") may not own devices: ", ref.Devices)
		return errors.New("gateway (" + id + ") may not own devices")
	}
	return
}

func (this *Persistence) gatewayMayOwnDevices(gatewayId string, devices []string) (result bool, err error) {
	for _, device := range devices {
		rel := model.DeviceGatewayRelation{Id: device}
		err = this.ordf.Select(&rel)
		if err != nil {
			return false, err
		}
		if rel.Gateway != gatewayId && rel.Gateway != GATEWAY_NONE {
			return false, err
		}
	}
	return true, err
}

func (this *Persistence) ProvideGateway(id string, owner string) (gateway model.Gateway, isNew bool, err error) {
	if id != "" {
		gateway, err = this.GetGateway(id)
		if err != nil {
			return gateway, isNew, err
		}
		if gateway.Name != "" {
			return gateway, false, err
		}
	}
	gateway.Name = "new-gateway-" + strconv.FormatInt(time.Now().Unix(), 36)
	err = this.ordf.SetIdDeep(&gateway)
	return gateway, true, err
}

func (this *Persistence) SetGatway(id string, name string, hash string, devices []string) (err error) {
	newGw := model.GatewayFlat{Id: id, Name: name, Hash: hash, Devices: devices}
	log.Println("DEBUG: set gateway: ", newGw)
	current := model.GatewayFlat{Id: id}
	err = this.ordf.Select(&current)
	if err != nil {
		return err
	}
	if current.Name == "" {
		_, err = this.ordf.Insert(newGw)
		if err != nil {
			return err
		}
		err = this.changeDeviceListGateway(devices, id)
	} else {
		_, err = this.ordf.Update(current, newGw)
		if err != nil {
			return err
		}
		add, remove := gatewayDeviceDiff(current.Devices, devices)
		err = this.changeDeviceListGateway(remove, GATEWAY_NONE)
		if err != nil {
			return err
		}
		err = this.changeDeviceListGateway(add, id)
	}

	return
}
