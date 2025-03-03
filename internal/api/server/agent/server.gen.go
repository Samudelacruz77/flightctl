// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/flightctl/flightctl/api/v1alpha1"
	externalRef0 "github.com/flightctl/flightctl/api/v1alpha1"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /api/v1/devices/{name}/rendered)
	GetRenderedDeviceSpec(w http.ResponseWriter, r *http.Request, name string, params GetRenderedDeviceSpecParams)

	// (PUT /api/v1/devices/{name}/status)
	ReplaceDeviceStatus(w http.ResponseWriter, r *http.Request, name string)

	// (POST /api/v1/enrollmentrequests)
	CreateEnrollmentRequest(w http.ResponseWriter, r *http.Request)

	// (GET /api/v1/enrollmentrequests/{name})
	ReadEnrollmentRequest(w http.ResponseWriter, r *http.Request, name string)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /api/v1/devices/{name}/rendered)
func (_ Unimplemented) GetRenderedDeviceSpec(w http.ResponseWriter, r *http.Request, name string, params GetRenderedDeviceSpecParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (PUT /api/v1/devices/{name}/status)
func (_ Unimplemented) ReplaceDeviceStatus(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /api/v1/enrollmentrequests)
func (_ Unimplemented) CreateEnrollmentRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /api/v1/enrollmentrequests/{name})
func (_ Unimplemented) ReadEnrollmentRequest(w http.ResponseWriter, r *http.Request, name string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetRenderedDeviceSpec operation middleware
func (siw *ServerInterfaceWrapper) GetRenderedDeviceSpec(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithOptions("simple", "name", chi.URLParam(r, "name"), &name, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetRenderedDeviceSpecParams

	// ------------- Optional query parameter "knownRenderedVersion" -------------

	err = runtime.BindQueryParameter("form", true, false, "knownRenderedVersion", r.URL.Query(), &params.KnownRenderedVersion)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "knownRenderedVersion", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRenderedDeviceSpec(w, r, name, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ReplaceDeviceStatus operation middleware
func (siw *ServerInterfaceWrapper) ReplaceDeviceStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithOptions("simple", "name", chi.URLParam(r, "name"), &name, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ReplaceDeviceStatus(w, r, name)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateEnrollmentRequest operation middleware
func (siw *ServerInterfaceWrapper) CreateEnrollmentRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateEnrollmentRequest(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ReadEnrollmentRequest operation middleware
func (siw *ServerInterfaceWrapper) ReadEnrollmentRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithOptions("simple", "name", chi.URLParam(r, "name"), &name, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ReadEnrollmentRequest(w, r, name)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/v1/devices/{name}/rendered", wrapper.GetRenderedDeviceSpec)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/api/v1/devices/{name}/status", wrapper.ReplaceDeviceStatus)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/v1/enrollmentrequests", wrapper.CreateEnrollmentRequest)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/v1/enrollmentrequests/{name}", wrapper.ReadEnrollmentRequest)
	})

	return r
}

type GetRenderedDeviceSpecRequestObject struct {
	Name   string `json:"name"`
	Params GetRenderedDeviceSpecParams
}

type GetRenderedDeviceSpecResponseObject interface {
	VisitGetRenderedDeviceSpecResponse(w http.ResponseWriter) error
}

type GetRenderedDeviceSpec200JSONResponse externalRef0.RenderedDeviceSpec

func (response GetRenderedDeviceSpec200JSONResponse) VisitGetRenderedDeviceSpecResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetRenderedDeviceSpec204Response struct {
}

func (response GetRenderedDeviceSpec204Response) VisitGetRenderedDeviceSpecResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type GetRenderedDeviceSpec401JSONResponse externalRef0.Error

func (response GetRenderedDeviceSpec401JSONResponse) VisitGetRenderedDeviceSpecResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetRenderedDeviceSpec404JSONResponse externalRef0.Error

func (response GetRenderedDeviceSpec404JSONResponse) VisitGetRenderedDeviceSpecResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetRenderedDeviceSpec409JSONResponse externalRef0.Error

func (response GetRenderedDeviceSpec409JSONResponse) VisitGetRenderedDeviceSpecResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type ReplaceDeviceStatusRequestObject struct {
	Name string `json:"name"`
	Body *ReplaceDeviceStatusJSONRequestBody
}

type ReplaceDeviceStatusResponseObject interface {
	VisitReplaceDeviceStatusResponse(w http.ResponseWriter) error
}

type ReplaceDeviceStatus200JSONResponse externalRef0.Device

func (response ReplaceDeviceStatus200JSONResponse) VisitReplaceDeviceStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ReplaceDeviceStatus400JSONResponse externalRef0.Error

func (response ReplaceDeviceStatus400JSONResponse) VisitReplaceDeviceStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type ReplaceDeviceStatus401JSONResponse externalRef0.Error

func (response ReplaceDeviceStatus401JSONResponse) VisitReplaceDeviceStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type ReplaceDeviceStatus404JSONResponse externalRef0.Error

func (response ReplaceDeviceStatus404JSONResponse) VisitReplaceDeviceStatusResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type CreateEnrollmentRequestRequestObject struct {
	Body *CreateEnrollmentRequestJSONRequestBody
}

type CreateEnrollmentRequestResponseObject interface {
	VisitCreateEnrollmentRequestResponse(w http.ResponseWriter) error
}

type CreateEnrollmentRequest201JSONResponse externalRef0.EnrollmentRequest

func (response CreateEnrollmentRequest201JSONResponse) VisitCreateEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type CreateEnrollmentRequest400JSONResponse externalRef0.Error

func (response CreateEnrollmentRequest400JSONResponse) VisitCreateEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type CreateEnrollmentRequest401JSONResponse externalRef0.Error

func (response CreateEnrollmentRequest401JSONResponse) VisitCreateEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type CreateEnrollmentRequest403JSONResponse externalRef0.Error

func (response CreateEnrollmentRequest403JSONResponse) VisitCreateEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	return json.NewEncoder(w).Encode(response)
}

type CreateEnrollmentRequest409JSONResponse externalRef0.Error

func (response CreateEnrollmentRequest409JSONResponse) VisitCreateEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type CreateEnrollmentRequest503JSONResponse externalRef0.Error

func (response CreateEnrollmentRequest503JSONResponse) VisitCreateEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

type ReadEnrollmentRequestRequestObject struct {
	Name string `json:"name"`
}

type ReadEnrollmentRequestResponseObject interface {
	VisitReadEnrollmentRequestResponse(w http.ResponseWriter) error
}

type ReadEnrollmentRequest200JSONResponse externalRef0.EnrollmentRequest

func (response ReadEnrollmentRequest200JSONResponse) VisitReadEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ReadEnrollmentRequest401JSONResponse externalRef0.Error

func (response ReadEnrollmentRequest401JSONResponse) VisitReadEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type ReadEnrollmentRequest403JSONResponse externalRef0.Error

func (response ReadEnrollmentRequest403JSONResponse) VisitReadEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	return json.NewEncoder(w).Encode(response)
}

