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
	"strconv"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"github.com/SmartEnergyPlatform/util/http/response"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/permission"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

func devicetype(router *jwt_http_router.Router, db interfaces.Persistence) {

	router.GET("/service/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		ids, err := permission.SelectFieldAll(jwt, util.Config.DeviceTypeTopic, util.Config.DeviceTypeServiceFieldSearchName, id, model.READ)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		if len(ids) > 1 {
			response.To(res).DefaultError("found multiple devicetypes with service", http.StatusInternalServerError)
			return
		}
		if len(ids) < 1 {
			response.To(res).DefaultError("found no accessible devicetypes with service", http.StatusUnauthorized)
			return
		}
		service, err := db.GetServiceById(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(service)
	})

	router.GET("/deviceTypes/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		ids, err := permission.List(jwt, util.Config.DeviceTypeTopic, model.READ, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		result := []model.DeviceType{}
		for _, id := range ids {
			dt, err := db.GetDeviceTypeById(id.Id, 1)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			if dt.Name != "" {
				result = append(result, dt)
			}
		}
		response.To(res).Json(result)
	})

	router.GET("/maintenance/deviceTypes/:maintenance", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		maintenance := ps.ByName("maintenance")
		ids, err := permission.SelectFieldAll(jwt, util.Config.DeviceTypeTopic, util.Config.DeviceTypeMaintenanceFieldSearchName, maintenance, model.WRITE)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		result := []model.DeviceType{}
		for _, id := range ids {
			dt, err := db.GetDeviceTypeById(id.Id, 1)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			if dt.Name != "" {
				result = append(result, dt)
			}
		}
		response.To(res).Json(result)
	})

	router.GET("/maintenance/deviceTypes/:maintenance/:limit/:offset/:sortby/:direction", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		maintenance := ps.ByName("maintenance")
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		sortby := ps.ByName("sortby")
		direction := ps.ByName("direction")
		ids, err := permission.SelectField(jwt, util.Config.DeviceTypeTopic, util.Config.DeviceTypeMaintenanceFieldSearchName, maintenance, model.WRITE, limit, offset, sortby, direction)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		result := []model.DeviceType{}
		for _, id := range ids {
			dt, err := db.GetDeviceTypeById(id.Id, 1)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			if dt.Name != "" {
				result = append(result, dt)
			}
		}
		response.To(res).Json(result)
	})

	router.GET("/maintenance/check/deviceTypes/:maintenance", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		maintenance := ps.ByName("maintenance")
		result, err := db.CheckDeviceTypeMaintenance(maintenance)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(result)
	})

	router.GET("/deviceType/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		err := permission.Check(jwt, util.Config.DeviceTypeTopic, id, model.READ)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		deviceType, err := db.GetDeepDeviceTypeById(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(deviceType)
	})

	router.GET("/deviceType/:id/:depth", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		depth, err := strconv.Atoi(ps.ByName("depth"))
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		id := ps.ByName("id")
		err = permission.Check(jwt, util.Config.DeviceTypeTopic, id, model.READ)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		deviceType, err := db.GetDeviceTypeById(id, depth)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(deviceType)
	})

	router.POST("/deviceType", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var deviceType model.DeviceType
		err := json.NewDecoder(r.Body).Decode(&deviceType)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		valid, validErr := deviceType.IsValid()
		if !valid {
			response.To(res).DefaultError("invalid deviceType: "+validErr, http.StatusBadRequest)
			return
		}

		ok, inconsistencies := db.DeviceTypeIsConsistent(deviceType)
		if !ok {
			response.To(res).Error(response.ErrorMessage{StatusCode: http.StatusBadRequest, Message: "inconsistencies found", ErrorCode: response.ERROR_INCONSISTENT_NEW_ELEMENT, Detail: []string{inconsistencies}})
			return
		}
		err = db.SetId(&deviceType)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		err = eventsourcing.PublishDeviceType(deviceType, jwt.UserId)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(deviceType)
	})

	router.POST("/import/deviceType", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		if jwt.UserId != "" && !contains(jwt.RealmAccess.Roles, "admin") {
			response.To(res).DefaultError("only for admins", http.StatusUnauthorized)
			return
		}
		var deviceType model.DeviceType
		err := json.NewDecoder(r.Body).Decode(&deviceType)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		if deviceType.Id == "" {
			response.To(res).DefaultError("missing id", http.StatusBadRequest)
			return
		}
		err = eventsourcing.PublishDeviceType(deviceType, jwt.UserId)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(deviceType)
	})

	router.POST("/deviceType/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		err := permission.Check(jwt, util.Config.DeviceTypeTopic, id, model.WRITE)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}

		var deviceType model.DeviceType
		err = json.NewDecoder(r.Body).Decode(&deviceType)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		valid, validErr := deviceType.IsValid()
		if !valid {
			response.To(res).DefaultError("invalid deviceType: "+validErr, http.StatusBadRequest)
			return
		}

		if !db.DeviceTypeIdExists(id) {
			response.To(res).DefaultError("unknown deviceType id", http.StatusBadRequest)
			return
		}

		ok, inconsistencies := db.DeviceTypeIsConsistent(deviceType)
		if !ok {
			response.To(res).Error(response.ErrorMessage{StatusCode: http.StatusBadRequest, Message: "inconsistencies found", ErrorCode: response.ERROR_INCONSISTENT_NEW_ELEMENT, Detail: []string{inconsistencies}})
			return
		}

		deviceType.Id = id
		err = db.SetId(&deviceType)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		err = eventsourcing.PublishDeviceType(deviceType, "")
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Json(deviceType)
	})

	router.DELETE("/deviceType/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		err := permission.Check(jwt, util.Config.DeviceTypeTopic, id, model.ADMINISTRATE)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		deviceInstancesIds, err := db.GetAllDeviceInstanceUsingDeviceTypes(id)
		if deviceInstancesIds != nil && len(deviceInstancesIds) > 0 {
			response.To(res).Error(response.ErrorMessage{StatusCode: http.StatusBadRequest, Message: "dependent device instances", ErrorCode: response.ERROR_DEPENDENT_DEVICE_INSTANCE, Detail: deviceInstancesIds})
			return
		}

		err = eventsourcing.PublishDeviceTypeRemove(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Text("ok")
	})

	router.GET("/ui/deviceType/allowedvalues", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		response.To(res).Json(db.GetAllowedValues())
	})

	router.POST("/query/deviceType", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var deviceType model.DeviceType
		err := json.NewDecoder(r.Body).Decode(&deviceType)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		exists, id, err := db.DeviceTypeQuery(deviceType)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Json(struct {
			Exists bool
			Id     string
		}{Exists: exists, Id: id})
	})

	router.POST("/query/service", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var service model.Service
		err := json.NewDecoder(r.Body).Decode(&service)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		ids, err := db.QueryServiceDeviceType(service)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		result := []string{}
		for _, id := range ids {
			err = permission.Check(jwt, util.Config.DeviceTypeTopic, id, model.READ)
			if err == nil {
				result = append(result, id)
			}
		}
		response.To(res).Json(result)
	})
}
