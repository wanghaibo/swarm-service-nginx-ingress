#!/bin/bash
set -e

# Warn if the DOCKER_HOST socket does not exist
if [[ $DOCKER_HOST == unix://* ]]; then
	socket_file=${DOCKER_HOST#unix://}
	if ! [ -S $socket_file ]; then
		cat >&2 <<-EOT
			ERROR: you need to share your Docker host socket with a volume at $socket_file
			Typically you should run your jwilder/nginx-proxy with: \`-v /var/run/docker.sock:$socket_file:ro\`
			See the documentation at http://git.io/vZaGJ
		EOT
		socketMissing=1
	fi
fi

if [ "$1" = "" ]; then
    if [ "$socketMissing" = 1 ]; then
        exit 1
    else
        swarm-service-nginx-ingress -tpl="/app/nginx.tpl" -dst="/etc/nginx/conf.d/default.conf"
        nginx
    fi
else 
    exec "$@"
fi

