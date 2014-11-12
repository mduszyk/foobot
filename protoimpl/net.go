package protoimpl

import(
	"io"
    "bufio"
	"net"
    "net/textproto"
    "container/list"
	"github.com/mduszyk/foobot/log"
	"github.com/mduszyk/foobot/conf"
	"github.com/mduszyk/foobot/proto"
)

type NetServerProto struct {
    handlers *list.List
    conns map[string]net.Conn
}

func NewNetServerProto() *NetServerProto {
	proto := &NetServerProto{
        handlers: list.New(),
        conns: make(map[string]net.Conn),
    }
    return proto
}

func (p *NetServerProto) Run() {
    log.INFO.Printf("Starting net server proto")

    typ := conf.Get("net.server.type")
    socket := conf.Get("net.server.socket")

	l, err := net.Listen(typ, socket)
	if err != nil {
        log.ERROR.Printf("Failed binding proto, type: %s, socket: %s, err: %s",
            typ, socket, err)
	}
	defer l.Close()

    log.INFO.Printf("Binded net server, type: %s, socket: %s", typ, socket)

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
            log.ERROR.Printf("Failed accepting, err: %s", err)
		} else {
            go p.handleConn(conn)
        }
	}
}

func (p *NetServerProto) handleConn(conn net.Conn) {
    p.conns[conn.RemoteAddr().String()] = conn
    reader := textproto.NewReader(bufio.NewReader(conn))
    for {
        line, err := reader.ReadLine()

        if err == io.EOF {
            log.INFO.Printf("Closing connection, conn: %s", conn)
            conn.Close()
            break
        }

        for e := p.handlers.Front(); e != nil; e = e.Next() {
            handler := e.Value.(func(net.Conn, string))
            handler(conn, line)
        }
    }
}

func (p *NetServerProto) Send(addr string, text string) {
    conn, ok := p.conns[addr]
    if !ok {
        log.ERROR.Printf("Connection not found, addr: %s", addr)
        return
    }
    conn.Write([]byte(text))
}

func (p *NetServerProto) Register(i proto.Interpreter) {
    handler := func(conn net.Conn, line string) {
        addr := conn.RemoteAddr().String()
        /* log.TRACE.Printf("Got message, addr: %s, line: %s", addr, line) */
        msg := proto.Parse(line)
        msg.Addr = addr
        msg.User = addr
        msg.Proto = p
        // pass message to agent
        rsp := i.Handle(msg)
        if rsp != "" {
            conn.Write([]byte(rsp))
        }
    }
    p.handlers.PushBack(handler)
}
