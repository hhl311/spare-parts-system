all: catalog orders packing-slips

catalog:
	cd modules/catalog-service && make clean deps docker-build

orders:
	cd modules/orders-service && make clean deps docker-build

packing-slips:
	cd modules/packing-slips-service && make clean deps docker-build

deploy:
	cd modules/catalog-service && make deploy
	cd modules/orders-service && make deploy
	cd modules/packing-slips-service && make deploy