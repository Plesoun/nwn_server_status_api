package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	ipaddr := "3.64.204.102" // Replace with the server IP address
	port := "5121"           // Replace with the server port

	nwnOnline, serverInfo := checkNWNServer(ipaddr, port)

	if nwnOnline {
		fmt.Println("Server online!")
		fmt.Println("Server name:", serverInfo.ServerName)
		fmt.Println("Module name:", serverInfo.ModuleName)
		fmt.Println("Players:", serverInfo.PlayersOnline, "/", serverInfo.PlayersMax)
		fmt.Println("Description:", serverInfo.Description)
		fmt.Println("Game Type:", serverInfo.GameType)
		fmt.Println("PvP:", serverInfo.PvP)
		fmt.Println("Version:", serverInfo.Version)
		// ... print other fields as needed ...
	} else {
		fmt.Println("Server offline or not responding.")
	}
}

// ServerInfo stores the information gathered about the NWN server.
type ServerInfo struct {
	ServerName    string
	ModuleName    string
	PlayersOnline int
	PlayersMax    int
	Description   string
	GameType      string
	PvP           string
	Version       string
	// ... other fields from the PHP code ...
}

// checkNWNServer checks the NWN server status using GameSpy and Eriniel methods.
func checkNWNServer(ipaddr, port string) (bool, *ServerInfo) {
	nwnOnline, serverInfo := checkNWNServerGameSpy(ipaddr, port)
	if nwnOnline {
		return true, serverInfo
	}

	// Fallback to Eriniel method if GameSpy fails
	return checkNWNServerEriniel(ipaddr, port)
}

// checkNWNServerGameSpy checks the NWN server status using the GameSpy protocol.
func checkNWNServerGameSpy(ipaddr, port string) (bool, *ServerInfo) {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("udp", net.JoinHostPort(ipaddr, port), timeout)
	if err != nil {
		fmt.Println("GameSpy: Error connecting:", err)
		return false, nil
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(timeout))

	// GameSpy query packet
	send := []byte{0xFE, 0xFD, 0x00, 0xE0, 0xEB, 0x2D, 0x0E, 0x14, 0x01, 0x0B, 0x01, 0x05, 0x08, 0x0A, 0x33, 0x34, 0x35, 0x13, 0x04, 0x36, 0x37, 0x38, 0x39, 0x14, 0x3A, 0x3B, 0x3C, 0x3D, 0x00, 0x00}

	_, err = conn.Write(send)
	if err != nil {
		fmt.Println("GameSpy: Error sending packet:", err)
		return false, nil
	}

	output := make([]byte, 5000)
	n, err := conn.Read(output)
	if err != nil {
		fmt.Println("GameSpy: Error reading response:", err)
		return false, nil
	}

	if n > 0 {
		// Server is online!
		lines := bytes.Split(output[:n], []byte{0x00})

		serverInfo := &ServerInfo{}
		if len(lines) > 4 {
			serverInfo.GameType = string(lines[2])
			serverInfo.ModuleName = strings.ReplaceAll(string(lines[4]), "_", " ")
			serverInfo.ServerName = string(lines[3])
			if serverInfo.ServerName == "" {
				serverInfo.ServerName = serverInfo.ModuleName
			}
		}

		if len(lines) > 6 {
			serverInfo.PlayersOnline = parseInt(string(lines[5]))
			serverInfo.PlayersMax = parseInt(string(lines[6]))
		}

		if len(lines) > 15 {
			serverInfo.Description = strings.ReplaceAll(string(lines[15]), "\n", "<br>\n")
		}

		if len(lines) > 9 {
			serverInfo.PvP = parsePVP(string(lines[9]))
		}

		if len(lines) > 14 {
			serverInfo.Version = parseVersion(string(lines[14]), string(lines[20]))
		}

		return true, serverInfo
	}

	return false, nil
}

// checkNWNServerEriniel checks the NWN server status using the Eriniel method.
func checkNWNServerEriniel(ipaddr, port string) (bool, *ServerInfo) {
	timeout := 5 * time.Second
	conn, err := net.DialTimeout("udp", net.JoinHostPort(ipaddr, port), timeout)
	if err != nil {
		fmt.Println("Eriniel: Error connecting:", err)
		return false, nil
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(timeout))

	// BNES packet
	send := []byte{0x42, 0x4e, 0x45, 0x53, 0x00, 0x14, 0x00}
	_, err = conn.Write(send)
	if err != nil {
		fmt.Println("Eriniel: Error sending BNES:", err)
		return false, nil
	}

	output := make([]byte, 500)
	n, err := conn.Read(output)
	if err != nil {
		fmt.Println("Eriniel: Error reading after BNES:", err)
		return false, nil
	}

	erinServer := ""
	if n > 0 {
		erinServer = string(output[9 : 9+bytes.IndexByte(output[9:], 0x00)]) // Extract server name, accounting for null terminator
	}

	// BNXI packet
	send = []byte{0x42, 0x4e, 0x58, 0x49, 0x00, 0x14, 0x00}
	_, err = conn.Write(send)
	if err != nil {
		fmt.Println("Eriniel: Error sending BNXI:", err)
		return false, nil
	}

	n, err = conn.Read(output)
	if err != nil {
		fmt.Println("Eriniel: Error reading after BNXI:", err)
		return false, nil
	}

	serverInfo := &ServerInfo{}
	online, err := strconv.Atoi(string(hex.EncodeToString(output)[20:22]))
	if err != nil {
		fmt.Println("error in conversion")
	}
	max, err := strconv.Atoi(string(hex.EncodeToString(output)[22:24]))
	if err != nil {
		fmt.Println("error in conversion")
	}

	if n > 0 {
		// Basic information extraction
		serverInfo.ServerName = erinServer
		if n > 10 {
			serverInfo.Version = parseVersion(string(output[10:11]), "") // Assuming version is in the 11th byte
		}
		if n > 21 {
			serverInfo.PlayersOnline = online
		}
		if n > 23 {
			serverInfo.PlayersMax = max
		}
		if n > 27 {
			serverInfo.PvP = parsePVP(string(hex.EncodeToString(output[26:28]))) // Assuming PvP info
		}
		if n > 41 {
			serverInfo.ModuleName = strings.ReplaceAll(string(output[40:40+bytes.IndexByte(output[40:], 0x00)]), "_", " ") // Assuming module name and replacing underscores
		}

		if serverInfo.ServerName == "" {
			serverInfo.ServerName = serverInfo.ModuleName
		}

		serverInfo.Description = "This server is not enabled with SkyWing's Gamespy replacement server, more info here: http://www.neverwinternights.info/builders_hosts.htm"

		return true, serverInfo
	}

	return false, nil
}

// parseInt converts a string to an integer, returning 0 if the conversion fails.
func parseInt(s string) int {
	var i int
	fmt.Sscan(s, &i)
	return i
}

// parsePVP parses the PvP setting from the server response.
func parsePVP(pvp string) string {
	switch pvp {
	case "0":
		return "None"
	case "1":
		return "Party"
	case "2":
		return "Full PvP"
	default:
		return "Unknown"
	}
}

// parseVersion parses the version from the server response.
func parseVersion(version string, expansion string) string {
	switch version {
	case "8109":
		ver := "NWN1"
		switch expansion {
		case "1":
			ver += "+SoU"
		case "2":
			ver += "+HotU"
		case "3":
			ver += "+SoU+HotU"
		}
		return ver
	case "1765":
		return "NWN2"
	default:
		return "NWN?"
	}
}
