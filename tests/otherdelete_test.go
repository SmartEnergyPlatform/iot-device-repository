/*
 * Copyright 2019 InfAI (CC SES)
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
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/api"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/persistence"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/util"
	"github.com/SmartEnergyPlatform/jwt-http-router"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/ory/dockertest"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestDeleteVendorAndDeviceClass(t *testing.T)  {
	impersonate := jwt_http_router.JwtImpersonate("Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIzaUtabW9aUHpsMmRtQnBJdS1vSkY4ZVVUZHh4OUFIckVOcG5CcHM5SjYwIn0.eyJqdGkiOiJiNjdjYjI2YS00OWFhLTQ0NTctYjE2NC0yNGJhMmNlMzE0YmYiLCJleHAiOjE1NDgwNzk2MTYsIm5iZiI6MCwiaWF0IjoxNTQ4MDc2MDE2LCJpc3MiOiJodHRwczovL2F1dGguc2VwbC5pbmZhaS5vcmcvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiZnJvbnRlbmQiLCJzdWIiOiJkZDY5ZWEwZC1mNTUzLTQzMzYtODBmMy03ZjQ1NjdmODVjN2IiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJmcm9udGVuZCIsIm5vbmNlIjoiNmI3YTJlNzQtYzAwNy00NmNjLWJhMWItMDIxYTRjYjY2NjE1IiwiYXV0aF90aW1lIjoxNTQ4MDU4MzM0LCJzZXNzaW9uX3N0YXRlIjoiNTA3NTMyMmMtNzc5ZC00NGZiLWExNDYtNDhmODYzODdiYmY5IiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJjcmVhdGUtcmVhbG0iLCJhZG1pbiIsImRldmVsb3BlciIsInVtYV9hdXRob3JpemF0aW9uIiwidXNlciJdfSwicmVzb3VyY2VfYWNjZXNzIjp7Im1hc3Rlci1yZWFsbSI6eyJyb2xlcyI6WyJ2aWV3LWlkZW50aXR5LXByb3ZpZGVycyIsInZpZXctcmVhbG0iLCJtYW5hZ2UtaWRlbnRpdHktcHJvdmlkZXJzIiwiaW1wZXJzb25hdGlvbiIsImNyZWF0ZS1jbGllbnQiLCJtYW5hZ2UtdXNlcnMiLCJxdWVyeS1yZWFsbXMiLCJ2aWV3LWF1dGhvcml6YXRpb24iLCJxdWVyeS1jbGllbnRzIiwicXVlcnktdXNlcnMiLCJtYW5hZ2UtZXZlbnRzIiwibWFuYWdlLXJlYWxtIiwidmlldy1ldmVudHMiLCJ2aWV3LXVzZXJzIiwidmlldy1jbGllbnRzIiwibWFuYWdlLWF1dGhvcml6YXRpb24iLCJtYW5hZ2UtY2xpZW50cyIsInF1ZXJ5LWdyb3VwcyJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwicm9sZXMiOlsidW1hX2F1dGhvcml6YXRpb24iLCJhZG1pbiIsImNyZWF0ZS1yZWFsbSIsImRldmVsb3BlciIsInVzZXIiLCJvZmZsaW5lX2FjY2VzcyJdLCJuYW1lIjoiU2VwbCBBZG1pbiIsInByZWZlcnJlZF91c2VybmFtZSI6InNlcGwiLCJnaXZlbl9uYW1lIjoiU2VwbCIsImZhbWlseV9uYW1lIjoiQWRtaW4iLCJlbWFpbCI6InNlcGxAc2VwbC5kZSJ9.UYFgPlwWTmXOLT4b3RNPlzlOXT-0vSzoivqPwtuC_7xlkIwt8vaHvPH8wPQvLHr52HspKvjoYZnjFVMEI4e-kzMfiBUvlY0dEq4LxTjFrI3hfLA1xkTSNyOVfgFkghXgd8o1UtCwoKj92X0IViawy_JEHbYzeXT3431rFlh0XwzP1fqVZSF7fMWiRHvdptMmYMSEtBUclPJ7uioX6yQjcY-IN5m71GcfxwcY2la45kj_DH3wNDG0iGXZEpNFzcHWlt5_46lVPusu15owIWRjASf-JDDg6ALCYKz_bh8Vfvz1XJrwlO-AHTRiOLX8BgAjVx2qGLCu1i1DLzvzbpXz-g")

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Error("Could not connect to docker:", err)
		return
	}
	dockerrdf, err := pool.Run("fgseitsrancher.wifa.intern.uni-leipzig.de:5000/iot-ontology", "unstable", []string{
		"DBA_PASSWORD=myDbaPassword",
		"DEFAULT_GRAPH=iot_test",
	})
	defer dockerrdf.Close()
	if err != nil {
		t.Error("Could not start rdf:", err)
		return
	}
	time.Sleep(20 * time.Second)

	err = util.LoadConfig("../config.json")
	util.Config.SparqlEndpoint = "http://localhost:" + dockerrdf.GetPort("8890/tcp") + "/sparql"
	util.Config.ForceAuth = "false"
	util.Config.ForceUser = "false"
	server := httptest.NewServer(logger.New(api.GetRoutes(persistence.New()), util.Config.LogLevel))
	defer server.Close()

	vendor_1 := model.Vendor{}
	vendor_1_result := api.Insert_OK{}
	err = impersonate.PostJSON(server.URL+"/other/vendor", vendor_1, &vendor_1_result)
	if err != nil {
		t.Error(err)
		return
	}
	if vendor_1_result.CreatedId != "" {
		t.Error("expect no result on request without name", err, vendor_1_result)
		return
	}

	vendor_list_result := []model.Vendor{}
	err = impersonate.GetJSON(server.URL+"/ui/search/others/vendors//20/0", &vendor_list_result)
	if err != nil {
		t.Error(err)
		return
	}

	if len(vendor_list_result) != 0 {
		t.Error("unexpected result: ", vendor_list_result)
		return
	}

	vendor_2 := model.Vendor{Name:"test_1"}
	vendor_2_result := api.Insert_OK{}
	err = impersonate.PostJSON(server.URL+"/other/vendor", vendor_2, &vendor_2_result)
	if err != nil || vendor_2_result.CreatedId == "" {
		t.Error(err, vendor_2_result)
		return
	}

	vendor_list_result = []model.Vendor{}
	err = impersonate.GetJSON(server.URL+"/ui/search/others/vendors//20/0", &vendor_list_result)
	if err != nil {
		t.Error(err)
		return
	}

	if len(vendor_list_result) != 1 || vendor_list_result[0].Name != vendor_2.Name || vendor_list_result[0].Id != vendor_2_result.CreatedId  {
		t.Error("unexpected result: ", vendor_list_result)
	}

	req, err := http.NewRequest("DELETE", server.URL+"/other/vendor/"+url.PathEscape(vendor_2_result.CreatedId), nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Authorization", string(impersonate))
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	vendor_list_result = []model.Vendor{}
	err = impersonate.GetJSON(server.URL+"/ui/search/others/vendors//20/0", &vendor_list_result)
	if err != nil {
		t.Error(err)
		return
	}

	if len(vendor_list_result) != 0 {
		t.Error("unexpected result: ", vendor_list_result)
	}
	
	// deviceClass
	deviceClass_1 := model.DeviceClass{}
	deviceClass_1_result := api.Insert_OK{}
	err = impersonate.PostJSON(server.URL+"/other/deviceclass", deviceClass_1, &deviceClass_1_result)
	if err != nil {
		t.Error(err)
		return
	}
	if deviceClass_1_result.CreatedId != "" {
		t.Error("expect no result on request without name", err, deviceClass_1_result)
		return
	}

	deviceClass_list_result := []model.DeviceClass{}
	err = impersonate.GetJSON(server.URL+"/ui/search/others/deviceClasses//20/0", &deviceClass_list_result)
	if err != nil {
		t.Error(err)
		return
	}

	if len(deviceClass_list_result) != 0 {
		t.Error("unexpected result: ", deviceClass_list_result)
	}

	deviceClass_2 := model.DeviceClass{Name:"test_1"}
	deviceClass_2_result := api.Insert_OK{}
	err = impersonate.PostJSON(server.URL+"/other/deviceclass", deviceClass_2, &deviceClass_2_result)
	if err != nil || deviceClass_2_result.CreatedId == "" {
		t.Error(err)
		return
	}

	deviceClass_list_result = []model.DeviceClass{}
	err = impersonate.GetJSON(server.URL+"/ui/search/others/deviceClasses//20/0", &deviceClass_list_result)
	if err != nil {
		t.Error(err)
		return
	}

	if len(deviceClass_list_result) != 1 || deviceClass_list_result[0].Name != deviceClass_2.Name || deviceClass_list_result[0].Id != deviceClass_2_result.CreatedId  {
		t.Error("unexpected result: ", deviceClass_list_result)
		return
	}

	req, err = http.NewRequest("DELETE", server.URL+"/other/deviceclass/"+url.PathEscape(deviceClass_2_result.CreatedId), nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Authorization", string(impersonate))
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	deviceClass_list_result = []model.DeviceClass{}
	err = impersonate.GetJSON(server.URL+"/ui/search/others/deviceClasses//20/0", &deviceClass_list_result)
	if err != nil {
		t.Error(err)
		return
	}

	if len(deviceClass_list_result) != 0 {
		t.Error("unexpected result: ", deviceClass_list_result)
	}
}