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

const SPARQL_UPDATE = `
DELETE DATA FROM <{{graph}}>
  {
    {{{remove}}}
  };
INSERT INTO <{{graph}}>
  {
    {{{add}}}
  }
`

const SPARQL_DELETE = `
DELETE DATA FROM <{{graph}}>
  {
    {{{turtle}}}
  }
`

const SPARQL_INSERT_TURTLE = `
INSERT INTO <{{graph}}>
  {
    {{{turtle}}}
  }
`

//find by entity id
const SPARQL_SELECT = `
SELECT <{{{id}}}> as ?s ?p ?o
FROM <{{{graph}}}>
WHERE { <{{{id}}}> ?p ?o }`

const SPARQL_DEEP = `
CONSTRUCT {
   <{{{id}}}> ?prop ?val .
   ?child ?childProp ?childPropVal .
   ?someSubj ?incomingChildProp ?child .
}
WHERE {
	GRAPH <{{{graph}}}> {
		 <{{{id}}}> ?prop ?val ;
			 (<overrides>|!<overrides>)+ ?child .
		 ?child ?childProp ?childPropVal.
		 ?someSubj ?incomingChildProp ?child.
     }
}
`

const SPARQL_LIST = `
SELECT ?s ?p ?o
FROM <{{{graph}}}>
WHERE {
    ?s ?p ?o.
    {
    	SELECT ?s WHERE {?s <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <{{{type}}}>.}
     	ORDER BY ?s
       	LIMIT {{limit}}
       	OFFSET {{offset}}
    }
}
ORDER BY ?s`

const SPARQL_SEARCH = `
SELECT ?s ?p ?o
FROM <{{{graph}}}>
WHERE {
    ?s ?p ?o.
    {
    	SELECT {{symbol}} as ?s WHERE {
    		{{{fields}}}
    	}
     	ORDER BY ?s
       	LIMIT {{limit}}
       	OFFSET {{offset}}
    }
}`

const SPARQL_SEARCH_ALL = `
SELECT ?s ?p ?o
FROM <{{{graph}}}>
WHERE {
    ?s ?p ?o.
    {
    	SELECT {{symbol}} as ?s WHERE {
    		{{{fields}}}
    	}
		ORDER BY ?s
    }
}
`

const SPARQL_SEARCH_VARIATION = `
SELECT DISTINCT ?s ?p ?o
FROM <{{{graph}}}>
WHERE {
    ?s ?p ?o.
    {
    	SELECT DISTINCT {{mainsymbol}} as ?s WHERE {
			{{{mainfields}}}
    		{{#variants}}
				{{^first}}UNION{{/first}}
				{
					SELECT {{symbol}} as {{mainsymbol}} WHERE {
						{{{fields}}}
					}
				}
			{{/variants}}
    	}
     	ORDER BY ?s
       	LIMIT {{limit}}
       	OFFSET {{offset}}
    }
}`

const SPARQL_SEARCH_ALL_VARIATION = `
SELECT DISTINCT ?s ?p ?o
FROM <{{{graph}}}>
WHERE {
    ?s ?p ?o.
    {
    	SELECT DISTINCT {{mainsymbol}} as ?s WHERE {
			{{{mainfields}}}
    		{{#variants}}
				{{^first}}UNION{{/first}}
				{
					SELECT {{symbol}} as {{mainsymbol}} WHERE {
						{{{fields}}}
					}
				}
			{{/variants}}
    	}
     	ORDER BY ?s
    }
}`

const SPARQL_TEXT_SEARCH_VARIATION = `
SELECT DISTINCT ?s ?p ?o
FROM <{{{graph}}}>
WHERE {
    ?s ?p ?o.
    {
    	SELECT DISTINCT {{mainsymbol}} as ?s WHERE {
			{{mainsymbol}} <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <{{{type}}}>.
			{
				{ {{mainsymbol}} <nope> <nope> }
				{{#textFields}}
					UNION
					{
						{{{subject}}} {{{textFieldName}}} ?o{{@index}}.
						FILTER regex(?o{{@index}}, {{{textFieldValue}}}, "i").
					}
				{{/textFields}}
			}
    		{{#variants}}
				{{^first}}UNION{{/first}}
				{
					SELECT {{symbol}} as {{mainsymbol}} WHERE {
						{{{fields}}}
					}
				}
			{{/variants}}
    	}
     	ORDER BY ?s
       	LIMIT {{limit}}
       	OFFSET {{offset}}
    }
}`

const SPARQL_TEXT_SEARCH = `
SELECT DISTINCT ?s ?p ?o
FROM <{{{graph}}}>
WHERE {
    ?s ?p ?o.
    {
    	SELECT DISTINCT {{symbol}} as ?s WHERE {
    		{{symbol}} <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <{{{type}}}>.
    		{
    			{ {{symbol}} <nope> <nope> }
    			{{#textFields}}
					UNION
					{
						{{{subject}}} {{{textFieldName}}} ?o{{@index}}.
						FILTER regex(?o{{@index}}, {{{textFieldValue}}}, "i").
					}
				{{/textFields}}
    		}
    	}
     	ORDER BY ?s
       	LIMIT {{limit}}
       	OFFSET {{offset}}
    }
}`

const SPARQL_ID_EXISTS = `
ASK
FROM <{{{graph}}}>
{ <{{{id}}}> a ?t }`

const SPARQL_ID_IS_OF_CLASS = `
ASK
FROM <{{{graph}}}>
{ <{{{id}}}> a <{{{type}}}> }`
