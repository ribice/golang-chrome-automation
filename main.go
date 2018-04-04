package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"
	"time"

	"github.com/chromedp/chromedp/kb"

	"github.com/chromedp/chromedp"
)

func main() {
	cfg, err := readConfig()
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	checkErr(err)

	// run task list
	checkErr(c.Run(ctxt, changeDrupalSettings(cfg)))

	// shutdown chrome
	checkErr(c.Shutdown(ctxt))

	// wait for chrome to finish
	checkErr(c.Wait())

	log.Println("Successfully changed Drupal settings")
}

func readConfig() (*config, error) {
	_, filePath, _, _ := runtime.Caller(0)
	pwd := filePath[:len(filePath)-7]
	txt, err := ioutil.ReadFile(pwd + "/config.json")
	if err != nil {
		return nil, err
	}
	var cfg = new(config)
	if err := json.Unmarshal(txt, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func changeDrupalSettings(cfg *config) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(cfg.URL + "user/login"),
		chromedp.WaitVisible(`#edit-name`, chromedp.ByID),
		chromedp.SendKeys(`#edit-name`, cfg.Username, chromedp.ByID),
		chromedp.SendKeys(`#edit-pass`, cfg.Password, chromedp.ByID),
		chromedp.Click("#edit-submit"),
		chromedp.Sleep(1 * time.Second),
		chromedp.Navigate(cfg.URL + "admin/appearance/settings/bootstrap#edit-advanced"),
		chromedp.WaitVisible(`#edit-cdn`, chromedp.ByID),
		chromedp.Click(`#edit-cdn`),
		chromedp.Click(`#edit-cdn-provider`),
		chromedp.SendKeys(`#edit-cdn-provider`, "c"+kb.Select, chromedp.ByID),
		chromedp.WaitVisible(`#edit-cdn-custom-css`, chromedp.ByID),
		chromedp.Clear(`#edit-cdn-custom-css`),
		chromedp.Clear(`#edit-cdn-custom-css-min`),
		chromedp.Clear(`#edit-cdn-custom-js`),
		chromedp.Clear(`#edit-cdn-custom-js-min`),
		chromedp.SendKeys(`#edit-cdn-custom-css`, cfg.BootstrapCSS, chromedp.ByID),
		chromedp.SendKeys(`#edit-cdn-custom-css-min`, cfg.BootstrapCSSMin, chromedp.ByID),
		chromedp.SendKeys(`#edit-cdn-custom-js`, cfg.BootstrapJS, chromedp.ByID),
		chromedp.SendKeys(`#edit-cdn-custom-js-min`, cfg.BootstrapJSMin, chromedp.ByID),
		chromedp.Click("#edit-submit"),
		chromedp.Sleep(1 * time.Second),
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type config struct {
	URL             string `json:"url"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	BootstrapCSS    string `json:"bootstrap_css"`
	BootstrapCSSMin string `json:"bootstrap_css_min"`
	BootstrapJS     string `json:"bootstrap_js"`
	BootstrapJSMin  string `json:"bootstrap_js_min"`
}
