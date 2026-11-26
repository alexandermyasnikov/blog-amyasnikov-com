package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type ScoreDTO struct {
	Key []int `json:"key,omitempty"`
}

const serverURL = "https://t-trickster-jbi8aw9z.spbctf.ru"

func start() (string, []byte, error) {
	uuid := uuid.NewString()
	flag, err := request("PUT", serverURL+"/start/"+uuid)
	return uuid, flag, err
}

func next(uuid string) ([]byte, error) {
	data, err := request("GET", serverURL+"/next/"+uuid)
	return data, err
}

func accept(uuid string, answer string) ([]byte, error) {
	data, err := request("GET", serverURL+"/accept/"+uuid+"/"+answer)
	return data, err
}

func request(method, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("request: invalid status: %v", resp.Status)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func calcAnswer(bs []byte, answer string) string {
	switch {
	case string(bs[len(bs)-3]) == answer:
		return "2"
	case string(bs[len(bs)-5]) == answer:
		return "1"
	case string(bs[len(bs)-7]) == answer:
		return "0"
	default:
		return "X"
	}
}

func encryptFlag(flag, data []byte) string {
	var score ScoreDTO
	err := json.Unmarshal(data, &score)
	if err != nil {
		return ""
	}

	if len(score.Key) == 0 {
		return ""
	}

	flagEncrypted := ""
	for i := range flag {
		flagEncrypted += string(byte(int(flag[i]) ^ score.Key[i%len(score.Key)]))
	}

	return flagEncrypted
}

func robotrickster() error {
	uuid, flagEncrypted, err := start()
	if err != nil {
		return err
	}

	fmt.Printf("uuid: %v\n", uuid)
	fmt.Printf("flag: %v\n", flagEncrypted)

	answer := "1"

	for {
		data, err := next(uuid)
		if err != nil {
			return err
		}

		fmt.Printf("positions: %v\n", string(data))

		answer = calcAnswer(data, answer)

		data, err = accept(uuid, answer)
		if err != nil {
			return err
		}

		fmt.Printf("score:    %v\n", string(data))

		flag := encryptFlag(flagEncrypted, data)
		if flag != "" {
			fmt.Printf("flag: %v\n", flag)
			break
		}
	}

	return nil
}

func main() {
	err := robotrickster()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}
