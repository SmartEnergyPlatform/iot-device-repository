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
	"errors"
	"reflect"

	"github.com/knakk/rdf"
)

func isEmpty(value reflect.Value) bool {
	return reflect.DeepEqual(reflect.Zero(value.Type()).Interface(), value.Interface())
}

func structToRdf(structValue reflect.Value, triples *[]map[string]rdf.Term, sideEffectsAllowed bool) (id rdf.Term, err error) {
	firstCall := len(*triples) == 0
	predicateMap, _ := getFieldMapping(structValue.Type(), RDF_FIELD_TAG_NAME)
	idFieldName, rdfType := GetFieldNameWithTag(structValue.Type(), RDF_ENTITY_TAG_NAME)
	idField, idFieldTypeExists := structValue.Type().FieldByName(idFieldName)
	if !idFieldTypeExists {
		err = errors.New("struct has no id field (missing tag: " + RDF_ENTITY_TAG_NAME + ")")
		return
	}
	idTemp, err2 := rdf.NewIRI(structValue.FieldByName(idFieldName).String())
	if err2 == nil {
		id = idTemp
	} else {
		id = newSymbol()
	}
	isRootEntity := "true" == idField.Tag.Get(RDF_ROOT_ENTITY_TAG_NAME)
	if sideEffectsAllowed || firstCall || !isRootEntity {
		addRdfType(triples, id, rdfType)
		for predicate, fieldName := range predicateMap {
			field, _ := structValue.Type().FieldByName(fieldName)
			isRef := "true" == field.Tag.Get(RDF_REF_TAG_NAME)
			isLend := "true" == field.Tag.Get(RDF_LEND_TAG_NAME)
			if sideEffectsAllowed || !isLend {
				err = addTriple(triples, id, predicate, structValue.FieldByName(fieldName), isRef, sideEffectsAllowed)
			}
		}
	}
	return
}

func addRdfType(triples *[]map[string]rdf.Term, subject rdf.Term, objectString string) (err error) {
	object, err := rdf.NewIRI(objectString)
	if err != nil {
		return
	}
	predicate, err := rdf.NewIRI(RDF_TYPE_PREDICATE_NAME)
	if err != nil {
		return
	}
	*triples = append(*triples, buildTriple(subject, predicate, object))
	return
}

func addTriple(triples *[]map[string]rdf.Term, id rdf.Term, predicate string, fieldValue reflect.Value, isRef bool, sideEffectsAllowed bool) (err error) {
	predicateIRI, err := rdf.NewIRI(predicate)
	if err != nil {
		return
	}
	if !isEmpty(fieldValue) {
		if fieldValue.Kind() == reflect.Slice {
			sliceToRdf(triples, fieldValue, id, predicateIRI, isRef, sideEffectsAllowed)
		} else {
			object, err := fieldToRdfTerm(triples, fieldValue, isRef, sideEffectsAllowed)
			if err == nil && isValidObject(object) {
				*triples = append(*triples, buildTriple(id, predicateIRI, object))
			}
		}
	}
	return
}

func isValidObject(term rdf.Term) (result bool) {
	return !(term.Type() == rdf.TermIRI && term.String() == "")
}

func buildTriple(subject rdf.Term, predicate rdf.Term, object rdf.Term) map[string]rdf.Term {
	zero := rdf.IRI{}
	if subject == zero {
		return map[string]rdf.Term{"p": predicate, "o": object}
	}
	return map[string]rdf.Term{"s": subject, "p": predicate, "o": object}
}

func sliceToRdf(triples *[]map[string]rdf.Term, field reflect.Value, subject rdf.Term, predicate rdf.Term, isRef bool, sideEffectsAllowed bool) (err error) {
	for i := 0; i < field.Len(); i++ {
		element := field.Index(i)
		object, err := fieldToRdfTerm(triples, element, isRef, sideEffectsAllowed)
		if err != nil {
			return err
		}
		if isValidObject(object) {
			*triples = append(*triples, buildTriple(subject, predicate, object))
		}
	}
	return
}

func fieldToRdfTerm(triples *[]map[string]rdf.Term, field reflect.Value, isRef bool, sideEffectsAllowed bool) (term rdf.Term, err error) {
	if field.Kind() == reflect.String {
		if isRef {
			term, err = rdf.NewIRI(field.String())
		} else {
			term, err = rdf.NewLiteral(field.String())
		}
	} else if field.Kind() == reflect.Bool {
		term, err = rdf.NewLiteral(field.Bool())
	} else if field.Kind() == reflect.Struct {
		term, err = structToRdf(field, triples, sideEffectsAllowed)
		if err != nil {
			return
		}
	}
	return
}

func StructToRdf(structure interface{}) (id rdf.Term, triples []map[string]rdf.Term, err error) {
	triples = []map[string]rdf.Term{}
	id, err = structToRdf(reflect.ValueOf(structure), &triples, true)
	return
}

func StructToRdfWithoutSideEffects(structure interface{}) (triples []map[string]rdf.Term, err error) {
	triples = []map[string]rdf.Term{}
	_, err = structToRdf(reflect.ValueOf(structure), &triples, false)
	return
}

func getId(structValue reflect.Value) string {
	structType := structValue.Type()
	fieldName, _ := GetFieldNameWithTag(structType, RDF_ENTITY_TAG_NAME)
	return structValue.FieldByName(fieldName).String()
}

func GetId(structure interface{}) string {
	structValue := reflect.Indirect(reflect.ValueOf(structure))
	return getId(structValue)
}
