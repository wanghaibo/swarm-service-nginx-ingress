## usage
```sh
//service with network will create a vip 
docker network create -d overlay --attachable test 
docker service  create -e "VIRTUAL_HOST=**.com" --network test ** 
docker run -d  -v /var/run/docker.sock:/tmp/docker.sock:ro -p 8095:80 --network test --name ingress wanghaibo/ingress 
```

## build
```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
docker build ./ -t wanghaibo/ingress -f Dockerfile.alpine  --no-cache 
```

## reload
```sh
docker exec ingress /app/docker-entrypoint.sh  reload
```
## todo
waiting for events of services https://github.com/moby/moby/pull/32421
