# db_to_d2

This workspace is populated with modules which are oriented around diagramming the details of databases using the d2 diagramming language. 

It is oriented in a hexagonal-like architecture. This is by no means an expert implementation, but I believe this is now complete enough to be useful, and structured well enough to extend by others.

## Usage

The following is assumed at this point in time: 
- You have a recent version of `Go` installed
- You have `git` installed
- You have `d2` installed

Currently there is no distributable, so you will need to:
- `git clone git@gitlab.com:JDOsborne1/db_to_d2.git`
- `cd cmd/`
- `go install`

Which will install the application in your GOBIN directory, which is assumed to be on your path for further instructions. 

> Alternatively, you can use the `go build` command to build the binary, and then run it directly from the `cmd/` directory.

The binary makes use of several environment variables to function, including for connection credentials. This is currently the only method supported, but please see the roadmap for future plans in this space.

An example of all the possible environment variables which you can use, along with possible options, can be found in `connection_vars.sh` in the root of the repo. At this stage, only the DB connection variables are required for the program to function, all the others have workable defaults. 

You will then want to run `db_to_d2` which will output the d2 of the database you've pointed it at. You can then pipe this `>` into a .d2 file which you can then render with the d2 commandline tool. 

It is recommended to use the -l "tala" option, which renders ERD diagrams the best, but this assumes that you also have `tala` installed, and that you either have a license, or aren't disturbed by the 'UNLICENSED' watermark.


## Features

### Basic ERD diagram
If you use the application with its defaults, and point at a MySQL database, you will get a basic ERD diagram of the database. This will include all tables, and all foreign keys which are recorded in the database information schema.

### Grouped Tables
As part of diagramming a large database, you may want to introduce some manual groupings of tables. This is supported by the application, and can be done by adding the groups into a .json file and pointing the application at it's path using the variable `TABLE_GROUPS_PATH`. An example of this can be found in `example_table_groups.json` in the `cmd/` directory.

This can be of special use when refactoring a database, as you can group tables together which are related, and the diagram will then illustrate the way that your logical groups remain interconnected.

You will need to enable the grouped tables behaviour by setting `TABLE_GROUPS` to 'true' in your environment variables.

### Virtual Links
Sometimes you may want to show a relationship between two tables which is not recorded in the database schema. This is supported by the application, and can be done by adding the virtual links into a .json file and pointing the application at it's path using the variable `VIRTUAL_LINKS_PATH`. An example of this can be found in `example_virtual_links.json` in the `cmd/` directory.

This feature can be of use when you have either a not very well formed database, where the logical links aren't captured in the structure (for one reason or another, many legitimate). It can also be of use when you have a database which is in the process of being refactored to support some microservices, and you still have multiple services writing to the same database.

You will need to enable the virtual links behaviour by setting `VIRTUAL_LINKS` to 'true' in your environment variables.

### Minimalist restrictions
When diagramming a large database, you may find yourself hitting the limitations of your chosen layout engine. This is especially common when you are working with a database where the tables have lots of columns. This is reasonably common if you have a database with a low level of normalisation, intentional or otherwise. 

You can simplify this by restricting the application to only display the columns which are used in the foreign keys of the tables. This will reduce the number of columns displayed, and will allow you to cleanly render diagrams of larger databases.

You can enable this by setting the environment variable `RESTRICTOR_TYPE` to 'minimal'. 

### User permission restrictions

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

- Find a way to run `d2 fmt` on the output of the program, for better formatting
- Extend docs to include walkthrough for using `incredible-cli` for env var management
- Optionally supply configuration options via commandline 
- Optionally supply configuration options via .yaml config
    - later on possibly extend to .json format also
- Support other database flavours
    - MSSQL
    - PostgreSQL
- Support Multi-schema diagrams (possibly using D2s multi-diagram options)
- (Possibly) break out `d2` elements of the core package to enable other diagramming tools to be used