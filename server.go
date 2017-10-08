package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

//URLs sdfsdf
var (
	mainURL = "https://getpocket.com"
	URLs    = map[string]string{
		"RequestUrl":      "%s/v3/oauth/request",
		"RequestTokenUrl": "%s/auth/authorize?request_token=%s&redirect_uri=%s",
		"RequestAuthUrl":  "%s/v3/oauth/authorize",
		"Redirect_uri":    "http://localhost:8080",
		"GetUrl":          "%s/v3/get",
	}
)

//Result can be structured with this tool https://mholt.github.io/json-to-go/
type Result struct {
	MsgID    string `json:"msg_id"`
	Text     string `json:"_text"`
	Entities struct {
		Subject []struct {
			Confidence int    `json:"confidence"`
			Value      string `json:"value"`
			Type       string `json:"type"`
		} `json:"subject"`
		Number []struct {
			Confidence int    `json:"confidence"`
			Value      int    `json:"value"`
			Type       string `json:"type"`
		} `json:"number"`
		Intent []struct {
			Confidence float64 `json:"confidence"`
			Value      string  `json:"value"`
		} `json:"intent"`
	} `json:"entities"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		//I receive a request from Telegram

		//I prepare the request to send to wit.ai to interpret what is being said in Telegram
		req, err := http.NewRequest("GET", "https://api.wit.ai/message", nil)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		res := Result{}   //Result has the structure of a JSON response from wit.ai
		doJSON(req, &res) //I send a resquest to wit.ai, the response is parsed trough Result
		//return c.JSON(http.StatusOK, res)
		return c.String(http.StatusOK, res.Entities.Intent[0].Value)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func doJSON(req *http.Request, res interface{}) error {
	q := req.URL.Query()
	q.Add("q", "Hello")
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	req.Header.Add("Authorization", "Bearer GPXW7M4BANRCYS2NDT2WYPUVX7ZOLOBS")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error when sending request to the server")
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(res)
}
