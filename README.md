# Integrand Website

Monilith Application to run the integrand.io marketing website. Very simple CMS with batteries included.


## How to Setup New Dev Environment
```bash
go get .
```

## How to Run
`go run .`

## Docker Instructions

### Build the Docker Image
`docker build -t integrand .`

### Run the Docker Container
`docker run -it -p 8000:8000 integrand:latest`

## Docker Tag and Push the Container
`docker tag integrand:latest registry.vineglue.com/integrand:latest`

`docker push registry.vineglue.com/providers:latest`

## Project Structure
TODO: Add this section

## Todos
Track this in project management

1. Add user ownership to features
2. Add organization(s)
3. Write test cases as curl(s)
4. Create bash script to be able to quickly run
5. Create integration test suite.
6. Run integration tests in docker compose file