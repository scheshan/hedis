package server

import (
	"bufio"
	"errors"
	"hedis/codec"
	"hedis/core"
	"io"
	"log"
	"net"
	"sync"
)

type ClientMessage struct {
	session *Session
	message codec.Message
}

type Server interface {
	Start() error
	Stop() error
}

type baseServer struct {
	commands      *core.Hash
	config        *ServerConfig
	listener      *net.TCPListener
	session       *Session
	running       bool
	clientId      int
	requests      chan *ClientMessage
	responses     chan *ClientMessage
	subscription  *core.Hash
	pSubscription *core.List
	decoder       *codec.Decoder
	encoder       *codec.Encoder
	mutex         *sync.Mutex
}

func newBaseServer(config *ServerConfig) *baseServer {
	s := &baseServer{}

	s.config = config
	s.commands = core.NewHash()
	s.requests = make(chan *ClientMessage, 10240)
	s.responses = make(chan *ClientMessage, 10240)
	s.subscription = core.NewHash()
	s.pSubscription = core.NewList()
	s.decoder = &codec.Decoder{}
	s.encoder = &codec.Encoder{}
	s.mutex = &sync.Mutex{}

	s.initCommands()

	return s
}

func (t *baseServer) bindAndListen() error {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:6379")
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	t.listener = listener
	return nil
}

func (t *baseServer) accept() {
	for t.running {
		conn, err := t.listener.Accept()
		if err != nil {
			log.Fatalf("连接创建出现错误:%v", err)
		}

		tcpConn := conn.(*net.TCPConn)
		if err = tcpConn.SetNoDelay(true); err != nil {
			log.Printf("配置tcp连接出错: %v", err)
		}

		session := t.initSession(tcpConn)
		t.addSession(session)

		go t.readLoop(session)
		go t.writeLoop(session)

		log.Printf("客户端 %s 连接成功", session)
	}
}

func (t *baseServer) initSession(conn *net.TCPConn) *Session {
	t.clientId++
	sessionId := t.clientId

	session := &Session{}
	session.id = sessionId
	session.conn = conn
	session.reader = bufio.NewReader(conn)
	session.writer = bufio.NewWriter(conn)
	session.messages = make(chan codec.Message, 1024)

	return session
}

func (t *baseServer) addSession(session *Session) {
	if t.session != nil {
		t.session.pre = session
	}
	session.next = t.session
	t.session = session
}

func (t *baseServer) removeSession(session *Session) {
	if session.pre != nil {
		session.pre.next = session.next
	}
	if session.next != nil {
		session.next.pre = session.pre
	}

	if session == t.session {
		t.session = session.next
	}

	session.pre = nil
	session.next = nil
}

func (t *baseServer) closeSession(session *Session) {
	session.state = SessionStateClosed
	close(session.messages)

	t.removeSession(session)
	session.conn.Close()
}

func (t *baseServer) onSessionError(session *Session, err error) {
	t.closeSession(session)

	if err == io.EOF {
		log.Printf("客户端 %s 主动关闭连接", session)
		return
	}
	log.Printf("客户端 %s 连接出错，关闭连接: %v", session, err)
}

func (t *baseServer) writeToSession(session *Session, msg codec.Message) error {
	err := t.encoder.Encode(session.writer, msg)
	if err != nil {
		return err
	}

	return session.writer.Flush()
}

func (t *baseServer) readLoop(session *Session) {
	for t.running && session.state == SessionStateConnected {
		msg, err := t.decoder.Decode(session.reader)
		if err != nil {
			t.onSessionError(session, err)
			return
		}

		cm := &ClientMessage{}
		cm.session = session
		cm.message = msg

		t.requests <- cm
	}
}

func (t *baseServer) writeLoop(session *Session) {
	for t.running && session.state == SessionStateConnected {
		msg, ok := <-session.messages
		if !ok {
			break
		}
		if msg == nil {
			continue
		}

		err := t.writeToSession(session, msg)
		if err != nil {
			t.onSessionError(session, err)
			return
		}
	}
}

func (t *baseServer) getCommand(key *core.String) (Command, bool) {
	v, find := t.commands.Get(key)
	if !find {
		return nil, find
	}

	return v.(Command), find
}

func (t *baseServer) addCommand(key string, cmd Command) {
	t.commands.Put(core.NewStringStr(key), cmd)
}

func (t *baseServer) processCommand() {
	for t.running {
		ctx, ok := <-t.requests
		if !ok {
			break
		}

		cmd, find := t.getCommand(ctx.message.Command())
		if !find {
			cmd = t.CommandNotFound
		}
		message := cmd(ctx.session, ctx.message.Args())

		ctx.session.messages <- message
	}
}

func (t *baseServer) Start() error {
	t.mutex.Lock()
	defer func() {
		t.mutex.Unlock()
	}()

	if t.running {
		log.Println("Server is already running")
		return errors.New("server is already running")
	}

	t.running = true
	err := t.bindAndListen()
	if err != nil {
		return err
	}

	go t.processCommand()
	go t.accept()

	return nil
}

func (t *baseServer) Stop() {

}

//#region commands

func (t *baseServer) initCommands() {
	t.addCommand("ping", t.CommandPing)
	t.addCommand("quit", t.CommandQuit)
	t.addCommand("echo", t.CommandEcho)
}

func (t *baseServer) CommandNotFound(s *Session, args []*core.String) codec.Message {
	return codec.ErrorCommandNotFound
}

func (t *baseServer) CommandPing(s *Session, args []*core.String) codec.Message {
	var msg codec.Message
	if len(args) == 1 {
		msg = codec.NewSimpleStr(args[0])
	} else {
		msg = codec.NewSimpleString("pong")
	}

	return msg
}

func (t *baseServer) CommandQuit(s *Session, args []*core.String) codec.Message {
	msg := codec.SimpleOK
	t.writeToSession(s, msg)

	t.closeSession(s)

	return nil
}

func (t *baseServer) CommandEcho(s *Session, args []*core.String) codec.Message {
	if len(args) != 1 {
		return MessageErrorInvalidArgNum
	}

	msg := codec.NewBulkStr(args[0])

	return msg
}

//endregion
