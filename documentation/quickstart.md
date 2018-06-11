# Quick start with docker-compose

You can locally build and deploy `gitbase-playground` and its dependencies using [`docker-compose`](https://docs.docker.com/compose/install/). Docker compose will run three different containers: for the playground frontend itself, gitbase and bblfsh services. It will be the [latest gitbase version](https://hub.docker.com/r/srcd/gitbase/tags/) and [latest bblfsh version](https://hub.docker.com/r/bblfsh/bblfshd/tags/).

If you preffer to run `gitbase-playground` dependencies manually, you can follow [the alternative playground quickstart](quickstart-manually.md)

## Download the project

```bash
$ git clone git@github.com:src-d/gitbase-playground.git gitbase-playground
$ cd gitbase-playground
```

This guide will assume you're running all commands from `gitbase-playground` sources directory

## Populate the database

It is needed to populate a directory with some git repositories to be served by [gitbase](https://github.com/src-d/gitbase) before running it.

example:

```bash
$ git clone git@github.com:src-d/gitbase-playground.git ./repos/gitbase-playground
$ git clone git@github.com:src-d/go-git-fixtures.git ./repos/go-git-fixtures
```

Everytime you want to add a new repository to gitbase, the application should be restarted.

## Run the application

Run the [latest released version of the frontend](https://hub.docker.com/r/srcd/gitbase-playground/tags/):

```bash
$ docker-compose pull --include-deps playground
$ GITBASEPG_REPOS_FOLDER=./repos docker-compose up --force-recreate
```

If you want to build and run the playground from sources instead of using the last released version you can do so:

```text

$ GITBASEPG_REPOS_FOLDER=./repos make compose-serve
```

## Stop the application

To kill the running containers just `Ctrl+C`

To delete the containers run `docker-compose rm -f`

## Access the playground and execute a query

Once the server is running –with its default values–, it will be accessible through: [http://localhost:8080](http://localhost:8080)

You will find more info about how to run queries using the playground API on the [rest-api guide](rest-api.md)

## More information

You have more information about the [playground architecture](contributing.md#architecture), [development guides](contributing.md#development) and [configuration options](contributing.md#configuration) in the [Contribution Guidelines](contributing.md).

