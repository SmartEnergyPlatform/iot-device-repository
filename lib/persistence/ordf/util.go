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
)

const (
	RDF_FIELD_TAG_NAME = "rdf_field" //add to field to describe predicateName -> `rdf:"www.example.com/predicateUri"`

	//add to field to inform that this field describes the id of the entity
	//informs about the type of the describe entity -> `rdf_entity:"www.example.com/ServiceTypeUri"`
	RDF_ENTITY_TAG_NAME = "rdf_entity"

	//add to string or []string field to declare that the value is a id of a entity and not a literal
	// --> `rdf_ref:"true"`
	RDF_REF_TAG_NAME = "rdf_ref"

	//add to id field (field with <<RDF_ENTITY_TAG_NAME>>) to indicate that the structure is a root entity
	//is needed to prevent the entities deletion on update off a depending entity
	//effects the result of StructToRdfWithoutSideEffects() but not StructToRdf()
	// --> `rdf_entity:"http://www.example.com/User" rdf_root:"true"`
	RDF_ROOT_ENTITY_TAG_NAME = "rdf_root"

	//similar to root but on the using side
	RDF_LEND_TAG_NAME = "rdf_lending"

	RDF_TYPE_PREDICATE_NAME = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"
)

func getFieldMapping(structType reflect.Type, tagName string) (tagValueToFieldName map[string]string, fieldNameToTagValue map[string]string) {
	tagValueToFieldName = map[string]string{}
	fieldNameToTagValue = map[string]string{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" {
			tagValueToFieldName[tag] = field.Name
			fieldNameToTagValue[field.Name] = tag
		}
	}
	return
}

func GetFieldNameWithTag(structType reflect.Type, tagName string) (fieldName string, tagValue string) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get(tagName)
		if tag != "" {
			fieldName = field.Name
			tagValue = tag
			return
		}
	}
	return
}

func (this *Persistence) formatId(id string) string {
	return this.Graph + "#" + id
}
