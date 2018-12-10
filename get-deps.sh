#!/usr/bin/env bash

( cd modules/catalog-service && make deps )
( cd modules/orders-service && make deps )
( cd modules/packing-slips-service && make deps )