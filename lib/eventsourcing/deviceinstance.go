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

type DeviceinstanceCommand struct {
	Command        string               `json:"command"`
	Id             string               `json:"id"`
	Owner          string               `json:"owner"`
	DeviceInstance model.DeviceInstance `json:"device_instance"`
}

func getDeviceInstanceCommandHandler(db interfaces.Persistence) amqp_wrapper_lib.ConsumerFunc {
	return func(msg []byte) (err error) {
		log.Println(util.Config.DeviceInstanceTopic, string(msg))
		command := DeviceinstanceCommand{}
		err = json.Unmarshal(msg, &command)
		if err != nil {
			return
		}
		switch command.Command {
		case "PUT":
			return db.SetDeviceInstance(command.DeviceInstance)
		case "DELETE":
			return db.DeleteDeviceInstance(command.Id)
		}
		return errors.New("unable to handle permission command: " + string(msg))
	}
}

func PublishDeviceInstance(instance model.DeviceInstance, creator string) (err error) {
	return sendEvent(util.Config.DeviceInstanceTopic, DeviceinstanceCommand{DeviceInstance: instance, Id: instance.Id, Command: "PUT", Owner: creator})
}

func PublishDeviceInstanceRemove(id string) (err error) {
	return sendEvent(util.Config.DeviceInstanceTopic, DeviceinstanceCommand{Id: id, Command: "DELETE"})
}
