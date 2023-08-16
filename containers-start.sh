#!/bin/bash

# inicia dois arquivos docker-compose
# corrigir o caminho dos arquivos
# verificar o arquivo .env
# melhor executar a partir do diretório do supabase

cd ~/supabase

docker compose \
    -f ~/smtp2api/smtp2api-compose.yml \
    -f ~/supabase/supabase-compose.yml \
    --env-file .env
    up -d