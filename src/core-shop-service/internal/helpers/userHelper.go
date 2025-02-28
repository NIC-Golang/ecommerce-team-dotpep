package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetName(email string) (string, error) {
	resp, err := http.Post("http://user-auth-service:8081/name-taking", "application/json", strings.NewReader(fmt.Sprintf(`{"email":"%s"}`, email)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {

		return "", fmt.Errorf("failed to decode response")
	}
	name, ok := result["name"]
	if !ok {
		return "", fmt.Errorf("failed to decode response")
	}
	return name, nil
}
func GetIdAndEmailFromToken(token string) (string, string, error) {
	resp, err := http.Post("http://user-auth-service:8081/validate-token/id-taking", "application/json", strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, token)))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {

		return "", "", fmt.Errorf("failed to decode response")
	}
	id, ok1 := result["id"].(string)
	email, ok2 := result["email"].(string)
	if !ok1 || !ok2 {
		return "", "", fmt.Errorf("failed to decode response")
	}
	return id, email, nil
}
