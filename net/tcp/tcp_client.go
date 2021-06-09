package tcp

import (
	"fmt"
	"net"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gcmd"
)

type Client interface {
	Address() string                             // TCP 服务器的地址
	Command() []byte                             // 需要发送的命令
	Check(data *gjson.Json) (*gjson.Json, error) // 对接受到的信息进行判断
}

func Send(client Client) (*gjson.Json, error) {
	command := client.Command()
	// 启动参数 "-c"，打印发送的信息
	if gcmd.ContainsOpt("c") {
		fmt.Printf("Send ==> \n%s\n", gjson.New(command).MustToJsonIndentString())
	}
	// 启动参数 "-n"，不请求服务器
	if gcmd.ContainsOpt("n") {
		return nil, nil
	}
	conn, err := net.Dial("tcp", client.Address())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if _, err = conn.Write(command); err != nil {
		return nil, err
	}
	var result []byte
	for {
		var buf [1024]byte
		var n int
		n, err = conn.Read(buf[:])
		if err != nil {
			return nil, err
		}
		result = append(result, buf[:n]...)
		if gjson.Valid(result) {
			break
		}
	}
	data := gjson.New(result)
	// 启动参数 "-r"，打印接受到的信息
	if gcmd.ContainsOpt("r") {
		fmt.Printf("Recv ==> \n%s\n", data.MustToJsonIndentString())
	}
	return client.Check(gjson.New(result))
}
