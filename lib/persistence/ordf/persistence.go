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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
	"github.com/satori/go.uuid"
)

type Persistence struct {
	Endpoint  string
	Graph     string
	User      string
	Pw        string
	SparqlLog string
	repo      *sparql.Repo
}

func (this *Persistence) Connect() (repo *sparql.Repo, err error) {
	if this.repo == nil {
		this.repo, err = sparql.NewRepo(this.Endpoint,
			sparql.DigestAuth(this.User, this.Pw),
			sparql.Timeout(time.Second*10),
		)
	}
	return this.repo, err
}

func (this *Persistence) CreateDeleteQuery(structure interface{}) (query string, err error) {
	rdf, err := StructToRdfWithoutSideEffects(structure)
	if err != nil {
		return
	}
	turtle := Turtle(rdf)
	query, err = mustache.Render(SPARQL_DELETE, map[string]string{
		"graph":  this.Graph,
		"turtle": turtle,
	})
	return
}

func (this *Persistence) CreateUpdateQuery(old interface{}, new interface{}) (query string, err error) {

	oldRdf, err := StructToRdfWithoutSideEffects(old)
	if err != nil {
		return
	}
	_, newRdf, err := StructToRdf(new)
	if err != nil {
		return
	}

	remove, add := RdfDiff(oldRdf, newRdf)
	query, err = mustache.Render(SPARQL_UPDATE, map[string]string{
		"graph":  this.Graph,
		"add":    Turtle(add),
		"remove": Turtle(remove),
	})

	return
}

func (this *Persistence) CreateInsertQuery(structure interface{}) (query string, err error) {
	_, rdf, err := StructToRdf(structure)
	if err != nil {
		return
	}
	turtle := Turtle(rdf)
	query, err = mustache.Render(SPARQL_INSERT_TURTLE, map[string]string{
		"graph":  this.Graph,
		"turtle": turtle,
	})
	return
}

func (this *Persistence) CreateSelectQuery(id string) (query string, err error) {
	query, err = mustache.Render(SPARQL_SELECT, map[string]string{
		"graph": this.Graph,
		"id":    id,
	})
	return
}

func (this *Persistence) CreateSelectDeepQuery(id string) (query string, err error) {
	query, err = mustache.Render(SPARQL_DEEP, map[string]string{
		"graph": this.Graph,
		"id":    id,
	})
	return
}

func (this *Persistence) CreateListQuery(rdfTypeName string, limit int, offset int) (query string, err error) {
	query, err = mustache.Render(SPARQL_LIST, map[string]interface{}{
		"graph":  this.Graph,
		"type":   rdfTypeName,
		"limit":  limit,
		"offset": offset,
	})
	return
}

func (this *Persistence) CreateSearchQuery(idSymbol rdf.Term, triple []map[string]rdf.Term, limit int, offset int) (query string, err error) {
	turtle := Turtle(triple)
	query, err = mustache.Render(SPARQL_SEARCH, map[string]interface{}{
		"graph":  this.Graph,
		"fields": turtle,
		"limit":  limit,
		"offset": offset,
		"symbol": idSymbol.Serialize(rdf.Turtle),
	})
	return
}

func (this *Persistence) CreateVariantSearchQuery(mainSymbol string, maintriple []map[string]rdf.Term, variants map[string][]map[string]rdf.Term, limit int, offset int) (query string, err error) {
	mainturtle := Turtle(maintriple)

	variantContext := []map[string]interface{}{}
	first := true
	for symbol, variant := range variants {
		context := map[string]interface{}{
			"symbol": symbol,
			"fields": Turtle(variant),
		}
		if first {
			first = false
			context["first"] = true
		}
		variantContext = append(variantContext, context)
	}

	query, err = mustache.Render(SPARQL_SEARCH_VARIATION, map[string]interface{}{
		"graph":      this.Graph,
		"mainfields": mainturtle,
		"limit":      limit,
		"offset":     offset,
		"mainsymbol": mainSymbol,
		"variants":   variantContext,
	})
	return
}

func (this *Persistence) CreateVariantSearchAllQuery(mainSymbol string, maintriple []map[string]rdf.Term, variants map[string][]map[string]rdf.Term) (query string, err error) {
	mainturtle := Turtle(maintriple)

	variantContext := []map[string]interface{}{}
	first := true
	for symbol, variant := range variants {
		context := map[string]interface{}{
			"symbol": symbol,
			"fields": Turtle(variant),
		}
		if first {
			first = false
			context["first"] = true
		}
		variantContext = append(variantContext, context)
	}

	query, err = mustache.Render(SPARQL_SEARCH_ALL_VARIATION, map[string]interface{}{
		"graph":      this.Graph,
		"mainfields": mainturtle,
		"mainsymbol": mainSymbol,
		"variants":   variantContext,
	})
	return
}

