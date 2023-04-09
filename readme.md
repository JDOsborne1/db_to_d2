# db_to_d2

This workspace is populated with modules which are oriented around diagramming the details of databases using the d2 diagramming language. 

It is oriented in a hexagonal-like architecture. This is by no means an expert implementation, but I believe this is now complete enough to be useful, and structured well enough to extend by others.

## Usage

The following is assumed at this point in time: 
- You have a recent version of `Go` installed
- You have `git` installed
- You have `d2` installed

Currently there is no distributable, so you will need to:
- `git clone git@github.com:JDOsborne1/rdb_to_rdf2.git`
- `cd cmd/`
- `go install`

Which will install the application in your GOBIN directory, which is assumed to be on your path for further instructions. 

The binary makes use of several evironment variables to function, including for connection credentials. This is currently the only method supported, but please see the roadmap for future plans in this space.

> It is possible to make use of a tool like `incredible-cli` to make this more straightforward with minimal security compromises. 

An example of all the possible environment variables which you can use, along with possible options, can be found in `connection_vars.sh` in the root of the repo. At this stage, only the DB connection variables are required for the program to function, all the others have workable defaults. 

You will then want to run `rdb_to_rdf2` which will output the d2 of the database you've pointed it at. You can then pipe this `>` into a .d2 file which you can then render with the d2 commandline tool. 

It is recommended to use the -l "tala" option, which renders ERD diagrams the best, but this assumes that you also have `tala` installed, and that you either have a license, or aren't disturbed by the 'UNLICENSED' watermark.


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