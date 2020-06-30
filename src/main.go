//package main
//author: Lubia Yang
//create: 2013-10-21
//about: www.lubia.me

package main

import (
	"github.com/CharlesDardaman/quic_modbus/src/app"
)

func main() {
	go app.QUICServer()
	app.QUICClient()

	//If you want to run the other examples uncomment the below code
	// app.RtuClient()
	// app.TcpClient()
	// app.TCPServer()
}
