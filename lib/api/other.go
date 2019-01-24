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

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/format"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"github.com/SmartEnergyPlatform/util/http/response"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/permission"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

type Insert_OK struct {
	CreatedId string `json:"created_id"`
}

func other(router *jwt_http_router.Router, db interfaces.Persistence) {

	router.POST("/other/vendor", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var element model.Vendor
		err := json.NewDecoder(r.Body).Decode(&element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		if element.Name == "" {
			response.To(res).DefaultError("missing name in request", http.StatusBadRequest)
			return
		}
		id, err := db.CreateVendor(element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Json(Insert_OK{CreatedId: id})
	})

	router.DELETE("/other/vendor/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		if !contains(jwt.RealmAccess.Roles, "admin") {
			response.To(res).DefaultError("access denied", http.StatusUnauthorized)
			return
		}
		err := db.DeleteVendor(ps.ByName("id"))
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Text("ok")
	})

	router.POST("/other/protocol", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var element model.Protocol
		err := json.NewDecoder(r.Body).Decode(&element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}

		id, err := db.CreateProtocol(element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Json(Insert_OK{CreatedId: id})
	})

	router.POST("/other/deviceclass", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var element model.DeviceClass
		err := json.NewDecoder(r.Body).Decode(&element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		if element.Name == "" {
			response.To(res).DefaultError("missing name in request", http.StatusBadRequest)
			return
		}
		id, err := db.CreateDeviceClass(element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Json(Insert_OK{CreatedId: id})
	})

	router.DELETE("/other/deviceclass/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		if !contains(jwt.RealmAccess.Roles, "admin") {
			response.To(res).DefaultError("access denied", http.StatusUnauthorized)
			return
		}
		err := db.DeleteDeviceClass(ps.ByName("id"))
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}

		response.To(res).Text("ok")
	})

	router.POST("/other/valueType", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		var element model.ValueType
		err := json.NewDecoder(r.Body).Decode(&element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		err = db.ValueTypeIsConsistent(element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		err = db.SetId(&element)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		err = eventsourcing.PublishValueType(element, jwt.UserId)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Json(Insert_OK{CreatedId: element.Id})
	})

	router.GET("/skeleton/:instance_id/:service_id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		instance := ps.ByName("instance_id")
		service := ps.ByName("service_id")
		err := permission.Check(jwt, util.Config.DeviceInstanceTopic, instance, model.EXECUTE)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		skeleton, err := format.GetBpmnSkeleton(db, instance, service)
		if err == nil {
			response.To(res).Json(skeleton)
		} else {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
		}
	})

	router.GET("/skeleton/:instance_id/:service_id/output/leaves", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		instance := ps.ByName("instance_id")
		service := ps.ByName("service_id")
		err := permission.Check(jwt, util.Config.DeviceInstanceTopic, instance, model.EXECUTE)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		skeleton, err := format.GetBpmnSkeletonOutputLeaves(db, instance, service)
		if err == nil {
			response.To(res).Json(skeleton)
		} else {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
		}
	})

	router.GET("/devicetype/skeleton/:type_id/:service_id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		devicetype := ps.ByName("type_id")
		service := ps.ByName("service_id")
		skeleton, err := format.GetBpmnSkeletonFromDeviceType(db, devicetype, service)
		if err == nil {
			response.To(res).Json(skeleton)
		} else {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
		}
	})

	router.POST("/format/example", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		msg := model.TypeAssignment{}
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
		}
		result, err := format.GetFormatExample(db, msg)
		if err != nil {
			response.To(res).DefaultError(err.Error(), 500)
		} else {
			response.To(res).Text(result)
		}
	})

	// evaluates request and responds on invalid requests with plain text messages and status code 200
	router.POST("/format/preview", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		msg := model.TypeAssignment{}
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
		}
		if msg.Name == "" {
			response.To(res).Text("error: missing name")
			return
		}
		if msg.Type.Id == "" {
			response.To(res).Text("error: missing ValueType")
			return
		}
		if msg.Format == "" {
			response.To(res).Text("error: missing Format")
			return
		}
		if msg.MsgSegment.Id == "" {
			response.To(res).Text("error: missing MessageSegment")
			return
		}
		result, err := format.GetFormatExample(db, msg)
		if err != nil {
			response.To(res).DefaultError(err.Error(), 500)
		} else {
			response.To(res).Text(result)
		}
	})
}
