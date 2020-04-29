package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Id          string
	Name        string
	Habits      string
	CreatedTime string
}

var tpl = `<html>
<head>
<title></title>
</head>
<body>
<form action="/info" method="post">
	用户名:<input type="text" name="username">
	兴趣爱好:<input type="text" name="habits">
	<input type="submit" value="提交">
</form>
</body>
</html>`

func connect(cName string) (*mgo.Session, *mgo.Collection) {
	session, err := mgo.Dial("mongodb://47.96.140.41:27017/") //Mongodb's connection
	checkErr(err)
	session.SetMode(mgo.Monotonic, true)
	//return a instantiated collect
	return session, session.DB("test").C(cName)
}

func queryByName(name string) []User {
	var user []User
	s, c := connect("user")
	defer s.Close()
	err := c.Find(bson.M{"name": name}).All(&user)
	checkErr(err)
	return user
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	var t *template.Template
	t = template.New("Products") //创建一个模板
	t, _ = t.Parse(tpl)
	log.Println(t.Execute(w, nil))
}

func store(user User) error {
	//插入数据
	s, c := connect("user")
	defer s.Close()
	user.Id = bson.NewObjectId().Hex()
	return c.Insert(&user)
}

func userInfo(w http.ResponseWriter, r *http.Request) {
	//请求的是登录数据，那么执行登录的逻辑判断
	_ = r.ParseForm()
	if r.Method == "POST" {
		user1 := User{Name: r.Form.Get("username"), Habits: r.Form.Get("habits")}
		store(user1)
		fmt.Fprintf(w, " %v", queryByName("aoho")) //这个写入到w的是输出到客户端的
	}
}

func main() {
	http.HandleFunc("/form", submitForm)     //设置访问的路由
	http.HandleFunc("/info", userInfo)       //设置访问的路由
	err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
