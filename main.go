package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"

	ginLogger "github.com/gin-contrib/logger"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	beApi "github.com/Sourceware-Lab/go-huma-gin-postgres-template/api"
	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/config"
	DBpostgres "github.com/Sourceware-Lab/go-huma-gin-postgres-template/database/postgres"
)

const apiVersion = "0.0.1"

type Options struct {
	Port int `help:"Port to listen on" short:"p"`
}

func (o *Options) loadFromViper() {
	o.Port = config.Config.Port
}

func initProvider() func() { //nolint:funlen
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("REPLACEME2"),
		),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create resource")
	}

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(config.Config.OTELExporterOTLPEndpoint))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create the collector metric exporter")
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExp,
				sdkmetric.WithInterval(2*time.Second), //nolint:mnd
			),
		),
	)
	otel.SetMeterProvider(meterProvider)

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(config.Config.OTELExporterOTLPEndpoint),
	)

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create the collector trace exporter")
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
		// pushes any last exports to the receiver
		if err := meterProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}

//nolint:ireturn
func getCli() humacli.CLI {
	log.Info().Msg("Creating CLI")

	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		log.Info().Msg("Starting server")
		options.loadFromViper()

		if config.Config.ReleaseMode {
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}

		gin.DisableConsoleColor()
		gin.DefaultWriter = log.Logger
		gin.DefaultErrorWriter = log.Logger

		// Create a new router & API
		router := gin.New()
		router.Use(otelgin.Middleware("REPLACEME2"))
		router.Use(ginLogger.SetLogger())
		api := humagin.New(router, huma.DefaultConfig("Example API", apiVersion))

		beApi.AddRoutes(api)

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			log.Info().Msg(fmt.Sprintf("Starting server on port %d...\n", options.Port))
			server := &http.Server{
				IdleTimeout:       300 * time.Second, //nolint:mnd
				ReadTimeout:       300 * time.Second, //nolint:mnd
				WriteTimeout:      300 * time.Second, //nolint:mnd
				ReadHeaderTimeout: 10 * time.Second,  //nolint:mnd
				Addr:              fmt.Sprintf(":%d", options.Port),
				Handler:           router,
			}
			_ = server.ListenAndServe()
		})
	})

	return cli
}

func main() {
	config.LoadConfig()
	config.InitLogger()

	shutdown := initProvider()
	defer shutdown()

	DBpostgres.Open(config.Config.DatabaseDSN)

	defer DBpostgres.Close()
	DBpostgres.RunMigrations()

	cli := getCli()
	cli.Run()
}
