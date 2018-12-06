#!/usr/bin/env bash

( cd modules/catalog-service && make clean docker-build )
( cd modules/orders-service && make clean docker-build )