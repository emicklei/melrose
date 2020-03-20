package main

import "fmt"

func printInfo(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;32mINFO:\033[0m"}, args...)...)
}

func printError(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;31mERROR:\033[0m"}, args...)...)
}

func printWarning(args ...interface{}) {
	fmt.Println(append([]interface{}{"\033[1;33mWARNING:\033[0m"}, args...)...)
}
