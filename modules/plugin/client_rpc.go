package plugin

//go:generate go run modules/plugin/interface_generator/main.go

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"reflect"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-plugin"
	"github.com/lib/pq"
	"github.com/sitename/sitename/model"
	"github.com/sitename/sitename/model/file"
	"github.com/sitename/sitename/model/plugins"
	"github.com/sitename/sitename/modules/json"
	"github.com/sitename/sitename/modules/slog"
)

var hookNameToId map[string]int = make(map[string]int)

type hooksRPCClient struct {
	client      *rpc.Client
	log         *slog.Logger
	muxBroker   *plugin.MuxBroker
	apiImpl     API
	driver      Driver
	implemented [TotalHooksID]bool
}

type hooksRPCServer struct {
	impl         interface{}
	muxBroker    *plugin.MuxBroker
	apiRPCClient *apiRPCClient
}

// Implements hashicorp/go-plugin/plugin.Plugin interface to connect the hooks of a plugin
type hooksPlugin struct {
	hooks      interface{}
	apiImpl    API
	driverImpl Driver
	log        *slog.Logger
}

func (p *hooksPlugin) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &hooksRPCServer{
		impl:      p.hooks,
		muxBroker: b,
	}, nil
}

func (p *hooksPlugin) Client(b *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &hooksRPCClient{
		client:    client,
		log:       p.log,
		muxBroker: b,
		apiImpl:   p.apiImpl,
		driver:    p.driverImpl,
	}, nil
}

type apiRPCClient struct {
	client    *rpc.Client
	muxBroker *plugin.MuxBroker
}

type apiRPCServer struct {
	impl      API
	muxBroker *plugin.MuxBroker
}

// ErrorString is a fallback for sending unregistered implementations of the error interface across
// rpc. For example, the errorString type from the github.com/pkg/errors package cannot be
// registered since it is not exported, but this precludes common error handling paradigms.
// ErrorString merely preserves the string description of the error, while satisfying the error
// interface itself to allow other registered types (such as model.AppError) to be sent unmodified.
type ErrorString struct {
	Code int // Code to map to various error variables
	Err  string
}

func (e ErrorString) Error() string {
	return e.Err
}

func encodableError(err error) error {
	if err == nil {
		return nil
	}

	switch err.(type) {
	case *model.AppError, *pq.Error, *mysql.MySQLError:
		return err
	}

	ret := &ErrorString{
		Err: err.Error(),
	}

	switch err {
	case io.EOF:
		ret.Code = 1
	case sql.ErrNoRows:
		ret.Code = 2
	case sql.ErrConnDone:
		ret.Code = 3
	case sql.ErrTxDone:
		ret.Code = 4
	case driver.ErrSkip:
		ret.Code = 5
	case driver.ErrBadConn:
		ret.Code = 6
	case driver.ErrRemoveArgument:
		ret.Code = 7
	}

	return ret
}

func decodableError(err error) error {
	if encErr, ok := err.(*ErrorString); ok {
		switch encErr.Code {
		case 1:
			return io.EOF
		case 2:
			return sql.ErrNoRows
		case 3:
			return sql.ErrConnDone
		case 4:
			return sql.ErrTxDone
		case 5:
			return driver.ErrSkip
		case 6:
			return driver.ErrBadConn
		case 7:
			return driver.ErrRemoveArgument
		}
	}
	return err
}

func init() {
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	gob.Register(&model.AppError{})
	gob.Register(&pq.Error{})
	gob.Register(&mysql.MySQLError{})
	gob.Register(&ErrorString{})
	gob.Register(&opengraph.OpenGraph{})
	// gob.Register(&model.AutocompleteDynamicListArg{})
	// gob.Register(&model.AutocompleteStaticListArg{})
	// gob.Register(&model.AutocompleteTextArg{})
}

// These enforce compile time checks to make sure types implement the interface
// If you are getting an error here, you probably need to run `make pluginapi` to
// autogenerate RPC glue code
var _ plugin.Plugin = &hooksPlugin{}
var _ Hooks = &hooksRPCClient{}

func (g *hooksRPCClient) Implemented() (impl []string, err error) {
	err = g.client.Call("Plugin.Implemented", struct{}{}, &impl)
	for _, hookName := range impl {
		if hookId, ok := hookNameToId[hookName]; ok {
			g.implemented[hookId] = true
		}
	}

	return
}

