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

	"encoding/json"

	"net/http"

	"time"

	"io/ioutil"

	"log"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
)

func TestValuetypeSearch(t *testing.T) {
	purge, _, err := InitTestContainer()
	defer purge(true)
	if err != nil {
		t.Fatal(err)
	}

	if Jwtuser == "" {
		t.Fatal("missiong jwt")
	}

	dtStr := `{  
   "id":"iot#f8b43fd0-6318-4cca-82d3-71eb8e6fce79",
   "name":"test",
   "description":"test",
   "device_class":{  
      "id":"iot#3e522022-38ee-4a8b-b5c7-dbcb54b887d1",
      "name":"test"
   },
   "services":[  
      {  
         "id":"iot#bd936af1-ad93-4dc9-b310-bf93264de0eb",
         "service_type":"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Actuator",
         "name":"test",
         "description":"test",
         "protocol":{  
            "id":"iot#d6a462c5-d4e0-4396-b3f3-28cd37b647a8",
            "protocol_handler_url":"connector",
            "name":"standard-connector",
            "description":"Generic protocol for transporting data and metadata.",
            "msg_structure":[  
               {  
                  "id":"iot#37ff5298-a7dd-4744-9080-7cfdbda5dc72",
                  "name":"metadata",
                  "constraints":null
               },
               {  
                  "id":"iot#88cd5b0e-a451-4070-a20d-464ee23742dd",
                  "name":"data",
                  "constraints":null
               }
            ]
         },
         "input":[  
            {  
               "id":"iot#7398cf74-2194-4399-841b-cf401dc8a67e",
               "name":"test",
               "msg_segment":{  
                  "id":"iot#88cd5b0e-a451-4070-a20d-464ee23742dd",
                  "name":"data",
                  "constraints":null
               },
               "type":{  
                  "id":"iot#e69373a9-2ab9-4dc4-b5d5-ff57aa742c3e",
                  "name":"test",
                  "description":"test",
                  "base_type":"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#structure",
                  "fields":[  
                     {  
                        "id":"iot#70908900-4acb-4b94-91ff-1c05b4f23c77",
                        "name":"a",
                        "type":{  
                           "id":"iot#e9104f3f-ffe1-410a-befa-dd68f0677ec6",
                           "name":"test_int",
                           "description":"test_int",
                           "base_type":"http://www.w3.org/2001/XMLSchema#integer",
                           "fields":null,
                           "literal":""
                        }
                     }
                  ],
                  "literal":""
               },
               "format":"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#json",
               "additional_formatinfo":[  
                  {  
                     "id":"iot#32f58890-8b8f-4e06-9082-e7848c116154",
                     "field":{  
                        "id":"iot#70908900-4acb-4b94-91ff-1c05b4f23c77",
                        "name":"a",
                        "type":{  
                           "id":"iot#e9104f3f-ffe1-410a-befa-dd68f0677ec6",
                           "name":"test_int",
                           "description":"test_int",
                           "base_type":"http://www.w3.org/2001/XMLSchema#integer",
                           "fields":null,
                           "literal":""
                        }
                     },
                     "format_flag":""
                  }
               ]
            }
         ],
         "url":"test"
      }
   ],
   "vendor":{  
      "id":"iot#91bff598-bd63-44ce-aa5e-f66e092b7279",
      "name":"test"
   }
}`
	dt := model.DeviceType{}
	json.Unmarshal([]byte(dtStr), &dt)

	//create devicetype and 2 value types

	err = Jwtuser.PostJSON("http://localhost:"+util.Config.ServerPort+"/import/deviceType", dt, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	// check if value type exists

	response, err := Jwtuser.Get(VtSearchUrl + "/search/valuetype/test/endpoint/20/0/name/asc")
	log.Println(err)
	b, err := ioutil.ReadAll(response.Body)
	log.Println(string(b), err)

	vts := []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/asc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 2 || vts[0]["name"] != "test" || vts[1]["name"] != "test_int" {
		t.Fatal(vts)
	}

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/desc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 2 || vts[1]["name"] != "test" || vts[0]["name"] != "test_int" {
		t.Fatal(vts)
	}

	//try deleting value type 2: should fail

	req, err := http.NewRequest("DELETE", "http://localhost:"+util.Config.ServerPort+"/valueType/iot%23e9104f3f-ffe1-410a-befa-dd68f0677ec6", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", string(Jwtuser))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) == "ok" {
		t.Fatal(string(msg))
	}

	//try deleting value type 1: should fail

	req, err = http.NewRequest("DELETE", "http://localhost:"+util.Config.ServerPort+"/valueType/iot%23e69373a9-2ab9-4dc4-b5d5-ff57aa742c3e", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", string(Jwtuser))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	msg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) == "ok" {
		t.Fatal(string(msg))
	}

	// delete device type

	req, err = http.NewRequest("DELETE", "http://localhost:"+util.Config.ServerPort+"/deviceType/iot%23f8b43fd0-6318-4cca-82d3-71eb8e6fce79", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", string(Jwtuser))

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	// value types should still exist

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/asc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 2 || vts[0]["name"] != "test" || vts[1]["name"] != "test_int" {
		t.Fatal(vts)
	}

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/desc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 2 || vts[1]["name"] != "test" || vts[0]["name"] != "test_int" {
		t.Fatal(vts)
	}

	//try deleting value type 2: should fail

	req, err = http.NewRequest("DELETE", "http://localhost:"+util.Config.ServerPort+"/valueType/iot%23e9104f3f-ffe1-410a-befa-dd68f0677ec6", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", string(Jwtuser))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	msg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) == "ok" {
		t.Fatal(string(msg))
	}

	time.Sleep(5 * time.Second)

	// value types should still exist

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/asc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 2 || vts[0]["name"] != "test" || vts[1]["name"] != "test_int" {
		t.Fatal(vts)
	}

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/desc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 2 || vts[1]["name"] != "test" || vts[0]["name"] != "test_int" {
		t.Fatal(vts)
	}

	//try deleting value type 1

	req, err = http.NewRequest("DELETE", "http://localhost:"+util.Config.ServerPort+"/valueType/iot%23e69373a9-2ab9-4dc4-b5d5-ff57aa742c3e", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", string(Jwtuser))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	msg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) != "ok" {
		t.Fatal(string(msg))
	}

	time.Sleep(5 * time.Second)

	// value type 2 should still exist, 1 not

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/asc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 1 || vts[0]["name"] != "test_int" {
		t.Fatal(vts)
	}

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/desc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 1 || vts[0]["name"] != "test_int" {
		t.Fatal(vts)
	}

	//try deleting value type 2

	req, err = http.NewRequest("DELETE", "http://localhost:"+util.Config.ServerPort+"/valueType/iot%23e9104f3f-ffe1-410a-befa-dd68f0677ec6", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", string(Jwtuser))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	msg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(msg) != "ok" {
		t.Fatal(string(msg))
	}

	// value type 2 should not exist

	time.Sleep(5 * time.Second)

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/asc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 0 {
		t.Fatal(vts)
	}

	vts = []map[string]interface{}{}
	err = Jwtuser.GetJSON(VtSearchUrl+"/search/valuetype/test/endpoint/20/0/name/desc", &vts)
	if err != nil {
		t.Fatal(err)
	}
	if len(vts) != 0 {
		t.Fatal(vts)
	}
}
