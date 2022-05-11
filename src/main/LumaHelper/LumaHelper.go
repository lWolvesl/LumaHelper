package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var client http.Client

type GAME struct {
	id   string
	name string
}

func getIndex() int {
	pwd, _ := os.Getwd()
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}

	return len(fileInfoList) - 1
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
	fmt.Println("写入文件" + filePath)
}

func setProxy(rawUrl string) http.Client {
	if rawUrl == "" {
		client := http.Client{}
		return client
	}
	ProxyUri, err := url.Parse(rawUrl)
	if err != nil {
		log.Fatal("parse url error: ", err)
	}
	client := http.Client{
		Transport: &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(ProxyUri),
		},
	}
	return client
}

func getUrl(SteamUrl string) io.Reader {
	res, err := client.Get(SteamUrl)
	if err != nil {
		fmt.Println("访问失败")
		os.Exit(0)
	}
	return res.Body
}

func getGames(SteamUrl string) []GAME {
	var games []GAME
	body := getUrl(SteamUrl)
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}
	text := doc.Find("a.search_result_row")
	for i := 0; i < text.Length() && i < 6; i++ {
		attr, _ := text.Eq(i).Attr("href")
		str := strings.Split(attr, "/")

		var game GAME
		game.id = str[4]
		game.name = str[5]
		games = append(games, game)
	}
	return games
}

func getDLCs(game string) []GAME {
	var DlCs []GAME
	SteamUrl := "https://store.steampowered.com/app/" + game
	body := getUrl(SteamUrl)
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}
	text := doc.Find("a.game_area_dlc_row")
	for i := 0; i < text.Length(); i++ {
		attr, _ := text.Eq(i).Attr("href")
		str := strings.Split(attr, "/")

		var game GAME
		game.id = str[4]
		game.name = str[5]
		DlCs = append(DlCs, game)
	}
	return DlCs
}

func SelectGame() string {
	var Game string
	var index int

	fmt.Print("请输入游戏名(英文)：")
	fmt.Scan(&Game)
	SteamUrl := "https://store.steampowered.com/search/?term=" + Game
	games := getGames(SteamUrl)
	for i := 0; i < len(games); i++ {
		fmt.Println(strconv.Itoa(i) + ":" + games[i].name)
	}

	fmt.Println("请选择你的游戏(若无选中游戏，请输入-1)：")
	fmt.Scan(&index)
	if index == -1 {
		return SelectGame()
	}
	if index < 0 || index > len(games) {
		fmt.Println("范围错误，请重新输入")
		return SelectGame()
	}

	return games[index].id
}

func SelectDLC(game string) {
	DLCs := getDLCs(game)
	fmt.Println("DLC列表：")
	for i := 0; i < len(DLCs); i++ {
		fmt.Println(strconv.Itoa(i) + ":" + DLCs[i].name)
	}
	var check int
	fmt.Println("选择全部添加还是选择添加：\n" +
		"1.全部添加\n" +
		"2.逐个添加\n" +
		"3.退出系统\n" +
		"0.重新选择游戏")
	fmt.Scan(&check)
	switch check {
	case 0:
		start()
	case 1:
		index := getIndex()
		for i := 0; i < len(DLCs); i++ {
			createFile(index, DLCs[i].id)
			index++
		}
	case 2:
		var checkX = 1
		for checkX != -1 {
			fmt.Println("输入序号选择游戏(输入-1退出)")
			fmt.Scan(&checkX)
		}
	case 3:
		os.Exit(0)
	}
}

func start() {
	//查找游戏
	game := SelectGame()

	//查找DLC
	SelectDLC(game)
}

func run() {
	//设置代理
	fmt.Print("是否启用代理(启用请输入1，不启用则输入0)：")
	var startProxy int
	fmt.Scan(&startProxy)
	if startProxy == 1 {
		var rawUrl string
		fmt.Print("请输入代理地址：")
		fmt.Scan(&rawUrl)
		setProxy(rawUrl)
	}

	//查找游戏
	game := SelectGame()

	//查找DLC
	SelectDLC(game)

	var check = 0
	fmt.Println("退出系统请输入0/继续添加请输入1")
	fmt.Scan(&check)
	for check != 0 {
		start()
		fmt.Println("退出系统请输入0/继续添加请输入1")
		fmt.Scan(&check)
	}
}

func main() {
	run()
}
