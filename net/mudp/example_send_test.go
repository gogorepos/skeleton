package mudp_test

import (
	"log"

	"github.com/gogorepos/skeleton/net/mudp"
)

func Example_send() {
	// 发送字节
	err := mudp.Send("224.224.1.1:35001", []byte("hello"))
	if err != nil {
		log.Printf("Send err: %s", err)
	}

	// 发送字符串
	err = mudp.SendString("224.224.1.1:35001", "hello")
	if err != nil {
		log.Printf("Send string err: %s", err)
	}

	// 将数据 JSON 序列化后再发送，本质发送的是 JSON 字符串
	data := map[string]int{"bob": 98, "小明": 99}
	err = mudp.SendJson("224.224.1.1:35001", data)
	if err != nil {
		log.Printf("Send json err: %s", err)
	}
}
