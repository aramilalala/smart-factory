package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tarm/serial"
)

func main() {
	// 创建文件夹及文件
	os.Mkdir(filepath.Join("./data"), os.ModeDir)
	f, err := os.OpenFile(filepath.Join("./data", "temp.log"), os.O_APPEND|os.O_WRONLY, 0600)
	if os.IsNotExist(err) {
		f, err = os.Create(filepath.Join("./data", "temp.log"))
	}
	cnt := true
	defer f.Close()

	// 配置并打开端口
	c := &serial.Config{Name: "COM4", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// 循环读出数据并写入文件
	for {
		buf := make([]byte, 128)
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		text := fmt.Sprintf("%s", string(buf[:n]))
		if cnt {
			text = time.Now().Format("2006-01-02 15:04:05") + " " + text
		}
		cnt = !cnt
		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}
		// fmt.Println(text)
	}
}
