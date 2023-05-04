# db_to_d2

This workspace is populated with modules which are oriented around diagramming the details of databases using the d2 diagramming language. 

It is oriented in a hexagonal-like architecture. This is by no means an expert implementation, but I believe this is now complete enough to be useful, and structured well enough to extend by others.


## Installation


### Released

The program was released in beta form here: https://github.com/JDOsborne1/db_to_d2/releases/download/v0.1-beta/db_to_d2 

### Latest

Below are the instructions if you want to use the latest version of the aplication available. 

The following is assumed at this point in time: 
- You have a recent version of `Go` installed
- You have `git` installed

Currently there is no distributable, so you will need to:
- `git clone git@gitlab.com:JDOsborne1/db_to_d2.git`
- `cd cmd/`
- `go install`

Which will install the application in your GOBIN directory, which is assumed to be on your path for further instructions. 

Alternatives: 
- Use the `go build` command to build the binary, and then run it directly from the `cmd/` directory.
- Use the `go run` command to run the application directly from the `cmd/` directory.
- Use any of the *passing* build artefacts in the Github Actions pipeline, which are available in the `Actions` tab of the repo.


## Usage 

The following is assumed at this point in time: 
- You have `d2` installed

The application makes use of the following variables, of which the *Bold* ones are required: 
    - *D2 target db host*
    - *D2 target db name*
    - *D2 target db password*
    - *D2 target db port*
    - *D2 target db type*
    - *D2 target db user*
    - Designated user
    - Restrictor type
    - Use table groups
    - Path to table groups file
    - Use virtual links
    - Path to virtual links file

These are currently supported to be supplied by any combination of _environment variables_ and _command line flags_. Please consult the relevant section below to see how to supply the information.

You will then want to run `db_to_d2` which will output the d2 of the database you've pointed it at. You can then pipe this `>` into a .d2 file which you can then render with the d2 commandline tool. 

When using `d2` it is recommended to use the -l "tala" option, which renders ERD diagrams the best, but this assumes that you also have `tala` installed, and that you either have a license, or aren't disturbed by the 'UNLICENSED' watermark.


### Environment Variables 

An example of all the possible environment variables which you can use, along with possible options, can be found in `connection_vars.sh` in the root of the repo. 


### Commandline Arguments

The below is taken from the help dialogue for the program, found by supplying `-h` or `--help` as a flag.

```
Usage of ./db_to_d2:
      --D2TargetDbHost string       D2 target db host
      --D2TargetDbName string       D2 target db name
      --D2TargetDbPassword string   D2 target db password
      --D2TargetDbPort string       D2 target db port
      --D2TargetDbType string       D2 target db type
      --D2TargetDbUser string       D2 target db user
      --DesignatedUser string       Designated user
      --RestrictorType string       Restrictor type
      --TableGroups string          Use table groups
      --TableGroupsPath string      Path to table groups file
      --VirtualLinks string         Use virtual links
      --VirtualLinksPath string     Path to virtual links file
pflag: help requested
```

## Features

### Basic ERD Diagram

If you use the application with its defaults, and point at a MySQL database, you will get a basic ERD diagram of the database. This will include all tables, and all foreign keys which are recorded in the database information schema.


### Grouped Tables


As part of diagramming a large database, you may want to introduce some manual groupings of tables. This is supported by the application, and can be done by adding the groups into a .json file and pointing the application at it's path using the variable `TABLE_GROUPS_PATH`. An example of this can be found in `example_table_groups.json` in the `cmd/` directory. Both relative and absolute paths are supported.

This can be of special use when refactoring a database, as you can group tables together which are related, and the diagram will then illustrate the way that your logical groups remain interconnected.

You will need to enable the grouped tables behaviour by setting `TABLE_GROUPS` to 'true' in your environment variables.


### Virtual Links


Sometimes you may want to show a relationship between two tables which is not recorded in the database schema. This is supported by the application, and can be done by adding the virtual links into a .json file and pointing the application at it's path using the variable `VIRTUAL_LINKS_PATH`. An example of this can be found in `example_virtual_links.json` in the `cmd/` directory. Both relative and absolute paths are supported.

This feature can be of use when you have either a not very well formed database, where the logical links aren't captured in the structure (for one reason or another, many legitimate). It can also be of use when you have a database which is in the process of being refactored to support some microservices, and you still have multiple services writing to the same database.

You will need to enable the virtual links behaviour by setting `VIRTUAL_LINKS` to 'true' in your environment variables.


### Minimalist Restrictions


When diagramming a large database, you may find yourself hitting the limitations of your chosen layout engine. This is especially common when you are working with a database where the tables have lots of columns. This is reasonably common if you have a database with a low level of normalisation, intentional or otherwise. 

You can simplify this by restricting the application to only display the columns which are used in the foreign keys of the tables. This will reduce the number of columns displayed, and will allow you to cleanly render diagrams of larger databases.

You can enable this by setting the environment variable `RESTRICTOR_TYPE` to 'minimal'. 


### User Permission Restrictions


For purposes of access management, it can be of use to create a view of the tables and columns which are accessible to a particular user. This can be done by setting the environment variable `RESTRICTOR_TYPE` to 'user'. 

You will need to also set the environment variable `DESIGNATED_USER` to the name of the user you want to restrict the diagram to. Currently this needs to be the full name of the user, including the host.

This can also be useful when diagnosing and preparing for a transformation to a microservices architecture, as you can use this to show the tables accessed by the various client users. 

## Local Environment

If you want to test the tool against a stable local system, you can make use of the distributed seeding query. 

In the root directory: 
- `docker compose up`
- navigate to localhost:8080
- Copy contents of `seeding_db.sql` into the phymyadmin query window
- Hit 'Go' and it should create a new database called `testdb`
- This db should then be suitable for the example setup contained in `connection_vars.sh`

## Roadmap

- Expand configuration options:
    - Support commandline flags
    - Support .yaml config
    - Support .json config
- Support other database flavours
    - MSSQL
    - PostgreSQL
- Support Multi-schema diagrams (possibly using D2s multi-diagram options)
- (Possibly) break out `d2` elements of the core package to enable other diagramming tools to be used

## Versioning

The project will be versioned with Semantic Versioning, as recommended by Github and Gitlab.

The Public API surface of this repo is the behaviour of the binary itself, as documented in this readme, and in the package documentation associated with that build. 

Special consideration will be given to the interfaces in the `core` package, for those who may eventually use it as a basis in other tooling. These are currently subject to change as part of the intial development, but are intended to be included as part of the public API for version 1.
