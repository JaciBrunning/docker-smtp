# Docker Postfix SMTP Server

A simple docker image containing the Postfix SMTP server, and OpenDKIM, designed to take TLS SSL certs from a Traefik Load Balancer, with the configuration stored in etcd.

## Environment
`DOMAIN`: The domain to send from
`ETCD`: The etcd host and port (e.g. `etcd:2379`)
`ACME_STORAGE`: The storage location of the acme in Traefik (configured in traefik.toml), e.g. `traefik/acme/account`

## Volumes / Secrets
This image can accept the DKIM private key as either a bindmount, or as a swarm secret. Swarm secret takes priority.

SASL keys require secrets.

### Secret
Secret name: `opendkim_private`
Value of the secret should contain the RSA private key generated for OpenDKIM.

Secret name: `smtp_passwd`
Value of the smtp user's password. Should be stored in form `user:pass`, separated by spaces for multiple

### Volume
Bind mount: `local_dkim_path:/dkim`
`local_dkim_path` should contain a single file - `mail.private`, containing the OpenDKIM private key.