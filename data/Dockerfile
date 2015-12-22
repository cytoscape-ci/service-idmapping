#
# Data volume container for ID Mapping service
#
#  This Container is only for building data sources from remote FTP server.
#
FROM python:3.5

RUN apt-get update
RUN pip install pandas

RUN mkdir -p /go/src/github.com/cytoscape-ci/service-idmapping/data
WORKDIR /go/src/github.com/cytoscape-ci/service-idmapping/data

ADD . /go/src/github.com/cytoscape-ci/service-idmapping/data

# This starts the actual data download and preprecessing.
RUN python ./data_table_generator.py

CMD /bin/true