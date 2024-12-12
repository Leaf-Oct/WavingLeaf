package sftp

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 用户名和密码的映射
var userPasswords = map[string]string{
	"leaf": "test",
}

func passwordAuth(sshConn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	// 获取用户名
	username := sshConn.User()
	// 验证用户名和密码
	if pwd, ok := userPasswords[username]; ok && pwd == string(password) {
		return nil, nil // 验证成功
	}
	return nil, ssh.ServerAuthError{}
}

// 处理 SSH 连接
func handleConn(conn net.Conn) {
	defer conn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, nil)
	if err != nil {
		log.Println("Failed to establish SSH connection:", err)
		return
	}
	log.Println("Logged in:", sshConn.User())

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() == "session" {
			print("是session")
			go handleSession(newChannel)
		} else if newChannel.ChannelType() == "sftp" {
			print("是sftp")
			go handleSFTP(newChannel)
		}
	}
}

// 处理 SSH 会话
func handleSession(newChannel ssh.NewChannel) {
	channel, _, err := newChannel.Accept()
	if err != nil {
		log.Println("Could not accept channel:", err)
		return
	}
	io.Copy(channel, channel) // Echo input
}

// 处理 SFTP
func handleSFTP(newChannel ssh.NewChannel) {
	channel, _, err := newChannel.Accept()
	if err != nil {
		log.Println("Could not accept channel:", err)
		return
	}

	sftpServer, err := sftp.NewServer(channel)
	if err != nil {
		log.Println("Failed to start SFTP server:", err)
		return
	}

	if err := sftpServer.Serve(); err != nil {
		log.Println("SFTP server exited:", err)
	}
}

func TestSFTP() {
	privateBytes, err := os.ReadFile("github") // 私钥文件
	if err != nil {
		log.Fatal("Failed to load private key:", err)
	}

	privateKey, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Unable to parse private key:", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: false,
		PasswordCallback: passwordAuth,
	}
	config.AddHostKey(privateKey)

	listener, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Fatal("Failed to listen for connection:", err)
	}

	log.Println("Listening on :2222...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleConn(conn)
	}
}
