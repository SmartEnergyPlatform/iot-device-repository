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
	"github.com/knakk/rdf"
	"fmt"
	"strings"
)

func TurtleTriple(triple map[string]rdf.Term) string{
	var s string
	subject, ok := triple["s"]
	if ok {
		s = subject.Serialize(rdf.Turtle)
	}else{
		s = "?s"
	}
	p := triple["p"].Serialize(rdf.Turtle)
	o := triple["o"].Serialize(rdf.Turtle)
	return fmt.Sprintf("%s %s %s .", s,p,o)
}

func TurtleList(triples []map[string]rdf.Term) (result []string){
	for _, triple := range triples {
		result = append(result, TurtleTriple(triple))
	}
	return
}

func Turtle(triples []map[string]rdf.Term) (result string){
	return strings.Join(TurtleList(triples), "\n")
}

func TextSearchTriples(triples []map[string]rdf.Term) (result []map[string]string) {
	for _, triple := range triples {
		if triple["o"].Type() == rdf.TermLiteral {
			var s string
			subject, ok := triple["s"]
			if ok {
				s = subject.Serialize(rdf.Turtle)
			}else{
				s = "?s"
			}
			result = append(result, map[string]string{
				"subject": s,
				"textFieldName": triple["p"].Serialize(rdf.Turtle),
				"textFieldValue": triple["o"].Serialize(rdf.Turtle),
			})
		}
	}
	return
}