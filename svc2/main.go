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
		"service1",
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport, nil),
	)

	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	var fn opengintracing.ParentSpanReferenceFunc
	fn = func(sc opentracing.SpanContext) opentracing.StartSpanOption {
		return opentracing.ChildOf(sc)
	}

	// Set up routes
	r := gin.Default()
	r.POST("",
		opengintracing.SpanFromHeadersHttpFmt(tracer, "service2", fn, false),
		handler)
	r.Run(":8002")
}


func handler(c *gin.Context) {
	_, found := opengintracing.GetSpan(c)
	if found == false {
		fmt.Println("Span not found")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println("Incoming Headers")
	for k, v := range c.Request.Header {
		fmt.Printf("%s: %s\n", k, v)
	}
	c.Status(http.StatusOK)
}
