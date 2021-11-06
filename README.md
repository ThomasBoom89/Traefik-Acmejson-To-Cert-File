# Traefik ACME.json to Certificate File

## Description

This Application does exactly what it's name describes. If you use Traefik V2 as your reverse proxy ,integrated certbot
for Let's Encrypt Certificates and save it into the ACME.json there may will be a time you need to extract the key and
crt file from json. This is where TATCF comes into play. The application will read the acme.json and create
Certificates (.key and .crt) for it. Also, it will watch the acme.json file for changes and automatic update the created
Certificates. 
