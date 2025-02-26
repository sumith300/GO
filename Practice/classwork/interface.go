package main


import "fmt"

type Ht interface{
	utility() int
}

type superHeroes struct{
	name string
	peopleSaved int
	power int
}

type superVillian struct{
	name string
	peopleHarmed int
	power int
}

func(s *superHeroes) utility() int{
	return s.peopleSaved * s.power
}

func (s *superVillian) utility() int{
	return s.peopleHarmed * s.power
}

func main(){	




	batman:=superHeroes{
		name: "batman",
		peopleSaved: 10,
		power : 5,
	}

	joker:=superVillian{
		name: "joker",
		peopleHarmed: 10,
		power:8,
	}

	hitsquad:=[]Ht{joker,batman}

	total:=0
	for _,value:=range hitsquad{
		total+=value.utility()
	}

	fmt.Println(total)
}