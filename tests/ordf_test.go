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

package tests

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/SmartEnergyPlatform/iot-device-repository/lib/persistence/ordf"

	"strings"

	"github.com/knakk/rdf"
)

type ClassB struct {
	Id      string `json:"id"          rdf_entity:"TypeB"`
	FieldB1 string `json:"field_1"     rdf_field:"bf_1"`
}

type ClassA struct {
	Id     string   `json:"id"              rdf_entity:"TypeA"`
	Field1 string   `json:"field_1"         rdf_field:"af_1"`
	Field2 []string `json:"field_2"         rdf_field:"af_2"`
	Field3 ClassB   `json:"field_3_abcasd"  rdf_field:"af_3"`
	Field4 []ClassB `json:"field_4"         rdf_field:"af_4"`
}

type ClassC struct {
	Id     string `json:"id"                              rdf_entity:"TypeA"`
	Field1 string `json:"field_1"     rdf_field:"af_3"    rdf_ref:"true"`
}

func createLiteral(value string) (result rdf.Literal) {
	result, _ = rdf.NewLiteral(value)
	return
}

func createIRI(iri string) (result rdf.IRI) {
	result, _ = rdf.NewIRI(iri)
	return
}

func getRdf() []map[string]rdf.Term {
	return []map[string]rdf.Term{
		//ClassB
		{
			"s": createIRI("B_1"),
			"p": createIRI("bf_1"),
			"o": createLiteral("b1"),
		}, {
			"s": createIRI("B_1"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeB"),
		}, {
			"s": createIRI("B_2"),
			"p": createIRI("bf_1"),
			"o": createLiteral("b2"),
		}, {
			"s": createIRI("B_2"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeB"),
		}, {
			"s": createIRI("B_3"),
			"p": createIRI("bf_1"),
			"o": createLiteral("b3"),
		}, {
			"s": createIRI("B_3"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeB"),
		},

		//ClassA
		{
			"s": createIRI("A_1"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeA"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_1"),
			"o": createLiteral("field_1_string"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_2"),
			"o": createLiteral("field_2_string_1"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_2"),
			"o": createLiteral("field_2_string_2"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_3"),
			"o": createIRI("B_1"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_4"),
			"o": createIRI("B_1"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_4"),
			"o": createIRI("B_2"),
		}, {
			"s": createIRI("A_2"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeA"),
		}, {
			"s": createIRI("A_2"),
			"p": createIRI("af_1"),
			"o": createLiteral("field_1_string_a2_1"),
		}, {
			"s": createIRI("A_2"),
			"p": createIRI("af_2"),
			"o": createLiteral("field_2_string_a2_1"),
		}, {
			"s": createIRI("A_2"),
			"p": createIRI("af_2"),
			"o": createLiteral("field_2_string_a2_2"),
		}, {
			"s": createIRI("A_2"),
			"p": createIRI("af_3"),
			"o": createIRI("B_3"),
		}, {
			"s": createIRI("A_2"),
			"p": createIRI("af_4"),
			"o": createIRI("B_2"),
		}, {
			"s": createIRI("A_2"),
			"p": createIRI("af_4"),
			"o": createIRI("B_3"),
		},
	}
}

func getSparseRdf() []map[string]rdf.Term {
	return []map[string]rdf.Term{

		//ClassA
		{
			"s": createIRI("A_1"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeA"),
		},
	}
}

func rdfContainsTriple(triples []map[string]rdf.Term, triple map[string]rdf.Term) bool {
	for _, element := range triples {
		if reflect.DeepEqual(element, triple) {
			return true
		}
	}
	return false
}

func equalRdf(a []map[string]rdf.Term, b []map[string]rdf.Term) bool {
	for _, triple := range a {
		if !rdfContainsTriple(b, triple) {
			return false
		}
	}
	for _, triple := range b {
		if !rdfContainsTriple(a, triple) {
			return false
		}
	}
	return true
}

func TestStruct(t *testing.T) {
	elemA := ClassA{}
	ordf.RdfToStruct(&elemA, "A_1", getRdf())
	_, newRdf, err := ordf.StructToRdf(elemA)
	newRdf2, err2 := ordf.StructToRdfWithoutSideEffects(elemA)

	targetRdf := []map[string]rdf.Term{
		{
			"s": createIRI("B_1"),
			"p": createIRI("bf_1"),
			"o": createLiteral("b1"),
		}, {
			"s": createIRI("B_1"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeB"),
		}, {
			"s": createIRI("B_2"),
			"p": createIRI("bf_1"),
			"o": createLiteral("b2"),
		}, {
			"s": createIRI("B_2"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeB"),
		},
		//ClassA
		{
			"s": createIRI("A_1"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeA"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_1"),
			"o": createLiteral("field_1_string"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_2"),
			"o": createLiteral("field_2_string_1"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_2"),
			"o": createLiteral("field_2_string_2"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_3"),
			"o": createIRI("B_1"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_4"),
			"o": createIRI("B_1"),
		}, {
			"s": createIRI("A_1"),
			"p": createIRI("af_4"),
			"o": createIRI("B_2"),
		},
	}

	if err != nil {
		t.Error(err)
	}

	if !equalRdf(newRdf, targetRdf) {
		t.Error("want: ", targetRdf, " got: ", newRdf)
	}

	if err2 != nil {
		t.Error(err2)
	}

	if !equalRdf(newRdf2, targetRdf) {
		t.Error("want: ", targetRdf, " got: ", newRdf2)
	}
}

type ById []ClassA

func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }

func ExampleRdfToStruct() {
	sliceA := []ClassA{}
	ordf.RdfToStructList(&sliceA, getRdf())
	sort.Sort(ById(sliceA))
	fmt.Println(sliceA)
	// Output:
	//[{A_1 field_1_string [field_2_string_1 field_2_string_2] {B_1 b1} [{B_1 b1} {B_2 b2}]} {A_2 field_1_string_a2_1 [field_2_string_a2_1 field_2_string_a2_2] {B_3 b3} [{B_2 b2} {B_3 b3}]}]
}

func ExampleSparseRdfToStruct() {
	elemA := ClassA{}
	ordf.RdfToStruct(&elemA, "A_1", getSparseRdf())
	fmt.Println(elemA)
	// Output:
	//{A_1  [] { } []}
}

func TestSparseStructToRdf(t *testing.T) {
	elemA := ClassA{}
	ordf.RdfToStruct(&elemA, "A_1", getSparseRdf())
	_, newRdf, err := ordf.StructToRdf(elemA)
	newRdf2, err2 := ordf.StructToRdfWithoutSideEffects(elemA)
	if err != nil {
		t.Error(err)
	}
	if err2 != nil {
		t.Error(err2)
	}

	targetRdf := []map[string]rdf.Term{
		{
			"s": createIRI("A_1"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeA"),
		},
	}

	if !equalRdf(newRdf, targetRdf) {
		t.Error("want: ", targetRdf, " got: ", newRdf)
	}
	if !equalRdf(newRdf2, targetRdf) {
		t.Error("want: ", targetRdf, " got: ", newRdf2)
	}
}

func TestRefStructToRdf(t *testing.T) {
	elem := ClassC{Id: "C", Field1: "B_1"}
	_, newRdf, err := ordf.StructToRdf(elem)
	newRdf2, err2 := ordf.StructToRdfWithoutSideEffects(elem)
	if err != nil {
		t.Error(err)
	}
	if err2 != nil {
		t.Error(err2)
	}

	targetRdf := []map[string]rdf.Term{
		{
			"s": createIRI("C"),
			"p": createIRI(ordf.RDF_TYPE_PREDICATE_NAME),
			"o": createIRI("TypeA"),
		}, {
			"s": createIRI("C"),
			"p": createIRI("af_3"),
			"o": createIRI("B_1"),
		},
	}

	if !equalRdf(newRdf, targetRdf) {
		t.Error("different rdf (ckeck type of o:B_1) \n want: ", targetRdf, " \n got: ", newRdf)
	}
	if !equalRdf(newRdf2, targetRdf) {
		t.Error("different rdf (ckeck type of o:B_1) \n want: ", targetRdf, " \n got: ", newRdf2)
	}

	elemA := ClassA{}
	elemA2 := ClassA{}
	elemA3 := ClassA{Id: "C", Field3: ClassB{Id: "B_1"}}
	ordf.RdfToStruct(&elemA, "C", newRdf)
	ordf.RdfToStruct(&elemA2, "C", newRdf2)

	if !reflect.DeepEqual(elemA, elemA3) {
		t.Error("different struct \n want: ", elemA3, " \n got: ", elemA)
	}

	if !reflect.DeepEqual(elemA2, elemA3) {
		t.Error("different struct \n want: ", elemA3, " \n got: ", elemA2)
	}

}

func TestInsert(t *testing.T) {
	elemA := ClassA{}
	ordf.RdfToStruct(&elemA, "A_1", getRdf())

	db := &ordf.Persistence{Endpoint: "http://localhost:8890/sparql", Graph: "test", User: "dba", Pw: "myDbaPassword"}
	_, err := db.Insert(elemA)
	if err != nil {
		t.Error(err)
	}

	elemB := ClassA{Id: "A_1"}
	err = db.SelectLevel(&elemB, -1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(elemA, elemB) {
		t.Error("different struct \n want: ", elemA, " \n got: ", elemB)
	}
}

func ExampleList() {
	db := &ordf.Persistence{Endpoint: "http://localhost:8890/sparql", Graph: "test", User: "dba", Pw: "myDbaPassword"}
	elements := []ClassB{}
	db.List(&elements, 2, 0)
	fmt.Println(elements)

	elements = []ClassB{}
	db.List(&elements, 1, 0)
	fmt.Println(elements)

	elements = []ClassB{}
	db.List(&elements, 1, 1)
	fmt.Println(elements)
	//Output:
	//[{B_1 b1} {B_2 b2}]
	//[{B_1 b1}]
	//[{B_2 b2}]
}

func ExampleSearch() {
	db := &ordf.Persistence{Endpoint: "http://localhost:8890/sparql", Graph: "test", User: "dba", Pw: "myDbaPassword"}
	result := []ClassA{}
	db.Search(&result, ClassA{Field1: "field_1_string", Field2: []string{"field_2_string_1"}}, 10, 0)
	fmt.Println(result)
	//Output:
	//[{A_1 field_1_string [field_2_string_1 field_2_string_2] {B_1 } [{B_1 } {B_2 }]}]
}

func ExampleIdExists() {
	db := &ordf.Persistence{Endpoint: "http://localhost:8890/sparql", Graph: "test", User: "dba", Pw: "myDbaPassword"}
	fmt.Println(db.IdExists("A_1"))
	fmt.Println(db.IdExists("X_1"))
	//Output:
	//true <nil>
	//false <nil>
}

func TestSetId(t *testing.T) {
	assert := Assertions(t)

	withoutId := ClassA{Field1: "abc", Field3: ClassB{FieldB1: "def"}}
	aid := ClassA{Id: "a", Field1: "abc", Field3: ClassB{FieldB1: "def"}}
	bid := ClassA{Field1: "abc", Field3: ClassB{FieldB1: "def", Id: "b"}}
	fullId := ClassA{Id: "a", Field1: "abc", Field3: ClassB{FieldB1: "def", Id: "b"}}

	list := ClassA{Id: "a", Field1: "abc", Field3: ClassB{FieldB1: "def", Id: "b"}, Field4: []ClassB{{FieldB1: "def", Id: "b2"}, {FieldB1: "def"}}}

	db := &ordf.Persistence{Endpoint: "http://localhost:8890/sparql", Graph: "test", User: "dba", Pw: "myDbaPassword"}

	db.SetIdDeep(&withoutId)
	db.SetIdDeep(&aid)
	db.SetIdDeep(&bid)
	db.SetIdDeep(&fullId)
	db.SetIdDeep(&list)

	assert.True(withoutId.Id != "", "missing a id")
	assert.Equal(withoutId.Field1, "abc")
	assert.Equal(withoutId.Field3.FieldB1, "def")
	assert.True(withoutId.Field3.Id != "", "missing b id")

	assert.Equal(aid.Id, "a")
	assert.Equal(aid.Field1, "abc")
	assert.Equal(aid.Field3.FieldB1, "def")
	assert.True(aid.Field3.Id != "", "missing b id")

	assert.True(bid.Id != "", "missing a id")
	assert.Equal(bid.Field1, "abc")
	assert.Equal(bid.Field3.FieldB1, "def")
	assert.Equal(bid.Field3.Id, "b")

	assert.Equal(fullId.Id, "a")
	assert.Equal(fullId.Field1, "abc")
	assert.Equal(fullId.Field3.FieldB1, "def")
	assert.Equal(fullId.Field3.Id, "b")

	assert.Equal(list.Id, "a")
	assert.Equal(list.Field1, "abc")
	assert.Equal(list.Field3.FieldB1, "def")
	assert.Equal(list.Field3.Id, "b")
	assert.Equal(list.Field4[0].Id, "b2")
	assert.Equal(list.Field4[0].FieldB1, "def")
	assert.Equal(list.Field4[1].FieldB1, "def")
	assert.UnEqual(list.Field4[1].Id, "")
}

func diffPrint(old []map[string]rdf.Term, new []map[string]rdf.Term) {
	remove, add := ordf.RdfDiff(old, new)
	removeTriples := ordf.TurtleList(remove)
	addTriples := ordf.TurtleList(add)
	sort.Strings(removeTriples)
	sort.Strings(addTriples)
	fmt.Println(removeTriples, addTriples)
}

func triplesPrint(triples []map[string]rdf.Term) {
	turtles := ordf.TurtleList(triples)
	sort.Strings(turtles)
	for _, turtle := range turtles {
		fmt.Println(turtle)
	}
}

func ExampleDiff() {
	_, a1, _ := ordf.StructToRdf(ClassA{Id: "1"})
	_, a1_1, _ := ordf.StructToRdf(ClassA{Id: "1", Field1: "abc"})

	_, a1_b1, _ := ordf.StructToRdf(ClassA{Id: "1", Field1: "abc", Field3: ClassB{Id: "2", FieldB1: "abc"}})
	_, a1_b2, _ := ordf.StructToRdf(ClassA{Id: "1", Field1: "abc", Field3: ClassB{Id: "2", FieldB1: "abc2"}})

	_, a1_b3, _ := ordf.StructToRdf(ClassA{Id: "1", Field3: ClassB{Id: "2", FieldB1: "abc2"}})

	diffPrint(a1, a1_1)
	diffPrint(a1_1, a1)

	diffPrint(a1_b1, a1_b2)
	diffPrint(a1_b2, a1_b1)

	diffPrint(a1_1, a1_b3)
	diffPrint(a1_b3, a1_1)

	//Output:
	//[] [<1> <af_1> "abc" .]
	//[<1> <af_1> "abc" .] []
	//[<2> <bf_1> "abc" .] [<2> <bf_1> "abc2" .]
	//[<2> <bf_1> "abc2" .] [<2> <bf_1> "abc" .]
	//[<1> <af_1> "abc" .] [<1> <af_3> <2> . <2> <bf_1> "abc2" . <2> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <TypeB> .]
	//[<1> <af_3> <2> . <2> <bf_1> "abc2" . <2> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <TypeB> .] [<1> <af_1> "abc" .]

}

type B struct {
	Id     string `rdf_entity:"EntityB"`
	Fieldb string `rdf_field:"fbb"`
	FieldC C      `rdf_field:"fbc"`
	SliceC []C    `rdf_field:"fbcs"`
}

type C struct {
	Id     string `rdf_entity:"RootC" rdf_root:"true"`
	Fieldc string `rdf_field:"fcc"`
}

type A struct {
	Id     string `rdf_entity:"RootA" rdf_root:"true"`
	Fielda string `rdf_field:"faa"`
	FieldB B      `rdf_field:"fab"`
	FieldC C      `rdf_field:"fac"`
}

func ExampleUpdateSparql() {
	orig := A{
		Id:     "A",
		Fielda: "fielda",
		FieldC: C{
			Id:     "A_C",
			Fieldc: "Fieldc_A_C",
		},
		FieldB: B{
			Id:     "B",
			Fieldb: "fieldb",
			FieldC: C{
				Id:     "B_C",
				Fieldc: "Fieldc_B_C",
			},
			SliceC: []C{
				{
					Id:     "B_SC1",
					Fieldc: "Fieldc_B_SC1",
				},
				{
					Id:     "B_SC2",
					Fieldc: "Fieldc_B_SC2",
				},
				{
					Id:     "B_SC3",
					Fieldc: "Fieldc_B_SC3",
				},
			},
		},
	}

	changeFa := orig
	changeFa.Fielda = "changedFielda"

	changeFb := orig
	changeFb.FieldB.Fieldb = "changedFieldb"

	changeFc := orig
	changeFc.FieldB.FieldC.Fieldc = "changedFieldc" //should not be in update or delete sparql

	changeC := orig
	changeC.FieldB.FieldC = C{
		Id:     "Changed_B_C",
		Fieldc: "changed_field_in_c", //should not be in update or delete sparql
	}

	changeCS := orig
	changeCS.FieldB.SliceC = []C{
		{
			Id:     "Changed_B_CS_1",
			Fieldc: "changed_field_in_cs1", //should not be in update or delete sparql
		}, {
			Id:     "Changed_B_CS_2",
			Fieldc: "changed_field_in_cs2", //should not be in update or delete sparql
		},
	}

	db := &ordf.Persistence{Endpoint: "http://localhost:8890/sparql", Graph: "test", User: "dba", Pw: "myDbaPassword"}

	noWhiteSpace := func(str string, err error) (result string) {
		result = strings.Join(strings.Fields(str), "")
		if err != nil {
			result += err.Error()
		}
		return
	}

	fmt.Println("===========changeFa===========")
	fmt.Println(noWhiteSpace(db.CreateUpdateQuery(orig, changeFa)))

	fmt.Println("===========changeFb===========")
	fmt.Println(noWhiteSpace(db.CreateUpdateQuery(orig, changeFb)))

	fmt.Println("===========changeFc===========")
	fmt.Println(noWhiteSpace(db.CreateUpdateQuery(orig, changeFc)))

	fmt.Println("===========changeC===========")
	fmt.Println(noWhiteSpace(db.CreateUpdateQuery(orig, changeC)))

	fmt.Println("===========changeCS===========")
	fmt.Println(noWhiteSpace(db.CreateUpdateQuery(orig, changeCS)))

	fmt.Println("===========remove_A===========")
	fmt.Println(noWhiteSpace(db.CreateDeleteQuery(orig)))

	//Output:
	//===========changeFa===========
	//DELETEDATAFROM<test>{<A><faa>"fielda".};INSERTINTO<test>{<A><faa>"changedFielda".}
	//===========changeFb===========
	//DELETEDATAFROM<test>{<B><fbb>"fieldb".};INSERTINTO<test>{<B><fbb>"changedFieldb".}
	//===========changeFc===========
	//DELETEDATAFROM<test>{};INSERTINTO<test>{}
	//===========changeC===========
	//DELETEDATAFROM<test>{<B><fbc><B_C>.};INSERTINTO<test>{<B><fbc><Changed_B_C>.}
	//===========changeCS===========
	//DELETEDATAFROM<test>{<B><fbcs><B_SC1>.<B><fbcs><B_SC2>.<B><fbcs><B_SC3>.};INSERTINTO<test>{<B><fbcs><Changed_B_CS_1>.<B><fbcs><Changed_B_CS_2>.}
	//===========remove_A===========
	//DELETEDATAFROM<test>{<A><http://www.w3.org/1999/02/22-rdf-syntax-ns#type><RootA>.<A><faa>"fielda".<B><http://www.w3.org/1999/02/22-rdf-syntax-ns#type><EntityB>.<B><fbb>"fieldb".<B><fbc><B_C>.<B><fbcs><B_SC1>.<B><fbcs><B_SC2>.<B><fbcs><B_SC3>.<A><fab><B>.<A><fac><A_C>.}

}

func ExampleStructToRdfWithRoot() {
	a := A{
		Id:     "A",
		Fielda: "fielda",
		FieldC: C{
			Id:     "A_C",
			Fieldc: "Fieldc_A_C",
		},
		FieldB: B{
			Id:     "B",
			Fieldb: "fieldb",
			FieldC: C{
				Id:     "B_C",
				Fieldc: "Fieldc_B_C",
			},
			SliceC: []C{
				{
					Id:     "B_SC1",
					Fieldc: "Fieldc_B_SC1",
				},
				{
					Id:     "B_SC2",
					Fieldc: "Fieldc_B_SC2",
				},
				{
					Id:     "B_SC3",
					Fieldc: "Fieldc_B_SC3",
				},
			},
		},
	}
	triples, err := ordf.StructToRdfWithoutSideEffects(a)
	if err != nil {
		panic(err)
	}
	triplesPrint(triples)
	fmt.Println("==========================")
	_, triples2, err2 := ordf.StructToRdf(a)
	if err2 != nil {
		panic(err2)
	}
	triplesPrint(triples2)

	//Output:
	//<A> <faa> "fielda" .
	//<A> <fab> <B> .
	//<A> <fac> <A_C> .
	//<A> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <RootA> .
	//<B> <fbb> "fieldb" .
	//<B> <fbc> <B_C> .
	//<B> <fbcs> <B_SC1> .
	//<B> <fbcs> <B_SC2> .
	//<B> <fbcs> <B_SC3> .
	//<B> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <EntityB> .
	//==========================
	//<A> <faa> "fielda" .
	//<A> <fab> <B> .
	//<A> <fac> <A_C> .
	//<A> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <RootA> .
	//<A_C> <fcc> "Fieldc_A_C" .
	//<A_C> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <RootC> .
	//<B> <fbb> "fieldb" .
	//<B> <fbc> <B_C> .
	//<B> <fbcs> <B_SC1> .
	//<B> <fbcs> <B_SC2> .
	//<B> <fbcs> <B_SC3> .
	//<B> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <EntityB> .
	//<B_C> <fcc> "Fieldc_B_C" .
	//<B_C> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <RootC> .
	//<B_SC1> <fcc> "Fieldc_B_SC1" .
	//<B_SC1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <RootC> .
	//<B_SC2> <fcc> "Fieldc_B_SC2" .
	//<B_SC2> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <RootC> .
	//<B_SC3> <fcc> "Fieldc_B_SC3" .
	//<B_SC3> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <RootC> .
}
