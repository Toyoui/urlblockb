package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url" // 新添加的库
	"os"
	"strconv"
	"time"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Page string `json:"page"`
		Data []struct {
			ID         int    `json:"id"`
			Title      string `json:"title"`
			Content    string `json:"content"`
			CreateTime string `json:"create_time"`
		} `json:"data"`
	} `json:"data"`
}
type Flash struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

var apiURL = "https://xizhi.qqoq.net/XZ3a8c491f49fb57b5ae792ab0f30e077b.channel"

func flashIDs() (map[int]bool, error) {
	file, err := os.Open("flashid.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ids := make(map[int]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		ids[id] = true
	}

	return ids, nil
}

func main() {
	tburl := "https://api.theblockbeats.news/v1/open-api/open-flash?size=3&page=1"
	duration := 2 * time.Minute

	file, err := os.OpenFile("flashid.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("打开flashid.txt文件出错：", err)
		return
	}
	defer file.Close()

	for {
		ids, err := flashIDs()
		if err != nil {
			fmt.Println("读取flashid.txt文件出错：", err)
			continue
		}

		response, err := http.Get(tburl)
		if err != nil {
			fmt.Println("请求出错：", err)
			continue
		}

		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println("读取响应出错：", err)
			continue
		}

		var data Response
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("解析JSON出错：", err)
			continue
		}

		var uniflashs []Flash
		for _, flash := range data.Data.Data {
			if ids[flash.ID] {
				fmt.Println("存在相同文章")
				continue
			}

			uniflashs = append(uniflashs, Flash{
				ID:         flash.ID,
				Title:      flash.Title,
				Content:    flash.Content,
				CreateTime: flash.CreateTime,
			})
		}
		fmt.Println("============================================")

		for _, flash := range uniflashs {
			fmt.Printf("ID: %d\n", flash.ID)
			fmt.Printf("Title: %s\n", flash.Title)
			fmt.Printf("Content: %s\n", flash.Content)
			timestamp, err := strconv.ParseInt(flash.CreateTime, 10, 64)
			if err != nil {
				fmt.Println("解析时间戳出错：", err)
				continue
			}
			createTime := time.Unix(timestamp, 0)
			createTimeFormatted := createTime.Format("2006-01-02 15:04:05")
			fmt.Printf("Create Time: %s\n", createTimeFormatted)

			title := flash.Title
			content := "发布时间: " + "\n" + createTimeFormatted + "\n" + "\n" + "概述: " + "\n" + "\n" + flash.Content
			formData := url.Values{
				"title":   {title},
				"content": {content},
			}

			response, err := http.PostForm(apiURL, formData)
			if err != nil {
				fmt.Println("发送请求出错：", err)
				continue
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				fmt.Println("请求未成功，状态码：", response.StatusCode)
				continue
			}

			// 如果需要，你可以像这样读取响应体
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("读取响应体出错：", err)
				continue
			}
			fmt.Println("请求成功，响应体：", string(responseBody))

			_, err = file.WriteString(strconv.Itoa(flash.ID) + "\n")
			if err != nil {
				fmt.Println("写入flashid.txt文件出错：", err)
			}
		}
		fmt.Println("============================================")

		time.Sleep(duration)
	}
}