// Implemented replies with the names of the hooks that are implemented.
func (s *hooksRPCServer) Implemented(args struct{}, reply *[]string) error {
	ifaceType := reflect.TypeOf((*Hooks)(nil)).Elem()
	implType := reflect.TypeOf(s.impl)
	selfType := reflect.TypeOf(s)
	var methods []string

	for i := 0; i < ifaceType.NumMethod(); i++ {
		method := ifaceType.Method(i)
		if m, ok := implType.MethodByName(method.Name); !ok {
			// implType HAS NOT implemented the method with certain name
			continue
		} else if m.Type.NumIn() != method.Type.NumIn()+1 {
			continue
		} else if m.Type.NumOut() != method.Type.NumOut() {
			continue
		} else {
			match := true
			for j := 0; j < method.Type.NumIn(); j++ {
				if m.Type.In(j+1) != method.Type.In(j) {
					match = false
					break
				}
			}
			for j := 0; j < method.Type.NumOut(); j++ {
				if m.Type.Out(j) != method.Type.Out(j) {
					match = false
					break
				}
			}

			if !match {
				continue
			}
		}
		if _, ok := selfType.MethodByName(method.Name); !ok {
			continue
		}
		methods = append(methods, method.Name)
	}
	*reply = methods
	return encodableError(nil)
}

type Z_OnActivateArgs struct {
	APIMuxId    uint32
	DriverMuxId uint32
}

type Z_OnActivateReturns struct {
	A error
}

func (g *hooksRPCClient) OnActivate() error {
	muxId := g.muxBroker.NextId()
	go g.muxBroker.AcceptAndServe(muxId, &apiRPCServer{
		impl:      g.apiImpl,
		muxBroker: g.muxBroker,
	})

	nextID := g.muxBroker.NextId()
	go g.muxBroker.AcceptAndServe(nextID, &dbRPCServer{
		dbImpl: g.driver,
	})

	_args := &Z_OnActivateArgs{
		APIMuxId:    muxId,
		DriverMuxId: nextID,
	}
	_returns := &Z_OnActivateReturns{}

	if err := g.client.Call("Plugin.OnActivate", _args, _returns); err != nil {
		g.log.Error("RPC call to OnActivate plugin failed.", slog.Err(err))
	}
	return _returns.A
}

func (s *hooksRPCServer) OnActivate(args *Z_OnActivateArgs, returns *Z_OnActivateReturns) error {
	connection, err := s.muxBroker.Dial(args.APIMuxId)
	if err != nil {
		return err
	}

	conn2, err := s.muxBroker.Dial(args.DriverMuxId)
	if err != nil {
		return err
	}

	s.apiRPCClient = &apiRPCClient{
		client:    rpc.NewClient(connection),
		muxBroker: s.muxBroker,
	}

	dbClient := &dbRPCClient{
		client: rpc.NewClient(conn2),
	}

	if snplugin, ok := s.impl.(interface {
		SetAPI(api API)
		SetHelpers(helpers Helpers)
		SetDriver(driver Driver)
	}); ok {
		snplugin.SetAPI(s.apiRPCClient)
		snplugin.SetHelpers(&HelpersImpl{
			API: s.apiRPCClient,
		})
		snplugin.SetDriver(dbClient)
	}

	if snplugin, ok := s.impl.(interface {
		OnConfigurationChange() error
	}); ok {
		if err := snplugin.OnConfigurationChange(); err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] call to OnConfigurationChange failed, error: %v", err)
		}
	}

	// Capture output of standard logger because go-plugin
	// redirects it.
	log.SetOutput(os.Stderr)

	if hook, ok := s.impl.(interface {
		OnActivate() error
	}); ok {
		returns.A = encodableError(hook.OnActivate())
	}

	return nil
}

type Z_LoadPluginConfigurationArgsArgs struct {
}

type Z_LoadPluginConfigurationArgsReturns struct {
	A []byte
}

func (g *apiRPCClient) LoadPluginConfiguration(dest interface{}) error {
	_args := &Z_LoadPluginConfigurationArgsArgs{}
	_returns := &Z_LoadPluginConfigurationArgsReturns{}
	if err := g.client.Call("Plugin.LoadPluginConfiguration", _args, _returns); err != nil {
		log.Printf("RPC call to LoadPluginConfiguration API failed: %s", err.Error())
	}
	if err := model.ModelFromJson(&dest, bytes.NewReader(_returns.A)); err != nil {
		log.Printf("LoadPluginConfiguration API failed to unmarshal: %s", err.Error())
	}
	return nil
}

