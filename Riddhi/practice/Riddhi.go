package main

import "fmt"

func main(){
	var n int
		
	fmt.Scan(&n)	
	
	arr := make([]int, n)

	for i:=0;i<n;i++{
		fmt.Scan(&arr[i])
	}

	// print elements 
	for i:=0;i<n;i++{
		fmt.Println(arr[i])
	}

	// print even numbers 
	var count int
    for i:=0;i<n;i++{
      if arr[i]%2==0{
		count ++;
	  }
	}
	fmt.Printf("%d\n", count)

	//find an element 
	var key int
	fmt.Scan(&key)
    for i:=0;i<n;i++{
		if key==arr[i]{
			fmt.Printf("Required element found at index %d\n", i)
			break;
		}
	}

	// print characters
	var word string 
	fmt.Scan(&word)
	
	for ch:=range word {
		fmt.Printf("%c\n", word[ch])
	}

	//odd or even 

    for num:= range arr{
		if arr[num]%2==0{
			fmt.Printf("%v = even\n",arr[num])
		}else{
			fmt.Printf("%v = odd\n",arr[num])
		}
	}

}