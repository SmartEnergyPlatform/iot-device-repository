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
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/persistence/ordf"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
)

type Persistence struct {
	ordf ordf.Persistence
}

func New() *Persistence {
	return &Persistence{
		ordf: ordf.Persistence{
			Endpoint:  util.Config.SparqlEndpoint,
			Graph:     util.Config.RdfGraph,
			User:      util.Config.RdfUser,
			Pw:        util.Config.RdfPW,
			SparqlLog: util.Config.SparqlLog,
		},
	}
}

func (this *Persistence) SetId(element interface{}) error {
	return this.ordf.SetIdDeep(element)
}

func (this *Persistence) GetOrdf() ordf.Persistence {
	return this.ordf
}
