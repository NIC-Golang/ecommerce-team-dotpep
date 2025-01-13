package helpers

import (
	"fmt"
	"net/http"
	"strings"

	"encoding/json"
)

func GetIdFromToken(token string) (string, error) {
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
		return "", fmt.Errorf("failed to decode response")
	}
	return id, nil
}

func HeaderTrimming(header string) (token string, msg error) {
	if header == "" {
		return "", fmt.Errorf("authorization header missing")
	}
	tokenTrim := strings.TrimPrefix(header, "Bearer ")

	if tokenTrim == "" {
		return "", fmt.Errorf("token missing")
	}
	return tokenTrim, nil
}
