#!/bin/bash

cd ~/supabase

docker compose \
    -f ~/supabase/supabase-compose.yml \
    --env-file .env \
    down

docker compose \
    -f ~/smtp2api/smtp2api-compose.yml \
    --env-file .env \
    down