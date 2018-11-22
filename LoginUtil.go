package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

//上网登录
func loginNet(account, password string) (sign int, err error) {
	const responseURL = "http://192.168.252.254"
	var client *http.Client
	client = http.DefaultClient

	postDict := map[string]string{
		"DDDDD":  account,
		"upass":  password,
		"0MKKey": "Login script written goalng",
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
		log.Fatal("post1 error is:", err1)
		return 0, err1
	}
	defer httpResp.Body.Close()
	flag, _ := checkLogin()
	if flag == false {
		fmt.Println("login fail! Please check your account or password")
		return 0, err
	}
	fmt.Printf("The account of %s ", account)
	return 1, err

}

//备用的登陆地址
func loginNetTwo(account, password string) (sign int, err error) {

	const responseURL = "http://172.31.252.71:801/eportal/?c=ACSetting&a=Login"
	var client *http.Client
	client = http.DefaultClient
	postDict := map[string]string{
		"DDDDD":  account,
		"upass":  password,
		"0MKKey": "Login script written goalng",
	}
	postValues := url.Values{}
	for postKey, PostValue := range postDict {
		postValues.Set(postKey, PostValue)
	}
	postDataStr := postValues.Encode()
	postDataBytes := []byte(postDataStr)
	postBytesReader := bytes.NewReader(postDataBytes)
	httpReq, _ := http.NewRequest("POST", responseURL, postBytesReader)

	httpReq.Header.Set("Accept", "*/*")
	httpReq.Header.Set("Accept-Encoding", " gzip, deflate")
	httpReq.Header.Set("Accept-Language", " zh-CN,zh;q=0.9")
	httpReq.Header.Set("Origin", "http://172.31.252.71")
	httpReq.Header.Set("Cache-Control", "no-cache")
	httpReq.Header.Set("Connection", " keep-alive")
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.92 Safari/537.36")

	httpResp, err1 := client.Do(httpReq)
	if err1 != nil {
		log.Fatal("post2 err:", err)
		return 0, err1
	}
	defer httpResp.Body.Close()
	flag, _ := checkLogin()
	if flag == false {
		fmt.Println("login fail! Please check your account or password.")
		return 0, err
	}
	fmt.Printf("The account of %s ", account)
	return 1, err

}

//登出网络
func loginOut() {
	_, err := http.Get("http://192.168.252.254/F.htm")
	if err != nil {
		fmt.Println("Log out error is ", err)
		return
	}
	flag, err := checkLogin()
	if err != nil {
		fmt.Println("Log out is fail ")
		return
	}
	if flag == false {
		fmt.Println("Log out success!")
		return
	}

	return
}

//检查是否登陆
func checkLogin() (flag bool, err error) {
	flag = false
	const URL = "http://192.168.252.254/"
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
		return flag, err
	}
	decoder := mahonia.NewDecoder("GBK")

	defer resp.Body.Close()

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		log.Fatal("error is ", err1)
		return flag, err1
	}
	result := decoder.ConvertString(string(body))

	flag = strings.Contains(result, "uid=")

	return flag, nil

}

//强制切换用户
func changeUser() (flag int, err error) {
	var account, password string

	fmt.Println("Please enter the new accout:")
	fmt.Printf("Please input your accunt:")
	fmt.Scanln(&account)
	fmt.Printf("Please input your password:")
	fmt.Scanln(&password)

	flag, err = loginNet(account, password)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	if flag != 1 {
		fmt.Println("Login fail!")
		return 0, err
	}
	fmt.Println("change login success !")
	return 1, nil
}

//初始化登录
func loginInit() (r int, err1 error) {

	var chose string
	var account, password string
	var times, timeout float64

	//检测账号是否已经登陆
	flag, err := checkLogin()
	if err != nil {
		log.Fatal("error is ", err)
		return 0, err
	}
	if flag == true {
		fmt.Println("YOU ALREADY LOGIN!")
		//询问是否退出登陆
		fmt.Printf("Do you want to login out now?(y or n):")
		fmt.Scanln(&chose)
		if chose == "n" {
			return 1, err
		}
		loginOut()
	}
	fmt.Printf("Do you want to login in now:(y or n)")
	fmt.Scanln(&chose)
	if chose == "n" {
		return 1, nil
	}

	fmt.Printf("Please input your accunt:")
	fmt.Scanln(&account)
	fmt.Printf("Please input your password:")
	fmt.Scanln(&password)

	sign, err := loginNet(account, password)
	if err != nil {
		log.Fatal("login error is ", err)
		return 0, err
	}
	if sign == 1 {
		fmt.Println("login success !")
		fmt.Printf("Please set time to check(second)(enter -1 to pass):")
		fmt.Scanln(&times)
		fmt.Printf("Please set time to keep(hour)(enter -1 to pass):")
		fmt.Scanln(&timeout)
		if times >= 0 {
			keepOnline(times, timeout, account, password)
		} else {
			fmt.Println("Pass this set up!")
		}
	}
	return 1, nil

}

