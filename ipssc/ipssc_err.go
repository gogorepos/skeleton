package ipssc

import "github.com/gogf/gf/errors/gerror"

var (
	ErrUsbOccupy  = gerror.New("usb occupy")
	ErrUsbFailed  = gerror.New("usb failed")
	ErrUsbDisable = gerror.New("usb disable")
)
