# opentracing-go-http-example

安装jaeger

```
docker run --name jaeger  -d \
-e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
-p 5775:5775/udp \
-p 6831:6831/udp \
-p 6832:6832/udp \
-p 5778:5778 \
-p 16686:16686 \
-p 14268:14268 \
-p 9411:9411 \
jaegertracing/all-in-one:latest
 ```

启动 api-gateway

 ```
 cd api-gateway
 go run main.go

 ```


启动 service1
 ```
 cd svc1
 go run main.go

 ```

启动 service2
 ```
 cd svc2
 go run main.go

 ```


测试
```
curl http://localhost:8000/service1 -X POST -d "hello"
curl http://localhost:8000/service2 -X POST -d "world"

```

