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

	"gitee.com/flexlb/flexlb-api/common"
	"gitee.com/flexlb/flexlb-api/config"
	"gitee.com/flexlb/flexlb-api/handlers"
	"gitee.com/flexlb/flexlb-api/restapi/operations"
	"gitee.com/flexlb/flexlb-api/wacher"
)

// flexlb command line options
var options = struct {
	ConfigFile string `short:"c" long:"config-file" description:"Configuration file" default:"/etc/flexlb/flexlb-api-config.yaml" group:"flexlb"`
	Debug      bool   `short:"d" long:"debug" description:"Debug mode, show verbose logs" group:"flexlb"`
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

func configureAPI(api *operations.FlexlbAPI) http.Handler {

	// show version and exit
	if options.Version {
		fmt.Printf("%s v%s\n", common.PROJECT, common.VERSION)
		os.Exit(0)
	}

	// set debug mode
	common.Debug = options.Debug

	// load config file
	config.LoadConfig(options.ConfigFile)

	// load saved instances
	config.LoadInstances()

	// start wachers
	go wacher.StartInstanceWatcher()

	// setup service handlers
	api.ServiceReadyzHandler = &handlers.ReadyzHandlerImpl{}

	// setup instance handlers
	api.InstanceCreateHandler = &handlers.InstanceCreateHandlerImpl{}
	api.InstanceListHandler = &handlers.InstanceListHandlerImpl{}
	api.InstanceGetHandler = &handlers.InstanceGetHandlerImpl{}
	api.InstanceModifyHandler = &handlers.InstanceModifyHandlerImpl{}
	api.InstanceDeleteHandler = &handlers.InstanceDeleteHandlerImpl{}
	api.InstanceStopHandler = &handlers.InstanceStopHandlerImpl{}
	api.InstanceStartHandler = &handlers.InstanceStartHandlerImpl{}

	api.PreServerShutdown = func() {
		wacher.StopInstanceWacher()
	}

	api.ServerShutdown = func() {}

	// Swagger generated code:

	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	// api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
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
