# README

[Babelfish](https://doc.bblf.sh) Go client library provides functionality to both connecting to the Babelfish server for parsing code \(obtaining an [UAST](https://doc.bblf.sh/uast/specification.html) as a result\) and for analysing UASTs with the functionality provided by [libuast](https://github.com/bblfsh/libuast).

## Installation

The recommended way to install _client-go_ is:

```bash
go get -d -u gopkg.in/bblfsh/client-go.v2/...
cd $GOPATH/src/gopkg.in/bblfsh/client-go.v2
make dependencies
```

Windows build is supported, provided by you have `make` and `curl` in your `%PATH%`. It is also possible to link against custom `libuast` on Windows, read [WINDOWS.md](windows.md).

## Example

This small example illustrates how to retrieve the [UAST](https://doc.bblf.sh/uast/specification.html) from a small Python script.

If you don't have a bblfsh server installed, please read the [getting started](https://doc.bblf.sh/user/getting-started.html) guide, to learn more about how to use and deploy a bblfsh server.

Go to the[quick start](https://github.com/bblfsh/bblfshd#quick-start) to discover how to run Babelfish with Docker.

```go
client, err := bblfsh.NewClient("0.0.0.0:9432")
if err != nil {
    panic(err)
}

python := "import foo"

res, err := client.NewParseRequest().Language("python").Content(python).Do()
if err != nil {
    panic(err)
}

query := "//*[@roleImport]"
nodes, _ := tools.Filter(res.UAST, query)
for _, n := range nodes {
    fmt.Println(n)
}
```

```text
Import {
.  Roles: Import,Declaration,Statement
.  StartPosition: {
.  .  Offset: 0
.  .  Line: 1
.  .  Col: 1
.  }
.  Properties: {
.  .  internalRole: body
.  }
.  Children: {
.  .  0: alias {
.  .  .  Roles: Import,Pathname,Identifier
.  .  .  TOKEN "foo"
.  .  .  Properties: {
.  .  .  .  asname: <nil>
.  .  .  .  internalRole: names
.  .  .  }
.  .  }
.  }
}

alias {
.  Roles: Import,Pathname,Identifier
.  TOKEN "foo"
.  Properties: {
.  .  asname: <nil>
.  .  internalRole: names
.  }
}

iter, err := tools.NewIterator(res.UAST)
if err != nil {
    panic(err)
}
defer iter.Dispose()

for node := range iter.Iterate() {
    fmt.Println(node)
}

// For XPath expressions returning a boolean/numeric/string value, you must
// use the right typed Filter function:

boolres, err := FilterBool(res.UAST, "boolean(//*[@strtOffset or @endOffset])")
strres, err := FilterString(res.UAST, "name(//*[1])")
numres, err := FilterNumber(res.UAST, "count(//*)")
```

Please read the [Babelfish clients](https://doc.bblf.sh/user/language-clients.html) guide section to learn more about babelfish clients and their query language.

## License

Apache License 2.0, see [LICENSE](https://github.com/src-d/gitbase-playground/tree/11b016900612037802e94240c5e73e9ea6770d36/vendor/gopkg.in/bblfsh/client-go.v2/LICENSE/README.md)

