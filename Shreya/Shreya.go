package main

import "fmt"

//Q1. Reverse an array
func revarray(arr []int) {
	if len(arr) == 0 {
		fmt.Println("Array is empty")
		return
	}
	var i = 0
	var j = len(arr) - 1
	for i < j {
		var temp = arr[i]
		arr[i] = arr[j]
		arr[j] = temp
		i++
		j--
	}
}

//Q2.Count even numbers
func countEven(arr []int) int {
	if len(arr) == 0 {
		fmt.Println("Array is empty")
		return 0
	}
	var count = 0
	for i := 0; i < len(arr); i++ {
		if arr[i]%2 == 0 {
			count++
		}
	}
	return count
}

//Q3. Linear Search
func linearsearch(arr []int, target int) int {
	for i, value := range arr {
		if value == target {
			return i
		}
	}
	return -1
}

//Q4.Reverse a string
func revString(s string) string {
	var revstring = ""
	for i := len(s) - 1; i >= 0; i-- {
		revstring = revstring + string(s[i])
	}
	return revstring
}

//Q.5Check for a prime number
func checkPrime(n int) int {
	if n <= 1 {
		return 0
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return 0
		}
	}
	return 1
}

func main() {
	var array = []int{1, 2, 3, 4, 5, 6}

	revarray(array)
	fmt.Println("Reversed array is:", array)

	var evenNum = countEven(array)
	fmt.Println("Total even numbers are:", evenNum)

	var target = 4
	var result = linearsearch(array, target)

	if result == -1 {
		fmt.Println("Element not found")
	} else {
		fmt.Printf("Element found at %d\n", result)
	}

	var str = "golang"
	var reversestr = revString(str)
	fmt.Println("Reversed string:", reversestr)

	var number = 19
	var answer = checkPrime(number)
	if answer == 1 {
		fmt.Printf("%d is a prime number\n", number)
	} else {
		fmt.Printf("%d is not a prime number\n", number)
	}
}
