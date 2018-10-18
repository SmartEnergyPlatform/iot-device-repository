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

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/api"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/persistence"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
)

func main() {
	defer fmt.Println("exit application")
	configLocation := flag.String("config", "config.json", "configuration file")
	flag.Parse()

	err := util.LoadConfig(*configLocation)
	if err != nil {
		log.Fatal("unable to load config", err)
	} else {
		log.Println("prepare connection to database")
		db := persistence.New()

		log.Println("init eventsourcing")
		eventsourcing.InitEventHandling(db)

		api.Init(db)
	}
}
