package complitions

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// https://platform.openai.com/docs/api-reference/completions
func CreateCompletion(c echo.Context) error {
	openAPIKey := viper.GetString("openAI.apiKey")

	var completionReq CompletionLegacyRequest
	err := c.Bind(&completionReq)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(completionReq)
	if err != nil {
		return err
	}

	url := "https://api.openai.com/v1/completions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+openAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var completionResp CompletionLegacyResponse
	err = json.Unmarshal(body, &completionResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, completionResp)
}
