# Cytoscape CI Service Sample Implementation for Go Developers

![](docs/cytoscape-flat-logo-orange.png) ![](docs/gopher-side_path.png)


## Introduction
This is a sample CI service implementation written in Go.  In this project, we use a simple ID mapping service as an example.


## Basic Design
_Any_ type of RESTful API server can work as a part of Cytoscape CI system.  This means if you have existing API server, 
the change required to use as a CI service is very small.

The followings are the building blocks of your CI service:

### 1. Actual Service Package

For all types of services, it is a good idea to create a single-function, command-line tool first.  You can start with 
 designing your function by defining the following:

* Function - What is the function of your service?
* Input - What are the parameters you need to pass for the tool?
* Output - What is the output data user get?


In this example, the goal is simple and clear:

* Function: For given list of gene IDs, generate mappings to other set of IDs.  e.g., Entrez Gene ID to Ensemble ID
* Input: List of IDs
* Output: List of maps from input ID to other set of IDs.

Once you decide, implement such function.  This portion of your code should work with or without RESTful server.


### 2. RESTful API Server
Once your function works, wrap it with an API server code.  In Go, there are some nice and thin frameworks for implementing 
RESTful API, like [Gin](https://gin-gonic.github.io/gin/).  You can choose any of those to do this job.  In this example, we simply use standard _net/http_ 
package to implement API server.

The requirement for this part is writing an API to consume and generate JSON. 


### 3. Registration to _Submit Agent_
When you finish writing the REST API server, now you are ready to register it to task submit agent for CI.  
Currently, we have an [Erlang](http://www.erlang.org/) implementation of Submit Agent called 
[elsa](https://github.com/cytoscape-ci/elsa).  To register your service, you can POST some required information 
to the agent.

In this example, there is a reusable module to do this task and you can use it if you just want to register single 
server.


# Sample Service Implementation

## ID Mapping Service

As a sample service for Go developers, this repository contains complete code for simple ID Mapping service, 
which loads all mapping data into memory and return result in JSON.

Currently, the original data sources are [NCBI Gene](ftp://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/) and 
[UniprotKB](ftp://ftp.uniprot.org/pub/databases/uniprot/current_release/knowledgebase/idmapping/by_organism/).


## How to Run

### 1. In single API server mode
You can run this service as a single RESTful API server independent from Cytoscape CI.

#### On your local machine

* Requirments
  * Go 1.5.x and later
  * Python 3.5.x or newer
    * pandas

* Clone this repository
* ```cd data```
* Run this command to download and create mapping tables from NCBI and Uniprot:  ```python ./data_table_generator.py```
* ```cd ..```
* ```go build app.go```
* ```./app```

Then access ```http://localhost:3000/``` to check the server is actually working or not.


#### Docker Container
The easiest way to run this application is using Docker.  Suppose you are using Docker host running on ```192.168.99.100```.

Make sure you have _docker-compose_ installed.
 
```bash
docker-compose build
docker-compose up 
```

Then access ```http://192.168.99.100:3000/map```.  You will see the following message:

```json
{
    "code": 405,
    "message": "You need to POST your data to use this service.",
    "error": "Invalid HTTP method used: GET"
}
```

(Since you need to _POST_ your data to use this service, you see this error message.)

To test, try the following command:

```bash
curl -H "Accept: application/json" -H "Content-type: application/json" -X \
    POST -d '{"ids": ["P53_HUMAN", "TP53", "P04637", "7157fdsfds"]}' \
    http://192.168.99.100:3000/map | jq .
```

(All in one line.)


### Register to _elsa_
This service works in two modes:

#### Single server mode
If you run this service without elsa instance, it works as a simple RESTful API server.  You can use your server by 
directly send your _POST_ requests to it.


#### CI service mode



## How to use the service

### REST API

#### _/_

* Function - Show service information
* Supported Methods - __GET__
* Output

```json
{
    name: "Gene ID Mapping service",
    version: "v1",
    description: "Converts list of IDs into other types of IDs.",
    documents: "https://github.com/cytoscape-ci/service-go"
}
```

#### _/map_

* Function - Gene ID mapping function
* Supported Methods - __POST__

##### Query

```json
{
    "ids": ["id1", "id2", ... ],
    "idTypes": ["GeneID", "Synonyms", ... ]
}
```

##### _ids_
A list of IDs as text.  You can mix multiple ID types, and multiple species in a query.  Matching is __NOT case sensitive__.

##### Supported input ID types and examples

* Entrez Gene ID - ```5594, 855005```
* UniProtKB-AC - ```P28482, P19880```
* UniProtKB-ID - ```MK01_HUMAN, YAP1_YEAST```
* Ensembl Gene ID - ```ENSG00000100030``` 
* Locus Tag - ```YML007W```
* Gene Symbol - ```MAPK1, YAP1```
* Synonyms - ```ERK, p38, p40, p41, PAR1, PDR4, SNQ3, etc...```

##### Supported Species (more coming soon...)

* Human (Homo sapiens: 9606)
* Yeast (Saccharomyces cerevisiae: 4932)
* Fly (Drosophila melanogaster: 7227)
* Mouse (Mus musculus: 10090)


#### _idTypes_ (Optional)
If you use this optional field, the mapper only returns IDs you specified.

__This parameter is case sensitive.__


##### Supported Types (as of 12/18/2015)

Although the main scope of this API is ID mapping, the data contains some additional information, such as gene ontology, 
chromosome ID, etc.

* UniProtKB-AC
* UniProtKB-ID
* GeneID
* RefSeq
* GI
* PDB
* GO
* UniRef100
* UniRef90
* UniRef50
* UniParc
* PIR
* MIM
* UniGene
* EMBL
* EMBL-CDS
* Ensembl
* Ensembl_TRS
* Ensembl_PRO
* Symbol
* LocusTag
* Synonyms
* dbXrefs
* chromosome
* map_location
* description
* Full_name_from_nomenclature_authority


### Result

#### Example

* Input:

```json
{
    "ids": ["Antp", "HOXA7"],
    "idTypes": ["GeneID", "Symbol", "UniProtKB-ID", "Synonyms"]
}
```

* Output:

```json
{
    "unmatched": [],
    "matched": [
        {
            "matches": {
                "UniProtKB-ID": "Q7KSY7_DROME",
                "Symbol": "Antp",
                "GeneID": "40835",
                "Synonyms": [
                    "3.4",
                    "ANT-C",
                    "ANT-P",
                    "ANTC",
                    "ANTP",
                    "Ant",
                    "AntP",
                    "AntP1",
                    "Antp P1",
                    "Antp P2",
                    "Antp1",
                    "Aus",
                    "BG:DS07700.1",
                    "CG1028",
                    "DMANTPE1",
                    "DRO15DC96Z",
                    "DmAntp",
                    "Dmel\\CG1028",
                    "Hu",
                    "Ns",
                    "Scx",
                    "antp",
                    "l(3)84Ba"
                ]
            },
            "in": "Antp",
            "inType": "Symbol",
            "species": "fly"
        },
        {
            "matches": {
                "UniProtKB-ID": "HXA7_HUMAN",
                "Symbol": "HOXA7",
                "GeneID": "3204",
                "Synonyms": [
                    "ANTP",
                    "HOX1",
                    "HOX1.1",
                    "HOX1A"
                ]
            },
            "in": "Antp",
            "inType": "Synonyms",
            "species": "human"
        },
        {
            "matches": {
                "UniProtKB-ID": "HXA7_HUMAN",
                "Symbol": "HOXA7",
                "GeneID": "3204",
                "Synonyms": [
                    "ANTP",
                    "HOX1",
                    "HOX1.1",
                    "HOX1A"
                ]
            },
            "in": "HOXA7",
            "inType": "Symbol",
            "species": "human"
        },
        {
            "matches": {
                "UniProtKB-ID": "Q8JZW2_MOUSE",
                "Symbol": "Hoxa7",
                "GeneID": "15404",
                "Synonyms": [
                    "AV118143",
                    "Hox-1.1",
                    "M6"
                ]
            },
            "in": "HOXA7",
            "inType": "Symbol",
            "species": "mouse"
        }
    ]
}
  
```

### Sample Client Code in Python

Since this is a RESTful API, you can use any client to access this service.  The following is an example code to convert 
list of IDs using Python.


```python

import json, requests

SERVICE_URL = "http://192.168.99.100:3000/map" # Replace this to your server location.

# Mixed species query allowed - human, mouse, yeast, and fly
query = {
    "ids": ["Antp", "HOXA7"],
    "idTypes": ["GeneID", "Symbol", "UniProtKB-ID", "Synonyms"]
}

res = requests.post(SERVICE_URL, json=query)
print(res.json())

```

## Errors

### Invalid user data

If the query does not contain required parameters (for this version, _ids_ is the only required field) 
the service returns this error.

* Code: 400
* Body:

  ```json
  {
    "code": 400,
    "message": "Invalid query.  Probably you missed ids parameter?",
    "error": "(Any error massage from the system.)"
  }
  ```
* Example

### Unsupported Method

You will see this when you use unsupported HTTP method for an endpoint.
For example you need to use __POST__ method to call this ID Mapping service, and you will get this error 
when you simply call the URL with GET method.

* Code: 405
* Body:

  ```json
  {
    "code": 405,
    "message": "Unsupported HTTP request type.",
    "error": "You need to use POST method to use this endpoint."
  }
  ```

* Example


### Resource Not found

If the service cannot find any result, it returns this error.

* Code: 404
* Body:

  ```json
  {
    "code": 404,
    "message": "No resource found for your query.",
    "error": "No maching IDs for your inputs"
  }
  ```

### Unexpected server side error

Usually this should not happen.  In Go, if critical panic happens due to bugs, this will be returned to the user.

* Code: 500
* Body:

  ```json
  {
    "code": 500,
    "message": "Something wrong happened to the service.  Now is the good time to call admin...",
    "error": "(stack trace, panic message, heap dump, etc.)"
  }
  ```
