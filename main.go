package main

import (
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/api"
	"log"
)

func main() {
	cfg, err := config.SetupEnv()

	if err != nil {
		log.Fatalf("config file is not loaded properly %v\n", err)
	}
	api.StartServer(cfg)
}

/*
func main() {
	fmt.Println("I am main function")

	app := fiber.New()

	//Basics Types: int, float64, string, bool
	//composite types: array, slice, map, struct
	//pointer types: *

	//Array (has fixed size)

	// var myFamily [3]string
	// myFamily[0] = "Manish"
	// myFamily[1] = "Ardent"
	// myFamily[2] = "FCB"

	myFamily := [3]string{"Manish", "Ardent", "FCB"}
	myFamily[2] = "Vas"

	myCourses := [3][2]string{
		{"Go", "NodeJS"},
		{"AWS", "GCP"},
		{"CDK", "Pulumi"},
	}

	fmt.Println("Available Courses %v", myCourses)

	fmt.Println("My Family: %v", myFamily)

	//Slice (dynamic)
	var myFriends []string
	myFriends = append(myFriends, "Teku", "KJ", "Tau")
	myFriends = append(myFriends, "Vasu")
	fmt.Println("My Friends: %v", myFriends)

	mySliceCourses := [][]string{
		{"Go", "NodeJS"},
		{"AWS", "GCP"},
		{"CDK", "Pulumi"},
	}

	course := []string{"IAC", "Cloud Formation"}

	mySliceCourses = append(mySliceCourses, course)

	fmt.Println("Available Courses %v", mySliceCourses)

	// var age int
	// var height float64
	// var firstName string
	// var isEmployed bool

	//NOTE: Another way to declare
	// age := 23
	// height := 111.11
	// firstName := "Manish"
	// isEmployed := true

	// // fmt.Println(age, height, firstName, isEmployed)

	// fmt.Printf("Age: %d\n", age)
	// fmt.Printf("Height: %f\n", height)
	// fmt.Printf("FirstName: %s\n", firstName)
	// fmt.Printf("Employed?: %t\n", isEmployed)

	// MyHelperFunction()        // i don't to import as both using same package
	// configs.LoadAppSettings() // now we need to import

	// myBeCourse := make([]int, 2, 10)

	myWishList := make(map[string]string)
	myWishList["First"] = "MacPro"
	myWishList["Second"] = "Trip to spain"
	myWishList["Third"] = "Coffee in Naples"

	delete(myWishList, "Third")
	firstWish := myWishList["first"]
	fmt.Println(firstWish)
	fmt.Printf("My wish list %v\n", myWishList)
	app.Listen("localhost:9000")

}
*/
