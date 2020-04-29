## start consul

```
sudo docker-compose -f docker/docker-compose.yml up
```


## start register-service

```
./register -consul.host localhost -consul.port 8500 -service.host 192.168.1.100 -service.port 9000
```

## start discover-service