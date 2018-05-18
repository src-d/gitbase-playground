# Quickstart

## Run bblfsh and gitbase Dependencies

It is recommended to read about `bblfsh` and `gitbase` from its own documentation, but here is a small guide about how to run both easily:

Launch [bblfshd](https://github.com/bblfsh/bblfshd) and install the drivers. More info in the [bblfshd documentation](https://doc.bblf.sh/user/getting-started.html):

```bash
$ docker run --privileged \
    --publish 9432:9432 \
    --volume /var/lib/bblfshd:/var/lib/bblfshd \
    --name bblfsh \
    bblfsh/bblfshd
$ docker exec -it bblfsh \
    bblfshctl driver install --recommended
```

[gitbase](https://github.com/src-d/gitbase) will serve git repositories, so it is needed to populate a directory with them:

```bash
$ mkdir -p ~/gitbase/repos
$ git clone git@github.com:src-d/go-git-fixtures.git ~/gitbase/repos/go-git-fixtures
```

Install and run [gitbase](https://github.com/src-d/gitbase):

```bash
$ docker run --publish 3306:3306 \
    --link bblfsh \
    --volume ~/gitbase/repos:/opt/repos \
    --env BBLFSH_ENDPOINT=bblfsh:9432 \
    --name gitbase \
    srcd/gitbase:latest
```


## Run gitbase-playground

Once bblfsh and gitbase are running and accessible, you can serve the playground:

```bash
$ docker run -d --publish 8080:8080 \
    --link gitbase \
    --env GITBASEPG_ENV=dev \
    --env GITBASEPG_DB_CONNECTION="gitbase@tcp(gitbase:3306)/none?maxAllowedPacket=4194304" \
    --name gitbase_playground \
   srcd/gitbase-playground:latest
```

Once the server is running &ndash;with its default values&ndash;, it will be accessible through: http://localhost:8080

You have more information about the [playground architecture](CONTRIBUTING.md#architecture), [development guides](CONTRIBUTING.md#development) and [configuration options](CONTRIBUTING.md#configuration) in the [CONTRIBUTING.md](CONTRIBUTING.md).


## Run a Query

You will find more info about how to run queries using the playground API on the [rest-api guide](rest-api.md)
