package main

import (
	"fmt"
	"reflect"
)

// 定义一个人的接口
type Person interface {

	// 和人说hello
	SayHello(name string)
	// 跑步
	Run() string
}

type Hero struct {
	Name string
	Age int
	Speed int
}
func (hero *Hero) SayHello(name string)  {
	fmt.Println("Hello " + name, ", I am " + hero.Name)
}

func (hero *Hero) Run() string {
	fmt.Println("I am running at speed " + string(hero.Speed))
	return "Running"
}

func main()  {


	typeOfHero := reflect.TypeOf(Hero{})
	fmt.Printf("Hero's type is %s, kind is %s", typeOfHero, typeOfHero.Kind())


	//fmt.Printf("*Hero's type is %s, kind is %s",reflect.TypeOf(&Hero{}), reflect.TypeOf(&Hero{}).Kind())


	//typeOfPtrHero := reflect.TypeOf(&Hero{})
	//fmt.Printf("*Hero's type is %s, kind is %s\n",typeOfPtrHero, typeOfPtrHero.Kind())
	//typeOfHero := typeOfPtrHero.Elem()
	//fmt.Printf(" typeOfPtrHero elem to typeOfHero, Hero's type is %s, kind is %s", typeOfHero, typeOfHero.Kind())


	//typeOfHero := reflect.TypeOf(Hero{})
	//
	//// 通过 #NumField 获取结构体字段的数量
	//for i := 0 ; i < typeOfHero.NumField(); i++{
	//	fmt.Printf("field' name is %s, type is %s, kind is %s\n",
	//		typeOfHero.Field(i).Name,
	//		typeOfHero.Field(i).Type,
	//		typeOfHero.Field(i).Type.Kind())
	//}
	//// 获取名称为 Name 的成员字段类型对象
	//nameField, _ := typeOfHero.FieldByName("Name")
	//fmt.Printf("field' name is %s, type is %s, kind is %s\n", nameField.Name, nameField.Type, nameField.Type.Kind())


	//// 声明一个 Person 接口，并用 Hero 作为接收器
	//var person Person = &Hero{}
	//// 获取接口Person的类型对象
	//typeOfPerson := reflect.TypeOf(person)
	//// 打印Person的方法类型和名称
	//for i := 0 ; i < typeOfPerson.NumMethod(); i++{
	//	fmt.Printf("method is %s, type is %s, kind is %s.\n",
	//		typeOfPerson.Method(i).Name,
	//		typeOfPerson.Method(i).Type,
	//		typeOfPerson.Method(i).Type.Kind())
	//}
	//method, _ := typeOfPerson.MethodByName("Run")
	//fmt.Printf("method is %s, type is %s, kind is %s.\n", method.Name, method.Type, method.Type.Kind())


	//name := "小明"
	//valueOfName := reflect.ValueOf(name)
	//fmt.Println(valueOfName.Interface())


	//name := "小明"
	//valueOfName := reflect.ValueOf(name)
	//fmt.Println(valueOfName.Bytes())


	//typeOfHero := reflect.TypeOf(Hero{})
	//heroValue := reflect.New(typeOfHero)
	//fmt.Printf("Hero's type is %s, kind is %s\n", heroValue.Type(), heroValue.Kind())


	//name := "小明"
	//valueOfName := reflect.ValueOf(&name)
	//valueOfName.Elem().Set(reflect.ValueOf("小红"))
	//fmt.Println(name)


	//name := "小明"
	//valueOfName := reflect.ValueOf(name)
	//fmt.Printf( "name can be address: %t\n", valueOfName.CanAddr())
	//valueOfName = reflect.ValueOf(&name)
	//fmt.Printf( "&name can be address: %t\n", valueOfName.CanAddr())
	//valueOfName = valueOfName.Elem()
	//fmt.Printf( "&name's Elem can be address: %t", valueOfName.CanAddr())


	//hero := &Hero{
	//	Name: "小白",
	//}
	//
	//valueOfHero := reflect.ValueOf(hero).Elem()
	//
	//valueOfName := valueOfHero.FieldByName("Name")
	//// 判断字段的 Value 是否可以设定变量值
	//if valueOfName.CanSet() {
	//	valueOfName.Set(reflect.ValueOf("小张"))
	//}
	//
	//fmt.Printf("hero name is %s", hero.Name)


	//var person Person = &Hero{
	//	Name: "小红",
	//	Speed: 100,
	//}
	//valueOfPerson := reflect.ValueOf(person)
	//// 获取SayHello 方法
	//sayHelloMethod := valueOfPerson.MethodByName("SayHello")
	//// 构建调用参数并通过 #Call 调用方法
	//sayHelloMethod.Call([]reflect.Value{reflect.ValueOf("小张")})
	//// 获取Run 方法
	//runMethod := valueOfPerson.MethodByName("Run")
	//// 通过 #Call 调用方法并获取结果
	//result := runMethod.Call([]reflect.Value{})
	//fmt.Printf("result of run method is %s.", result[0])


	//var person Person = &Hero{
	//	Name: "小红",
	//}
	//// 获取接口Person的类型对象
	//typeOfPerson := reflect.TypeOf(person)
	//// 打印Person的方法类型和名称
	//sayHelloMethod, _ := typeOfPerson.MethodByName("Run")
	//// 将 person 接收器放在参数的第一位
	//result:=sayHelloMethod.Func.Call([]reflect.Value{reflect.ValueOf(person)})
	//fmt.Printf("result of run method is %s.", result[0])


	//methodOfHello := reflect.ValueOf(hello)
	//methodOfHello.Call([]reflect.Value{})

}

func hello() {
	fmt.Print("Hello World！")
}








