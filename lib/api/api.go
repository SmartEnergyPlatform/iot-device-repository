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
	"log"
	"net/http"
	"net/url"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"

	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"

	"errors"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

//creates all rest-endpoints, initiates the request-logger
func Init(db interfaces.Persistence) {
	if util.Config.FlushOnStartup == "true" {
		log.Println("flush database to event")
		db.FlushDevicetypes()
		db.FlushDevices()
		db.FlushGateways()
		db.FlushValueTypes()
	}
	log.Println("start server on port: ", util.Config.ServerPort)
	httpHandler := getRoutes(db)
	corsHandler := cors.New(httpHandler)
	logger := logger.New(corsHandler, util.Config.LogLevel)
	if util.Config.DecodeUrlFix == "true" {
		log.Println("Use DecodeUrlFix")
		urldecodefix := NewUrlDecodeMiddleWare(logger)
		log.Println(http.ListenAndServe(":"+util.Config.ServerPort, urldecodefix))
	} else {
		log.Println(http.ListenAndServe(":"+util.Config.ServerPort, logger))
	}
}

func getRoutes(db interfaces.Persistence) *jwt_http_router.Router {

	router := jwt_http_router.New(jwt_http_router.JwtConfig{PubRsa: util.Config.JwtPubRsa, ForceAuth: util.Config.ForceAuth == "true", ForceUser: util.Config.ForceUser == "true"})

	devicetype(router, db)
	deviceinstance(router, db)
	valuetype(router, db)
	search(router, db)
	other(router, db)
	gateway(router, db)
	intern(router, db)

	return router
}

func NewUrlDecodeMiddleWare(handler http.Handler) *UrlDecodeMiddleWare {
	return &UrlDecodeMiddleWare{handler: handler}
}

type UrlDecodeMiddleWare struct {
	handler http.Handler
}

func (this *UrlDecodeMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this.handler != nil {
		r.URL.Path, _ = url.QueryUnescape(r.URL.Path)
		r.URL.RawPath, _ = url.QueryUnescape(r.URL.RawPath)
		this.handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", 403)
	}
}

func AuthActionFromStr(str string) (action model.AuthAction, err error) {
	switch str {
	case "read":
		action = model.READ
	case "write":
		action = model.WRITE
	case "execute":
		action = model.EXECUTE
	default:
		err = errors.New("unknown resource action")
	}
	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
