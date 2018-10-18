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
	"log"

	"github.com/knakk/rdf"
)

func groupByPredicate(triples []map[string]rdf.Term) (result map[string][]map[string]rdf.Term) {
	result = map[string][]map[string]rdf.Term{}
	for _, triple := range triples {
		predId := triple["p"].String()
		if subjTriples, ok := result[predId]; ok {
			result[predId] = append(subjTriples, triple)
		} else {
			result[predId] = []map[string]rdf.Term{triple}
		}
	}
	return
}

func groupByObject(triples []map[string]rdf.Term) (result map[string][]map[string]rdf.Term) {
	result = map[string][]map[string]rdf.Term{}
	for _, triple := range triples {
		predId := triple["o"].String()
		if subjTriples, ok := result[predId]; ok {
			result[predId] = append(subjTriples, triple)
		} else {
			result[predId] = []map[string]rdf.Term{triple}
		}
	}
	return
}

func RdfHash(triple map[string]rdf.Term) string {
	subject := ""
	predicate := ""
	object := ""
	if s, ok := triple["s"]; ok {
		subject = s.String()
	} else {
		log.Println("RdfHash with unknown subject", triple)
	}
	if p, ok := triple["p"]; ok {
		predicate = p.String()
	} else {
		log.Println("RdfHash with unknown predicate", triple)
	}
	if o, ok := triple["o"]; ok {
		object = o.String()
	} else {
		log.Println("RdfHash with unknown object", triple)
	}
	return subject + "::" + predicate + "::" + object
}

func createRdfIndex(triples []map[string]rdf.Term) (result map[string]map[string]rdf.Term) {
	result = map[string]map[string]rdf.Term{}
	for _, triple := range triples {
		result[RdfHash(triple)] = triple
	}
	return result
}

func RdfDiff(old []map[string]rdf.Term, new []map[string]rdf.Term) (remove []map[string]rdf.Term, add []map[string]rdf.Term) {
	indexOld := createRdfIndex(old)
	indexNew := createRdfIndex(new)
	for _, triple := range new {
		_, ok := indexOld[RdfHash(triple)]
		if !ok {
			add = append(add, triple)
		}
	}

	for _, triple := range old {
		_, ok := indexNew[RdfHash(triple)]
		if !ok {
			remove = append(remove, triple)
		}
	}
	return
}
