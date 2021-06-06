package mudp_test

import (
	"fmt"

	"github.com/gogf/gf/net/gudp"
	"github.com/gogorepos/skeleton/net/mudp"
)

func Example_server() {
	// 创建 UDP 组播
	s := mudp.NewServer("224.224.1.1:35001", func(conn *gudp.Conn) {
		// 不要忘记关闭连接
		defer conn.Close()
		for {
			// 接受数据示例
			if data, err := conn.Recv(-1); err == nil {
				if len(data) > 0 {
					// 发送数据示例
					_ = conn.Send([]byte("> " + string(data)))
					fmt.Printf("Recv: %s", string(data))
				}
			}
		}
	})
	// 启动服务
	s.Run()

	// Output:
	// Recv: xxx
}
