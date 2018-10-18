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

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func (this *Persistence) ValueTypeQuery(valueType model.ValueType) (exists bool, id string, err error) {
	if valueType.Id != "" {
		return true, valueType.Id, nil
	}
	found := []model.ValueType{}
	valueType.Id = ""
	err = this.ordf.Search(&found, valueType, 1, 0)
	if err != nil {
		return exists, id, err
	}
	exists = len(found) > 0
	if exists {
		id = found[0].Id
	}
	return
}

func (this *Persistence) CreateValueType(element model.ValueType) (err error) {
	_, err = this.ordf.Insert(element)
	return
}

func (this *Persistence) GetValueTypeList(limit int, offset int) (valueTypes []model.ValueType, err error) {
	err = this.ordf.List(&valueTypes, limit, offset)
	return
}

func (this *Persistence) GetValueTypeById(id string) (valueType model.ValueType, err error) {
	valueType.Id = id
	err = this.ordf.SelectLevel(&valueType, -1)
	return
}

func (this *Persistence) CheckValueTypeDelete(id string) (err error) {
	valuetype := model.ValueType{Fields: []model.FieldType{{Type: model.ValueType{Id: id}}}}
	valueTypes := []model.ValueType{}
	err = this.ordf.Search(&valueTypes, valuetype, 1, 0)
	if err != nil {
		return
	}
	if len(valueTypes) > 0 {
		err = errors.New("dependent valueTypes found: " + valueTypes[0].Name + ", " + valueTypes[0].Id)
		return
	}

	deviceTypes := []model.DeviceType{}
	deviceTypeInp := model.DeviceType{Services: []model.Service{{Input: []model.TypeAssignment{{Type: model.ValueType{Id: id}}}}}}
	err = this.ordf.Search(&deviceTypes, deviceTypeInp, 1, 0)
	if err != nil {
		return
	}
	if len(deviceTypes) > 0 {
		err = errors.New("dependent devicetype found: " + deviceTypes[0].Name + ", " + deviceTypes[0].Id)
		return
	}

	deviceTypes = []model.DeviceType{}
	deviceTypeOutp := model.DeviceType{Services: []model.Service{{Output: []model.TypeAssignment{{Type: model.ValueType{Id: id}}}}}}
	err = this.ordf.Search(&deviceTypes, deviceTypeOutp, 1, 0)
	if err != nil {
		return
	}
	if len(deviceTypes) > 0 {
		err = errors.New("dependent devicetype found: " + deviceTypes[0].Name + ", " + deviceTypes[0].Id)
		return
	}
	return nil
}

func (this *Persistence) DeleteValueType(id string) (err error) {
	vt, err := this.GetValueTypeById(id)
	log.Println("DEBUG: delete valuetype", id, err, vt)
	if err != nil {
		return err
	}
	_, err = this.ordf.Delete(vt)
	return
}

func (this *Persistence) ValueTypeIdExists(id string) (exists bool, err error) {
	return this.ordf.IdExists(id)
}