type ReadEnrollmentRequest404JSONResponse externalRef0.Error

func (response ReadEnrollmentRequest404JSONResponse) VisitReadEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type ReadEnrollmentRequest503JSONResponse externalRef0.Error

func (response ReadEnrollmentRequest503JSONResponse) VisitReadEnrollmentRequestResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (GET /api/v1/devices/{name}/rendered)
	GetRenderedDeviceSpec(ctx context.Context, request GetRenderedDeviceSpecRequestObject) (GetRenderedDeviceSpecResponseObject, error)

	// (PUT /api/v1/devices/{name}/status)
	ReplaceDeviceStatus(ctx context.Context, request ReplaceDeviceStatusRequestObject) (ReplaceDeviceStatusResponseObject, error)

	// (POST /api/v1/enrollmentrequests)
	CreateEnrollmentRequest(ctx context.Context, request CreateEnrollmentRequestRequestObject) (CreateEnrollmentRequestResponseObject, error)

	// (GET /api/v1/enrollmentrequests/{name})
	ReadEnrollmentRequest(ctx context.Context, request ReadEnrollmentRequestRequestObject) (ReadEnrollmentRequestResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// GetRenderedDeviceSpec operation middleware
func (sh *strictHandler) GetRenderedDeviceSpec(w http.ResponseWriter, r *http.Request, name string, params GetRenderedDeviceSpecParams) {
	var request GetRenderedDeviceSpecRequestObject

	request.Name = name
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetRenderedDeviceSpec(ctx, request.(GetRenderedDeviceSpecRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetRenderedDeviceSpec")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetRenderedDeviceSpecResponseObject); ok {
		if err := validResponse.VisitGetRenderedDeviceSpecResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// ReplaceDeviceStatus operation middleware
func (sh *strictHandler) ReplaceDeviceStatus(w http.ResponseWriter, r *http.Request, name string) {
	var request ReplaceDeviceStatusRequestObject

	request.Name = name

	var body ReplaceDeviceStatusJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ReplaceDeviceStatus(ctx, request.(ReplaceDeviceStatusRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ReplaceDeviceStatus")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ReplaceDeviceStatusResponseObject); ok {
		if err := validResponse.VisitReplaceDeviceStatusResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// CreateEnrollmentRequest operation middleware
func (sh *strictHandler) CreateEnrollmentRequest(w http.ResponseWriter, r *http.Request) {
	var request CreateEnrollmentRequestRequestObject

	var body CreateEnrollmentRequestJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreateEnrollmentRequest(ctx, request.(CreateEnrollmentRequestRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateEnrollmentRequest")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreateEnrollmentRequestResponseObject); ok {
		if err := validResponse.VisitCreateEnrollmentRequestResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// ReadEnrollmentRequest operation middleware
func (sh *strictHandler) ReadEnrollmentRequest(w http.ResponseWriter, r *http.Request, name string) {
	var request ReadEnrollmentRequestRequestObject

	request.Name = name

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ReadEnrollmentRequest(ctx, request.(ReadEnrollmentRequestRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ReadEnrollmentRequest")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ReadEnrollmentRequestResponseObject); ok {
		if err := validResponse.VisitReadEnrollmentRequestResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
