package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/ssh"
)

const (
	sshUser    = "<user>"
	sshHost    = "<host>"
	sshPort    = 22
	sshKeyPath = "<key>"
	httpAddr   = "localhost:8000"
)

type SSHTransport struct {
	sshClient *ssh.Client
	httpAddr  string
}

func (t *SSHTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	session, err := t.sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	var stdout bytes.Buffer
	session.Stdout = &stdout

	var stderr bytes.Buffer
	session.Stderr = &stderr

	if err := session.Start("socat - TCP:" + t.httpAddr); err != nil {
		return nil, fmt.Errorf("failed to start socat: %w", err)
	}

	if err := req.Write(stdin); err != nil {
		return nil, fmt.Errorf("failed to write request: %v", err)
	}

	if err := session.Wait(); err != nil {
		return nil, fmt.Errorf("socat session error: %w", err)
	}

	if stderr.Len() > 0 {
		log.Printf("stderr: %s", stderr.String())
		return nil, fmt.Errorf("socat session error")
	}

	resp, err := http.ReadResponse(bufio.NewReader(&stdout), req)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return resp, nil
}

func NewSshClient() (*ssh.Client, error) {
	key, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshAddr := fmt.Sprintf("%s:%d", sshHost, sshPort)

	client, err := ssh.Dial("tcp", sshAddr, sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	log.SetFlags(log.Lmicroseconds)

	client, err := NewSshClient()
	if err != nil {
		panic(err)
	}

	defer client.Close()

	for _, cmd := range []string{"hostname", "uptime", "whoami"} {
		session, err := client.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()

		output, err := session.CombinedOutput(cmd)
		if err != nil {
			panic(err)
		}

		log.Printf("%s -> %s", cmd, output)
	}

	httpClient := &http.Client{
		Transport: &SSHTransport{
			sshClient: client,
			httpAddr:  httpAddr,
		},
	}

	for _, arg := range []string{"my-arg-1"} {
		req, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
		if err != nil {
			panic(err)
		}

		req.Header.Add("Custom-Header", arg)

		q := req.URL.Query()
		q.Add("lang", "go")
		q.Add("query-arg", arg)
		req.URL.RawQuery = q.Encode()

		resp, err := httpClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		log.Printf("Body: %s - %s", resp.Status, body)
	}
}
