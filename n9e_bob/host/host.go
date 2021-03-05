package host

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var N9e_server = "72.160.3.56"

const GetAllHostsApi  string = "/api/ams-ce/hosts"


type HostInfo struct{
	Id int `json:"id"`
	Sn string `json:"sn"`
	Ip string `json:"ip"`
	Ident string `json:"ident"`
	Note string `json:"note"`
	Cpu string	`json:"cpu"`
	Name string	`json:"name"`
	Mem string `json:"mem"`
	Disk string `json:"disk"`
	Cate string	`json:"cate"`
	Clock string `json:"clock"`
	Tenant string `json:"tenant"`
}

type Hosts struct{
	Dat `json:"dat"`
	Err string `json:"err"`
}

type Dat struct {
	List []HostInfo `json:"list"`
	Total string `json:"total"`
}


func GetHosts()*Hosts{

	var hosts Hosts
	url := "https://" + N9e_server + GetAllHostsApi
	resp , err := http.Get(url)
	if err != nil{
		fmt.Println("GetHost Get method Error:",err)
		return nil
	}
	defer resp.Body.Close()
	resp_splitbytes , err  := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("read from resp body error:",err)
		return nil
	}
	err = json.Unmarshal(resp_splitbytes , &hosts)
	fmt.Println(hosts)
	return &hosts
}