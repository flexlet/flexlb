// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	flags "github.com/jessevdk/go-flags"

	"gitee.com/flexlb/flexlb-api/common"
	"gitee.com/flexlb/flexlb-api/config"
	"gitee.com/flexlb/flexlb-api/handlers"
	"gitee.com/flexlb/flexlb-api/restapi/operations"
	"gitee.com/flexlb/flexlb-api/wacher"

	"github.com/00ahui/utils"
)

// flexlb command line options
var options = struct {
	ConfigFile string `short:"c" long:"config-file" env:"FLEXLB_CONFIG_FILE" description:"Configuration file" default:"/etc/flexlb/flexlb-api-config.yaml" group:"flexlb"`
	Version    bool   `short:"v" long:"version" description:"Show version" group:"flexlb"`
}{}

// configure flexlb command line options
func configureFlags(api *operations.FlexlbAPI) {
	optionsGroup := swag.CommandLineOptionsGroup{
		ShortDescription: "FlexLB options",
		LongDescription:  "Flexible load balancer options",
		Options:          &options,
	}
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		optionsGroup,
	}
}

func configureAPI(s *Server) http.Handler {

	// show version and exit
	if options.Version {
		fmt.Printf("%s v%s\n", common.PROJECT, common.VERSION)
		os.Exit(0)
	}

	// load config file
	config.LoadConfig(options.ConfigFile)

	// set log level
	utils.LogLevel = int(config.LB.LogLevel)

	// set server params
	s.Host = config.LB.Host
	s.Port = int(config.LB.Port)
	s.TLSHost = config.LB.TLSHost
	s.TLSPort = int(config.LB.TLSPort)
	s.TLSCertificate = flags.Filename(config.LB.TLSCert)
	s.TLSCertificateKey = flags.Filename(config.LB.TLSKey)
	s.TLSCACertificate = flags.Filename(config.LB.TLSCACert)

	// load saved instances
	config.LoadInstances()

	// update cluster isntance
	config.UpdateClusterInstance()

	// start cluster gossip
	config.StartClusterGossip()

	// start wachers
	go wacher.StartInstanceWatcher()

	// setup service handlers
	s.api.ServiceReadyzHandler = &handlers.ReadyzHandlerImpl{}

	// setup instance handlers
	s.api.InstanceCreateHandler = &handlers.InstanceCreateHandlerImpl{}
	s.api.InstanceListHandler = &handlers.InstanceListHandlerImpl{}
	s.api.InstanceGetHandler = &handlers.InstanceGetHandlerImpl{}
	s.api.InstanceModifyHandler = &handlers.InstanceModifyHandlerImpl{}
	s.api.InstanceDeleteHandler = &handlers.InstanceDeleteHandlerImpl{}
	s.api.InstanceStopHandler = &handlers.InstanceStopHandlerImpl{}
	s.api.InstanceStartHandler = &handlers.InstanceStartHandlerImpl{}

	s.api.PreServerShutdown = func() {
		wacher.StopInstanceWacher()
	}

	s.api.ServerShutdown = func() {}

	// Swagger generated code:

	s.api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	// api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	s.api.JSONConsumer = runtime.JSONConsumer()

	s.api.JSONProducer = runtime.JSONProducer()

	return setupGlobalMiddleware(s.api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
