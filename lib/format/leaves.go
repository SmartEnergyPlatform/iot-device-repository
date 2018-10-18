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

package format

import (
	"reflect"

	"strconv"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
)

type Leaf struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type LeavesResult struct {
	Example map[string]interface{} `json:"example"`
	Leaves  []Leaf                 `json:"leaves"`
}

func GetBpmnSkeletonOutputLeaves(db interfaces.Persistence, deviceInstanceId string, serviceId string) (result LeavesResult, err error) {
	instance, err := db.GetDeviceInstanceById(deviceInstanceId)
	if err != nil {
		return
	}
	skeleton := BpmnValueSkeleton{}
	skeleton, err = GetBpmnSkeletonFromDeviceType(db, instance.DeviceType, serviceId)
	if err != nil {
		return
	}
	result.Example = map[string]interface{}{"value": skeleton.Outputs, "device_id": deviceInstanceId, "service_id": serviceId, "source_topic": "topic"}
	result.Leaves, err = GetLeaves(result.Example)
	return
}

func GetLeaves(value map[string]interface{}) (result []Leaf, err error) {
	return getLeavesMap(value, "$")
}

func getLeavesMap(value map[string]interface{}, prefix string) (result []Leaf, err error) {
	for key, val := range value {
		path := prefix + "." + key
		leaves, err := getLeavesInterface(val, path)
		if err != nil {
			return result, err
		}
		result = append(result, leaves...)
	}
	return
}

func getLeavesArray(value []interface{}, prefix string) (result []Leaf, err error) {
	for index, val := range value {
		path := prefix + "[" + strconv.Itoa(index) + "]"
		leaves, err := getLeavesInterface(val, path)
		if err != nil {
			return result, err
		}
		result = append(result, leaves...)
	}
	return
}

func getLeavesInterface(val interface{}, path string) (result []Leaf, err error) {
	switch v := val.(type) {
	case map[string]interface{}:
		leaves, err := getLeavesMap(v, path)
		if err != nil {
			return result, err
		}
		result = append(result, leaves...)
	case []interface{}:
		leaves, err := getLeavesArray(v, path)
		if err != nil {
			return result, err
		}
		result = append(result, leaves...)
	case nil:
		result = append(result, Leaf{
			Path: path,
			Type: "null",
		})
	default:
		leaves, err := getLeavesPrimitives(val, path)
		if err != nil {
			return result, err
		}
		result = append(result, leaves...)
	}
	return
}

func getLeavesPrimitives(value interface{}, path string) (result []Leaf, err error) {
	result = []Leaf{{
		Path: path + "+",
		Type: reflect.TypeOf(value).Name(),
	}}
	return
}
