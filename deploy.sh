timestamp=$(date +"%Y%m%d%H%M%S")

docker build --platform linux/amd64 -t emma-site-htmx:${timestamp} .

docker tag emma-site-htmx:${timestamp} 198576290984.dkr.ecr.eu-central-1.amazonaws.com/emma-site-htmx:${timestamp}

docker push 198576290984.dkr.ecr.eu-central-1.amazonaws.com/emma-site-htmx:${timestamp}