// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	Api_keyScopes = "api_key.Scopes"
)

// Defines values for SearchKindleParamsRegion.
const (
	Au SearchKindleParamsRegion = "au"
	Ca SearchKindleParamsRegion = "ca"
	De SearchKindleParamsRegion = "de"
	Es SearchKindleParamsRegion = "es"
	Fr SearchKindleParamsRegion = "fr"
	In SearchKindleParamsRegion = "in"
	It SearchKindleParamsRegion = "it"
	Jp SearchKindleParamsRegion = "jp"
	Uk SearchKindleParamsRegion = "uk"
	Us SearchKindleParamsRegion = "us"
)

// BookMetadata defines model for BookMetadata.
type BookMetadata struct {
	Asin   *string `json:"asin,omitempty"`
	Author *string `json:"author,omitempty"`

	// Cover URL to the cover image
	Cover       *string `json:"cover,omitempty"`
	Description *string `json:"description,omitempty"`

	// Duration Duration in seconds
	Duration      *int              `json:"duration,omitempty"`
	Genres        *[]string         `json:"genres,omitempty"`
	Isbn          *string           `json:"isbn,omitempty"`
	Language      *string           `json:"language,omitempty"`
	Narrator      *string           `json:"narrator,omitempty"`
	PublishedYear *string           `json:"publishedYear,omitempty"`
	Publisher     *string           `json:"publisher,omitempty"`
	Series        *[]SeriesMetadata `json:"series,omitempty"`
	Subtitle      *string           `json:"subtitle,omitempty"`
	Tags          *[]string         `json:"tags,omitempty"`
	Title         string            `json:"title"`
}

// SeriesMetadata defines model for SeriesMetadata.
type SeriesMetadata struct {
	Sequence *string `json:"sequence,omitempty"`
	Series   string  `json:"series"`
}

// Author defines model for author.
type Author = string

// Query defines model for query.
type Query = string

// N200 defines model for 200.
type N200 struct {
	Matches *[]BookMetadata `json:"matches,omitempty"`
}

// N400 defines model for 400.
type N400 struct {
	Error *string `json:"error,omitempty"`
}

// N401 defines model for 401.
type N401 struct {
	Error *string `json:"error,omitempty"`
}

// N500 defines model for 500.
type N500 struct {
	Error *string `json:"error,omitempty"`
}

// SearchGoodreadsParams defines parameters for SearchGoodreads.
type SearchGoodreadsParams struct {
	Query  Query   `form:"query" json:"query"`
	Author *Author `form:"author,omitempty" json:"author,omitempty"`
}

// SearchKindleParams defines parameters for SearchKindle.
type SearchKindleParams struct {
	Query  Query   `form:"query" json:"query"`
	Author *Author `form:"author,omitempty" json:"author,omitempty"`
}

// SearchKindleParamsRegion defines parameters for SearchKindle.
type SearchKindleParamsRegion string

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Search for books using goodreads
	// (GET /goodreads/search)
	SearchGoodreads(w http.ResponseWriter, r *http.Request, params SearchGoodreadsParams)
	// Search for books using kindle
	// (GET /kindle/{region}/search)
	SearchKindle(w http.ResponseWriter, r *http.Request, region SearchKindleParamsRegion, params SearchKindleParams)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Search for books using goodreads
// (GET /goodreads/search)
func (_ Unimplemented) SearchGoodreads(w http.ResponseWriter, r *http.Request, params SearchGoodreadsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Search for books using kindle
// (GET /kindle/{region}/search)
func (_ Unimplemented) SearchKindle(w http.ResponseWriter, r *http.Request, region SearchKindleParamsRegion, params SearchKindleParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// SearchGoodreads operation middleware
func (siw *ServerInterfaceWrapper) SearchGoodreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params SearchGoodreadsParams

	// ------------- Required query parameter "query" -------------

	if paramValue := r.URL.Query().Get("query"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "query"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "query", r.URL.Query(), &params.Query)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "query", Err: err})
		return
	}

	// ------------- Optional query parameter "author" -------------

	err = runtime.BindQueryParameter("form", true, false, "author", r.URL.Query(), &params.Author)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "author", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SearchGoodreads(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SearchKindle operation middleware
func (siw *ServerInterfaceWrapper) SearchKindle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "region" -------------
	var region SearchKindleParamsRegion

	err = runtime.BindStyledParameterWithOptions("simple", "region", chi.URLParam(r, "region"), &region, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: false})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "region", Err: err})
		return
	}

	ctx = context.WithValue(ctx, Api_keyScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params SearchKindleParams

	// ------------- Required query parameter "query" -------------

	if paramValue := r.URL.Query().Get("query"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "query"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "query", r.URL.Query(), &params.Query)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "query", Err: err})
		return
	}

	// ------------- Optional query parameter "author" -------------

	err = runtime.BindQueryParameter("form", true, false, "author", r.URL.Query(), &params.Author)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "author", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SearchKindle(w, r, region, params)
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
		r.Get(options.BaseURL+"/goodreads/search", wrapper.SearchGoodreads)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/kindle/{region}/search", wrapper.SearchKindle)
	})

	return r
}

type N200JSONResponse struct {
	Matches *[]BookMetadata `json:"matches,omitempty"`
}

type N400JSONResponse struct {
	Error *string `json:"error,omitempty"`
}

type N401JSONResponse struct {
	Error *string `json:"error,omitempty"`
}

type N500JSONResponse struct {
	Error *string `json:"error,omitempty"`
}

type SearchGoodreadsRequestObject struct {
	Params SearchGoodreadsParams
}

