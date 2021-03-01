package main

import (
	"database/sql"
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

type Mysql struct{
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Network string `yaml:"network"`
	Server string `yaml:"server"`
	Port int `yaml:"port"`
	Database string `yaml:"database"`
}

func init()  {
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
	for _ , username := range userList {
		var uid int
		r := db.QueryRow("select uid from dzz_user where username=?", username)
		err := r.Scan(&uid)
		if err != nil && err == sql.ErrNoRows{
			fmt.Println("没有在数据库中查询到用户:",username,"跳过对该用户的授权")
			continue
		}else if err !=nil{
			fmt.Println("在数据库中查询用户",username,"时候发生错误:",err)
			os.Exit(1)
		}
		if uid == 0 {
			os.Exit(0)
		}
		date := time.Now().Unix()
		for id := range tempId {
			_ , err = db.Exec("insert into dzz_organization_user values(? , ? , 0 , ?)", id, uid, date)
			if err != nil{
				fmt.Println("授权失败：",err)
				os.Exit(1)
			}
		}
	}


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


}