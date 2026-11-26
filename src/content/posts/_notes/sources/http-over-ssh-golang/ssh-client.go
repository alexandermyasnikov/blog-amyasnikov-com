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