package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// Conectar al servidor
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	var windowSize, packetSize int
	fmt.Print("Tamaño de la ventana: ")
	fmt.Scan(&windowSize)
	fmt.Print("Tamaño de los paquetes: ")
	fmt.Scan(&packetSize)
	var fileName string
	fmt.Print("Nombre del archivo: ")
	fmt.Scan(&fileName)

	// Enviar tamaño de la ventana, de los paquetes y nombre del archivo al servidor
	err = binary.Write(conn, binary.LittleEndian, int32(windowSize))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	err = binary.Write(conn, binary.LittleEndian, int32(packetSize))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Enviar el archivo
	fileUpload(conn, windowSize, packetSize, fileName)
}

func fileUpload(conn net.Conn, windowSize, packetSize int, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	fmt.Println("Enviando archivo...")
	buf := make([]byte, packetSize)
	for {
		bytesRead, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error: ", err)
			}
			break
		}

		_, err = conn.Write(buf[:bytesRead])
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if bytesRead%windowSize == 0 {
			ack := make([]byte, 4)
			_, err = conn.Read(ack)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
	}

}
