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

package api

import (
	"encoding/json"
	"net/http"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/SmartEnergyPlatform/util/http/response"
)

func intern(router *jwt_http_router.Router, db interfaces.Persistence) {
	router.POST("/intern/gatewaynames", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		ids := []string{}
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		result := map[string]string{}
		for _, id := range ids {
			name, err := db.GetGatewayName(id)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			result[id] = name
		}
		response.To(res).Json(result)
	})

	router.POST("/intern/device/gatewaynames", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		ids := []string{}
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		result := map[string]string{}
		for _, id := range ids {
			name, err := db.GetGatewayNameByDevice(id)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			result[id] = name
		}
		response.To(res).Json(result)
	})
}
