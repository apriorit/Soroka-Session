#!/usr/bin/env sh

echo "Adding keys to k/v storage..."

# add connection string to consul
DB_CONNECTION="stub"
curl -X PUT --data-binary "$DB_CONNECTION" http://$CONSUL_HOST_ADDR/v1/kv/$SESSIONS_DB_KEY

#Add public key to consul
cat /serv.crt
curl -X PUT --data-binary @/serv.crt http://$CONSUL_HOST_ADDR/v1/kv/$PUBLIC_KEY

#Add private key
cat /serv.key
curl -X PUT --data-binary @/serv.key http://$CONSUL_HOST_ADDR/v1/kv/$PRIVATE_KEY

#Add
SECRET="secret"
echo $SECRET
curl -X PUT --data-binary "$SECRET" http://$CONSUL_HOST_ADDR/v1/kv/$SESSIONS_SECRET

chmod +x /bin/sessionssvc

/bin/sessionssvc -consul.address $CONSUL_HOST_ADDR \
 - consul.tls.pubkey $PUBLIC_KEY \
 - consul.tls.privkey $PRIVATE_KEY \
 - consul.service.signingKey $SESSIONS_SECRET