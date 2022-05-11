package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func help() {
	var number string
	//控制台提示语句
	fmt.Print("请输入一个整数：")
	//控制台的输出
	fmt.Scan(&number)
	fmt.Println("数值是：", number)
	fmt.Printf("数据类型是：%T\n", number)
	//数据类型转换string---> int
	value, _ := strconv.Atoi(number)
	fmt.Printf("转换后的数据类型是：%T\n", value)
	//数值判断
	if value > 100 {
		fmt.Println("数值较大")
	} else {
		fmt.Println("数值较小")
	}
}

func createFile(start int, value string) {
	filePath := strconv.Itoa(start) + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file err:", err)
		return
	}
	//关闭文件
	defer file.Close()
	//写入的内容
	str := value
	//带缓存写入文件
	writer := bufio.NewWriter(file)
	writer.WriteString(str)
	//内容是写入到缓存中的，需要Flush()
	writer.Flush()
}

func addList() {
	createFile(1, "123456")
}

func run() {
	var number string
	fmt.Print("请输入一个整数：")
	fmt.Scan(&number)
	addList()
}

func main() {

}
