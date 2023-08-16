#!/bin/bash

# inicia dois arquivos docker-compose
# corrigir o caminho dos arquivos
# verificar o arquivo .env
# melhor executar a partir do diretório do supabase

cd ~/supabase

# Nome do container Docker
CONTAINER_NAME="your_container_name_here"

# Obtém a rede que o container Docker está utilizando
NETWORK_NAME=$(docker inspect $CONTAINER_NAME --format '{{range .NetworkSettings.Networks}}{{.NetworkID}}{{end}}')

# Exporta a rede como uma variável de ambiente
export DOCKER_CONTAINER_NETWORK=$NETWORK_NAME

# Imprime a variável para verificar
echo $DOCKER_CONTAINER_NETWORK



docker compose \
    -f ~/smtp2api/smtp2api-compose.yml \
    --env-file .env \
    up -d

docker compose \
    -f ~/supabase/supabase-compose.yml \
    --env-file .env \
    up -d