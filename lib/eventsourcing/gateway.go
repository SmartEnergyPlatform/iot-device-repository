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

package eventsourcing

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/SmartEnergyPlatform/amqp-wrapper-lib"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
)

type GatewayCommand struct {
	Command string   `json:"command"`
	Id      string   `json:"id"`
	Owner   string   `json:"owner"`
	Name    string   `json:"name"`
	Hash    string   `json:"hash"`
	Devices []string `json:"devices"`
}

func getGatewayCommandHandler(db interfaces.Persistence) amqp_wrapper_lib.ConsumerFunc {
	return func(msg []byte) (err error) {
		log.Println(util.Config.GatewayTopic, string(msg))
		command := GatewayCommand{}
		err = json.Unmarshal(msg, &command)
		if err != nil {
			return
		}
		switch command.Command {
		case "PUT":
			return db.SetGatway(command.Id, command.Name, command.Hash, command.Devices)
		case "DELETE":
			return db.DeleteGateway(command.Id)
		}
		return errors.New("unable to handle permission command: " + string(msg))
	}
}

func PublishGateway(gw model.Gateway, owner string) (err error) {
	devices := []string{}
	for _, device := range gw.Devices {
		devices = append(devices, device.Id)
	}
	return sendEvent(util.Config.GatewayTopic, GatewayCommand{Command: "PUT", Id: gw.Id, Name: gw.Name, Hash: gw.Hash, Owner: owner, Devices: devices})
}

func PublishGatewayRef(gw model.GatewayRef, name string, owner string) (err error) {
	return sendEvent(util.Config.GatewayTopic, GatewayCommand{Command: "PUT", Id: gw.Id, Name: name, Hash: gw.Hash, Owner: owner, Devices: gw.Devices})
}

func PublishGatewayCommand(gw GatewayCommand) (err error) {
	return sendEvent(util.Config.GatewayTopic, gw)
}

func PublisGatewayRemove(id string) (err error) {
	return sendEvent(util.Config.GatewayTopic, GatewayCommand{Id: id, Command: "DELETE"})
}
