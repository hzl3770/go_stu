package custom_rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_stu/custom_rpc/codec"
	"go_stu/custom_rpc/log"
	"io"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int        // MagicNumber marks this's a geerpc request
	CodecType   codec.Type // client may choose different Codec to encode body
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (s *Server) Accept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("rpc server: accept error:", err)
			return
		}
		go s.serveConn(conn)
	}
}

func (s *Server) serveConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()

	var option Option
	if err := json.NewDecoder(conn).Decode(&option); err != nil {
		log.Error("rpc server: options error: ", err)
		return
	}

	if option.MagicNumber != MagicNumber {
		log.Error("rpc server: invalid magic number %x", option.MagicNumber)
		return
	}

	codecFunc := codec.NewCodecFuncMap[option.CodecType]
	if codecFunc == nil {
		log.Error("rpc server: invalid codec type", option.CodecType)
		return
	}

	s.serveCodec(codecFunc(conn))
}

var invalidRequest = struct{}{}

func (s *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := s.readRequest(cc)
		if err != nil {
			if req == nil {
				break // it's not possible to recover, so close the connection
			}

			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}

		wg.Add(1)
		go s.handleRequest(cc, req, sending, wg)
	}
}

type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
}

func (s *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := s.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
	req.argv = reflect.New(reflect.TypeOf(""))
	if err := cc.ReadBody(req.argv.Interface()); err != nil {
		log.Error("rpc server: read body error: ", err)
		return nil, err
	}

	return req, nil
}

func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			log.Error("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (s *Server) sendResponse(cc codec.Codec, header *codec.Header, reply any, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()

	if err := cc.Write(header, reply); err != nil {
		log.Error("rpc server: write response error:", err)
	}
}

func (s *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Infof("rpc server: request: %v\n", req.h)
	log.Infof("rpc server: request argv: %v\n", req.argv.Elem())

	req.replyv = reflect.ValueOf(fmt.Sprintf("sbydx resp %d", req.h.Seq))
	s.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}
