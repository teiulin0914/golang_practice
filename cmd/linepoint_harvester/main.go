package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {

	codes := readPointCodes("./codes")

	if len(codes) == 0 {
		return
	}

	var id, pwd string
	fmt.Print("請輸入 LINE 帳號: ")
	fmt.Scanln(&id)
	fmt.Print("請輸入 LINE 密碼: ")
	password, err := terminal.ReadPassword(0)
	if err == nil {
		pwd = string(password)
	}

	fmt.Println()
	log.Println("---start---")

	failedCodes, errorCodes := takePoints(id, pwd, codes)

	if len(failedCodes) > 0 {
		writePointCodes("./failed_codes", failedCodes)
	}
	if len(errorCodes) > 0 {
		writePointCodes("./error_codes", errorCodes)
	}

	log.Println("---end---")
}

func takePoints(id string, pwd string, codes []string) ([]string, []string) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var failedCodes []string = make([]string, 0)
	var errorCodes []string = make([]string, 0)
	for i, code := range codes {
		fmt.Printf("%d,%s", i+1, code)

		var text string
		var err error
		if i == 0 {
			err = chromedp.Run(chromeCtx,
				chromedp.Navigate("https://points.line.me/pointcode?pincode="+code),
				chromedp.Sleep(100*time.Millisecond),
				chromedp.WaitEnabled(".MdBtn01"),
				chromedp.Click(".MdBtn01"),
				chromedp.Sleep(100*time.Millisecond),
				chromedp.WaitVisible("#id"),
				chromedp.SetValue("#id", id),
				chromedp.SetValue("#passwd", pwd),
				chromedp.WaitEnabled(".MdBtn03Login"),
				chromedp.Click(".MdBtn03Login"),
				chromedp.Sleep(1000*time.Millisecond),
				chromedp.WaitReady("body"),
				chromedp.Location(&text),
			)
		} else {
			err = chromedp.Run(chromeCtx,
				chromedp.Navigate("https://points.line.me/pointcode?pincode="+code),
				chromedp.Sleep(100*time.Millisecond),
				chromedp.WaitEnabled(".MdBtn01"),
				chromedp.Click(".MdBtn01"),
				chromedp.Sleep(1000*time.Millisecond),
				chromedp.WaitReady("body"),
				chromedp.Location(&text),
			)
		}

		if err != nil {
			errorCodes = append(errorCodes, code)
			fmt.Println(",[Error]")

			if i == 0 {
				log.Fatal("首筆失敗！請至 https://points.line.me/pointcode/ 手動領取，確認輸入資料無誤。")
			}

		} else if !strings.Contains(text, "complete") && !strings.Contains(text, "compleate") { // LINE 網址打錯字
			failedCodes = append(failedCodes, code)
			fmt.Println(",[Failed]")

			if i == 0 {
				log.Fatal("首筆失敗！請至 https://points.line.me/pointcode/ 手動領取，確認輸入資料無誤。")
			}

		} else {
			fmt.Println(",[Succeed]")
		}

		time.Sleep(500 * time.Millisecond)
	}

	return failedCodes, errorCodes
}

func readPointCodes(filePath string) []string {
	var codes []string = make([]string, 0)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > 0 {
			codes = append(codes, text)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return codes
}

func writePointCodes(filePath string, codes []string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, code := range codes {
		writer.WriteString(code + "\n")
	}
	writer.Flush()
}
