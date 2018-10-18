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

package gen

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"sort"

	"strings"

	"github.com/bouk/monkey"
)

func ExampleFormatCheck() {
	jsonval := `{
	"a": "a",
	"b": 1,
	"c": [1,2,3],
    "d": ["a", "b", "c"]
}`
	format, struca := getFormatStruct(jsonval)
	var strucb interface{}
	err := json.Unmarshal([]byte(jsonval), &strucb)
	fmt.Println(format == JSON_FORMAT, err, reflect.DeepEqual(strucb, struca))

	// Output:
	// true <nil> true
}

type dbMockCheck func(model.ValueType) (bool, string, error)

type DbMock struct {
	ValueTypeQueryMock dbMockCheck
}

func (this DbMock) ValueTypeQuery(valueType model.ValueType) (exists bool, id string, err error) {
	if this.ValueTypeQueryMock == nil {
		return false, "", nil
	}

	return this.ValueTypeQueryMock(valueType)
}

type ByName []model.FieldType

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) < 0 }

func ExampleGenerateValueType() {
	count = 0

	//set time: ultra dirty
	wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	jsonval := `{
	"a": "a",
	"b": 1,
	"c": [1,2,3],
    "d": ["a", "b", "c"]
}`
	_, struc := getFormatStruct(jsonval)
	valueType, err := interfaceToValueType(DbMock{}, struc)
	sort.Sort(ByName(valueType.Fields))
	vtJson, err2 := json.Marshal(valueType)
	fmt.Println(string(vtJson), err, err2)

	// Output:
	//{"name":"138157323_3","description":"generated","base_type":"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#structure","fields":[{"name":"a","type":{"id":"iot#c8c36810-c8e0-403e-b00f-187414a84ccd","fields":null,"literal":""}},{"name":"b","type":{"id":"iot#cb0dc896-6d89-4e0c-ac59-33eceed512b0","fields":null,"literal":""}},{"name":"c","type":{"name":"138157323_1","description":"generated","base_type":"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#list","fields":[{"type":{"id":"iot#cb0dc896-6d89-4e0c-ac59-33eceed512b0","fields":null,"literal":""}}],"literal":""}},{"name":"d","type":{"name":"138157323_2","description":"generated","base_type":"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#list","fields":[{"type":{"id":"iot#c8c36810-c8e0-403e-b00f-187414a84ccd","fields":null,"literal":""}}],"literal":""}}],"literal":""} <nil> <nil>
}

func ExampleGenerateValueType2() {
	count = 0

	//set time: ultra dirty
	wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	jsonval := `{
	"a": "a",
	"b": 1,
	"c": [1,2,3],
    "d": ["a", "b", "c"]
}`
	_, struc := getFormatStruct(jsonval)
	valueType, err := interfaceToValueType(DbMock{
		ValueTypeQueryMock: func(valueType model.ValueType) (bool, string, error) {
			if valueType.BaseType == model.ListBaseType {
				return true, "thisismyid", nil
			}
			return false, "", nil
		},
	}, struc)
	sort.Sort(ByName(valueType.Fields))
	vtJson, err2 := json.Marshal(valueType)
	fmt.Println(string(vtJson), err, err2)

	// Output:
	//{"name":"138157323_1","description":"generated","base_type":"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#structure","fields":[{"name":"a","type":{"id":"iot#c8c36810-c8e0-403e-b00f-187414a84ccd","fields":null,"literal":""}},{"name":"b","type":{"id":"iot#cb0dc896-6d89-4e0c-ac59-33eceed512b0","fields":null,"literal":""}},{"name":"c","type":{"id":"thisismyid","fields":null,"literal":""}},{"name":"d","type":{"id":"thisismyid","fields":null,"literal":""}}],"literal":""} <nil> <nil>
}
