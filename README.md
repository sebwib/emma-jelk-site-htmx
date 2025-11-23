# htmx-experiments

Experimenting with HTMX

## Authentication

The `/edit` route is protected and requires login.

### Setup

1. Copy `.env.example` to `.env`
2. Set your admin credentials in `.env`:
   ```
   ADMIN_USERNAME=your_username
   ADMIN_PASSWORD=your_password
   ```
3. If no `.env` file is provided, the default credentials are:
   - Username: `admin`
   - Password: `password`

### Login

- Navigate to `/login` to access the login page
- Sessions last 24 hours
- Use `/logout` to end your session (POST request)

docker build --platform linux/amd64 -t emma-site-htmx:latest .
docker tag emma-site-htmx:latest 198576290984.dkr.ecr.eu-central-1.amazonaws.com/emma-site-htmx:latest
docker push 198576290984.dkr.ecr.eu-central-1.amazonaws.com/emma-site-htmx:latest