func (this *Persistence) CreateSearchAllQuery(idSymbol rdf.Term, triple []map[string]rdf.Term) (query string, err error) {
	turtle := Turtle(triple)
	query, err = mustache.Render(SPARQL_SEARCH_ALL, map[string]interface{}{
		"graph":  this.Graph,
		"fields": turtle,
		"symbol": idSymbol.Serialize(rdf.Turtle),
	})
	return
}

func (this *Persistence) CreateTextVariantSearchQuery(entityType string, mainSymbol string, textTriples []map[string]rdf.Term, variants map[string][]map[string]rdf.Term, limit int, offset int) (query string, err error) {
	textFields := TextSearchTriples(textTriples)
	variantContext := []map[string]interface{}{}
	first := true
	for symbol, variant := range variants {
		context := map[string]interface{}{
			"symbol": symbol,
			"fields": Turtle(variant),
		}
		if first {
			first = false
			context["first"] = true
		}
		variantContext = append(variantContext, context)
	}
	query, err = mustache.Render(SPARQL_TEXT_SEARCH_VARIATION, map[string]interface{}{
		"graph":      this.Graph,
		"limit":      limit,
		"offset":     offset,
		"textFields": textFields,
		"mainsymbol": mainSymbol,
		"variants":   variantContext,
		"type":       entityType,
	})
	return
}

func (this *Persistence) CreateTextSearchQuery(entityType string, idSymbol rdf.Term, textTriples []map[string]rdf.Term, limit int, offset int) (query string, err error) {
	textFields := TextSearchTriples(textTriples)
	query, err = mustache.Render(SPARQL_TEXT_SEARCH, map[string]interface{}{
		"graph":      this.Graph,
		"limit":      limit,
		"offset":     offset,
		"textFields": textFields,
		"symbol":     idSymbol.Serialize(rdf.Turtle),
		"type":       entityType,
	})
	return
}

func (this *Persistence) CreateIdExistsQuery(id string) (query string, err error) {
	query, err = mustache.Render(SPARQL_ID_EXISTS, map[string]interface{}{
		"graph": this.Graph,
		"id":    id,
	})
	return
}

func (this *Persistence) CreateIdIsOfClassQuery(id string, entityType string) (query string, err error) {
	query, err = mustache.Render(SPARQL_ID_IS_OF_CLASS, map[string]interface{}{
		"graph": this.Graph,
		"id":    id,
		"type":  entityType,
	})
	return
}

func (this *Persistence) Insert(structure interface{}) (results []map[string]rdf.Term, err error) {
	query, err := this.CreateInsertQuery(structure)
	if err != nil {
		return
	}
	resp, err := this.Request(query)
	if err == nil {
		results = resp.Solutions()
	}
	return
}

func (this *Persistence) Delete(structure interface{}) (results []map[string]rdf.Term, err error) {
	query, err := this.CreateDeleteQuery(structure)
	if err != nil {
		return
	}
	resp, err := this.Request(query)
	if err == nil {
		results = resp.Solutions()
	}
	return
}

func (this *Persistence) Update(old interface{}, new interface{}) (results []map[string]rdf.Term, err error) {
	query, err := this.CreateUpdateQuery(old, new)
	if err != nil {
		return
	}
	resp, err := this.Request(query)
	if err == nil {
		results = resp.Solutions()
	}
	return
}

func (this *Persistence) selectDeep(structureValue reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			err = errors.New(fmt.Sprint(r))
		}
	}()

	id := getId(structureValue)
	query, err := this.CreateSelectDeepQuery(id)
	if err != nil {
		return
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	rdfToStruct(&structureValue, id, groupBySubject(resp.Solutions()))
	return
}

func (this *Persistence) selectValue(structureValue reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			err = errors.New(fmt.Sprint(r))
		}
	}()

	id := getId(structureValue)
	query, err := this.CreateSelectQuery(id)
	if err != nil {
		return
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	rdfToStruct(&structureValue, id, groupBySubject(resp.Solutions()))
	return
}

func (this *Persistence) Select(structure interface{}) (err error) {
	return this.selectValue(reflect.Indirect(reflect.ValueOf(structure)))
}

func (this *Persistence) SelectDeep(structure interface{}) (err error) {
	return this.selectDeep(reflect.Indirect(reflect.ValueOf(structure)))
}

