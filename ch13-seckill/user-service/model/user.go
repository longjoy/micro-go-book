package model

import (
	"github.com/gohouse/gorose/v2"
	"github.com/longjoy/micro-go-book/ch13-seckill/pkg/mysql"
	"log"
)

type User struct {
	UserId     int64    `json:"user_id"`     //Id
	UserName   string `json:"user_name"`   //用户名称
	Password   string `json:"password"`    //密码
	Age        int    `json:"age"`         //年龄

}

type UserModel struct {
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (p *UserModel) getTableName() string {
	return "user"
}

func (p *UserModel) GetUserList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Get()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	return list, nil
}

/*func (p *UserModel) GetUserByUsername(username string) (*User,  error)  {

	conn := mysql.DB()
	if result, err := conn.Table(p.getTableName()).Where(map[string]interface{}{"username": username}).First(); err == nil{

	}else {
		return nil, err
	}

}*/

func (p *UserModel) CheckUser(username string, password string) (*User, error) {
	conn := mysql.DB()
	data, err := conn.Table(p.getTableName()).Where(map[string]interface{}{"user_name": username, "password": password}).First()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	user := &User{
		UserId:data["user_id"].(int64),
		UserName:data["user_name"].(string),
		Password:data["password"].(string),
		Age:int(data["age"].(int64)),


	}
	return user, nil
}



func (p *UserModel) CreateUser(user *User) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"user_id":     user.UserId,
		"user_name":   user.UserName,
		"password":    user.Password,
		"age":         user.Age,

	}).Insert()
	if err != nil {
		log.Printf("Error : %v", err)
		return err
	}
	return nil
}
