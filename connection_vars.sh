#!/bin/bash

# connection vars
export D2_TARGET_DB_USER="root"
export D2_TARGET_DB_PASSWORD="example_password"
export D2_TARGET_DB_HOST="localhost"
export D2_TARGET_DB_PORT="3306"
export D2_TARGET_DB_NAME="testdb"

# override vars
export TABLE_GROUPS_PATH="example_table_groups.json"
export VIRTUAL_LINKS_PATH="example_virtual_links.json"
export DESIGNATED_USER="'testuser'@'%'"

# config vars
export D2_TARGET_DB_TYPE="mysql" # accepts mysql
export VIRTUAL_LINKS="true"
export TABLE_GROUPS="true"
export RESTRICTOR_TYPE="minimal" # accepts user or minimal