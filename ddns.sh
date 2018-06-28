#!/bin/bash

DDNS_ACCESS_KEY=$DDNS_ACCESS_KEY
DDNS_SECRET_KEY=$DDNS_SECRET_KEY
IP=`curl --ipv4 -s http://icanhazip.com/`
DOMAIN=y00ns.com
SUBDOMAIN=dev

curl -H "Authorization: sso-key $DDNS_ACCESS_KEY:$DDNS_SECRET_KEY" \
     -H "Content-Type: application/json" \
     -H "Accept: application/json" \
     -X "PUT" \
     -i "https://api.godaddy.com/v1/domains/$DOMAIN/records/A/$SUBDOMAIN" \
     -d "[{\"data\": \"$IP\", \"ttl\": 1800}]"