func (s *apiRPCServer) LoadPluginConfiguration(args *Z_LoadPluginConfigurationArgsArgs, returns *Z_LoadPluginConfigurationArgsReturns) error {
	var config interface{}
	if hook, ok := s.impl.(interface {
		LoadPluginConfiguration(dest interface{}) error
	}); ok {
		if err := hook.LoadPluginConfiguration(&config); err != nil {
			return err
		}
	}
	b, err := json.JSON.Marshal(config)
	if err != nil {
		return err
	}
	returns.A = b
	return nil
}

func init() {
	hookNameToId["ServeHTTP"] = ServeHTTPID
}

type Z_ServeHTTPArgs struct {
	ResponseWriterStream uint32
	Request              *http.Request
	Context              *Context
	RequestBodyStream    uint32
}

func (g *hooksRPCClient) ServeHTTP(c *Context, w http.ResponseWriter, r *http.Request) {
	if !g.implemented[ServeHTTPID] {
		http.NotFound(w, r)
		return
	}

	serveHTTPStreamId := g.muxBroker.NextId()
	go func() {
		connection, err := g.muxBroker.Accept(serveHTTPStreamId)
		if err != nil {
			g.log.Error("Plugin failed to ServeHTTP, muxBroker couldn't accept connection", slog.Uint32("serve_http_stream_id", serveHTTPStreamId), slog.Err(err))
			return
		}
		defer connection.Close()

		rpcServer := rpc.NewServer()
		if err := rpcServer.RegisterName("Plugin", &httpResponseWriterRPCServer{w: w, log: g.log}); err != nil {
			g.log.Error("Plugin failed to ServeHTTP, couldn't register RPC name", slog.Err(err))
			return
		}
		rpcServer.ServeConn(connection)
	}()

	requestBodyStreamId := uint32(0)
	if r.Body != nil {
		requestBodyStreamId = g.muxBroker.NextId()
		go func() {
			bodyConnection, err := g.muxBroker.Accept(requestBodyStreamId)
			if err != nil {
				g.log.Error("Plugin failed to ServeHTTP, muxBroker couldn't Accept request body connection", slog.Err(err))
				return
			}
			defer bodyConnection.Close()
			serveIOReader(r.Body, bodyConnection)
		}()
	}

	forwardedRequest := &http.Request{
		Method:     r.Method,
		URL:        r.URL,
		Proto:      r.Proto,
		ProtoMajor: r.ProtoMajor,
		ProtoMinor: r.ProtoMinor,
		Header:     r.Header,
		Host:       r.Host,
		RemoteAddr: r.RemoteAddr,
		RequestURI: r.RequestURI,
	}

	if err := g.client.Call("Plugin.ServeHTTP", Z_ServeHTTPArgs{
		Context:              c,
		ResponseWriterStream: serveHTTPStreamId,
		Request:              forwardedRequest,
		RequestBodyStream:    requestBodyStreamId,
	}, nil); err != nil {
		g.log.Error("Plugin failed to ServeHTTP, RPC call failed", slog.Err(err))
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
	}
}

func (s *hooksRPCServer) ServeHTTP(args *Z_ServeHTTPArgs, returns *struct{}) error {
	connection, err := s.muxBroker.Dial(args.ResponseWriterStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Can't connecti to remote response writer stream, error: %v", err.Error())
		return err
	}
	w := connectHTTPResponseWriter(connection)
	defer w.Close()

	r := args.Request
	if args.RequestBodyStream != 0 {
		connection, err := s.muxBroker.Dial(args.RequestBodyStream)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Can't connect to remote request body stream, error: %v", err.Error())
			return err
		}
		r.Body = connectIOReader(connection)
	} else {
		r.Body = ioutil.NopCloser(&bytes.Buffer{})
	}
	defer r.Body.Close()

	if hook, ok := s.impl.(interface {
		ServeHTTP(c *Context, w http.ResponseWriter, r *http.Request)
	}); ok {
		hook.ServeHTTP(args.Context, w, r)
	} else {
		http.NotFound(w, r)
	}

	return nil
}

type Z_PluginHTTPArgs struct {
	Request     *http.Request
	RequestBody []byte
}

