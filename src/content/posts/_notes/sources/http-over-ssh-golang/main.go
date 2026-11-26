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

		log.Printf("Status: %s", resp.Status)
		log.Printf("Body: %s", body)
	}
}