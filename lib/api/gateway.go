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
	"github.com/SmartEnergyPlatform/util/http/response"

	"net/http"

	"encoding/json"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"log"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/permission"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

func gateway(router *jwt_http_router.Router, db interfaces.Persistence) {

	router.GET("/gateways/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		limit := ps.ByName("limit")
		offset := ps.ByName("offset")
		ids, err := permission.List(jwt, util.Config.GatewayTopic, model.READ, limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		gateways := []model.Gateway{}
		for _, id := range ids {
			gw, err := db.GetGateway(id.Id)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			if gw.Name != "" {
				gateways = append(gateways, gw)
			}
		}
		response.To(res).Json(gateways)
	})

	router.GET("/gateway/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		err := permission.Check(jwt, util.Config.GatewayTopic, id, model.READ)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		gateway, err := db.GetGateway(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(gateway)
	})

	router.DELETE("/gateway/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		err := permission.Check(jwt, util.Config.GatewayTopic, id, model.ADMINISTRATE)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		err = eventsourcing.PublisGatewayRemove(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Text("ok")
	})

	router.POST("/gateway/:id/clear", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		ready, gw, err := gatewayIsReady(jwt, db, id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		if !ready {
			log.Println(db.GetGateway(id))
			response.To(res).DefaultError("gateway not ready", http.StatusPreconditionFailed)
			return
		}
		err = permission.Check(jwt, util.Config.GatewayTopic, id, model.WRITE)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		gw, err = db.GetGateway(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		if gw.Name == "" {
			response.To(res).Text("ok")
			return
		}
		gw.Devices = []model.DeviceInstance{}
		gw.Hash = ""
		err = eventsourcing.PublishGateway(gw, "")
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Text("ok")
	})

	router.POST("/gateway/:id/name/:name", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		name := ps.ByName("name")
		err := permission.Check(jwt, util.Config.GatewayTopic, id, model.WRITE)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		gw, err := db.GetGateway(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		gw.Name = name
		err = eventsourcing.PublishGateway(gw, "")
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Text("ok")
	})

	router.GET("/gateway/:id/name", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		err := permission.Check(jwt, util.Config.GatewayTopic, id, model.READ)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		name, err := db.GetGatewayName(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Text(name)
	})

	//ignores model.GatewayRef.Id and ueses :id
	router.POST("/gateway/:id/commit", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		ready, gw, err := gatewayIsReady(jwt, db, id)
		if err != nil {
			log.Println("DEBUG: gatewayIsReady() error:", err)
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		if !ready {
			log.Println("DEBUG: gateway is not ready")
			response.To(res).DefaultError("gateway not ready", http.StatusPreconditionFailed)
			return
		}
		err = permission.Check(jwt, util.Config.GatewayTopic, id, model.WRITE)
		if err != nil {
			log.Println("DEBUG: permission check error:", err)
			response.To(res).DefaultError(err.Error(), http.StatusUnauthorized)
			return
		}
		var gateway model.GatewayRef
		err = json.NewDecoder(r.Body).Decode(&gateway)
		if err != nil {
			log.Println("DEBUG: json error:", err)
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		err = db.GatewayCheckCommit(id, gateway)
		if err != nil {
			log.Println("DEBUG: GatewayCheckCommit() error:", err)
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		gateway.Id = id
		err = eventsourcing.PublishGatewayRef(gateway, gw.Name, "")
		if err != nil {
			log.Println("DEBUG: eventsourcing.PublishGatewayRef() error:", err)
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("DEBUG: gateway commited", err)
		response.To(res).Text("ok")
	})

	router.GET("/gateway/:id/provide", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		if id != "" {
			gw, err := db.GetGateway(id)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
			if gw.Name != "" {
				allowed, err := permission.CheckBool(jwt, util.Config.GatewayTopic, id, model.EXECUTE)
				if err != nil {
					response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
					return
				}
				if !allowed {
					id = ""
				}
			}
		}
		gateway, isNew, err := db.ProvideGateway(id, jwt.UserId)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		if isNew {
			err = eventsourcing.PublishGateway(gateway, jwt.UserId)
			if err != nil {
				response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
				return
			}
		}
		response.To(res).Json(gateway)
	})
}

func gatewayIsReady(jwt jwt_http_router.Jwt, db interfaces.Persistence, id string) (exists bool, gw model.Gateway, err error) {
	gw, err = db.GetGateway(id)
	if err != nil || gw.Name == "" {
		exists = false
		log.Println("DEBUG gateway does not exist", id)
		return
	}
	exists, err = permission.Exists(jwt, util.Config.GatewayTopic, id)
	return
}
