//package app
//author: Lubia Yang
//create: 2013-10-21
//about: www.lubia.me

package app

import (
	"fmt"
	"log"
	"os"

	"github.com/CharlesDardaman/quic_modbus/src/modbusquic"
	"github.com/CharlesDardaman/quic_modbus/src/modbusrtu"
	"github.com/CharlesDardaman/quic_modbus/src/modbustcp"
)

func RtuClient() {
	fd, err := os.Open("/dev/ttyAM0")
	if err != nil {
		log.Println("unable to open rs485")
		return
	}
	b, err := modbusrtu.Read(fd, 0x03, 1, 3, 1)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(b)
	err = modbusrtu.Write(fd, 0x03, 1, 3, 1, []byte{0, 1})
	if err != nil {
		log.Println(err.Error())
	}
}

//TCPClient starts the TCP Client
func TCPClient() {
	mt := new(modbustcp.MbTcp)
	mt.Addr = 1
	mt.Code = 0x03
	mt.Data = []byte{0, 1}
	res, err := mt.Send("127.0.0.1:80")
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}

//QUICClient starts up the QUIC Client
func QUICClient() {
	log.Println("Starting client")
	fmt.Println("Starting client")
	mq := new(modbusquic.MbQUIC)
	mq.Addr = 1
	mq.Code = 0x03
	mq.Data = []byte{0, 1}
	res, err := mq.Send("127.0.0.1:4443")
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
