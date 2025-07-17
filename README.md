# Packaging service

Service provides calculations for packages to amount fo goods that should be packed

## Commands:
```make test``` - test with coverage

```make service-build``` - build service

```make service-up``` - start service

```make service-down``` - stop the service

```make service-logs``` - observe the logs from service

```make atlas-init``` - generate initial database schema

```make atlas-diff``` - generate migrations

```make seed``` - for seeding the database with initial data. Should be used when service is up.

```make test``` - tests for business logic of packs

```make coverage``` - tests with coverage report

## Description

Service is hosted on 8080 port and for localhost can be accessed via localhost:8080.

**/index** is just a welcome page

**/calculate** is endpoint responsible for calculations of packages amount

**/packs** is for managing packs and their sizes. This is protected endpoint. Admin creds are in *.env* file

Login / Logout operations may be done via form on the webpage

There is an API possibility for interactions and jwt token is required for them. Example of token generation is 
in *internal/auth/token.go* - ```GenerateDevToken```

*isAdmin* Claim is needed for /packs

Endpoints:
 - /packs GET POST DELETE
 - /calculate POST

Example:
```
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{"quantity": 42}'
```

```
curl -X GET http://localhost:8080/api/packs \
  -H "Authorization: Bearer your-token"
```
