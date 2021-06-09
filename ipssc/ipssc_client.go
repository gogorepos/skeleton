package ipssc

import (
	"net/http"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogorepos/skeleton/net/tcp"
)

type Ipssc struct {
	Header string      `json:"cmd_header"`
	Body   interface{} `json:"cmd_body"`
}

func (i Ipssc) Address() string {
	return "169.254.0.253:6000"
}

func (i Ipssc) Command() []byte {
	return gjson.New(i).MustToJson()
}

func (i Ipssc) Check(data *gjson.Json) (*gjson.Json, error) {
	if code := data.GetString("cmd_end.Err_code"); code != "0" && code != "" {
		err := gerror.New(http.StatusText(http.StatusBadRequest))
		// USB 切换将错误码放到了 Err_str 中
		if str := data.GetString("cmd_end.Err_str"); str != "" {
			switch str {
			case "-2":
				err = ErrUsbOccupy
			case "-3":
				err = ErrUsbFailed
			case "-4":
				err = ErrUsbDisable
			}
		}
		return nil, err
	}
	return data.GetJson("cmd_body"), nil
}

func Send(header string, body interface{}) (*gjson.Json, error) {
	if g.IsNil(body) {
		body = struct{}{}
	}
	return tcp.Send(Ipssc{
		Header: header,
		Body:   body,
	})
}
