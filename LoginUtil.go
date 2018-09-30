// LoginUtil
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
    "github.com/axgle/mahonia"
	"net/http"
	"net/url"
	"strings"
)

//登陆
var count = 0 //记录登陆次数

func loginNet() (sign int, err error) {
	const responseURL = "http://192.168.252.254"
	var account, password string

	fmt.Printf("Please input your accunt:")
	fmt.Scanln(&account)
	fmt.Printf("Please input your password:")
	fmt.Scanln(&password)

	var client *http.Client
	client = http.DefaultClient

	postDict := map[string]string{
		"DDDDD":  account,
		"upass":  password,
		"0MKKey": "Login script written by python3",
	}
	postValues := url.Values{}
	for postKey, PostValue := range postDict {
		postValues.Set(postKey, PostValue)
	}
	postDataStr := postValues.Encode()
	postDataBytes := []byte(postDataStr)
	postBytesReader := bytes.NewReader(postDataBytes)
	httpReq, _ := http.NewRequest("POST", responseURL, postBytesReader)
	httpResp, err1 := client.Do(httpReq)
	if err1 != nil {
		return 0 , err1
	}
	defer httpResp.Body.Close()

	flag, _ := checkLogin()
	if flag == false {
		if count < 3 {
			fmt.Println("login fail! Please check your account or password.")
			loginNet()

		} else {
			fmt.Println("You have login four times!")
			return 0, err
		}
	}
	fmt.printf("The account of %s ",account)
	return 1, err

}

//登出
func loginOut() {
	_, err := http.Get("http://192.168.252.254/F.htm")
	if err != nil {
		fmt.Println("Login out error is ", err)
		return
	}
	flag, err := checkLogin()
	if err != nil {
		fmt.Println("Login out is fail ")
		return
	}
	if flag == false {
		fmt.Println("Login out success!")
		return
	}

	return
}
func checkLogin() (flag bool, err error) {
	flag = false
	const URL = "http://192.168.252.254/"
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("checklogin error is ", err)
		return flag, err
	}
	decoder := mahonia.NewDecoder("gb2312")

	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println("error is ", err1)
		return flag, err
	}
	result := decoder.ConvertString(string(body))

	flag = strings.Contains(result, "uid=")

	return flag, nil

}
func main() {
	var chose string
	//检测账号是否已经登陆
	flag, err := checkLogin()
	if err != nil {
		fmt.Println("error is ", err)
		return
	}
	if flag == true {
		fmt.Println("YOU ALREADY LOGIN!")
		//询问是否退出登陆
		fmt.Printf("Do you want to login out now?(y or n):")
		fmt.Scanln(&chose)
		if chose == "n" {
			return
		}
		loginOut()
	}
	fmt.Printf("Do you want to login in now:(y or n)")
	fmt.Scanln(&chose)
	if chose == "n" {
		return
	}
	sign, err := loginNet()
	if err != nil{
	    fmt.Println("login error is ",err)
	    return 
	}
	if sign == 1 {
		fmt.Println("login success !")
	}
	return
}
