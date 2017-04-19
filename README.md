## usage
```sh
//service with network will create a vip 
docker network create -d overlay test
docker service  create -e "VIRTUAL_HOST=**.com" --network test ** 
docker run -d  -v /var/run/docker.sock:/tmp/docker.sock:ro -p 8095:80 --network test wanghaibo/ingress 
```

## build
```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
docker build ./ -t wanghaibo/ingress -f Dockerfile.alpine  --no-cache 
```
