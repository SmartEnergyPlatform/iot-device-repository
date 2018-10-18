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
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/model"

	"log"

	"github.com/knakk/rdf"
	"github.com/SmartEnergyPlatform/iot-device-repository/lib/interfaces"
)

type Format string

const (
	UNKNOWN_FORMAT Format = ""
	JSON_FORMAT           = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#json"
	PLAIN_FORMAT          = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#PlainText"
)

func getFormatStruct(value string) (format Format, structure interface{}) {
	err := json.Unmarshal([]byte(value), &structure)
	if err == nil {
		format = JSON_FORMAT
		return
	}

	structure = value
	format = PLAIN_FORMAT
	return
}

func interfaceToValueType(db interfaces.Persistence, value interface{}) (result model.ValueType, err error) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice:
		result, err = sliceToValueType(db, v)
	case reflect.Map:
		result, err = mapToValueType(db, v)
	case reflect.Int64:
		result.Id = "iot#01190060-db2e-4ed0-a424-c82b60f981e4"
	case reflect.String:
		result.Id = "iot#c8c36810-c8e0-403e-b00f-187414a84ccd"
	case reflect.Bool:
		result.Id = "iot#939963e5-1ab0-44e0-8fb4-5235fd6f5363"
	case reflect.Float64:
		result.Id = "iot#cb0dc896-6d89-4e0c-ac59-33eceed512b0"
	default:
		err = errors.New("unknown kind of value: " + v.Kind().String())
	}
	if result.Id == "" {
		result.Name = generateValueTypeName(result)
		result.Description = "generated"
	}
	return
}

var count = 0

const maxCount = 10000000
const TermSymbol = rdf.TermLiteral + 1

func getNextCount() int {
	count = (count % maxCount) + 1
	return count
}

func generateValueTypeName(valueType model.ValueType) (result string) {
	return strconv.FormatInt(time.Now().Unix(), 10) + "_" + strconv.Itoa(getNextCount())
}

func mapToValueType(db interfaces.Persistence, mapVal reflect.Value) (result model.ValueType, err error) {
	prevIsNew := false
	result.BaseType = model.StructBaseType
	for _, key := range mapVal.MapKeys() {
		val := mapVal.MapIndex(key)
		kind := reflect.ValueOf(val.Interface()).Kind()
		if kind != reflect.Invalid {
			subType, err := interfaceToValueType(db, val.Interface())
			if err != nil {
				log.Println("ERROR: mapToValueType() field ", key.String(), kind.String(), err)
				return result, err
			}
			if subType.Id == "" && !prevIsNew {
				prevIsNew = true
			}
			if subType.Id != "" || subType.BaseType != "" {
				result.Fields = append(result.Fields, model.FieldType{Name: key.String(), Type: subType})
			}
		}
	}
	if !prevIsNew {
		result, err = checkForExistingValueType(db, result)
	}
	return
}

func sliceToValueType(db interfaces.Persistence, slice reflect.Value) (result model.ValueType, err error) {
	length := slice.Len()
	prevIsNew := false
	if length > 0 {
		result.BaseType = model.ListBaseType
		subType, err := interfaceToValueType(db, slice.Index(0).Interface())
		if err != nil {
			return result, err
		}
		prevIsNew = subType.Id == ""
		result.Fields = append(result.Fields, model.FieldType{Type: subType})
	}
	if !prevIsNew {
		result, err = checkForExistingValueType(db, result)
	}
	return
}

func checkForExistingValueType(db interfaces.Persistence, vt model.ValueType) (result model.ValueType, err error) {
	exists, id, err := db.ValueTypeQuery(vt)
	if err != nil {
		return result, err
	}
	if exists {
		result.Id = id
	} else {
		result = vt
	}
	return
}

func ValueTypeFromMessage(db interfaces.Persistence, msg string) (valueType model.ValueType, format Format, struc interface{}, err error) {
	format, struc = getFormatStruct(msg)
	if format == UNKNOWN_FORMAT {
		err = errors.New("unknown message format; able to interprete the following formats: " + JSON_FORMAT)
	} else {
		valueType, err = interfaceToValueType(db, struc)
	}
	return
}
