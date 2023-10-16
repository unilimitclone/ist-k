#!/bin/bash
# Set permissions and owners.
chown -R ${PUID}:${PGID} /opt/alist/
umask ${UMASK}
# Create data folder.
mkdir -p /opt/alist/data
# Config Alist
/main
# main
exec su-exec ${PUID}:${PGID} /usr/bin/supervisord -c /etc/supervisord.conf