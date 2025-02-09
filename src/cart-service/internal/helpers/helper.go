package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func IdAuthorization(bearerToken string) (string, error) {
	if bearerToken == "" {
		return "", fmt.Errorf("header is missing")
	}
	token := strings.TrimPrefix("Bearer ", bearerToken)
	if token == "" {
		return "", fmt.Errorf("token is missing")
	}
	resp, err := http.Post("http://user-auth-service:8081/validate-token/id-taking", "application/json", strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, token)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to decode response")
	}
	id, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("error with converting id")
	}
	if id == "" {
		return "", fmt.Errorf("id is missing")
	}
	return id, nil

}
