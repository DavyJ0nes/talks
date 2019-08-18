package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	msg, err := getMessage()
	if err != nil {
		return err
	}

	c.Set("message", msg)
	return c.Render(200, r.HTML("index.html"))
}

func getMessage() (string, error) {
	client := http.DefaultClient
	apiAddr := envy.Get("API_ADDR", "http://localhost:8080")
	resp, err := client.Get(apiAddr + "/goodvibes")
	if err != nil {
		return "", err
	}

	msg := struct {
		Message string `json:"message"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&msg)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned %d response. Message: %s", resp.StatusCode, msg.Message)
	}

	return msg.Message, nil

}