type SearchGoodreadsResponseObject interface {
	VisitSearchGoodreadsResponse(w http.ResponseWriter) error
}

type SearchGoodreads200JSONResponse struct{ N200JSONResponse }

func (response SearchGoodreads200JSONResponse) VisitSearchGoodreadsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type SearchGoodreads400JSONResponse struct{ N400JSONResponse }

func (response SearchGoodreads400JSONResponse) VisitSearchGoodreadsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type SearchGoodreads401JSONResponse struct{ N401JSONResponse }

func (response SearchGoodreads401JSONResponse) VisitSearchGoodreadsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type SearchGoodreads500JSONResponse struct{ N500JSONResponse }

func (response SearchGoodreads500JSONResponse) VisitSearchGoodreadsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type SearchKindleRequestObject struct {
	Region SearchKindleParamsRegion `json:"region,omitempty"`
	Params SearchKindleParams
}

type SearchKindleResponseObject interface {
	VisitSearchKindleResponse(w http.ResponseWriter) error
}

type SearchKindle200JSONResponse struct{ N200JSONResponse }

func (response SearchKindle200JSONResponse) VisitSearchKindleResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type SearchKindle400JSONResponse struct{ N400JSONResponse }

func (response SearchKindle400JSONResponse) VisitSearchKindleResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type SearchKindle401JSONResponse struct{ N401JSONResponse }

func (response SearchKindle401JSONResponse) VisitSearchKindleResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type SearchKindle500JSONResponse struct{ N500JSONResponse }

func (response SearchKindle500JSONResponse) VisitSearchKindleResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Search for books using goodreads
	// (GET /goodreads/search)
	SearchGoodreads(ctx context.Context, request SearchGoodreadsRequestObject) (SearchGoodreadsResponseObject, error)
	// Search for books using kindle
	// (GET /kindle/{region}/search)
	SearchKindle(ctx context.Context, request SearchKindleRequestObject) (SearchKindleResponseObject, error)
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

// SearchGoodreads operation middleware
func (sh *strictHandler) SearchGoodreads(w http.ResponseWriter, r *http.Request, params SearchGoodreadsParams) {
	var request SearchGoodreadsRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.SearchGoodreads(ctx, request.(SearchGoodreadsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SearchGoodreads")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(SearchGoodreadsResponseObject); ok {
		if err := validResponse.VisitSearchGoodreadsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// SearchKindle operation middleware
func (sh *strictHandler) SearchKindle(w http.ResponseWriter, r *http.Request, region SearchKindleParamsRegion, params SearchKindleParams) {
	var request SearchKindleRequestObject

	request.Region = region
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.SearchKindle(ctx, request.(SearchKindleRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "SearchKindle")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(SearchKindleResponseObject); ok {
		if err := validResponse.VisitSearchKindleResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xWTXPjNgz9Kxq0R43ldLcX3XbbTuvJbtNJsoc24+nQEiIzlkgGBDPjZvTfO6CkyB+y",
	"07SHvfRiS8QD+AjiAXqGwjbOGjTsIX8Gp0g1yEjxTQVeW5InbSCHx4C0hRSMahDywZqCL9bYKIHx1onF",
	"M2lTQdumvc+JCMMr4WPQhCXkTAHPBWwF7J01HiPD7+Zz+SusYTQcOTtX60KxtiZ78NbI2hjPkXVIrDvv",
	"RnGx7h41YxMfviW8hxy+ycbEZJ2/zz5au/mMrErFCtp0YKeI1Daetl+wqwcsuKNboi9IO+EDOVxdit/7",
	"/8QaibpLOc72q/t/VGVyjY8BPXdELr4SkS+mKx/9F5bC5PuvlpKFYSSj6uQG6Qkp+SnGOnYcyjJuuFcJ",
	"R3SU15HuvaVGsVSILKSH7NIdgR2ZCvuE0XKQt+tPCduE15hERKIbVeFU7D2/iQ3KQGow7u/xY29JtEk8",
	"FtaUHtLxMNrwuJ82jBWSBKzQ0IGYjjbdV0wK2q/2MxUXJk5TK1MFOelUVCPx+EQiXVjV2q+x/B3VecS0",
	"1SPpN/SImwg/3SVS8GHFmuvpo7Cq3pjCU7Ha3b5618OWR3WdwgHjo2r20i9Mga9k5/z2PW45qSssAmne",
	"3kgKewU5/ecGx8GxRlUijZPjw5fbX66uF398uF1c/TrWi3L6EredyrW5t+Jf6wKNj+x758+LW0ghUC2B",
	"mZ3Ps8w6NN4GKnBmqcp6J58Jdkwy/BA82yYZkpX8RvZJd8yekHynn/nsYjYXLwmqnIYc3s3mszmk4BSv",
	"4wGzytqSUJU+86ioWMtihXysxptoTu4tJStrNz4JXpsqefGHuE+n2EX54vDzjn13qt9Nl+8Iybqp3Kav",
	"Avvu1S6nh/KU9wsuE9A4Cs9j3w/Yi3+CvdiZJ+exAooFGJpGyWfK68kWeLbRpqwxeyastDXtv7zBLsqJ",
	"67scjAd3F9UgVTRqoWMh0KEYcxPqevczCk1oRIUqQAqFApkOkALGvi7FG8eTlr7+4EQbG/nZleso6//r",
	"5+310191u9vt4nW+9Lm7paTBx0+Q7qb3a+iTLVQ92Xe6PlaLfW095+/mwmrZ/h0AAP//HMviW9sLAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
