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

type DeviceTypeCommand struct {
	Command    string           `json:"command"`
	Id         string           `json:"id"`
	Owner      string           `json:"owner"`
	DeviceType model.DeviceType `json:"device_type"`
}

func getDeviceTypeCommandHandler(db interfaces.Persistence) amqp_wrapper_lib.ConsumerFunc {
	return func(msg []byte) (err error) {
		log.Println(util.Config.DeviceTypeTopic, string(msg))
		command := DeviceTypeCommand{}
		err = json.Unmarshal(msg, &command)
		if err != nil {
			return
		}
		switch command.Command {
		case "PUT":
			err = publishMissingValueTypes(db, command.DeviceType, command.Owner)
			if err != nil {
				return err
			}
			return db.SetDeviceType(command.DeviceType)
		case "DELETE":
			return db.DeleteDeviceType(command.Id)
		}
		return errors.New("unable to handle permission command: " + string(msg))
	}
}

func publishMissingValueTypes(db interfaces.Persistence, deviceType model.DeviceType, owner string) (err error) {
	for _, service := range deviceType.Services {
		for _, assignment := range service.Input {
			err = publishMissingValueType(db, assignment.Type, owner)
			if err != nil {
				return err
			}
		}

		for _, assignment := range service.Output {
			err = publishMissingValueType(db, assignment.Type, owner)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func publishMissingValueType(db interfaces.Persistence, valueType model.ValueType, owner string) (err error) {
	exists, err := db.ValueTypeIdExists(valueType.Id)
	if err != nil {
		return err
	}
	if !exists {
		return PublishValueType(valueType, owner)
	}
	return nil
}

func PublishDeviceType(dt model.DeviceType, owner string) (err error) {
	return sendEvent(util.Config.DeviceTypeTopic, DeviceTypeCommand{DeviceType: dt, Id: dt.Id, Command: "PUT", Owner: owner})
}

func PublishDeviceTypeRemove(id string) (err error) {
	return sendEvent(util.Config.DeviceTypeTopic, DeviceTypeCommand{Id: id, Command: "DELETE"})
}
