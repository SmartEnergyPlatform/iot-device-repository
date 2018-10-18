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
	"regexp"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
)

func (this *Persistence) SearchValueType(query string, limit int, offset int) (valueTypes []model.ValueType, err error) {
	if query == "" {
		return this.GetValueTypeList(limit, offset)
	}
	query = regexp.QuoteMeta(query)
	err = this.ordf.SearchText(&valueTypes, model.ValueType{Name: query}, limit, offset)
	return
}

func (this *Persistence) SearchProtocol(query string, limit int, offset int) (protocols []model.Protocol, err error) {
	query = regexp.QuoteMeta(query)
	if query != "" {
		err = this.ordf.SearchText(&protocols, model.Protocol{Name: query}, limit, offset)
	} else {
		err = this.ordf.List(&protocols, limit, offset)
	}
	for index, protocol := range protocols {
		if err == nil {
			newProtocol := model.Protocol{Id: protocol.Id}
			err = this.ordf.SelectLevel(&newProtocol, -1)
			protocols[index] = newProtocol
		}
	}
	return
}

func (this *Persistence) SearchText(resultList interface{}, queryStruct interface{}, limit int, offset int) (err error) {
	err = this.ordf.SearchText(resultList, queryStruct, limit, offset)
	return
}
