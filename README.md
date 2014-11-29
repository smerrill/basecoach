## Base Coach

The Base Coach project aims to allow simple alerting based on log data in ElasticSearch, most likely put there by Logstash.

Base Coach is inspired by Heroku's excellent Umpire project, but is designed to be useful if you only have the ELK stack in place.

Base Coach is written in the Go language and is deployed as a single binary. The executable is called `basecoach`.

### Running `basecoach`

The `basecoach` daemon knows two commands, `run`, and `configtest`. In config test mode, it will try to parse all checks and stop before starting its HTTP server.

Both modes accept the following command-line flags with the following defaults:

#### -bindaddr

The address that basecoach should bind to.

Defaults to '0.0.0.0:8080'.

#### -checkfiles

The path to search for check definition files. Globs are allowed.

Defaults to '/opt/basecoach/config/\*.yml'.

#### -esurl

The URL to the ElasticSearch server that will be checked. Note that currently all queries are done as POST requests to the \_count endpoint.

Defaults to 'http://localhost:9200'.

#### -timeout

The timeout value for HTTP requests to ElasticSearch in milliseconds.

Defaults to 2000.

### Configuring Checks

Base Coach uses simple YAML files to define checks and thresholds.

Base Coach will create a route for each check file at path 

It is planned to allow using Go's text.Template system to further templatize the ElasticSearch queries in the future. This should allow for more dynamic query elements such as time intervals to be injected.

Although Base Coach does not yet use any templating, it is still possible to do dynamic queries with date fields. See examples/check.yml for a simple example looking for more than 10 entries of type 'php-error' in the last 15 minutes.

### HTTP Endpoints

#### GET /checks/:check

These are the main check endpoints. One will be created for each check file that is parsed, based on the basename of the file (without an extension.)

#### GET /ping

This endpoint will return a 200 response with the text "pong".

#### GET /

This endpoint currently returns a 200 response with the text "Base Coach is waving you home on $bindaddr.".

### HTTP Response Codes

Just like Umpire, Base Coach is designed to be used in a `check_http` check as part of an internal system like Icinga or an external system like Pingdom.

It will normally return a 200 or a 500.

#### 200

The search returned results that were inside of the acceptable thresholds.

#### 500

The search returned results that were outside of the acceptable thresholds.

#### 503

Base Coach encountered an error while connecting to ElasticSearch, which could also mean a timeout above the configured connection timeout.
