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

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/api"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/eventsourcing"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/persistence"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
)

//adjusted variant of /go-1.6/src/testing/testing.go
//with lvl parameter
func decorate(s string, lvl int) string {
	_, file, line, ok := runtime.Caller(lvl) // decorate + log + public function.
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteByte('\t')
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}

type AssertionsHandler struct {
	t *testing.T
}

func (this *AssertionsHandler) Log(lvl int, args ...interface{}) {
	this.t.Log("\n" + decorate(fmt.Sprintln(args...), 2+lvl))
}

func (this *AssertionsHandler) Error(lvl int, args ...interface{}) {
	this.Log(1+lvl, args)
	this.t.Fail()
}

func Assertions(t *testing.T) *AssertionsHandler {
	return &AssertionsHandler{t: t}
}

func (this *AssertionsHandler) True(assertion bool, msg string) {
	if !assertion {
		this.Error(1, msg)
	}
}

func (this *AssertionsHandler) Equal(a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		this.Error(1, "unequal: ", a, b)
	}
}

func (this *AssertionsHandler) UnEqual(a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) {
		this.Error(1, "equal: ", a, b)
	}
}

func testGetFreePort() string {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	parts := strings.Split(l.Addr().String(), ":")
	return parts[len(parts)-1]
}

var Jwtuser = jwt_http_router.JwtImpersonate("Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIzaUtabW9aUHpsMmRtQnBJdS1vSkY4ZVVUZHh4OUFIckVOcG5CcHM5SjYwIn0.eyJqdGkiOiJiOGUyNGZkNy1jNjJlLTRhNWQtOTQ4ZC1mZGI2ZWVkM2JmYzYiLCJleHAiOjE1MzA1MzIwMzIsIm5iZiI6MCwiaWF0IjoxNTMwNTI4NDMyLCJpc3MiOiJodHRwczovL2F1dGguc2VwbC5pbmZhaS5vcmcvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiZnJvbnRlbmQiLCJzdWIiOiJkZDY5ZWEwZC1mNTUzLTQzMzYtODBmMy03ZjQ1NjdmODVjN2IiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJmcm9udGVuZCIsIm5vbmNlIjoiMjJlMGVjZjgtZjhhMS00NDQ1LWFmMjctNGQ1M2JmNWQxOGI5IiwiYXV0aF90aW1lIjoxNTMwNTI4NDIzLCJzZXNzaW9uX3N0YXRlIjoiMWQ3NWE5ODQtNzM1OS00MWJlLTgxYjktNzMyZDgyNzRjMjNlIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJjcmVhdGUtcmVhbG0iLCJhZG1pbiIsImRldmVsb3BlciIsInVtYV9hdXRob3JpemF0aW9uIiwidXNlciJdfSwicmVzb3VyY2VfYWNjZXNzIjp7Im1hc3Rlci1yZWFsbSI6eyJyb2xlcyI6WyJ2aWV3LWlkZW50aXR5LXByb3ZpZGVycyIsInZpZXctcmVhbG0iLCJtYW5hZ2UtaWRlbnRpdHktcHJvdmlkZXJzIiwiaW1wZXJzb25hdGlvbiIsImNyZWF0ZS1jbGllbnQiLCJtYW5hZ2UtdXNlcnMiLCJxdWVyeS1yZWFsbXMiLCJ2aWV3LWF1dGhvcml6YXRpb24iLCJxdWVyeS1jbGllbnRzIiwicXVlcnktdXNlcnMiLCJtYW5hZ2UtZXZlbnRzIiwibWFuYWdlLXJlYWxtIiwidmlldy1ldmVudHMiLCJ2aWV3LXVzZXJzIiwidmlldy1jbGllbnRzIiwibWFuYWdlLWF1dGhvcml6YXRpb24iLCJtYW5hZ2UtY2xpZW50cyIsInF1ZXJ5LWdyb3VwcyJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwicm9sZXMiOlsidW1hX2F1dGhvcml6YXRpb24iLCJhZG1pbiIsImNyZWF0ZS1yZWFsbSIsImRldmVsb3BlciIsInVzZXIiLCJvZmZsaW5lX2FjY2VzcyJdLCJuYW1lIjoiZGYgZGZmZmYiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJzZXBsIiwiZ2l2ZW5fbmFtZSI6ImRmIiwiZmFtaWx5X25hbWUiOiJkZmZmZiIsImVtYWlsIjoic2VwbEBzZXBsLmRlIn0.eOwKV7vwRrWr8GlfCPFSq5WwR_p-_rSJURXCV1K7ClBY5jqKQkCsRL2V4YhkP1uS6ECeSxF7NNOLmElVLeFyAkvgSNOUkiuIWQpMTakNKynyRfH0SrdnPSTwK2V1s1i4VjoYdyZWXKNjeT2tUUX9eCyI5qOf_Dzcai5FhGCSUeKpV0ScUj5lKrn56aamlW9IdmbFJ4VwpQg2Y843Vc0TqpjK9n_uKwuRcQd9jkKHkbwWQ-wyJEbFWXHjQ6LnM84H0CQ2fgBqPPfpQDKjGSUNaCS-jtBcbsBAWQSICwol95BuOAqVFMucx56Wm-OyQOuoQ1jaLt2t-Uxtr-C9wKJWHQ")

