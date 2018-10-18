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
	"strconv"

	"github.com/knakk/rdf"
)

var count = 0

const maxCount = 10000000

const TermSymbol = rdf.TermLiteral + 1

func getNextCount() int {
	count = (count % maxCount) + 1
	return count
}

type Symbol struct {
	Id int
}

func (this Symbol) Serialize(rdf.Format) string {
	return "?s" + strconv.Itoa(this.Id)
}

func (this Symbol) String() string {
	return "?s" + strconv.Itoa(this.Id)
}

func (this Symbol) Type() rdf.TermType {
	return TermSymbol
}

func newSymbol() Symbol {
	return Symbol{Id: getNextCount()}
}
