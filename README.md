# Docker Postfix SMTP Server

A simple docker image containing the Postfix SMTP server, and OpenDKIM

This doesn't contain any auth, since it should only be sending emails from the internal docker network it's bound to.

## Environment
`DOMAIN`: The domain to send from

## Volumes / Secrets
This image can accept the DKIM private key as either a bindmount, or as a swarm secret. Swarm secret takes priority.

### Secret
Secret name: `opendkim_private`
Value of the secret should contain the RSA private key generated for OpenDKIM.

### Volume
Bind mount: `local_dkim_path:/dkim`
`local_dkim_path` should contain a single file - `mail.private`, containing the OpenDKIM private key.