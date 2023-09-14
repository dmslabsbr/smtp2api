#!/bin/bash

# v3 - 18-08-2023
# inicia dois arquivos docker-compose
# corrigir o caminho dos arquivos
# verificar o arquivo .env
# melhor executar a partir do diretório do supabase

cd ~/supabase

# Nome do container Docker
CONTAINER_NAME="supabase-db"
DEFAULT_NETWORK="supabase-default"
ENV_NAME=".env.tmp"

# Obtém a rede que o container Docker está utilizando
#NETWORK_NAME=$(docker inspect $CONTAINER_NAME --format '{{range .NetworkSettings.Networks}}{{.NetworkID}}{{end}}')
# Exporta a rede como uma variável de ambiente
#export DOCKER_CONTAINER_NETWORK=$NETWORK_NAME

# Imprime a variável para verificar
echo -e '\033[0;32m$DOCKER_CONTAINER_NETWORK:\033[0m'  $DOCKER_CONTAINER_NETWORK

docker compose \
    -f ~/supabase/supabase-compose.yml \
    --env-file .env \
    up -d

# cria novo .env
cp .env $ENV_NAME


# verifica o nome da rede
source check-docker-network.sh $CONTAINER_NAME $DEFAULT_NETWORK

# Imprime a variável para verificar
echo -e 'Novo: \033[0;32m$DOCKER_CONTAINER_NETWORK:\033[0m'  $DOCKER_CONTAINER_NETWORK

# adiciona var no novo arquivo

echo "" >> $ENV_NAME
echo "# Define Network" >> $ENV_NAME
echo "DOCKER_CONTAINER_NETWORK="'"'$DOCKER_CONTAINER_NETWORK'"' >> $ENV_NAME

pause

docker compose \
    -f ~/smtp2api/smtp2api-compose.yml \
    --env-file $ENV_NAME \
    up -d