func (this *Persistence) selectLevel(structValue reflect.Value, lvl int) (err error) {
	if lvl != 0 {
		err := this.selectValue(structValue)
		if err != nil {
			return err
		}
		_, fieldNameToTagValue := getFieldMapping(structValue.Type(), RDF_FIELD_TAG_NAME)
		for fieldName, tagValue := range fieldNameToTagValue {
			if tagValue != "" {
				fieldValue := structValue.FieldByName(fieldName)
				if fieldValue.Type().Kind() == reflect.Struct && !isEmpty(fieldValue) {
					err = this.selectLevel(fieldValue, lvl-1)
					if err != nil {
						return err
					}
				} else if fieldValue.Type().Kind() == reflect.Slice && !isEmpty(fieldValue) && fieldValue.Type().Elem().Kind() == reflect.Struct {
					for index := 0; index < fieldValue.Len(); index++ {
						element := fieldValue.Index(index)
						if element.Type().Kind() == reflect.Struct && !isEmpty(element) {
							err = this.selectLevel(element, lvl-1)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return
}

func (this *Persistence) SelectLevel(structure interface{}, lvl int) (err error) {
	return this.selectLevel(reflect.Indirect(reflect.ValueOf(structure)), lvl)
}

func (this *Persistence) List(structure interface{}, limit int, offset int) (err error) {
	elementType := reflect.Indirect(reflect.ValueOf(structure)).Type().Elem()
	_, rdfTypeName := GetFieldNameWithTag(elementType, RDF_ENTITY_TAG_NAME)
	query, err := this.CreateListQuery(rdfTypeName, limit, offset)
	if err != nil {
		return
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	RdfToStructList(structure, resp.Solutions())
	return
}

func (this *Persistence) Search(resultList interface{}, queryStruct interface{}, limit int, offset int) (err error) {
	id, triple, err := StructToRdf(queryStruct)
	if err != nil {
		return err
	}
	triple = filterNonQueryTrible(triple)
	query, err := this.CreateSearchQuery(id, triple, limit, offset)
	if err != nil {
		return err
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	RdfToStructList(resultList, resp.Solutions())
	return
}

func (this *Persistence) VariationSearch(resultList interface{}, limit int, offset int, mandatory interface{}, andOneOf ...interface{}) (err error) {
	mainId, mainTriples, err := StructToRdf(mandatory)
	if err != nil {
		return err
	}
	mainTriples = filterNonQueryTrible(mainTriples)

	variants := map[string][]map[string]rdf.Term{}
	for _, one := range andOneOf {
		id, triple, err := StructToRdf(one)
		if err != nil {
			return err
		}
		triple = filterNonQueryTrible(triple)
		variants[id.Serialize(rdf.Turtle)] = triple
	}

	query, err := this.CreateVariantSearchQuery(mainId.Serialize(rdf.Turtle), mainTriples, variants, limit, offset)
	if err != nil {
		return err
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	RdfToStructList(resultList, resp.Solutions())
	return
}

func (this *Persistence) VariationSearchAll(resultList interface{}, mandatory interface{}, andOneOf ...interface{}) (err error) {
	mainId, mainTriples, err := StructToRdf(mandatory)
	if err != nil {
		return err
	}
	mainTriples = filterNonQueryTrible(mainTriples)

	variants := map[string][]map[string]rdf.Term{}
	for _, one := range andOneOf {
		id, triple, err := StructToRdf(one)
		if err != nil {
			return err
		}
		triple = filterNonQueryTrible(triple)
		variants[id.Serialize(rdf.Turtle)] = triple
	}

	query, err := this.CreateVariantSearchAllQuery(mainId.Serialize(rdf.Turtle), mainTriples, variants)
	if err != nil {
		return err
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	RdfToStructList(resultList, resp.Solutions())
	return
}

func (this *Persistence) VariationTextSearch(resultList interface{}, limit int, offset int, queryStruct interface{}, andOneOf ...interface{}) (err error) {
	id, queryTriples, err := StructToRdf(queryStruct)
	if err != nil {
		return err
	}
	variants := map[string][]map[string]rdf.Term{}
	for _, one := range andOneOf {
		id, triple, err := StructToRdf(one)
		if err != nil {
			return err
		}
		triple = filterNonQueryTrible(triple)
		variants[id.Serialize(rdf.Turtle)] = triple
	}

	valueType := reflect.ValueOf(queryStruct).Type()
	_, entityType := GetFieldNameWithTag(valueType, RDF_ENTITY_TAG_NAME)

	query, err := this.CreateTextVariantSearchQuery(entityType, id.Serialize(rdf.Turtle), queryTriples, variants, limit, offset)
	if err != nil {
		return err
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	RdfToStructList(resultList, resp.Solutions())
	return
	err = errors.New("not implemented")
	return
}

func filterNonQueryTrible(terms []map[string]rdf.Term) (result []map[string]rdf.Term) {
	for _, trible := range terms {
		if trible["s"].Type() == TermSymbol || trible["p"].Type() == TermSymbol || trible["o"].Type() == TermSymbol {
			result = append(result, trible)
		}
	}
	return
}

func (this *Persistence) SearchAll(resultList interface{}, queryStruct interface{}) (err error) {
	id, triple, err := StructToRdf(queryStruct)
	if err != nil {
		return err
	}
	triple = filterNonQueryTrible(triple)
	query, err := this.CreateSearchAllQuery(id, triple)
	if err != nil {
		return err
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	RdfToStructList(resultList, resp.Solutions())
	return
}

func (this *Persistence) SearchText(resultList interface{}, queryStruct interface{}, limit int, offset int) (err error) {
	id, queryTriples, err := StructToRdf(queryStruct)
	if err != nil {
		return err
	}
	valueType := reflect.ValueOf(queryStruct).Type()
	_, entityType := GetFieldNameWithTag(valueType, RDF_ENTITY_TAG_NAME)
	query, err := this.CreateTextSearchQuery(entityType, id, queryTriples, limit, offset)
	if err != nil {
		return err
	}
	resp, err := this.Request(query)
	if err != nil {
		return
	}
	RdfToStructList(resultList, resp.Solutions())
	return
}

func (this *Persistence) IdExists(id string) (result bool, err error) {
	query, err := this.CreateIdExistsQuery(id)
	if err != nil {
		return false, err
	}
	return this.Ask(query)
}

func (this *Persistence) IdIsOfClass(structure interface{}) (result bool, err error) {
	value := reflect.ValueOf(structure)
	valueType := value.Type()
	fieldName, entityType := GetFieldNameWithTag(valueType, RDF_ENTITY_TAG_NAME)
	id := value.FieldByName(fieldName).String()
	query, err := this.CreateIdIsOfClassQuery(id, entityType)
	if err != nil {
		return false, err
	}
	return this.Ask(query)
}

func (this *Persistence) Ask(query string) (result bool, err error) {
	form := url.Values{}
	form.Set("query", query)
	b := form.Encode()

	req, err := http.NewRequest(
		"POST",
		this.Endpoint,
		bytes.NewBufferString(b))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(b)))
	req.Header.Set("Accept", "application/sparql-results+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		var msg string
		if err != nil {
			msg = "Failed to read response body"
		} else {
			if strings.TrimSpace(string(b)) != "" {
				msg = "Response body: \n" + string(b)
			}
		}
		return false, fmt.Errorf("Query: SPARQL request failed: %s. "+msg, resp.Status)
	}

	type BoolWraper struct {
		Boolean bool `json:"boolean"`
	}
	resultWrapper := &BoolWraper{}
	json.NewDecoder(resp.Body).Decode(&resultWrapper)
	return resultWrapper.Boolean, err
}

func (this *Persistence) setId(structValue reflect.Value) (err error) {
	idFieldName, _ := GetFieldNameWithTag(structValue.Type(), RDF_ENTITY_TAG_NAME)
	id := structValue.FieldByName(idFieldName).String()
	exists := id == ""
	if !isEmpty(structValue) {
		for exists && err == nil {
			id = this.formatId(uuid.NewV4().String())
			exists, err = this.IdExists(id)
		}
		if err == nil {
			structValue.FieldByName(idFieldName).SetString(id)
		}
	}
	return
}

func (this *Persistence) SetId(entity interface{}) (err error) {
	value := reflect.Indirect(reflect.ValueOf(entity))
	return this.setId(value)
}

func (this *Persistence) setIdDeep(structValue reflect.Value) (err error) {
	this.setId(structValue)
	mapToFieldName, _ := getFieldMapping(structValue.Type(), RDF_FIELD_TAG_NAME)
	for _, fieldName := range mapToFieldName {
		field := structValue.FieldByName(fieldName)
		if field.Kind() == reflect.Struct {
			err = this.setIdDeep(field)
			if err != nil {
				break
			}
		} else if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
			for i := 0; i < field.Len(); i++ {
				this.setIdDeep(field.Index(i))
			}
		}
	}
	return
}

func (this *Persistence) SetIdDeep(entity interface{}) (err error) {
	value := reflect.Indirect(reflect.ValueOf(entity))
	return this.setIdDeep(value)
}

func (this *Persistence) Request(query string) (result *sparql.Results, err error) {
	if this.SparqlLog == "true" {
		log.Println(this.Endpoint, query)
	}
	repo, err := this.Connect()
	if err != nil {
		return result, err
	}
	result, err = repo.Query(query)
	if this.SparqlLog == "true" {
		if err != nil {
			log.Println(err)
		} else {
			log.Println(result.Solutions())
		}
	}
	return
}
