package modbusquic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"

	quic "github.com/lucas-clemente/quic-go"
)

var _server func(req []byte) (res []byte)
var _fault func(detail string)

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

type Handler interface {
	Server(req []byte) (res []byte)
	Fault(detail string)
}

type MbQUIC struct {
	Addr byte
	Code byte
	Data []byte
}

func (m MbQUIC) generate() []byte {
	head := make([]byte, 8, 8)
	l := byte(len(m.Data) + 2)
	head[0] = 0x00
	head[1] = 0x00
	head[2] = 0x00
	head[3] = 0x00
	head[4] = 0x00
	head[5] = l
	head[6] = m.Addr
	head[7] = m.Code
	body := make([]byte, 260)
	body = append(body, head...)
	body = append(body, m.Data...)
	return body
}

//Send the data to the server
func (m *MbQUIC) Send(addr string) ([]byte, error) {
	req := m.generate()
	return send(addr, req)
}

func send(addr string, d []byte) ([]byte, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-modbus-echo-example"},
	}
	session, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Printf("Client: Sending '%s'\n", d)
	_, err = stream.Write(d)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, len(d))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		return buf, err
	}
	fmt.Printf("Client: Got '%s'\n", buf)

	return buf, nil
}

func SetHandler(h Handler) {
	_server = h.Server
	_fault = h.Fault
}

//ServerCreate creates the QUIC server
func ServerCreate(port string) error {
	listener, err := quic.ListenAddr("127.0.0.1:"+port, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}
		go handle(sess)
	}
}

//Currently handle just echos the connection but it could be repurposed
func handle(sess quic.Session) error {

	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}
	// Echo through the loggingWriter
	_, err = io.Copy(loggingWriter{stream}, stream)
	return err
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-modbus-echo-example"},
	}
}
