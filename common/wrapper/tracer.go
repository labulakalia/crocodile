package wrapper

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/metadata"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const contextTracerKey = "Tracer-context"

// sf sampling frequency
var sf = 100

func init() {
	rand.Seed(time.Now().Unix())
}

// SetSamplingFrequency 设置采样频率
// 0 <= n <= 100
func SetSamplingFrequency(n int) {
	sf = n
}

// TracerWrapper tracer 中间件
func TracerWrapper(c *gin.Context) {
	md := make(map[string]string)
	spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	sp := opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
	defer sp.Finish()

	if err := opentracing.GlobalTracer().Inject(sp.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(md)); err != nil {
		logging.Error(err)
	}

	ctx := context.TODO()
	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	c.Set(contextTracerKey, ctx)

	c.Next()

	statusCode := c.Writer.Status()
	ext.HTTPStatusCode.Set(sp, uint16(statusCode))
	ext.HTTPMethod.Set(sp, c.Request.Method)
	ext.HTTPUrl.Set(sp, c.Request.URL.EscapedPath())
	if statusCode >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	} else if rand.Intn(100) > sf {
		ext.SamplingPriority.Set(sp, 0)
	}

}

// ContextWithSpan 返回context
func ContextWithSpan(c *gin.Context) (ctx context.Context, ok bool) {
	v, exist := c.Get(contextTracerKey)
	if exist == false {
		ok = false
		ctx = context.TODO()
		return
	}

	ctx, ok = v.(context.Context)
	return
}

// NewTracer 创建一个jaeger Tracer
func NewTracer(servicename string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	return tracer, closer, err
}
