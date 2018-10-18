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

package gen

import (
	"encoding/json"
	"errors"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
)

func CreateNewEndpoint(db interfaces.Persistence, endpointMsg EndpointGenMsg, owner string) (endpoint model.Endpoint, err error) {
	protocol, err := db.GetProtocolByUri(endpointMsg.ProtocolHandler)
	if err != nil {
		return endpoint, err
	}
	dt, err := generateDeviceType(db, protocol, endpointMsg)
	if dt.Id == "" {
		err = db.SetId(&dt)
		if err != nil {
			return endpoint, err
		}
		err = eventsourcing.PublishDeviceType(dt, owner)
		if err != nil {
			return endpoint, err
		}
	}
	instance := model.DeviceInstance{
		DeviceType: dt.Id,
		Url:        endpointMsg.Endpoint,
		Tags:       []string{"generated:generated"},
		Name:       endpointMsg.Endpoint,
	}
	err = db.SetId(&instance)
	if err != nil {
		return endpoint, err
	}
	err = eventsourcing.PublishDeviceInstance(instance, owner)
	if err != nil {
		return endpoint, err
	}
	endpoint = model.Endpoint{Device: instance.Id, Endpoint: endpointMsg.Endpoint, ProtocolHandler: protocol.ProtocolHandlerUrl, Service: dt.Services[0].Id}
	return
}

func generateDeviceType(db interfaces.Persistence, protocol model.Protocol, endpointMsg EndpointGenMsg) (result model.DeviceType, err error) {
	msgDesc := map[string]interface{}{}
	outputs := []model.TypeAssignment{}
	for _, part := range endpointMsg.Parts {
		vt, format, struc, err := ValueTypeFromMessage(db, part.Msg)
		if err != nil {
			return result, err
		}
		msgDesc[part.MsgSegmentName] = struc
		if vt.Id == "" {
			vt.Description = "generated valuetype"
		}
		msgSegmentId := ""
		for _, msgSegment := range protocol.MsgStructure {
			if msgSegment.Name == part.MsgSegmentName {
				msgSegmentId = msgSegment.Id
				break
			}
		}
		if msgSegmentId == "" {
			return result, errors.New("unknown msgSegment name: " + part.MsgSegmentName)
		}
		outputs = append(outputs, model.TypeAssignment{
			Name:   part.MsgSegmentName,
			Format: string(format),
			Type:   vt,
			MsgSegment: model.MsgSegment{
				Id: msgSegmentId,
			},
		})
	}

	result.Services = []model.Service{{
		Name:           "get",
		Description:    "data send by " + protocol.Name,
		EndpointFormat: "{{device_uri}}",
		Url:            "get",
		Protocol:       model.Protocol{Id: protocol.Id},
		ServiceType:    "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Sensor",
		Output:         outputs,
	}}
	result.Vendor = model.Vendor{Id: util.Config.GeneratVendor}
	result.DeviceClass = model.DeviceClass{Id: util.Config.GeneratDeviceClass}
	result.Generated = true

	exists, id, err := db.DeviceTypeQuery(result)
	if err != nil {
		return result, err
	}
	if exists {
		result, err = db.GetDeviceTypeById(id, -1)
	} else {
		msgDescStr, _ := json.Marshal(msgDesc)
		result.Name = protocol.ProtocolHandlerUrl + "_" + endpointMsg.Endpoint
		result.Description = "generated devicetype for " + protocol.Name + "\n" + string(msgDescStr)
		result.Maintenance = []string{"rename"}
	}

	return
}

type EndpointGenMsgPart struct {
	Msg            string `json:"msg"`
	MsgSegmentName string `json:"msg_segment_name"`
}

type EndpointGenMsg struct {
	Endpoint        string               `json:"endpoint"`
	ProtocolHandler string               `json:"protocol_handler"`
	Parts           []EndpointGenMsgPart `json:"parts"`
}
