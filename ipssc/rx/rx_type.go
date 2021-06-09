package rx

type Entity struct {
	Online         string `json:"online"`
	Area           string `json:"area"`
	Name           string `json:"name"`
	Id             string `json:"id"`
	Ip             string `json:"ip"`
	DevType        string `json:"dev_type"`
	Version        string `json:"version"`
	FpgaVersion    string `json:"fpga_version"`
	GroupO         string `json:"group_name_o"`
	GroupS         string `json:"group_name_s"`
	GroupT         string `json:"group_name_t"`
	GroupF         string `json:"group_name_f"`
	Blink          string `json:"blink"`
	Multicast      string `json:"multicast"`
	Port           string `json:"port"`
	Display        string `json:"display"`
	Bind           string `json:"bind"`
	Hdcp           string `json:"hdcp"`
	DisableUsbType string `json:"disable_usb_type"`
	VideoFormat    string `json:"video_format"`
	Resolution     string `json:"resolution"`
	BaudRate       string `json:"baud_rate"`
	AnalogVol      string `json:"analog_vol"`
	HardDisk       string `json:"hard_disk"`
	DataBit        string `json:"data_bit"`
	ParityBit      string `json:"parity_bit"`
	StopBit        string `json:"stop_bit"`
	Link           string `json:"link"`
	Mac            string `json:"mac"`
	Conflict       bool   `json:"conflict"`
}
