package iswitch

type IfUnit struct {
	ID     string `json:"id"`
	IP     string `json:"ip"`
	Mac    string `json:"mac"`
	Status int    `json:"status"`
	Desc   string `json:"desc"`
	Speed  int    `json:"speed"`
}
