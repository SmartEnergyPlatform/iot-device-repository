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
	"log"

	"github.com/SmartEnergyPlatform/amqp-wrapper-lib"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
)

var AmqpConn *amqp_wrapper_lib.Connection

func InitEventHandling(db interfaces.Persistence) (err error) {
	AmqpConn, err = amqp_wrapper_lib.Init(util.Config.AmqpUrl, []string{util.Config.DeviceInstanceTopic, util.Config.DeviceTypeTopic, util.Config.GatewayTopic, util.Config.ValueTypeTopic}, util.Config.AmqpReconnectTimeout)
	if err != nil {
		log.Fatal("ERROR: while initializing amqp connection", err)
		return
	}

	log.Println("init deviceinstance event handler")
	err = AmqpConn.Consume(util.Config.AmqpConsumerName+"_"+util.Config.DeviceInstanceTopic, util.Config.DeviceInstanceTopic, getDeviceInstanceCommandHandler(db))
	if err != nil {
		log.Fatal("ERROR: while initializing deviceinstance consumer", err)
		return
	}

	log.Println("init devicetype event handler")
	err = AmqpConn.Consume(util.Config.AmqpConsumerName+"_"+util.Config.DeviceTypeTopic, util.Config.DeviceTypeTopic, getDeviceTypeCommandHandler(db))
	if err != nil {
		log.Fatal("ERROR: while initializing devicetype consumer", err)
		return
	}

	log.Println("init gateway event handler")
	err = AmqpConn.Consume(util.Config.AmqpConsumerName+"_"+util.Config.GatewayTopic, util.Config.GatewayTopic, getGatewayCommandHandler(db))
	if err != nil {
		log.Fatal("ERROR: while initializing event consumer", err)
		return
	}

	log.Println("init valuetype event handler")
	err = AmqpConn.Consume(util.Config.AmqpConsumerName+"_"+util.Config.ValueTypeTopic, util.Config.ValueTypeTopic, getValueTypeCommandHandler(db))
	if err != nil {
		log.Fatal("ERROR: while initializing event consumer", err)
		return
	}

	return
}

func sendEvent(topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Println("ERROR: event marshaling:", err)
		return err
	}
	log.Println("DEBUG: send amqp event: ", topic, string(payload))
	return AmqpConn.Publish(topic, payload)
}
