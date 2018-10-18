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
	"testing"

	"time"

	"net/url"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/format"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/gen"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
)

func TestLeaves(t *testing.T) {
	purge, _, err := InitTestContainer()
	defer purge(true)
	if err != nil {
		t.Fatal(err)
	}

	if Jwtuser == "" {
		t.Fatal("missiong jwt")
	}

	msg := `{
		"foo": "bar",
		"batz": [2, 3],
		"bla": {"name": "test", "something": 42},
		"list": [{"element":1}, {"element":2}],
		"n": null,
		"b": true
	}`

	newEndpoints := []model.Endpoint{}
	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/endpoint/generate", gen.EndpointGenMsg{
		ProtocolHandler: "mqtt",
		Endpoint:        "foo/bar/1",
		Parts:           []gen.EndpointGenMsgPart{{MsgSegmentName: "payload", Msg: msg}},
	}, &newEndpoints)
	if err != nil {
		t.Fatal(err)
	}

	if len(newEndpoints) != 1 || newEndpoints[0].Endpoint != "foo/bar/1" || newEndpoints[0].ProtocolHandler != "mqtt" || newEndpoints[0].Device == "" || newEndpoints[0].Service == "" {
		t.Fatal(newEndpoints)
	}

	time.Sleep(5 * time.Second)

	instanceId := newEndpoints[0].Device
	serviceId := newEndpoints[0].Service
	leavesResult := format.LeavesResult{}
	err = Jwtuser.GetJSON("http://localhost:"+util.Config.ServerPort+"/skeleton/"+url.PathEscape(instanceId)+"/"+url.PathEscape(serviceId)+"/output/leaves", &leavesResult)
	if err != nil {
		t.Fatal(err)
	}

	if len(leavesResult.Leaves) != 9 {
		t.Fatal("not matching leaves count:", leavesResult.Leaves)
	}

	leaf := format.Leaf{Path: "$.source_topic+", Type: "string"}
	if !test_sliceContainsInterface(leavesResult.Leaves, leaf) {
		t.Fatal("missing leaf", leaf)
	}

	leaf = format.Leaf{Path: "$.value.payload.foo+", Type: "string"}
	if !test_sliceContainsInterface(leavesResult.Leaves, leaf) {
		t.Fatal("missing leaf", leaf)
	}

	leaf = format.Leaf{Path: "$.value.payload.batz[0]+", Type: "float64"}
	if !test_sliceContainsInterface(leavesResult.Leaves, leaf) {
		t.Fatal("missing leaf", leaf)
	}

	leaf = format.Leaf{Path: "$.value.payload.bla.name+", Type: "string"}
	if !test_sliceContainsInterface(leavesResult.Leaves, leaf) {
		t.Fatal("missing leaf", leaf)
	}

	leaf = format.Leaf{Path: "$.value.payload.bla.something+", Type: "float64"}
	if !test_sliceContainsInterface(leavesResult.Leaves, leaf) {
		t.Fatal("missing leaf", leaf)
	}

	leaf = format.Leaf{Path: "$.value.payload.list[0].element+", Type: "float64"}
	if !test_sliceContainsInterface(leavesResult.Leaves, leaf) {
		t.Fatal("missing leaf", leaf)
	}

	leaf = format.Leaf{Path: "$.value.payload.b+", Type: "bool"}
	if !test_sliceContainsInterface(leavesResult.Leaves, leaf) {
		t.Fatal("missing leaf", leaf)
	}

	leaf = format.Leaf{Path: "$.service_id+", Type: "string"}
	if !test_sliceContainsInterface(leavesResult.Leaves, leaf) {
		t.Fatal("missing leaf", leaf)
	}

}
