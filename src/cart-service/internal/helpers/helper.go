package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Order struct {
	Id string `json:"id"`
}

func IdAuthorization(bearerToken string) (string, error) {
	if bearerToken == "" {
		return "", fmt.Errorf("header is missing")
	}
	token := strings.TrimPrefix(bearerToken, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("token is missing")
	}

	resp, err := http.Post("http://user-auth-service:8081/validate-token/id-taking", "application/json", strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, token)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var result Order
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	if result.Id == "" {
		return "", fmt.Errorf("id is missing! cannot create a cart")
	}

	return result.Id, nil
}
