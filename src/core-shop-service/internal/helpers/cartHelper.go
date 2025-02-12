package helpers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

func HeaderTrimming(header string) (token string, msg error) {
	tokenTrim := strings.TrimPrefix(header, "Bearer ")

	if tokenTrim == "" {
		return "", fmt.Errorf("token missing")
	}
	return tokenTrim, nil
}

func SendWithHeaders(header string, orderJSON []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", "http://cart-service:8083/cart/orders", bytes.NewReader(orderJSON))
	if err != nil {

		return nil, fmt.Errorf("error creating request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", header)

	client := &http.Client{}
	resp1, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to cart-service")
	}
	return resp1, nil
}