type Z_PluginHTTPReturns struct {
	Response     *http.Response
	ResponseBody []byte
}

func (g *apiRPCClient) PluginHTTP(request *http.Request) *http.Response {
	forwardedRequest := &http.Request{
		Method:     request.Method,
		URL:        request.URL,
		Proto:      request.Proto,
		ProtoMajor: request.ProtoMajor,
		ProtoMinor: request.ProtoMinor,
		Header:     request.Header,
		Host:       request.Host,
		RemoteAddr: request.RemoteAddr,
		RequestURI: request.RequestURI,
	}

	_args := &Z_PluginHTTPArgs{
		Request: forwardedRequest,
	}

	if request.Body != nil {
		requestBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Printf("RPC call to PluginHTTP API failed: %s", err.Error())
			return nil
		}
		request.Body.Close()
		request.Body = nil

		_args.RequestBody = requestBody
	}

	_returns := &Z_PluginHTTPReturns{}
	if err := g.client.Call("Plugin.PluginHTTP", _args, _returns); err != nil {
		log.Printf("RPC call to PluginHTTP API failed: %s", err.Error())
		return nil
	}

	_returns.Response.Body = ioutil.NopCloser(bytes.NewBuffer(_returns.ResponseBody))

	return _returns.Response
}

func (s *apiRPCServer) PluginHTTP(args *Z_PluginHTTPArgs, returns *Z_PluginHTTPReturns) error {
	args.Request.Body = ioutil.NopCloser(bytes.NewBuffer(args.RequestBody))

	if hook, ok := s.impl.(interface {
		PluginHTTP(request *http.Request) *http.Response
	}); ok {
		response := hook.PluginHTTP(args.Request)

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return encodableError(fmt.Errorf("RPC call to PluginHTTP API failed: %s", err.Error()))
		}
		response.Body.Close()
		response.Body = nil

		returns.Response = response
		returns.ResponseBody = responseBody
	} else {
		return encodableError(fmt.Errorf("API PluginHTTP called but not implemented"))
	}
	return nil
}

func init() {
	hookNameToId["FileWillBeUploaded"] = FileWillBeUploadedID
}

type Z_FileWillBeUploadedArgs struct {
	A                     *Context
	B                     *file.FileInfo
	UploadedFileStream    uint32
	ReplacementFileStream uint32
}

type Z_FileWillBeUploadedReturns struct {
	A *file.FileInfo
	B string
}

func (g *hooksRPCClient) FileWillBeUploaded(c *Context, info *file.FileInfo, file io.Reader, output io.Writer) (*file.FileInfo, string) {
	if !g.implemented[FileWillBeUploadedID] {
		return info, ""
	}

	uploadedFileStreamId := g.muxBroker.NextId()
	go func() {
		uploadedFileConnection, err := g.muxBroker.Accept(uploadedFileStreamId)
		if err != nil {
			g.log.Error("Plugin failed to serve upload file stream. MuxBroker could not Accept connection", slog.Err(err))
			return
		}
		defer uploadedFileConnection.Close()
		serveIOReader(file, uploadedFileConnection)
	}()

	replacementDone := make(chan bool)
	replacementFileStreamID := g.muxBroker.NextId()
	go func() {
		defer close(replacementDone)

		replacementFileConnection, err := g.muxBroker.Accept(replacementFileStreamID)
		if err != nil {
			g.log.Error("Plugin failed to server replacement file stream. MuxBroker could not Accept connection", slog.Err((err)))
			return
		}
		defer replacementFileConnection.Close()
		if _, err := io.Copy(output, replacementFileConnection); err != nil {
			g.log.Error("Error reading replacement file.", slog.Err(err))
		}
	}()

	_args := &Z_FileWillBeUploadedArgs{c, info, uploadedFileStreamId, replacementFileStreamID}
	_returns := &Z_FileWillBeUploadedReturns{A: _args.B}
	if err := g.client.Call("Plugin.FileWillBeUploaded", _args, _returns); err != nil {
		g.log.Error("RPC call FileWillBeUploaded to plugin failed.", slog.Err(err))
	}

	// Ensure the io.Copy from the replacementFileConnection above completes.
	<-replacementDone

	return _returns.A, _returns.B
}

