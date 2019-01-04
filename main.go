package main

import (
	"fmt"
	"github.com/go-xorm/core"
	"github.com/tarm/serial"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func main() {
	// 配置并打开端口
	c := &serial.Config{Name: "COM10", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// 连接数据库
	engine, err := xorm.NewEngine("mysql", "root:root@peng@tcp(222.200.184.65:13306)/tb_database?charset=utf8")
	handleError(err)
	engine.Logger().SetLevel(core.LOG_DEBUG)

	// 循环读出数据并写入数据库
	for {
		// 读出数据
		buf := make([]byte, 128)
		n, err := s.Read(buf)
		handleError(err)
		if n == 0 {
			continue
		}
		// D2:T&Ph:25C  08.0
		deviceText := fmt.Sprintf("%s", string(buf[:n]))
		fmt.Printf("New data received: %s\n", deviceText)

		// 解析数据
		temp := deviceText[8:10]
		temperature, err := strconv.Atoi(temp)
		if err != nil {
			log.Println(err)
			temperature = -1000
			fmt.Println("Invalid temperature in given data, set to -1000 as default.")
		}
		temp = deviceText[13:]
		ph, err := strconv.ParseFloat(temp, 32)
		if err != nil {
			log.Println(err)
			ph = -100.0
			fmt.Println("Invalid ph in given message, set to -100.0 as default.")
		}
		dataToBeStored := TemperaturePHAndTime{
			Temperature: temperature,
			PH:          float32(ph),
			Datatime:    time.Now(),
		}

		// 写入数据库
		affected, err := engine.Insert(&dataToBeStored)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%d data was inserted.\n", affected)
	}
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type TemperaturePHAndTime struct {
	Temperature int       `xorm:"'temperature'"`
	PH          float32   `xorm:"'PH'"`
	Datatime    time.Time `xorm:"'datatime'"`
}

func (t TemperaturePHAndTime) TableName() string {
	return "tb_T&PH"
}
