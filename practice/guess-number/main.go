package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Welcome! Let's guess the number!")

	maxNum := 100
	rand.Seed(time.Now().UnixNano())
	target := rand.Intn(maxNum)

	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again.")
			continue
		}

		input = strings.TrimSuffix(input, "\r\n")
		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("The input is not a number. Please try again.")
			continue
		}

		fmt.Printf("Your input is %d. ", number)

		if number > target {
			fmt.Println("The number you input is bigger.")
		} else if number < target {
			fmt.Println("The number you input is smaller.")
		} else {
			fmt.Println("The number is right!")
			return
		}
	}
}
