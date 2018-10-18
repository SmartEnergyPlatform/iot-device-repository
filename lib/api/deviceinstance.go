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

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/gen"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/permission"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

func deviceinstance(router *jwt_http_router.Router, db interfaces.Persistence) {

	router.GET("/deviceInstances/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")

		ids, err := permission.List(jwt, util.Config.DeviceInstanceTopic, model.READ, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		deviceInstances := []model.DeviceInstance{}
		for _, id := range ids {
			instance, err := db.GetDeviceInstanceById(id.Id)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			if instance.Name != "" {
				deviceInstances = append(deviceInstances, instance)
			}
		}
		response.To(res).Json(deviceInstances)
	})

	router.GET("/deviceInstances/:limit/:offset/:action", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		actionName := ps.ByName("action")
		action, err := AuthActionFromStr(actionName)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		ids, err := permission.List(jwt, util.Config.DeviceInstanceTopic, action, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		deviceInstances := []model.DeviceInstance{}
		for _, id := range ids {
			instance, err := db.GetDeviceInstanceById(id.Id)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			if instance.Name != "" {
				deviceInstances = append(deviceInstances, instance)
			}
		}

		response.To(res).Json(deviceInstances)
	})

	router.GET("/deviceInstance/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		if err := permission.Check(jwt, util.Config.DeviceInstanceTopic, id, model.READ); err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		deviceInstance, err := db.GetDeviceInstanceById(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(deviceInstance)
	})

	router.POST("/deviceInstance", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var deviceInstance model.DeviceInstance
		err := json.NewDecoder(r.Body).Decode(&deviceInstance)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		valid, validErr := deviceInstance.IsValid()
		if !valid {
			response.To(res).DefaultError("invalid deviceInstance: "+validErr, http.StatusBadRequest)
			return
		}

		ok, inconsistencies := db.DeviceInstanceIsConsistent(deviceInstance)
		if !ok {
			response.To(res).Error(response.ErrorMessage{StatusCode: http.StatusBadRequest, Message: "inconsistencies found", ErrorCode: response.ERROR_INCONSISTENT_NEW_ELEMENT, Detail: []string{inconsistencies}})
			return
		}

		err = db.SetId(&deviceInstance)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		if deviceInstance.ImgUrl == "" {
			dt, err := db.GetDeviceTypeById(deviceInstance.DeviceType, 1)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			}
			deviceInstance.ImgUrl = dt.ImgUrl
		}

		err = eventsourcing.PublishDeviceInstance(deviceInstance, jwt.UserId)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(deviceInstance)
	})

	router.POST("/deviceInstance/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		if err := permission.Check(jwt, util.Config.DeviceInstanceTopic, id, model.WRITE); err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}

		var deviceInstance model.DeviceInstance
		err := json.NewDecoder(r.Body).Decode(&deviceInstance)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		valid, validErr := deviceInstance.IsValid()
		if !valid {
			response.To(res).DefaultError("invalid deviceInstance: "+validErr, http.StatusBadRequest)
			return
		}

		if !db.DeviceInstanceIdExists(id) {
			response.To(res).DefaultError("unknown deviceInstance id", http.StatusBadRequest)
			return
		}

		deviceInstance.Id = id
		ok, inconsistencies := db.DeviceInstanceIsConsistent(deviceInstance)
		if !ok {
			response.To(res).Error(response.ErrorMessage{StatusCode: http.StatusBadRequest, Message: "inconsistencies found", ErrorCode: response.ERROR_INCONSISTENT_NEW_ELEMENT, Detail: []string{inconsistencies}})
			return
		}

		err = eventsourcing.PublishDeviceInstance(deviceInstance, "")
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Json(deviceInstance)
	})

	router.DELETE("/deviceInstance/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		if err := permission.Check(jwt, util.Config.DeviceInstanceTopic, id, model.ADMINISTRATE); err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		deviceInstance, err := db.GetDeviceInstanceById(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		err = eventsourcing.PublishDeviceInstanceRemove(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(deviceInstance)
	})

	router.GET("/ui/deviceInstance/resourceSkeleton/:deviceTypeId", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		instance := model.DeviceInstance{}
		deviceType, err := db.GetDeepDeviceTypeById(ps.ByName("deviceTypeId"))
		instance.DeviceType = deviceType.Id
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		instance.Config = createParameterSkeleton(deviceType.Config)

		response.To(res).Json(instance)
	})

	router.GET("/byservice/deviceInstances/:serviceid", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		serviceid := ps.ByName("serviceid")
		deviceTypeId, err := db.GetDeviceTypeIdByServiceId(serviceid)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		deviceInstances, err := bydevicetype(deviceTypeId, db, jwt, model.READ)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(deviceInstances)
	})

	router.GET("/bydevicetype/deviceInstances/:devicetype", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		devicetype := ps.ByName("devicetype")
		deviceInstances, err := bydevicetype(devicetype, db, jwt, model.READ)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(deviceInstances)
	})

	router.GET("/byservice/deviceInstances/:serviceid/:action", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		serviceid := ps.ByName("serviceid")
		actionName := ps.ByName("action")
		action, err := AuthActionFromStr(actionName)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		deviceTypeId, err := db.GetDeviceTypeIdByServiceId(serviceid)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		deviceInstances, err := bydevicetype(deviceTypeId, db, jwt, action)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(deviceInstances)
	})

	router.GET("/bydevicetype/deviceInstances/:devicetype/:action", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		devicetype := ps.ByName("devicetype")
		actionName := ps.ByName("action")
		action, err := AuthActionFromStr(actionName)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		deviceInstances, err := bydevicetype(devicetype, db, jwt, action)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(deviceInstances)
	})

	router.GET("/url_to_devices/:device_url", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		device_url := ps.ByName("device_url")
		entities := []model.DeviceServiceEntity{}
		ids, err := permission.SelectFieldAll(jwt, util.Config.DeviceInstanceTopic, util.Config.DeviceInstanceUrlFieldSearchName, device_url, model.READ)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		for _, id := range ids {
			entity, err := db.GetDeviceServiceEntity(id.Id)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			entities = append(entities, entity)
		}
		response.To(res).Json(entities)
	})

	router.POST("/endpoint/in", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var e model.Endpoint
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		endpoints, err := db.GetEndpoints(e.Endpoint, e.ProtocolHandler)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		filteredEndpoints := []model.Endpoint{}
		for _, endpoint := range endpoints {
			service, err := db.GetServiceById(endpoint.Service)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			if err := permission.Check(jwt, util.Config.DeviceInstanceTopic, endpoint.Device, model.EXECUTE); err == nil && service.ServiceType == "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Sensor" {
				filteredEndpoints = append(filteredEndpoints, endpoint)
			}
		}
		response.To(res).Json(filteredEndpoints)
	})

	router.POST("/endpoint/listen/auth/check/:handler/:endpoint", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		endpoint := ps.ByName("endpoint")
		handler := ps.ByName("handler")
		endpoints, err := db.GetEndpoints(endpoint, handler)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		for _, endpoint := range endpoints {
			if err := permission.Check(jwt, util.Config.DeviceInstanceTopic, endpoint.Device, model.EXECUTE); err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
				return
			}
		}
		response.To(res).Text("ok")
	})

	router.POST("/endpoint/out", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var e model.Endpoint
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		if err := permission.Check(jwt, util.Config.DeviceInstanceTopic, e.Device, model.EXECUTE); err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		endpoint, err := db.GetEndpointByDeviceAndService(e.Device, e.Service)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(endpoint)
	})

	router.POST("/endpoint/generate", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var e gen.EndpointGenMsg
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		endpoint, err := gen.CreateNewEndpoint(db, e, jwt.UserId)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json([]model.Endpoint{endpoint})
	})

	router.GET("/endpoints/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		if !contains(jwt.RealmAccess.Roles, "admin") {
			response.To(res).DefaultError("only for admins", http.StatusUnauthorized)
			return
		}
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
		endpoints, err := db.GetEndpointsList(limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(endpoints)
	})

}

func createParameterSkeleton(fieldTypes []model.ConfigFieldType) (parameters []model.ConfigField) {
	for _, fieldType := range fieldTypes {
		parameters = append(parameters, model.ConfigField{
			Name: fieldType.Name,
		})
	}
	return
}

func bydevicetype(devicetype string, db interfaces.Persistence, jwt jwt_http_router.Jwt, action model.AuthAction) (deviceInstances []model.DeviceInstance, err error) {
	ids, err := permission.SelectFieldAll(jwt, util.Config.DeviceInstanceTopic, util.Config.DeviceInstanceDtFieldSearchName, devicetype, action)
	if err != nil {
		return deviceInstances, err
	}
	for _, id := range ids {
		instance, err := db.GetDeviceInstanceById(id.Id)
		if err != nil {
			return deviceInstances, err
		}
		deviceInstances = append(deviceInstances, instance)
	}
	return
}
