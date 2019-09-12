FROM alpine

RUN apk add --no-cache curl

ENV CONSUL_HOST_ADDR sessions.consul:8500
ENV PUBLIC_KEY tls/pubKey
ENV PRIVATE_KEY tls/privKey
ENV SESSIONS_SECRET service/secret

ENTRYPOINT "bin/run_service.sh"

EXPOSE 443