//下一步操作
func nextAction() (sign int) {
	var chose int
	fmt.Println("Please enter next operate:")
	fmt.Printf("1、change user  2、exit script: ")
	fmt.Scanln(&chose)
	switch chose {
	case 1:
		flag, err := changeUser()
		if err != nil {
			log.Fatal("change error is ", err)
		}
		if flag != 1 {
			fmt.Println("change login fail!")
		}
	case 2:
		return 1
	default:
		fmt.Println("Please enter the right chose!")
	}
	return 0
}

//后台保护
func keepOnline(times, until float64, account, password string) {

	until = until * 3600 //数据换算

	tc := time.Tick(time.Duration(times) * time.Second)

	timeout := time.NewTimer(time.Duration(until) * time.Second)

	go func() {
		for _ = range tc {
			flag, err := checkLogin()
			if err != nil {
				log.Fatal(err)
				return
			}
			if flag == false {
				_, err1 := loginNet(account, password)
				if err1 != nil {
					log.Fatal(err1)
					return
				}

				fmt.Println("return login!")
			}

		}

	}()
	<-timeout.C
	loginOut()
	return
}

//命令行函数
func command() (chose bool) {

	pattern := flag.Bool("i", false, "Enter the interactive interface(true or false):")
	loginleft := flag.Bool("l", false, "Disconnect the internet(true or false)")
	account := flag.String("a", "", "Input your accout")
	password := flag.String("p", "", "Input your password")
	compulsive := flag.Bool("c", false, "Compulsive login")
	times := flag.Float64("t", -1, "Set time(second)")
	until := flag.Float64("u", -1, "keep online until(hour)")
	spare := flag.Bool("s", false, "Use the spare URL to login in")

	flag.Parse()
	sign, err := checkLogin()
	if err != nil {
		log.Fatal("check login error is:", err)
		return false
	}

	if *loginleft == true {
		loginOut()
		return false
	}
	if *pattern == true {
		return *pattern
	}
	if *account != " " && *password != " " && (sign == false || *compulsive == true) {
		var sign int
		var err error

		if *spare == true {
			sign, err = loginNetTwo(*account, *password)
		} else {
			sign, err = loginNet(*account, *password)
		}
		if err != nil {
			log.Fatal("login error is ", err)
			return false
		}
		if sign == 1 {
			fmt.Println("login success !")
			//执行后台保护
			if *times >= 0 {
				keepOnline(*times, *until, *account, *password)
			}
		}
	} else if sign == true && len(*account) > 1 {
		fmt.Println("You have Login,don't login once more!")
	} else {
		fmt.Println("optional arguments:")
		fmt.Println("-i , Interactive interface ,   change into interactive interface(bool) ")
		fmt.Println("-l , Loginout ,  Disconnect the internet(enter -l=true/false)")
		fmt.Println("-a , Account ,  Input your login account")
		fmt.Println("-p , Password , Input your login password")
		fmt.Println("-c , Compulsive login , Compulsive login the Internet(bool)")
		fmt.Println("-t , set time(second), set time to keep online")
		fmt.Println("-u , set time(hour), to make sure on line")
		fmt.Println("-s , use other url to login(bool)")

	}
	return false
}
func main() {
	//默认为参数行执行
	flag := command()
	if flag == true {
		fmt.Println("Welcome to the interactive interface:")
		//执行初始化函数
		result, err := loginInit()
		if err != nil {
			log.Fatal("error is ", err)
			return
		}
		if result != 1 {
			return
		}
		//执行下一步操作
		sign := nextAction()
		for sign == 0 {
			sign = nextAction()
		}
	}
	return

}
