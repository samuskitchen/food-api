#!/bin/bash

echo Inicia proceso de documentacion con Swagger

. ~/.bashrc

killall -9 swagger

cd food-api/

$GOBIN/swagger serve -F=redoc --host=0.0.0.0 --port=8082 --no-open swagger.json

echo "Se ha completado la documentacion correctamente"