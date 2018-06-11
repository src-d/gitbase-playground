# REST API

## GET /schema

Returns the database schema as a list of tables and columns for each table. See below for more details.

```bash
curl -X GET http://localhost:8080/schema
```

```javascript
{
    "status": 200,
    "data": [
        {
            "table": "tree_entries",
            "columns": [
                { "name": "tree_hash", "type": "TEXT" },
                { "name": "entry_hash", "type": "TEXT" },
                { "name": "mode", "type": "TEXT" },
                { "name": "name", "type": "TEXT" }
            ]
        },
        {
            "table": "blobs",
            "columns": [
                { "name": "hash", "type": "TEXT" },
                { "name": "size", "type": "INT64" },
                { "name": "content", "type": "BLOB" }
            ]
        },
        [...]
    }
}
```

## POST /query

Receives an SQL query and forwards it to the `gitbase` server.

The request body can have:

* `query`: An SQL statement string. Do not include `LIMIT` here.
* `limit`: Number, will be added as SQL `LIMIT` to the query. Optional.

The success response will contain:

* `status`: HTTP status code.
* `data`: Rows, array of JSON objects.
* `meta`: JSON object, with these fields:
  * `headers`: Array of strings with the names of the requested columns.
  * `types`: Array of strings with the types of each column. Note: these are the types reported by MySQL, so for example a type `BIT` will be a boolean in the `data` JSON.

A failure response will contain:

* `status`: HTTP status code.
* `errors`: Array of JSON objects, with the fields below:
  * `status`: HTTP status code.
  * `title`: Error description.
  * `mysqlCode`: Error code reported by MySQL. May not be present for some errors.

Some examples follow. A basic query:

```bash
curl -X POST \
  http://localhost:8080/query \
  -H 'content-type: application/json' \
  -d '{
  "query": "SELECT name, hash, is_remote(name) AS is_remote FROM refs",
  "limit": 20
}'
```

```javascript
{
    "status": 200,
    "data": [
        {
            "hash": "66fd81178abfa342f873df5ab66639cca43f5104",
            "is_remote": false,
            "name": "HEAD"
        },
        {
            "hash": "66fd81178abfa342f873df5ab66639cca43f5104",
            "is_remote": false,
            "name": "refs/heads/master"
        },
        {
            "hash": "66fd81178abfa342f873df5ab66639cca43f5104",
            "is_remote": true,
            "name": "refs/remotes/origin/master"
        }
    ],
    "meta": {
        "headers": [
            "name",
            "hash",
            "is_remote"
        ],
        "types": [
            "TEXT",
            "TEXT",
            "BIT"
        ]
    }
}
```

A failure:

```bash
curl -X POST \
  http://localhost:8080/query \
  -H 'content-type: application/json' \
  -d '{
  "query": "SELECT * FROM somewhere",
  "limit": 20
}'
```

```javascript
{
    "status": 400,
    "errors": [
        {
            "status": 400,
            "title": "unknown error: table not found: somewhere",
            "mysqlCode": 1105
        }
    ]
}
```

For a query with uast nodes the protobuf response is unmarshalled and the json is returned:

```bash
curl -X POST \
  http://localhost:8080/query \
  -H 'content-type: application/json' \
  -d '{
  "query": "SELECT hash, content, uast(blobs.content, 'go') FROM blobs WHERE hash='fd30cea52792da5ece9156eea4022bdd87565633'",
  "limit": 20
}'
```

```javascript
{
    "status": 200,
    "data": [
        {
            "content": "package server\n\nimport (\n\t\"net/http\"\n\n\t\"github.com/src-d/gitbase-playground/server/handler\"\n\n\t\"github.com/go-chi/chi\"\n\t\"github.com/go-chi/chi/middleware\"\n\t\"github.com/pressly/lg\"\n\t\"github.com/rs/cors\"\n\t\"github.com/sirupsen/logrus\"\n)\n\n// Router returns a Handler to serve the backend\nfunc Router(\n\tlogger *logrus.Logger,\n\tstatic *handler.Static,\n\tversion string,\n) http.Handler {\n\n\t// cors options\n\tcorsOptions := cors.Options{\n\t\tAllowedOrigins:   []string{\"*\"},\n\t\tAllowedMethods:   []string{\"GET\", \"POST\", \"PUT\", \"OPTIONS\"},\n\t\tAllowedHeaders:   []string{\"Location\", \"Authorization\", \"Content-Type\"},\n\t\tAllowCredentials: true,\n\t}\n\n\tr := chi.NewRouter()\n\n\tr.Use(middleware.Recoverer)\n\tr.Use(cors.New(corsOptions).Handler)\n\tr.Use(lg.RequestLogger(logger))\n\n\tr.Get(\"/version\", handler.APIHandlerFunc(handler.Version(version)))\n\n\tr.Get(\"/static/*\", static.ServeHTTP)\n\tr.Get(\"/*\", static.ServeHTTP)\n\n\treturn r\n}\n",
            "hash": "fd30cea52792da5ece9156eea4022bdd87565633",
            "uast(blobs.content, \"go\")": [
                {
                    "InternalType": "File",
                    "Properties": {
                        "Package": "1"
                    },
                    "Children": [
                        {
                            "InternalType": "Ident",
                            "Properties": {
                                "internalRole": "Name"
                            },
                            "Token": "server",
                            "StartPosition": {
                                "Offset": 9,
                                "Line": 1,
                                "Col": 10
                            },
                            "Roles": [
                                18,
                                1
                            ]
                        },
                        {
                            "InternalType": "GenDecl",
                            "Properties": {
                                "Lparen": "24",
                                "Tok": "import",
                                "internalRole": "Decls"
                            },
                            "Children": [
                                {
                                    "InternalType": "ImportSpec",
                                    "Properties": {
                                        "EndPos": "0",
                                        "internalRole": "Specs"
                                    },

        [...]

        }
    ],
    "meta": {
        "headers": [
            "hash",
            "content",
            "uast(blobs.content, \"go\")"
        ],
        "types": [
            "TEXT",
            "TEXT",
            "JSON"
        ]
    }
}
```

## POST /parse

Receives a file content and returns UAST.

```bash
curl -X POST \
  http://localhost:8080/parse \
  -H 'content-type: application/json' \
  -d '{
  "language": "javascript",
  "content": "console.log(test)"
}'
```

```javascript
{
    "status": 200,
    "data": {
        "InternalType": "File",
        "Children": [
            {
                "InternalType": "Program",
                "Properties": {
                    "internalRole": "program",
                    "sourceType": "module"
                },
                [...]
            }
        ]
    }
}
```

The endpoint also receives additional parameters:

* `serverUrl` - allows to override bblfsh server url.
* `filename` - can be used instead of language. Then the bblfsh server would try to guess the language.
* `filter` - [xpath query](https://doc.bblf.sh/user/uast-querying.html) to filter the results.

## GET /export

This endpoint is similar to `/query` but returns results as CSV file without LIMIT.

```bash
curl -X GET http://localhost:8080/export?query=select+*+from+repositories
```

```javascript
id
/opt/repos/gitbase-playground
/opt/repos/go-git-fixtures
```

