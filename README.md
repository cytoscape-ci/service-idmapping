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

##### 1. Prepare Docker Host
You need a Docker host to run the application using Docker Compose.

* System Requirements
    * 1 or 2 processors
    * 3+ GB of memory
    * Fast, stable network connection

In general, this application is I/O intensive, and you do not need a very fast processors to run this application.

I recommend to use [Docker Machine](https://docs.docker.com/machine/) to setup your host.

----

###### For the first-time Docker Machine users

Here is how to create new Docker host using VirtualBox:

* Install everything using [Docker Toolbox](https://www.docker.com/docker-toolbox)
* Open terminal
* Make sure you have docker-machine and docker-compose on your machine  

```bash
> docker-machine -v                                                        ✱
docker-machine version 0.5.5, build 02c4254

> docker-compose -v                                                        ✱
docker-compose version 1.5.2, build 7240ff3
```

* Create a new Docker host using VirtualBox driver

```bash
> docker-machine create -d virtualbox dev2                                 ✱
Running pre-create checks...
Creating machine...
(dev2) Copying /Users/kono/.docker/machine/cache/boot2docker.iso to /Users/kono/.docker/machine/machines/dev2/boot2docker.iso...
(dev2) Creating VirtualBox VM...
(dev2) Creating SSH key...
(dev2) Starting VM...
Waiting for machine to be running, this may take a few minutes...
Machine is running, waiting for SSH to be available...
Detecting operating system of created instance...
Detecting the provisioner...
Provisioning with boot2docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect Docker to this machine, run: docker-machine env dev2
```

* Specify the machie as your target Docker host


```bash
> eval "$(docker-machine env dev2)"
```

* Make sure you have access to that Docker host

```bash
> docker info                                                            ⏎ ✱
Containers: 0
Images: 0
Server Version: 1.9.1
Storage Driver: aufs
 Root Dir: /mnt/sda1/var/lib/docker/aufs
 Backing Filesystem: extfs
 Dirs: 0
 Dirperm1 Supported: true
Execution Driver: native-0.2
Logging Driver: json-file
Kernel Version: 4.1.13-boot2docker
Operating System: Boot2Docker 1.9.1 (TCL 6.4.1); master : cef800b - Fri Nov 20 19:33:59 UTC 2015
CPUs: 1
Total Memory: 996.2 MiB
Name: dev2
ID: QQER:IUDG:JNRA:YCAD:K4LS:QHNT:U5X6:WHGS:KRZM:DDKD:J5RP:ZYSY
Debug mode (server): true
 File Descriptors: 11
 Goroutines: 18
 System Time: 2016-01-04T19:14:07.938448467Z
 EventsListeners: 0
 Init SHA1:
 Init Path: /usr/local/bin/docker
 Docker Root Dir: /mnt/sda1/var/lib/docker
Username: keiono
Registry: https://index.docker.io/v1/
Labels:
 provider=virtualbox
```

Now you are ready to use _compose_.

----

##### 2. Run the Service
Make sure you have _docker-compose_ and _git_ installed.

The following commands automatically build and run new service in daemon mode.

```bash
git clone git@github.com:cytoscape-ci/service-idmapping.git
cd service-idmapping

docker-compose build
docker-compose up -d
```

Then access the endpoint, e.g. ```http://192.168.99.100:3000/```.  You will see the following message:

```json
{
    name: "Gene ID Mapping service",
    version: "v1",
    description: "Converts list of IDs into other types of IDs.",
    documents: "https://github.com/cytoscape-ci/service-go"
}
```


To test, try the following command (you need curl and jq to run this):

```bash
curl -H "Accept: application/json" -H "Content-type: application/json" -X \
    POST -d '{"ids": ["P53_HUMAN", "TP53", "P04637", "7157"]}' \
    http://192.168.99.100:3000/map | jq .
```

(All in one line.)


### Register to _elsa_
This service works in two modes:

#### Single server mode
If you run this service without elsa instance, it works as a simple RESTful API server.  You can use your server by 
directly send your _POST_ requests to it.


#### CI service mode
If you have a running instance of elsa, the registration process is (semi) automatic.  You can specify the following 
command-line options to register the service to _elsa_.

(TBD)

Options

* id - ID of this service (e.g., _idmapping_)
* ip - IP address of this server
* ver - version of this API (e.g., _v1_)
* port - Port number of this service
* agent - Location of elsa instance


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
