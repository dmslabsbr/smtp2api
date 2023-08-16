docker run --name smtp2api \
    -e "BREVO_APIKEY=12343" \
    --network supabase_default \
    -d smtp2api