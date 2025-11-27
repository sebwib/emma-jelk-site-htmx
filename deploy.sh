API_TOKEN=$(grep '^API_TOKEN=' .env | cut -d '=' -f2)

HOST=https://www.emmajelk.se

curl -H "Authorization: Bearer $API_TOKEN" $HOST/api/db --output ./down/database.db
curl -H "Authorization: Bearer $API_TOKEN" $HOST/api/db-shm --output ./down/database.db-shm
curl -H "Authorization: Bearer $API_TOKEN" $HOST/api/db-wal --output ./down/database.db-wal

timestamp=$(date +"%Y%m%d%H%M%S")

docker build --platform linux/amd64 -t emma-site-htmx:${timestamp} .

docker tag emma-site-htmx:${timestamp} 198576290984.dkr.ecr.eu-central-1.amazonaws.com/emma-site-htmx:${timestamp}

docker push 198576290984.dkr.ecr.eu-central-1.amazonaws.com/emma-site-htmx:${timestamp}