var VtSearchUrl = ""

func InitTestContainer() (purge func(active bool), db *persistence.Persistence, err error) {
	config := util.ConfigStruct{}
	err = json.Unmarshal([]byte(testConfigStr), &config)
	if err != nil {
		log.Fatalf("Could not unmarshal config: %s", err)
	}
	util.Config = &config
	util.Config.ServerPort = testGetFreePort()
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	dockerrabbitmq, err := pool.Run("rabbitmq", "3-management", []string{})
	if err != nil {
		log.Fatalf("Could not start dockerrabbitmq: %s", err)
	}
	dockerelastic, err := pool.Run("elasticsearch", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start dockerelastic: %s", err)
	}
	dockerrdf, err := pool.Run("fgseitsrancher.wifa.intern.uni-leipzig.de:5000/iot-ontology", "unstable", []string{
		"DBA_PASSWORD=myDbaPassword",
		"DEFAULT_GRAPH=iot",
	})
	if err != nil {
		log.Fatalf("Could not start rdf: %s", err)
	}
	time.Sleep(5 * time.Second)

	dockersearch, err := pool.Run("fgseitsrancher.wifa.intern.uni-leipzig.de:5000/permissionsearch", "unstable", []string{
		"AMQP_URL=" + "amqp://guest:guest@" + dockerrabbitmq.Container.NetworkSettings.IPAddress + ":5672/",
		"ELASTIC_URL=" + "http://" + dockerelastic.Container.NetworkSettings.IPAddress + ":9200",
	})

	time.Sleep(5 * time.Second)
	dockervtsearch, err := pool.Run("fgseitsrancher.wifa.intern.uni-leipzig.de:5000/valuetypesearch", "test", []string{
		"AMQP_URL=" + "amqp://guest:guest@" + dockerrabbitmq.Container.NetworkSettings.IPAddress + ":5672/",
		"ELASTIC_URL=" + "http://" + dockerelastic.Container.NetworkSettings.IPAddress + ":9200",
	})
	if err != nil {
		log.Fatalf("Could not start search: %s", err)
	}
	purge = func(active bool) {
		eventsourcing.AmqpConn.Close()
		if active {
			pool.Purge(dockersearch)
			pool.Purge(dockervtsearch)
			pool.Purge(dockerelastic)
			pool.Purge(dockerrdf)
			pool.Purge(dockerrabbitmq)
		}
	}

	time.Sleep(10 * time.Second)

	util.Config.SparqlEndpoint = "http://localhost:" + dockerrdf.GetPort("8890/tcp") + "/sparql"
	util.Config.AmqpUrl = "amqp://guest:guest@localhost:" + dockerrabbitmq.GetPort("5672/tcp") + "/"
	util.Config.PermissionsUrl = "http://localhost:" + dockersearch.GetPort("8080/tcp")
	VtSearchUrl = "http://localhost:" + dockervtsearch.GetPort("8080/tcp")
	fmt.Println(util.Config)
	db = persistence.New()
	eventsourcing.InitEventHandling(db)
	go api.Init(db)
	time.Sleep(2 * time.Second)
	return
}

var testConfigStr = `{
    "ServerPort" : "8080",
    "LogLevel" : "CALL",
    "SparqlEndpoint": "http://iot-ontology:8890/sparql",
    "RdfGraph": "iot",
    "RdfUser": "dba",
    "RdfPW": "myDbaPassword",
    "GeneratVendor": "iot#24e5bb75-6d18-4e4e-87eb-ea4e554a14fb",
    "GeneratDeviceClass": "iot#042fffe4-ecb6-42a0-8d9b-09c68804baca",
    "ForceUser": "true",
    "ForceAuth": "true",

    "AmqpUrl": "amqp://guest:guest@rabbitmq:5672/",
    "AmqpReconnectTimeout ": 10,
    "AmqpConsumerName": "iotrepo",
    "DeviceInstanceTopic": "deviceinstance",
    "DeviceTypeTopic": "devicetype",
	"ValueTypeTopic": "valuetype",
    "GatewayTopic": "gateway",
    "DeviceInstanceDtFieldSearchName": "devicetype",
    "DeviceInstanceUrlFieldSearchName": "uri",
    "DeviceTypeServiceFieldSearchName": "service",
    "DeviceTypeMaintenanceFieldSearchName": "maintenance",

    "PermissionsUrl": "http://permissionsearch:8080",
    "DefaultPermissionsUser": "336508f3-e2ef-4aff-9627-e844a4c2de51",

    "FlushOnStartup": "false"
}`

func test_sliceContainsInterface(slice interface{}, element interface{}) (result bool) {
	arrV := reflect.ValueOf(slice)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {

			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if reflect.DeepEqual(arrV.Index(i).Interface(), element) {
				return true
			}
		}
	}
	return false
}
