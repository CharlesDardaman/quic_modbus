package app

import (
	"log"

	"github.com/CharlesDardaman/quic_modbus/src/modbusquic"
	"github.com/CharlesDardaman/quic_modbus/src/modbustcp"
)

var h *handler

type handler struct {
}

func (h *handler) Server(req []byte) []byte {
	return []byte{}
}

func (h *handler) Fault(detail string) {

}

//TCPServer starts the TCP server
func TCPServer() {
	modbustcp.SetHandler(h)
	err := modbustcp.ServerCreate(80)
	if err != nil {
		log.Println(err.Error())
	}
}

//QUICServer starts the QUIC Server
func QUICServer() {
	modbustcp.SetHandler(h)
	err := modbusquic.ServerCreate("4443")
	if err != nil {
		log.Println(err.Error())
	}
}
