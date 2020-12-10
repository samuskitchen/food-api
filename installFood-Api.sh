#!/bin/bash

echo Inicia proceso de instalacion de Food Api

. ~/.bashrc

docker image prune -a

cd installers/

rm -rf food-api/

git clone https://github.com/samuskitchen/food-api.git -b master food-api

cd food-api/

docker-compose down --remove-orphans --volumes

docker-compose up -d --build

echo "Se ha completado la instalación correctamente"