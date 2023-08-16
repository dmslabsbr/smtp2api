#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Uso: $0 nome_do_container nome_da_rede"
    return 1
fi

CONTAINER_NAME=$1
NETWORK_NAME=$2

# Verifica se o container existe
docker inspect $CONTAINER_NAME > /dev/null 2>&1
CONTAINER_EXISTS=$?

# Se o container existir, obtenha a rede associada a ele
if [ $CONTAINER_EXISTS -eq 0 ]; then
    NETWORK_USED=$(docker inspect $CONTAINER_NAME --format '{{range .NetworkSettings.Networks}}{{.NetworkID}}{{end}}')
    echo -e "O container" "\033[0;33m${CONTAINER_NAME}\033[0m" está usando a rede "\033[0;33m${NETWORK_USED}\033[0m"
    # pega o nome da rede
    NETWORK_NAME=$(docker network inspect $NETWORK_USED --format '{{.Name}}')
    NETWORK_EXISTS=$?
    if [ $NETWORK_EXISTS -eq 0 ]; then
      echo "O nome da rede é $NETWORK_NAME"
      echo -e "O nome da rede é " "\033[0;35m${NETWORK_NAME}\033[0m"
      # Salva o nome da rede em uma variável
      export DOCKER_CONTAINER_NETWORK=$NETWORK_NAME
      echo -e '\033[0;32m$DOCKER_CONTAINER_NETWORK:\033[0m'  $DOCKER_CONTAINER_NETWORK
      return 0
    else
        echo -e '\033[0;31mOcorreu um erro !!\033[0m' 
        return 1
    fi
else
    # Se o container não existir, crie uma rede com o nome fornecido
    echo -e "O container" "\033[0;33m${CONTAINER_NAME}\033[0m" não encontrado. Criando rede "\033[0;36m${NETWORK_NAME}\033[0m"
    docker network create $NETWORK_NAME
    export DOCKER_CONTAINER_NETWORK=$NETWORK_NAME
    echo -e '\033[0;32m$DOCKER_CONTAINER_NETWORK:\033[0m'  $DOCKER_CONTAINER_NETWORK
    return 1
fi
