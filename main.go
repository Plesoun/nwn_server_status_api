package main

import (
	"fmt"
	"nwn_server_info/udp"
)

func main() {
	ipaddr := "3.64.204.102"
	port := "5121"

	nwnOnline, serverInfo := udp.CheckNWNServer(ipaddr, port)

	if nwnOnline {
		fmt.Println("Server online!")
		fmt.Println("Server name:", serverInfo.ServerName)
		fmt.Println("Module name:", serverInfo.ModuleName)
		fmt.Println("Players:", serverInfo.PlayersOnline, "/", serverInfo.PlayersMax)
		fmt.Println("Description:", serverInfo.Description)
		fmt.Println("Game Type:", serverInfo.GameType)
		fmt.Println("PvP:", serverInfo.PvP)
		fmt.Println("Version:", serverInfo.Version)
	} else {
		fmt.Println("Server offline or not responding.")
	}
}