func (s *hooksRPCServer) FileWillBeUploaded(args *Z_FileWillBeUploadedArgs, returns *Z_FileWillBeUploadedReturns) error {
	uploadFileConnection, err := s.muxBroker.Dial(args.UploadedFileStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Can't connect to remote upload file stream, error: %v", err.Error())
		return err
	}
	defer uploadFileConnection.Close()
	fileReader := connectIOReader(uploadFileConnection)
	defer fileReader.Close()

	replacementFileConnection, err := s.muxBroker.Dial(args.ReplacementFileStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Can't connect to remote replacement file stream, error: %v", err.Error())
		return err
	}
	defer replacementFileConnection.Close()
	returnFileWriter := replacementFileConnection

	if hook, ok := s.impl.(interface {
		FileWillBeUploaded(c *Context, info *file.FileInfo, file io.Reader, output io.Writer) (*file.FileInfo, string)
	}); ok {
		returns.A, returns.B = hook.FileWillBeUploaded(args.A, args.B, fileReader, returnFileWriter)
	} else {
		return fmt.Errorf("hook FileWillBeUploaded called but not implemented")
	}
	return nil
}

// MessageWillBePosted is in this file because of the difficulty of identifying which fields need special behaviour.
// The special behaviour needed is decoding the returned post into the original one to avoid the unintentional removal
// of fields by older plugins.
func init() {
	hookNameToId["MessageWillBePosted"] = MessageWillBePostedID
}

// type Z_MessageWillBePostedArgs struct {
// 	A *Context
// 	B *model.Post
// }

// type Z_MessageWillBePostedReturns struct {
// 	A *model.Post
// 	B string
// }

// func (g *hooksRPCClient) MessageWillBeUpdated(c *Context, newPost, oldPost *model.Post) (*model.Post, string) {
// 	_args := &Z_MessageWillBeUpdatedArgs{c, newPost, oldPost}
// 	_returns := &Z_MessageWillBeUpdatedReturns{A: _args.B}
// 	if g.implemented[MessageWillBeUpdatedID] {
// 		if err := g.client.Call("Plugin.MessageWillBeUpdated", _args, _returns); err != nil {
// 			g.log.Error("RPC call MessageWillBeUpdated to plugin failed.", mlog.Err(err))
// 		}
// 	}
// 	return _returns.A, _returns.B
// }

// func (s *hooksRPCServer) MessageWillBeUpdated(args *Z_MessageWillBeUpdatedArgs, returns *Z_MessageWillBeUpdatedReturns) error {
// 	if hook, ok := s.impl.(interface {
// 		MessageWillBeUpdated(c *Context, newPost, oldPost *model.Post) (*model.Post, string)
// 	}); ok {
// 		returns.A, returns.B = hook.MessageWillBeUpdated(args.A, args.B, args.C)

// 	} else {
// 		return encodableError(fmt.Errorf("hook MessageWillBeUpdated called but not implemented"))
// 	}
// 	return nil
// }

type Z_LogDebugArgs struct {
	A string
	B []interface{}
}

type Z_LogDebugReturns struct {
}

func (g *apiRPCClient) LogDebug(msg string, keyValuePairs ...interface{}) {
	stringifiedPairs := stringifyToObjects(keyValuePairs)
	_args := &Z_LogDebugArgs{msg, stringifiedPairs}
	_returns := &Z_LogDebugReturns{}
	if err := g.client.Call("Plugin.LogDebug", _args, _returns); err != nil {
		log.Printf("RPC call to LogDebug API failed: %s", err.Error())
	}

}

func (s *apiRPCServer) LogDebug(args *Z_LogDebugArgs, returns *Z_LogDebugReturns) error {
	if hook, ok := s.impl.(interface {
		LogDebug(msg string, keyValuePairs ...interface{})
	}); ok {
		hook.LogDebug(args.A, args.B...)
	} else {
		return encodableError(fmt.Errorf("API LogDebug called but not implemented"))
	}
	return nil
}

type Z_LogInfoArgs struct {
	A string
	B []interface{}
}

type Z_LogInfoReturns struct {
}

func (g *apiRPCClient) LogInfo(msg string, keyValuePairs ...interface{}) {
	stringifiedPairs := stringifyToObjects(keyValuePairs)
	_args := &Z_LogInfoArgs{msg, stringifiedPairs}
	_returns := &Z_LogInfoReturns{}
	if err := g.client.Call("Plugin.LogInfo", _args, _returns); err != nil {
		log.Printf("RPC call to LogInfo API failed: %s", err.Error())
	}
}

