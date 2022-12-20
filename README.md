# cozy-doc-api
REST api exposing an endpoint to insert documents in a couchDB database

## Couch DB doc
https://docs.couchdb.org/en/3.2.2-docs/ . 

## Endpoints :

- health :  Return status = 200 if the server is up
```
GET /health
```
- docs : Insert an array of documents in a given database. If the database doesn't exist it will be created. 

```
POST /docs
```
with body :

```
{
    "database": "test",
    "docs": [
        {
            "servings": 2,
            "subtitle": "Delicious with  bread",
            "title": "Fish Stew"
        },
        {
            "servings": 4,
            "subtitle": "Delicious with fresh bread",
            "title": "TEst Stew"
        },
        {
            "servings": 4,
            "subtitle": "Delic fresh bread",
            "title": "Fish "
        }
    ]
}
```
## Dev : 

- To run the server

```sh
make serve
```
- to run tests

In some cases you will need to regenerate mock files for tests. This is done with the following command:

```sh
make generate-mocks
```

You will also need to have mockery (https://github.com/vektra/mockery) in your path.


```sh
make test
```

- to install docker containers needed for dev en (redis)
```sh
make docker-dev
```

- to run api as docker container 
```sh
make docker-run
```