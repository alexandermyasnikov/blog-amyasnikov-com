package main

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
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
	if pattern == "" {
		pattern = ".*"
	}

	err := conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	checkErr(err)

	var buf bytes.Buffer
	bufTmp := make([]byte, 1024)

	defer func() {
		fmt.Printf("%s", buf.String())
	}()

	for {
		w, err := conn.Read(bufTmp)
		if err != nil {
			break
		}

		buf.Write(bufTmp[:w])

		re := regexp.MustCompile(pattern)
		res := re.FindStringSubmatch(buf.String())

		if len(res) > 0 {
			return res[1:]
		}
	}

	return nil
}

func write(conn net.Conn, msgs ...string) {
	for _, msg := range msgs {
		fmt.Printf("%s", msg)
		_, err := conn.Write([]byte(msg))
		checkErr(err)
	}
}

func main() {
	seeds := 15

	conn, err := net.Dial("tcp", "0.0.0.0:420")
	checkErr(err)
	defer conn.Close()

	_ = readAll(conn, ">")
	write(conn, "1", "\n")

	_ = readAll(conn, "Gimme the index of a sentence:")
	write(conn, "2", "\n")

	_ = readAll(conn, "Gimme the number of seeds:")
	write(conn, fmt.Sprintf("%d", seeds), "\n")

	seed := big.NewInt(0xFFFFFFFFFF)
	for range seeds {
		_ = readAll(conn, "Seed of the number")
		write(conn, seed.String(), "\n")

		seed.Lsh(seed, 5*8)
	}

	result := readAll(conn, `Result: \[(.*)\]`)
	write(conn, "2", "\n")

	resultList := strings.Split(result[0], ", ")
	numResult := big.NewInt(0)
	for _, num := range resultList {
		bigNum, _ := new(big.Int).SetString(num, 10)
		numResult.Or(numResult, bigNum)
	}

	log.Printf("sequence: %q", numResult.Bytes())
}
