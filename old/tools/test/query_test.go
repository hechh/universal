package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestQuery(t *testing.T) {
	data, err := ioutil.ReadFile("ttt.html")
	if err != nil {
		t.Log(err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		t.Log("Error loading HTML:", err)
		return
	}

	t.Log("--->", doc.Find("tbody").Text())
}

func TestUrl(t *testing.T) {
	rsp, err := http.Get("https://quote.eastmoney.com/ztb/detail#type=ztgc")
	if err != nil {
		t.Log("--1--->", err)
	}
	doc, err := goquery.NewDocumentFromResponse(rsp)
	if err != nil {
		t.Log("--2--->", err)
	}

	// 使用CSS选择器来获取页面中的元素（示例）
	t.Log("---------->", doc.Find("tbody").Text())

	/*
		doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
			// 输出元素文本内容
			fmt.Println(s.Text())
		})
	*/
}
