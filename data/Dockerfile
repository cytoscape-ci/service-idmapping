FROM python:3.5

RUN apt-get update

RUN pip install pandas

RUN mkdir /tables
WORKDIR /tables

ADD . /tables

CMD python ./data_table_generator.py