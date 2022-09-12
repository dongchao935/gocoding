package main

// ssh远程连接Linux主机

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

const sshPassword = "******"
const userName = "******"
const sshHost = "******"
const sshPort = 22
const sshKeyFile = "E:/xshell_download/id_rsa"

func runCmd(cfg *ssh.ClientConfig, sshAddr, cmdLine string) error {
	client, err := ssh.Dial("tcp", sshAddr, cfg)
	if err != nil {
		return fmt.Errorf("ssh连接目标%s失败:%v", sshAddr, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("开启session addr:%s失败:%v", sshAddr, err)
	}
	// 设置 session的 tty 配置
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	err = session.RequestPty("xterm", 24, 80, modes)
	if err != nil {
		return fmt.Errorf("设置TTY:%s失败:%v", sshAddr, err)
	}
	defer session.Close()

	// stdout stderr
	var b, eb bytes.Buffer
	session.Stdout = &b
	session.Stderr = &eb
	err = session.Run(cmdLine)
	log.Printf("HOST:[%s]  CMD:[%s] Err:[%v] OUT:[%s] OUT_Err:[%s]\n", sshAddr, cmdLine, err, b.String(), eb.String())
	if err != nil {
		return fmt.Errorf("ssh执行cmd:[ %s ]失败:%v", cmdLine, err)
	}
	return nil
}

func main() {

	key, err := ioutil.ReadFile(sshKeyFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	singer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sshConfig := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
		// 同时设置了密码和密钥文件
			ssh.Password(sshPassword),
			ssh.PublicKeys(singer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshAddr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	cmdLine := "pwd;ping baidu.com -c 6"

	err = runCmd(sshConfig, sshAddr, cmdLine)
	if err != nil {
		fmt.Println(err)
	}
}