func (s *apiRPCServer) LogInfo(args *Z_LogInfoArgs, returns *Z_LogInfoReturns) error {
	if hook, ok := s.impl.(interface {
		LogInfo(msg string, keyValuePairs ...interface{})
	}); ok {
		hook.LogInfo(args.A, args.B...)
	} else {
		return encodableError(fmt.Errorf("API LogInfo called but not implemented"))
	}
	return nil
}

type Z_LogWarnArgs struct {
	A string
	B []interface{}
}

type Z_LogWarnReturns struct {
}

func (g *apiRPCClient) LogWarn(msg string, keyValuePairs ...interface{}) {
	stringifiedPairs := stringifyToObjects(keyValuePairs)
	_args := &Z_LogWarnArgs{msg, stringifiedPairs}
	_returns := &Z_LogWarnReturns{}
	if err := g.client.Call("Plugin.LogWarn", _args, _returns); err != nil {
		log.Printf("RPC call to LogWarn API failed: %s", err.Error())
	}

}

func (s *apiRPCServer) LogWarn(args *Z_LogWarnArgs, returns *Z_LogWarnReturns) error {
	if hook, ok := s.impl.(interface {
		LogWarn(msg string, keyValuePairs ...interface{})
	}); ok {
		hook.LogWarn(args.A, args.B...)
	} else {
		return encodableError(fmt.Errorf("API LogWarn called but not implemented"))
	}
	return nil
}

type Z_LogErrorArgs struct {
	A string
	B []interface{}
}

type Z_LogErrorReturns struct {
}

func (g *apiRPCClient) LogError(msg string, keyValuePairs ...interface{}) {
	stringifiedPairs := stringifyToObjects(keyValuePairs)
	_args := &Z_LogErrorArgs{msg, stringifiedPairs}
	_returns := &Z_LogErrorReturns{}
	if err := g.client.Call("Plugin.LogError", _args, _returns); err != nil {
		log.Printf("RPC call to LogError API failed: %s", err.Error())
	}
}

func (s *apiRPCServer) LogError(args *Z_LogErrorArgs, returns *Z_LogErrorReturns) error {
	if hook, ok := s.impl.(interface {
		LogError(msg string, keyValuePairs ...interface{})
	}); ok {
		hook.LogError(args.A, args.B...)
	} else {
		return encodableError(fmt.Errorf("API LogError called but not implemented"))
	}
	return nil
}

type Z_InstallPluginArgs struct {
	PluginStreamID uint32
	B              bool
}

type Z_InstallPluginReturns struct {
	A *plugins.Manifest
	B *model.AppError
}

func (g *apiRPCClient) InstallPlugin(file io.Reader, replace bool) (*plugins.Manifest, *model.AppError) {
	pluginStreamID := g.muxBroker.NextId()

	go func() {
		uploadPluginConnection, err := g.muxBroker.Accept(pluginStreamID)
		if err != nil {
			log.Print("Plugin failed to upload plugin. MuxBroker could not Accept connection", slog.Err(err))
			return
		}
		defer uploadPluginConnection.Close()
		serveIOReader(file, uploadPluginConnection)
	}()

	_args := &Z_InstallPluginArgs{pluginStreamID, replace}
	_returns := &Z_InstallPluginReturns{}
	if err := g.client.Call("Plugin.InstallPlugin", _args, _returns); err != nil {
		log.Print("RPC call InstallPlugin to plugin failed.", slog.Err(err))
	}

	return _returns.A, _returns.B
}

func (s *apiRPCServer) InstallPlugin(args *Z_InstallPluginArgs, returns *Z_InstallPluginReturns) error {
	hook, ok := s.impl.(interface {
		InstallPlugin(file io.Reader, replace bool) (*plugins.Manifest, *model.AppError)
	})
	if !ok {
		return encodableError(fmt.Errorf("API InstallPlugin called but not implemented"))
	}

	receivePluginConnection, err := s.muxBroker.Dial(args.PluginStreamID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Can't connect to remote plugin stream, error: %v", err.Error())
		return err
	}
	pluginReader := connectIOReader(receivePluginConnection)
	defer pluginReader.Close()

	returns.A, returns.B = hook.InstallPlugin(pluginReader, args.B)
	return nil
}