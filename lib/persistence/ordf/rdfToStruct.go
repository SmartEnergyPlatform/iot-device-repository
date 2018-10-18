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

package ordf

import (
	"reflect"

	"github.com/knakk/rdf"
)

func rdfToStruct(structType *reflect.Value, id string, triples map[string][]map[string]rdf.Term) {
	predicateMap, _ := getFieldMapping(structType.Type(), RDF_FIELD_TAG_NAME)
	idFieldName, _ := GetFieldNameWithTag(structType.Type(), RDF_ENTITY_TAG_NAME)
	if idFieldName != "" {
		structType.FieldByName(idFieldName).SetString(id)
	}
	for _, triple := range triples[id] {
		predicate, predicateExists := triple["p"]
		obj, objExists := triple["o"]
		if predicateExists && objExists {
			if fieldName, ok := predicateMap[predicate.String()]; ok {
				field := structType.FieldByName(fieldName)
				if field.Kind() == reflect.String {
					field.SetString(obj.String())
				} else if field.Kind() == reflect.Bool {
					field.SetBool(obj.String() == "1")
				} else if field.Kind() == reflect.Struct && obj.Type() == rdf.TermIRI {
					rdfToStruct(&field, obj.String(), triples)
				} else if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct && obj.Type() == rdf.TermIRI {
					element := reflect.Indirect(reflect.New(field.Type().Elem()))
					rdfToStruct(&element, obj.String(), triples)
					field.Set(reflect.Append(field, element))
				} else if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.String {
					field.Set(reflect.Append(field, reflect.ValueOf(obj.String())))
				}
			}
		}
	}
}

//sorts triples by subject; deserializeType: n^2 --> n
func groupBySubject(triples []map[string]rdf.Term) (result map[string][]map[string]rdf.Term) {
	result = map[string][]map[string]rdf.Term{}
	for _, triple := range triples {
		subjectId := triple["s"].String()
		if subjTriples, ok := result[subjectId]; ok {
			result[subjectId] = append(subjTriples, triple)
		} else {
			result[subjectId] = []map[string]rdf.Term{triple}
		}
	}
	return
}

func getEntitysOfType(entityType reflect.Type, triples map[string][]map[string]rdf.Term) (result []string) {
	_, rdfType := GetFieldNameWithTag(entityType, RDF_ENTITY_TAG_NAME)
	for id, tripleList := range triples {
		if isRdfType(tripleList, rdfType) {
			result = append(result, id)
		}
	}
	return
}

func isRdfType(triples []map[string]rdf.Term, rdfTypeName string) bool {
	for _, triple := range triples {
		if triple["p"].String() == RDF_TYPE_PREDICATE_NAME {
			if triple["o"].String() == rdfTypeName {
				return true
			}
		}
	}
	return false
}

func RdfToStructList(result interface{}, triples []map[string]rdf.Term) {
	resultValue := reflect.Indirect(reflect.ValueOf(result))
	resultElemType := resultValue.Type().Elem()
	preparedTriples := groupBySubject(triples)
	entitys := getEntitysOfType(resultElemType, preparedTriples)
	for _, entity := range entitys {
		value := reflect.Indirect(reflect.New(resultElemType))
		rdfToStruct(&value, entity, preparedTriples)
		resultValue.Set(reflect.Append(resultValue, value))
	}
}

func RdfToStruct(result interface{}, id string, triples []map[string]rdf.Term) {
	resultValue := reflect.Indirect(reflect.ValueOf(result))
	rdfToStruct(&resultValue, id, groupBySubject(triples))
}
