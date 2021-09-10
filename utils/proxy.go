package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func LocalSocket5(sshAddr, username, password, localAddr string) {
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         10 * time.Second,
	}

	//监听本地映射端口
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//客户端连接
		go forward(conn, sshAddr, sshConfig)
	}

}

func forward(conn net.Conn, sshAddr string, sshConfig *ssh.ClientConfig) {
	defer conn.Close()
	var b [1024]byte
	_, err := conn.Read(b[:])
	if err != nil {
		log.Println("读取请求数据包失败", err)
	}
	if b[0] == 0x05 {
		log.Println("只处理Socket5协议")
		//客户端回应：Socket服务端不需要验证方式
		conn.Write([]byte{0x05, 0x00})
		n, err := conn.Read(b[:])
		if err != nil {
			log.Println("2次读取请求数据包失败", err)
		}
		var host, port string
		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度
		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))
		targetAddr := net.JoinHostPort(host, port)
		log.Println("target:", targetAddr)

		//建立本地与SSH服务器的连接
		sshClient, err := ssh.Dial("tcp", sshAddr, sshConfig)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer sshClient.Close()
		log.Printf("%s  ===>  %s 连接建立成功", sshClient.Conn.LocalAddr(), sshAddr)

		//建立ssh服务器 到 后端服务的连接
		forwardConn, err := sshClient.Dial("tcp", targetAddr)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Printf("%s  ===>  %s 连接建立成功", sshAddr, targetAddr)
		defer forwardConn.Close()
		conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发
		var wait2close = sync.WaitGroup{}
		wait2close.Add(1)

		go func() {
			n, err := io.Copy(forwardConn, conn)
			if err != nil {
				log.Println("write", err.Error())
				wait2close.Done()
			}
			log.Printf("入流量共%s", formatFlowSize(int64(n)))
		}()

		go func() {
			n, err := io.Copy(conn, forwardConn)
			if err != nil {
				log.Println("read", err.Error())
				wait2close.Done()
			}
			log.Printf("出流量共%s", formatFlowSize(n))
		}()
		wait2close.Wait()
	}
}

// 字节的单位转换 保留两位小数
func formatFlowSize(s int64) (size string) {
	if s < 1024 {
		return fmt.Sprintf("%.2fB", float64(s)/float64(1))
	} else if s < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(s)/float64(1024))
	} else if s < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(s)/float64(1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(s)/float64(1024*1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(s)/float64(1024*1024*1024*1024))
	} else { //if s < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(s)/float64(1024*1024*1024*1024*1024))
	}
}
