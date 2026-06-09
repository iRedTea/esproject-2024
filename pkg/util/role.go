package util

import (
	"encoding/json"
	"esproject"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func MatchRole(c *gin.Context, role string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://auth.easystartup.su/api/user/self"), nil)
	if err != nil {
		return false
	}

	req.Header.Add("Authorization", c.GetHeader("Authorization"))
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var result esproject.AuthorizedUser
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false
	}

	return result.Access == role
}
