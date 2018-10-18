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

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func (this *Persistence) ListDeviceClass(limit int, offset int) (result []model.DeviceClass, err error) {
	err = this.ordf.List(&result, limit, offset)
	return
}

func (this *Persistence) ListVendor(limit int, offset int) (result []model.Vendor, err error) {
	err = this.ordf.List(&result, limit, offset)
	return
}

func (this *Persistence) CreateVendor(element model.Vendor) (id string, err error) {
	err = this.ordf.SetIdDeep(&element)
	if err == nil {
		id = element.Id
		_, err = this.ordf.Insert(element)
	}
	return
}

func (this *Persistence) CreateProtocol(element model.Protocol) (id string, err error) {
	err = this.ordf.SetIdDeep(&element)
	if err == nil {
		id = element.Id
		_, err = this.ordf.Insert(element)
	}
	return
}

func (this *Persistence) CreateDeviceClass(element model.DeviceClass) (id string, err error) {
	err = this.ordf.SetIdDeep(&element)
	if err == nil {
		id = element.Id
		_, err = this.ordf.Insert(element)
	}
	return
}

func (this *Persistence) GetProtocolByUri(uri string) (result model.Protocol, err error) {
	results := []model.Protocol{}
	err = this.ordf.Search(&results, model.Protocol{ProtocolHandlerUrl: uri}, 1, 0)
	if err != nil {
		return
	}
	if len(results) == 0 {
		err = errors.New("no protocol with matching uri found")
		return
	}
	result = results[0]
	err = this.ordf.SelectDeep(&result)
	return
}
