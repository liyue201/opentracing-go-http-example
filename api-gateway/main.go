package main

import (
	"fmt"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"log"
	"net/http"
)

func main() {

	transport, err := zipkin.NewHTTPTransport(
		"http://localhost:9411/api/v1/spans",
		zipkin.HTTPBatchSize(10),
		zipkin.HTTPLogger(jlog.StdLogger),
	)
	if err != nil {
		log.Fatalf("Cannot initialize Zipkin HTTP transport: %v", err)
		panic(err.Error())
	}
	tracer, closer := jaeger.NewTracer(
		"api-gateway",
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport),
	)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// Set up routes
	r := gin.Default()
	r.POST("/service1",
		opengintracing.NewSpan(tracer, "forward to service 1"),
		service1handler)
	r.POST("/service2",
		opengintracing.NewSpan(tracer, "forward to service 2"),
		service2handler)
	r.Run(":8000")
}

func printHeaders(message string, header http.Header) {
	fmt.Println(message)
	for k, v := range header {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func service1handler(c *gin.Context) {
	span, found := opengintracing.GetSpan(c)
	if found == false {
		fmt.Println("Span not found")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	req, _ := http.NewRequest("POST", "http://localhost:8001", nil)

	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	printHeaders("Incoming Headers", c.Request.Header)
	printHeaders("Outgoing Headers", req.Header)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func service2handler(c *gin.Context) {
	span, found := opengintracing.GetSpan(c)
	if found == false {
		fmt.Println("Span not found")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	req, _ := http.NewRequest("POST", "http://localhost:8002", nil)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	printHeaders("Incoming Headers", c.Request.Header)
	printHeaders("Outgoing Headers", req.Header)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
