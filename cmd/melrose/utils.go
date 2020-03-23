package main

import "fmt"

func printInfo(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;32minfo:\033[0m"}, args...)...)
}

func printError(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;31merror:\033[0m"}, args...)...)
}

func printWarning(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;33mwarning:\033[0m"}, args...)...)
}

func shortTypename(v interface{}) string {
	return ""
}
