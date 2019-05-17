#!/bin/bash
set -e

postconf -e myhostname=$DOMAIN
postconf -e mynetworks=0.0.0.0/0  # Note: don't expose port 25 to the outside world!

# OpenDKIM
postconf -e milter_protocol=2
postconf -e milter_default_action=accept
postconf -e smtpd_milters=inet:localhost:12301
postconf -e non_smtpd_milters=inet:localhost:12301

export SOCKET="inet:12301@localhost"

cat > /etc/opendkim/opendkim.conf <<EOF
AutoRestart             Yes
AutoRestartRate         10/1h
UMask                   002
Syslog                  yes
SyslogSuccess           Yes
LogWhy                  Yes

Canonicalization        relaxed/simple

ExternalIgnoreList      refile:/etc/opendkim/TrustedHosts
InternalHosts           refile:/etc/opendkim/TrustedHosts
KeyTable                refile:/etc/opendkim/KeyTable
SigningTable            refile:/etc/opendkim/SigningTable

Mode                    sv
PidFile                 /var/run/opendkim/opendkim.pid
SignatureAlgorithm      rsa-sha256

UserID                  opendkim:opendkim

Socket                  inet:12301@localhost
EOF

cat > /etc/opendkim/TrustedHosts <<EOF
127.0.0.1
localhost

*.$DOMAIN
EOF

if [ -f /run/secrets/opendkim_private ]; then
cp /run/secrets/opendkim_private /opendkim.private
else
cp /dkim/mail.private /opendkim.private
fi

cat > /etc/opendkim/KeyTable <<EOF
mail._domainkey.$DOMAIN $DOMAIN:mail:/opendkim.private
EOF

cat > /etc/opendkim/SigningTable <<EOF
*@$DOMAIN mail._domainkey.$DOMAIN
EOF

chown opendkim:opendkim /opendkim.private
chmod 440 /opendkim.private

exec "$@"