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
	"time"

	"testing"

	"net/url"

	"net/http"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/gen"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
)

func TestEndpointGenerated(t *testing.T) {
	purge, db, err := InitTestContainer()
	defer purge(true)
	if err != nil {
		t.Fatal(err)
	}

	if Jwtuser == "" {
		t.Fatal("missiong jwt")
	}

	newEndpoints := []model.Endpoint{}
	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/endpoint/generate", gen.EndpointGenMsg{
		ProtocolHandler: "mqtt",
		Endpoint:        "foo/bar/1",
		Parts:           []gen.EndpointGenMsgPart{{MsgSegmentName: "payload", Msg: `{"foo": "bar"}`}},
	}, &newEndpoints)
	if err != nil {
		t.Fatal(err)
	}

	if len(newEndpoints) != 1 || newEndpoints[0].Endpoint != "foo/bar/1" || newEndpoints[0].ProtocolHandler != "mqtt" || newEndpoints[0].Device == "" || newEndpoints[0].Service == "" {
		t.Fatal(newEndpoints)
	}

	time.Sleep(5 * time.Second)

	endpoints := []model.Endpoint{}

	err = Jwtuser.GetJSON("http://localhost:"+util.Config.ServerPort+"/endpoints/10/0", &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 ||
		endpoints[0].Endpoint != newEndpoints[0].Endpoint ||
		endpoints[0].ProtocolHandler != newEndpoints[0].ProtocolHandler ||
		endpoints[0].Service != newEndpoints[0].Service ||
		endpoints[0].Device != newEndpoints[0].Device {
		t.Fatal(endpoints, newEndpoints)
	}

	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/endpoint/in", model.Endpoint{ProtocolHandler: "mqtt", Endpoint: "foo/bar/1"}, &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 ||
		endpoints[0].Endpoint != newEndpoints[0].Endpoint ||
		endpoints[0].ProtocolHandler != newEndpoints[0].ProtocolHandler ||
		endpoints[0].Service != newEndpoints[0].Service ||
		endpoints[0].Device != newEndpoints[0].Device {
		t.Fatal(endpoints, newEndpoints)
	}

	deviceid := endpoints[0].Device
	serviceid := endpoints[0].Service

	device := model.DeviceInstance{}
	err = Jwtuser.GetJSON("http://localhost:"+util.Config.ServerPort+"/deviceInstance/"+url.PathEscape(deviceid), &device)
	if err != nil {
		t.Fatal(err)
	}
	if device.Name == "" {
		t.Fatal(device)
	}

	dt := model.DeviceType{}
	err = Jwtuser.GetJSON("http://localhost:"+util.Config.ServerPort+"/deviceType/"+url.PathEscape(device.DeviceType), &dt)
	if err != nil {
		t.Fatal(err)
	}
	if dt.Name == "" || len(dt.Services) != 1 {
		t.Fatal(dt)
	}

	//remove generated flag
	withoutGenFlag := dt
	copy(withoutGenFlag.Services, dt.Services)
	withoutGenFlag.Generated = false
	ordf := db.GetOrdf()
	_, err = ordf.Update(dt, withoutGenFlag)
	if err != nil {
		t.Fatal(err)
	}

	dt.Services[0].EndpointFormat = "prefix,path::{{device_uri}},service::{{service_uri}}"
	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/deviceType/"+url.PathEscape(device.DeviceType), dt, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	dt = model.DeviceType{}
	err = Jwtuser.GetJSON("http://localhost:"+util.Config.ServerPort+"/deviceType/"+url.PathEscape(device.DeviceType), &dt)
	if err != nil {
		t.Fatal(err)
	}
	if dt.Name == "" || len(dt.Services) != 1 || dt.Services[0].EndpointFormat != "prefix,path::{{device_uri}},service::{{service_uri}}" {
		t.Fatal(dt)
	}

	endpoints = []model.Endpoint{}

	err = Jwtuser.GetJSON("http://localhost:"+util.Config.ServerPort+"/endpoints/10/0", &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 ||
		endpoints[0].Endpoint != "prefix,path::foo/bar/1,service::get" ||
		endpoints[0].ProtocolHandler != "mqtt" ||
		endpoints[0].Service != serviceid ||
		endpoints[0].Device != deviceid {
		t.Fatal(endpoints)
	}

	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/endpoint/in", model.Endpoint{ProtocolHandler: "mqtt", Endpoint: "prefix,path::foo/bar/1,service::get"}, &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 ||
		endpoints[0].Endpoint != "prefix,path::foo/bar/1,service::get" ||
		endpoints[0].ProtocolHandler != "mqtt" ||
		endpoints[0].Service != serviceid ||
		endpoints[0].Device != deviceid {
		t.Fatal(endpoints)
	}

	device.Url = "foo/bar/2"
	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/deviceInstance/"+url.PathEscape(deviceid), device, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	endpoints = []model.Endpoint{}

	err = Jwtuser.GetJSON("http://localhost:"+util.Config.ServerPort+"/endpoints/10/0", &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 ||
		endpoints[0].Endpoint != "prefix,path::foo/bar/2,service::get" ||
		endpoints[0].ProtocolHandler != "mqtt" ||
		endpoints[0].Service != serviceid ||
		endpoints[0].Device != deviceid {
		t.Fatal(endpoints)
	}

	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/endpoint/in", model.Endpoint{ProtocolHandler: "mqtt", Endpoint: "prefix,path::foo/bar/1,service::get"}, &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 0 {
		t.Fatal(endpoints)
	}

	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/endpoint/in", model.Endpoint{ProtocolHandler: "mqtt", Endpoint: "prefix,path::foo/bar/2,service::get"}, &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 1 ||
		endpoints[0].Endpoint != "prefix,path::foo/bar/2,service::get" ||
		endpoints[0].ProtocolHandler != "mqtt" ||
		endpoints[0].Service != serviceid ||
		endpoints[0].Device != deviceid {
		t.Fatal(endpoints)
	}

	req, err := http.NewRequest("DELETE", "http://localhost:"+util.Config.ServerPort+"/deviceInstance/"+url.PathEscape(deviceid), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", string(Jwtuser))

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	endpoints = []model.Endpoint{}

	err = Jwtuser.GetJSON("http://localhost:"+util.Config.ServerPort+"/endpoints/10/0", &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 0 {
		t.Fatal(endpoints)
	}

	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/endpoint/in", model.Endpoint{ProtocolHandler: "mqtt", Endpoint: "prefix,path::foo/bar/1,service::get"}, &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 0 {
		t.Fatal(endpoints)
	}

	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/endpoint/in", model.Endpoint{ProtocolHandler: "mqtt", Endpoint: "prefix,path::foo/bar/2,service::get"}, &endpoints)
	if err != nil {
		t.Fatal(err)
	}
	if len(endpoints) != 0 {
		t.Fatal(endpoints)
	}

	//TODO other tests with multiple existsing other endpoints (no side effects)

	//TODO: test migration
}
