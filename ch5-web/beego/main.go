package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver
)

// Model Struct
type User struct {
	UserId int    `orm:"pk"`
	Name   string `orm:"size(100)"`
}

func init() {
	// set default database
	orm.RegisterDataBase("default", "mysql", "root:root_test@tcp(114.67.98.210:3396)/user_tmp?charset=utf8", 30)

	// register model
	orm.RegisterModel(new(User))

	// create table
	orm.RunSyncdb("default", false, true)
}

func main() {
	o := orm.NewOrm()

	user := User{Name: "aoho"}

	// insert
	id, err := o.Insert(&user)
	fmt.Printf("ID: %d, ERR: %v\n", id, err)

	// update
	user.Name = "boho"
	num, err := o.Update(&user)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)

	// read one
	u := User{UserId: user.UserId}
	err = o.Read(&u)
	fmt.Printf("ERR: %v\n", err)

	var maps []orm.Params

	res, err := o.Raw("SELECT * FROM user").Values(&maps)
	fmt.Printf("NUM: %d, ERR: %v\n", res, err)
	for _, term := range maps {
		fmt.Println(term["user_id"], ":", term["name"])
	}
	// delete
	num, err = o.Delete(&u)
	fmt.Printf("NUM: %d, ERR: %v\n", num, err)
}
