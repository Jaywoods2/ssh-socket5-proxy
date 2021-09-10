package main

import "ssh-socket5-proxy/utils"

func main()  {
	utils.LocalSocket5("192.168.33.14:22", "root", "vagrant", "localhost:10111")
}
