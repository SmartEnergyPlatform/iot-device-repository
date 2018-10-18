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
	"net/http"
	"strconv"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/gen"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"github.com/SmartEnergyPlatform/util/http/response"

	"github.com/SmartEnergyPlatform/jwt-http-router"

	"bytes"

	"log"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
)

func valuetype(router *jwt_http_router.Router, db interfaces.Persistence) {

	router.GET("/valueTypes/:limit/:offset", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
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

		valueTypes, err := db.GetValueTypeList(limit, offset)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		model.SortValueTypes(&valueTypes)
		response.To(res).Json(valueTypes)
	})

	router.GET("/valueType/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		id := ps.ByName("id")
		valueType, err := db.GetValueTypeById(id)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(valueType)
	})

	router.POST("/valueType/generate", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		msg := buf.String()
		vt, format, _, err := gen.ValueTypeFromMessage(db, msg)
		if err != nil {
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		response.To(res).Json(struct {
			ValueType model.ValueType `json:"value_type"`
			Format    gen.Format      `json:"format"`
		}{
			ValueType: vt,
			Format:    format,
		})
	})

	router.DELETE("/valueType/:id", func(res http.ResponseWriter, r *http.Request, ps jwt_http_router.Params, jwt jwt_http_router.Jwt) {
		if !contains(jwt.RealmAccess.Roles, "admin") {
			response.To(res).DefaultError("only for admins", http.StatusUnauthorized)
			return
		}
		id := ps.ByName("id")
		err := db.CheckValueTypeDelete(id)
		if err != nil {
			log.Println("ERROR: ", err)
			response.To(res).DefaultError(err.Error(), http.StatusBadRequest)
			return
		}
		err = eventsourcing.PublishValueTypeRemove(id)
		if err != nil {
			log.Println("ERROR:", err)
			response.To(res).DefaultError(err.Error(), http.StatusInternalServerError)
			return
		}
		response.To(res).Text("ok")
	})
}
