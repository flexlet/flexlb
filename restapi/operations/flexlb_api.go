// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/flexlet/flexlb/restapi/operations/instance"
	"github.com/flexlet/flexlb/restapi/operations/service"
)

// NewFlexlbAPI creates a new Flexlb instance
func NewFlexlbAPI(spec *loads.Document) *FlexlbAPI {
	return &FlexlbAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		InstanceCreateHandler: instance.CreateHandlerFunc(func(params instance.CreateParams) middleware.Responder {
			return middleware.NotImplemented("operation instance.Create has not yet been implemented")
		}),
		InstanceDeleteHandler: instance.DeleteHandlerFunc(func(params instance.DeleteParams) middleware.Responder {
			return middleware.NotImplemented("operation instance.Delete has not yet been implemented")
		}),
		InstanceGetHandler: instance.GetHandlerFunc(func(params instance.GetParams) middleware.Responder {
			return middleware.NotImplemented("operation instance.Get has not yet been implemented")
		}),
		InstanceListHandler: instance.ListHandlerFunc(func(params instance.ListParams) middleware.Responder {
			return middleware.NotImplemented("operation instance.List has not yet been implemented")
		}),
		InstanceModifyHandler: instance.ModifyHandlerFunc(func(params instance.ModifyParams) middleware.Responder {
			return middleware.NotImplemented("operation instance.Modify has not yet been implemented")
		}),
		ServiceReadyzHandler: service.ReadyzHandlerFunc(func(params service.ReadyzParams) middleware.Responder {
			return middleware.NotImplemented("operation service.Readyz has not yet been implemented")
		}),
		InstanceStartHandler: instance.StartHandlerFunc(func(params instance.StartParams) middleware.Responder {
			return middleware.NotImplemented("operation instance.Start has not yet been implemented")
		}),
		InstanceStopHandler: instance.StopHandlerFunc(func(params instance.StopParams) middleware.Responder {
			return middleware.NotImplemented("operation instance.Stop has not yet been implemented")
		}),
	}
}

/*FlexlbAPI Flexible load balancer API to control keepalived and haproxy
 */
type FlexlbAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// InstanceCreateHandler sets the operation handler for the create operation
	InstanceCreateHandler instance.CreateHandler
	// InstanceDeleteHandler sets the operation handler for the delete operation
	InstanceDeleteHandler instance.DeleteHandler
	// InstanceGetHandler sets the operation handler for the get operation
	InstanceGetHandler instance.GetHandler
	// InstanceListHandler sets the operation handler for the list operation
	InstanceListHandler instance.ListHandler
	// InstanceModifyHandler sets the operation handler for the modify operation
	InstanceModifyHandler instance.ModifyHandler
	// ServiceReadyzHandler sets the operation handler for the readyz operation
	ServiceReadyzHandler service.ReadyzHandler
	// InstanceStartHandler sets the operation handler for the start operation
	InstanceStartHandler instance.StartHandler
	// InstanceStopHandler sets the operation handler for the stop operation
	InstanceStopHandler instance.StopHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *FlexlbAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *FlexlbAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *FlexlbAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *FlexlbAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *FlexlbAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *FlexlbAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *FlexlbAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *FlexlbAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *FlexlbAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the FlexlbAPI
func (o *FlexlbAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.InstanceCreateHandler == nil {
		unregistered = append(unregistered, "instance.CreateHandler")
	}
	if o.InstanceDeleteHandler == nil {
		unregistered = append(unregistered, "instance.DeleteHandler")
	}
	if o.InstanceGetHandler == nil {
		unregistered = append(unregistered, "instance.GetHandler")
	}
	if o.InstanceListHandler == nil {
		unregistered = append(unregistered, "instance.ListHandler")
	}
	if o.InstanceModifyHandler == nil {
		unregistered = append(unregistered, "instance.ModifyHandler")
	}
	if o.ServiceReadyzHandler == nil {
		unregistered = append(unregistered, "service.ReadyzHandler")
	}
	if o.InstanceStartHandler == nil {
		unregistered = append(unregistered, "instance.StartHandler")
	}
	if o.InstanceStopHandler == nil {
		unregistered = append(unregistered, "instance.StopHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *FlexlbAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *FlexlbAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	return nil
}

// Authorizer returns the registered authorizer
func (o *FlexlbAPI) Authorizer() runtime.Authorizer {
	return nil
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *FlexlbAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *FlexlbAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *FlexlbAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the flexlb API
func (o *FlexlbAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *FlexlbAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/instances"] = instance.NewCreate(o.context, o.InstanceCreateHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/instances/{name}"] = instance.NewDelete(o.context, o.InstanceDeleteHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/instances/{name}"] = instance.NewGet(o.context, o.InstanceGetHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/instances"] = instance.NewList(o.context, o.InstanceListHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/instances"] = instance.NewModify(o.context, o.InstanceModifyHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/readyz"] = service.NewReadyz(o.context, o.ServiceReadyzHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/instances/{name}/start"] = instance.NewStart(o.context, o.InstanceStartHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/instances/{name}/stop"] = instance.NewStop(o.context, o.InstanceStopHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *FlexlbAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *FlexlbAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *FlexlbAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *FlexlbAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *FlexlbAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[method][path] = builder(h)
	}
}
