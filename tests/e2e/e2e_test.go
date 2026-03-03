//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func baseURL() string {
	if url := os.Getenv("BASE_URL"); url != "" {
		return url
	}

	return "http://localhost:8081"
}

func post(t *testing.T, url string, body []byte) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func get(t *testing.T, url string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	require.NoError(t, err)

	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func TestShorten_Success(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "https://example.com"})
	resp := post(t, baseURL()+"/", body)

	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.Len(t, result["shortUrl"], 10)
}

func TestShorten_Idempotent(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "https://idempotent-test.com"})

	resp1 := post(t, baseURL()+"/", body)
	defer resp1.Body.Close()

	var result1 map[string]string
	require.NoError(t, json.NewDecoder(resp1.Body).Decode(&result1))

	resp2 := post(t, baseURL()+"/", body)
	defer resp2.Body.Close()

	var result2 map[string]string
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&result2))

	require.Equal(t, result1["shortUrl"], result2["shortUrl"])
}

func TestResolve_Success(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "https://resolve-test.com"})

	shortenResp := post(t, baseURL()+"/", body)
	defer shortenResp.Body.Close()

	var shortened map[string]string
	require.NoError(t, json.NewDecoder(shortenResp.Body).Decode(&shortened))

	resolveResp := get(t, baseURL()+"/"+shortened["shortUrl"])
	defer resolveResp.Body.Close()

	require.Equal(t, http.StatusOK, resolveResp.StatusCode)

	var resolved map[string]string
	require.NoError(t, json.NewDecoder(resolveResp.Body).Decode(&resolved))
	require.Equal(t, "https://resolve-test.com", resolved["originalUrl"])
}

func TestResolve_NotFound(t *testing.T) {
	t.Parallel()

	resp := get(t, baseURL()+"/notexists1")
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestShorten_InvalidURL(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "not-a-url"})

	resp := post(t, baseURL()+"/", body)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShorten_EmptyURL(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": ""})

	resp := post(t, baseURL()+"/", body)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShorten_InvalidScheme(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "ht://example.com"})

	resp := post(t, baseURL()+"/", body)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
