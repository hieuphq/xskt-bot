package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/hieuphq/xskt-bot/errors"
	"github.com/hieuphq/xskt-bot/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	FACEBOOK_API   = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
	IMAGE          = "http://37.media.tumblr.com/e705e901302b5925ffb2bcf3cacb5bcd/tumblr_n6vxziSQD11slv6upo3_500.gif"
	VISIT_SHOW_URL = "http://labouardy.com"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/webhook", VerifyEndpointHandler)
	e.POST("/webhook", MessageEndpointHandler)
	e.Logger.Fatal(e.Start(":8084"))
}

// VerifyEndpointHandler ...
func VerifyEndpointHandler(c echo.Context) error {
	ch := c.QueryParam("hub.challenge")
	m := c.QueryParam("hub.mode")
	token := c.QueryParam("hub.verify_token")

	if ch == "" || m == "" || token == "" {

		return errors.ErrInvalidRequest
	}

	if m != "" && token == os.Getenv("VERIFY_TOKEN") {
		return c.String(http.StatusOK, ch)
	}

	return c.String(http.StatusNotFound, "Error, wrong validation token")
}

// MessageEndpointHandler ...
func MessageEndpointHandler(c echo.Context) error {
	var callback models.Callback

	err := c.Bind(&callback)
	if err != nil {
		return err
	}

	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				if !reflect.DeepEqual(event.Message, models.Message{}) && event.Message.Text != "" {
					ProcessMessage(event)
				}
			}
		}
		c.String(http.StatusOK, "Got your message")
	}

	return c.String(http.StatusNotFound, "Message not supported")
}

// ProcessMessage ...
func ProcessMessage(event models.Messaging) {
	client := &http.Client{}
	response := models.Response{
		Recipient: models.User{
			ID: event.Sender.ID,
		},
		Message: models.Message{
			Attachment: &models.Attachment{
				Type: "image",
				Payload: models.Payload{
					URL: IMAGE,
				},
			},
		},
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&response)

	url := fmt.Sprintf(FACEBOOK_API, os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
