package main

// 交叉编译
//CGO_ENABLED=0 GOOS=windows  GOARCH=amd64 go build test.go

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)


var mysqlConfig Mysql
var db *sql.DB
var help *bool

type Mysql struct{
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Network string `yaml:"network"`
	Server string `yaml:"server"`
	Port int `yaml:"port"`
	Database string `yaml:"database"`
}

func init()  {
	usage := "1.第一个参数输入all,将dzzoffice所有用户执行授权操作\n2.第一个参数为用户名，例如:user01 则表示给user01用户授权\n3.第一个参数为 user01,user02,user03... 表示给多个用户授权"

	help = flag.Bool("h", false, usage)
	flag.Parse()
	if *help{
		fmt.Println(usage)
		os.Exit(0)
	}
	GetMysqlInfo()
	GetMysqlConn()
}

func GetUserList()[]string{
	username := os.Args[1]
	if strings.Contains(username,","){
		userSplit := strings.Split(username , ",")
		for index, user := range userSplit{
			if user == ""{
				//[1,2,3,4,5]
				userSplit = append(userSplit[:index],userSplit[index+1:]...)
			}
		}
		if len(userSplit) ==  0{
			fmt.Println("错误的参数输入,",username)
			os.Exit(1)
		}else {
			return userSplit
		}
	}else if username == "all"{
		return GetAllUsers(db)
	}
	return []string{username}

}

func GetAllUsers(db *sql.DB)[]string{
	out , err := db.Prepare("select username from dzz_user where uid!=?")
	if err !=nil{
		fmt.Println("select all users error:",err)
		os.Exit(1)
	}
	rows , err := out.Query(1)
	var UserList []string
	var username string
	for rows.Next(){
		err := rows.Scan(&username)
		if err != nil{
			fmt.Println("select all user scan error:",err)
		}
		UserList = append(UserList , username)
	}
	return UserList

}


func GetMysqlInfo(){
	var mysqlInfo map[string]Mysql
	// 读取文件
	conf , err := ioutil.ReadFile("config.yml")
	if err != nil{
		fmt.Println("read file config.yml faild :",err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(conf , &mysqlInfo)
	if err != nil{
		fmt.Println("解析文件错误",err)
		os.Exit(1)
	}
	mysqlConfig = mysqlInfo["mysql"]

}

func GetMysqlConn(){
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Network, mysqlConfig.Server, mysqlConfig.Port, mysqlConfig.Database)
	//fmt.Println(dsn)
	dbcon , err := sql.Open("mysql" , dsn)
	if err != nil{
		fmt.Println("open connection err:",err)
		os.Exit(1)
	}
	err = dbcon.Ping()
	if err == nil{
		fmt.Println("数据库连接成功✌️")
	}else{
		fmt.Println("连接数据库失败🙅‍♂️")
	}
	db = dbcon
}



func main() {
	// 异常处理
	defer func() {
		err := recover()
		errstr := fmt.Sprintf("%s" , err)
		if err != nil{

			if errstr == "runtime error: index out of range [1] with length 1"{
				fmt.Println("请输入正确的参数，例如bob|bob,tom,jerry")
			}
		}
	}()
	// 获取用户名
	userList := GetUserList()

	defer db.Close()
	rows , _ := db.Query("select pathkey from dzz_organization")
	colums , _ := rows.Columns()

	tempId := []int{}
	for  rows.Next(){
		rows.Scan(&colums[0])
		temp := strings.Split(colums[0] , "_")
		orgid , err := strconv.Atoi(temp[len(temp)-2:len(temp)-1][0])
		if err != nil {
			fmt.Println("strcov err:",err)
			continue
		}
		tempId = append(tempId , orgid)
	}
	//fmt.Println(tempId)
	// 获取forigid
	rows , _ = db.Query("select forgid from dzz_organization")
	colums , _ = rows.Columns()
	for rows.Next(){
		rows.Scan(&colums[0])
		temp ,err := strconv.Atoi(colums[0])
		if err != nil{
			fmt.Println(err)
			continue
		}
		for i , v := range tempId{
			if v == temp{
				tempId = append(tempId[:i] , tempId[i+1:]...)
			}
		}
	}
	//fmt.Println(tempId)
	for _ , username := range userList {
		r := db.QueryRow("select uid from dzz_user where username=?", username)
		var uid int
		err := r.Scan(&uid)
		if err != nil && err == sql.ErrNoRows{
			fmt.Println("没有在数据库中查询到用户:",username,"跳过对该用户的授权")
			continue
		}else if err !=nil{
			fmt.Println("在数据库中查询用户",username,"时候发生错误:",err)
			os.Exit(1)
		}
		// 此处判断用户是否为管理员， 如果是管理员不进行操作
		var adminid , groupid int
		row := db.QueryRow("select adminid , groupid from dzz_user where uid=?" , uid)
		if err != nil{
			fmt.Println("select adminid & groupid error :",err)
			os.Exit(1)
		}

		row.Scan(&adminid , &groupid)
		if adminid == 1 && groupid == 1{
			fmt.Println(username,"用户是管理员,跳过授权操作")
			continue
		}

		db.Exec("delete from dzz_organization_user where uid=?", uid)
		date := time.Now().Unix()
		for _, id := range tempId {
			// 先从将之前的授权信息删除
			_ , err = db.Exec("insert into dzz_organization_user values(? , ? , 0 , ?)", id, uid, date)
			if err != nil{
				fmt.Println("授权失败：",err)
				continue
			}
		}
		fmt.Println("用户:",username,"授权成功!")
	}


}