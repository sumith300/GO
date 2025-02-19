package main;
import	"fmt"

type Person struct{
	Name string
	Age int
}

type emp struct{
	Person
	job string
}
func main(){


	// s:=[]int{1,2,3,4,5,6,7,8,9,10}
	// sub:=s[1:3]
	// sub=append(sub,11)
	// index:=2
	// sub=append(s[:index],s[index+1:]...)

	// s1:=make([]int ,5,10)

	// fmt.Println(len(s1))
	// fmt.Println(cap(s1))

	// matrix:=[][]int{
	// 	{1,2,3},
	// 	{4,5,6},
	// 	{7,8,9},
	// }
	// s:=matrix[1][:]

	// p1:=Person{
	// 	Name:"sumith",
	// 	Age:20,
	// }
	// fmt.Println(p1.Name)


	// user:=struct{
	// 	Name string
	// 	Age int
	// }{"sumith",20}
	// fmt.Println(user.Name)

	m:=map[string]int{
		"sumith":20,
		"kumar":21,
		"raja":22,
	}

	m["raja"]=23
	m["raja"]=24

	_,exists:=m["raja"]
	if exists{
		fmt.Println("raja exists")
	}else{
		fmt.Println("raja does not exist")
	}
	delete(m,"raja")
	for k,v :=range m{
		fmt.Println(k,v)
	}

	m2:=make(map[string]int)

	m2["sumith"]=20
	for key,value:=range m2{
		fmt.Println(key,value)
	}
	
}
