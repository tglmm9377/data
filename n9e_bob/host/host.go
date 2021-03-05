package host

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var N9e_server = "72.160.3.56"

const GetAllHostsApi  string = "/api/ams-ce/hosts"

const token = "5db38d09ff83ec7c39856c5e2b822f5e"


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

var Client *http.Client
var url string = "http://"+N9e_server + GetAllHostsApi


func init(){
	SetClient(url)
}

func SetClient(url  string)error{
	client := &http.Client{}
	req , err := http.NewRequest("GET" , url , nil)
	if err != nil{
		fmt.Println(err)
		return err
	}
	req.Header.Set("X-User-Token" , token)
	_ , err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	Client = client
	return nil
}


func GetHosts()*Hosts{
	//X-User-Token: xxxx"
	var hosts Hosts
	url := "http://" + N9e_server + GetAllHostsApi
	SetClient(url)
	resp , err := Client.Get(url)
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