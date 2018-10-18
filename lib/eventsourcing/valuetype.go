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

type ValueTypeCommand struct {
	Command   string          `json:"command"`
	Id        string          `json:"id"`
	Owner     string          `json:"owner"`
	ValueType model.ValueType `json:"value_type"`
}

func getValueTypeCommandHandler(db interfaces.Persistence) amqp_wrapper_lib.ConsumerFunc {
	return func(msg []byte) (err error) {
		log.Println(util.Config.ValueTypeTopic, string(msg))
		command := ValueTypeCommand{}
		err = json.Unmarshal(msg, &command)
		if err != nil {
			return
		}
		switch command.Command {
		case "PUT":
			return recursiveValueTypeCreation(db, command.ValueType, command.Owner)
		case "DELETE":
			return db.DeleteValueType(command.Id)
		}
		return errors.New("unable to handle permission command: " + string(msg))
	}
}

func recursiveValueTypeCreation(db interfaces.Persistence, vt model.ValueType, owner string) (err error) {
	for _, field := range vt.Fields {
		err = PublishValueType(field.Type, owner)
		if err != nil {
			return err
		}
	}
	return db.CreateValueType(vt)
}

func PublishValueType(vt model.ValueType, owner string) (err error) {
	if vt.Id == "" {
		log.Println("WARNING: missing id in valuetype --> no publish")
		return nil
	}
	return sendEvent(util.Config.ValueTypeTopic, ValueTypeCommand{ValueType: vt, Id: vt.Id, Command: "PUT", Owner: owner})
}

func PublishValueTypeRemove(id string) (err error) {
	return sendEvent(util.Config.ValueTypeTopic, ValueTypeCommand{Id: id, Command: "DELETE"})
}
