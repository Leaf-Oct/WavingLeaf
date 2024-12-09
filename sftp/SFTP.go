package sftp

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func setupSSHServer() *ssh.ServerConfig {
	// 创建一个新的 ServerConfig 并添加主机密钥
	privateBytes, err := os.ReadFile("id_rsa") // 读取私钥文件
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(private)

	return config
}

func handleConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	// 协商 SSH 连接
	newConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}
	defer newConn.Close()

	log.Printf("SSH connection from %s (%s)", newConn.RemoteAddr(), newConn.ClientVersion())

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "subsystem":
					if string(req.Payload[4:]) == "sftp" {
						req.Reply(true, nil)
						server, err := sftp.NewServer(channel, &SimpleFileSystem{})
						if err != nil {
							log.Printf("Failed to create SFTP server: %v", err)
							return
						}
						err = server.Serve()
						if err != nil && err != io.EOF {
							log.Printf("SFTP server error: %v", err)
						}
						return
					}
				default:
					req.Reply(false, nil)
				}
			}
		}(requests)
	}
}

type SimpleFileSystem struct{}

func (fs *SimpleFileSystem) Lstat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

func (fs *SimpleFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (fs *SimpleFileSystem) OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flag, perm)
}

func (fs *SimpleFileSystem) Remove(path string) error {
	return os.Remove(path)
}

func (fs *SimpleFileSystem) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (fs *SimpleFileSystem) Mkdir(path string, perm os.FileMode) error {
	return os.Mkdir(path, perm)
}

func (fs *SimpleFileSystem) Rmdir(path string) error {
	return os.Remove(path)
}

func (fs *SimpleFileSystem) Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

func (fs *SimpleFileSystem) Chown(path string, uid, gid int) error {
	return os.Chown(path, uid, gid)
}

func (fs *SimpleFileSystem) CreateDirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

func TestSFTP() {
	listener, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Fatalf("Failed to listen on port 2222: %v", err)
	}
	defer listener.Close()

	log.Println("Listening for SFTP connections on :2222")

	config := setupSSHServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %v", err)
			continue
		}

		go handleConnection(conn, config)
	}
}
