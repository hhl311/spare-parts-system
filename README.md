# Spare Parts System

[![img](https://travis-ci.org/AntoineAube/spare-parts-system.svg?branch=master)](https://travis-ci.org/AntoineAube/spare-parts-system/)

A minimalistic system for selling spare parts.

## What is it?

Here is the exercise:
> ### Exercise goal: 
> Build a spare parts catalog, an order and a “packing slip” services following microservice architecture. 
> ### Implementation constraints: 
> - Using dgraph (https://dgraph.io/) build a graph db of a spare parts catalog for Car or Appliance 
> - Create an order and “packing slip” services that are scalable and decoupled following the microservice architecture paradigm. You can store the order in any db you like or in dgraph. Services must be REST API or GRPC services written in golang. (packing slip are used in the warehouse in order to prepare the command)
> - You are free to choose the structure of the database and the structure of API responses  
> - Build a docker image for all the services.
> ### Optional: 
> - Make it all run in Kubernetes 
> - Create an asynchronous communication between services to create an invoice as soon as a command is validated 
> - Support product bundle (one product in the command but several products in the packing slip)

As the constraints are somehow ambiguous, I made choices and the resulting system is described in [the `docs` folder](./docs).
There can be found a components diagram, a sequence diagram and a Postman collection with all possible requests.

There is three Docker images for the whole system:
- `antoineaube/sps-catalog-service`
- `antoineaube/sps-orders-service`
- `antoineaube/sps-packing-slips-service`

There is a [docker-compose.yml](docker-compose.yml) which setups and runs the whole system.

## How use it?

1. Clone this repository in `$GOPATH/src/github.com/AntoineAube/`.
1. To build the Docker images of the three services, run `make`.
1. Run the system with `docker-compose up`. It is not mandatory to locally build the images because they are on DockerHub.