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
	"errors"
	"net/http"
	"strconv"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"github.com/SmartEnergyPlatform/util/http/response"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/permission"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

func search(router *jwt_http_router.Router, db interfaces.Persistence) {
	router.GET("/ui/search/deviceTypes/:query/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		ids, err := permission.Search(jwt, util.Config.DeviceTypeTopic, query, model.READ, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		result := []model.DeviceType{}
		for _, id := range ids {
			dt, err := db.GetDeviceTypeById(id.Id, 1)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			result = append(result, dt)
		}
		response.To(res).Json(result)
	})

	router.GET("/ui/search/deviceInstances/:query/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		ids, err := permission.Search(jwt, util.Config.DeviceInstanceTopic, query, model.READ, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		result := []model.DeviceInstance{}
		for _, id := range ids {
			instance, err := db.GetDeviceInstanceById(id.Id)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			result = append(result, instance)
		}
		response.To(res).Json(result)
	})

	router.GET("/ui/search/deviceInstances/:query/:limit/:offset/:action", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		actionName := ps.ByName("action")
		action, err := AuthActionFromStr(actionName)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		ids, err := permission.Search(jwt, util.Config.DeviceInstanceTopic, query, action, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		result := []model.DeviceInstance{}
		for _, id := range ids {
			instance, err := db.GetDeviceInstanceById(id.Id)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			result = append(result, instance)
		}
		response.To(res).Json(result)
	})

	router.GET("/ui/search/valueTypes/:query/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit, err := strconv.Atoi(ps.ByName("limit"))
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(ps.ByName("offset"))
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		query := ps.ByName("query")
		valueTypes, err := db.SearchValueType(query, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(valueTypes)
	})

	router.GET("/ui/search/others/:type/:query/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		query := ps.ByName("query")
		typeString := ps.ByName("type")

		limit, err := strconv.Atoi(ps.ByName("limit"))
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(ps.ByName("offset"))
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		err = errors.New("unknown type")
		switch typeString {
		case "vendors":
			result := []model.Vendor{}
			if query != "" {
				err = db.SearchText(&result, model.Vendor{Name: query}, limit, offset)
			} else {
				result, err = db.ListVendor(limit, offset)
			}
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
				return
			}
			response.To(res).Json(result)
		case "deviceClasses":
			result := []model.DeviceClass{}
			if query != "" {
				err = db.SearchText(&result, model.DeviceClass{Name: query}, limit, offset)
			} else {
				result, err = db.ListDeviceClass(limit, offset)
			}
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
				return
			}
			response.To(res).Json(result)
		case "protocols":
			result, err := db.SearchProtocol(query, limit, offset)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
				return
			}
			response.To(res).Json(result)
		}

	})

	router.GET("/ui/search/gateways/:query/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		query := ps.ByName("query")
		ids, err := permission.Search(jwt, util.Config.GatewayTopic, query, model.READ, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		result := []model.Gateway{}
		for _, id := range ids {
			gateway, err := db.GetGateway(id.Id)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			result = append(result, gateway)
		}
		response.To(res).Json(result)
	})
}
