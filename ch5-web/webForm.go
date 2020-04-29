package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.tpl")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断
		_ = r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		if pwd := r.Form.Get("password"); pwd == "123456" { // 验证密码是否正确
			fmt.Fprintf(w, "欢迎登陆，Hello %s!", r.Form.Get("username")) //这个写入到w的是输出到客户端的
		} else {
			fmt.Fprintf(w, "密码错误，请重新输入!")
		}
	}
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello aoho!") //这个写入到w的是输出到客户端的
}

func main() {
	http.HandleFunc("/", sayHelloName)       //设置访问的路由
	http.HandleFunc("/login", login)         //设置访问的路由
	err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
