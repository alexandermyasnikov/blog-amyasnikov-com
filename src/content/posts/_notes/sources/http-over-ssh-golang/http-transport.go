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