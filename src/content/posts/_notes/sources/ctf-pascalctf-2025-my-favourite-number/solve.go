package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.Lmsgprefix)

	readFile, err := os.Open("output.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)

	scanner.Scan() // read: hi, i'm Alice, my public parameters are:
	scanner.Scan() // read: n=...
	scanner.Scan() // read: e=...
	scanner.Scan() // read: hi Alice! i'm Bob, my public parameters are:

	scanner.Scan() // read: n=...
	bobPNumStr := strings.TrimPrefix(scanner.Text(), "n=")

	scanner.Scan() // read: e=...
	bobENumStr := strings.TrimPrefix(scanner.Text(), "e=")

	scanner.Scan() // read: let's play a game, you have to guess my favourite number

	bobPNum, _ := new(big.Int).SetString(bobPNumStr, 10)
	bobENum, _ := new(big.Int).SetString(bobENumStr, 10)

	lowerbound := new(big.Int).SetInt64(0)
	upperbound := new(big.Int).Exp(big.NewInt(2), big.NewInt(501), nil)

	for new(big.Int).Sub(upperbound, lowerbound).Cmp(big.NewInt(1)) > 0 {
		scanner.Scan() // read: Is your number greater than {mid}?
		scanner.Scan() // read: Yes... or No...
		aliceMsgEnc := strings.TrimPrefix(scanner.Text(), "alice: ")

		mid := new(big.Int).Div(new(big.Int).Add(upperbound, lowerbound), big.NewInt(2))

		msgYes := fmt.Sprintf("Yes!, my number is greater than %s", mid)
		msgEncYes := new(big.Int).Exp(new(big.Int).SetBytes([]byte(msgYes)), bobENum, bobPNum)

		msgNo := fmt.Sprintf("No!, my number is lower or equal to %s", mid)
		msgEncNo := new(big.Int).Exp(new(big.Int).SetBytes([]byte(msgNo)), bobENum, bobPNum)

		log.Printf("mid: %126x", mid.Bytes())

		switch {
		case aliceMsgEnc == msgEncYes.String():
			lowerbound = mid
		case aliceMsgEnc == msgEncNo.String():
			upperbound = mid
		default:
			fmt.Println("unknown answer")
			upperbound = lowerbound
		}
	}

	log.Printf("answer: %s", upperbound.Bytes())
}
