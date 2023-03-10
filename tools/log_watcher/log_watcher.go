package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	logPath = "dev"
)

func main() {
	sshHost := "10.24.13.177"
	sshUser := "overlord_test"
	sshPassword := "overlord999"
	sshPort := 22

	config := &ssh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		logError("ssh.Dial() error: %v", err)
		return
	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		logError("sshClient.NewSession() error: %v", err)
		return
	}

	defer session.Close()

	now := time.Now()
	cmd := fmt.Sprintf("tail -n 100 /app/logs/%s/game/runtime/%s/%d.log", logPath, now.Format("2006-01-02"), now.Hour())
	ret, err := session.Output(cmd)
	if err != nil {
		logError("Not found log")
		return
	}

	logInfo(string(ret))

	pressAnyKeyExit()
}
