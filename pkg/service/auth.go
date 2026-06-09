package service

import (
	"bytes"
	"encoding/json"
	"esproject"
	"fmt"
	"io"
	"net/http"
	"os"
)

const AUTH_URL = "https://auth.easystartup.su/api"

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Authorize(header string) (*esproject.AuthorizedUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", AUTH_URL+"/user/self", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", header)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authozire: %d %s\nheader: %s", resp.StatusCode, resp.Status, header)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result esproject.AuthorizedUser
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) AuthorizeById(id int64) (*esproject.AuthorizedUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(AUTH_URL+"/user/%d", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("ADMIN_TOKEN"))
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authozire: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result esproject.AuthorizedUser
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *AuthService) AuthorizeAndUpdatePicture(header string, picUrl string) (*esproject.AuthorizedUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", AUTH_URL+"/user/self", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", header)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authozire: %d %s\nheader: %s", resp.StatusCode, resp.Status, header)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result esproject.AuthorizedUser
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	result.PictureURL = picUrl
	requestBody, _ := json.Marshal(result)

	req, err = http.NewRequest("POST", AUTH_URL+"/user/self", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", header)
	req.Header.Add("Accept", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
