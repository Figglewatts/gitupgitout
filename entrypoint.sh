#!/usr/bin/env sh

set -e

USERNAME="${USERNAME:-default}"
GROUPNAME="${GROUPNAME:-${USERNAME}}"
UID="${UID:-1000}"
GID="${GID:-${UID}}"

if grep -q -E "^${GROUPNAME}:" /etc/group > /dev/null 2>&1; then
  echo "INFO: Group exists; skipping creation"
else
  addgroup -g "${GID}" "${GROUPNAME}"
fi

if id -u "${USERNAME}" > /dev/null 2>&1; then
  echo "INFO: User exists; skipping creation"
else
  adduser -u "${UID}" -G "${GROUPNAME}" -h "/home/${USERNAME}" -s /bin/sh -D "${USERNAME}"
  mkdir -p "/home/${USERNAME}"
  chown "${USERNAME}:${GROUPNAME}" "/home/${USERNAME}"
fi

mkdir -p "/app"
chown -R "${USERNAME}:${GROUPNAME}" "/app"

exec su-exec "${USERNAME}" "${@}"