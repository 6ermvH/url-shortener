//go:build e2e

package e2e_test

import (
	"bytes"
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

func TestShorten_Success(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "https://example.com"})
	resp, err := http.Post(baseURL()+"/", "application/json", bytes.NewReader(body)) //nolint:noctx

	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.Len(t, result["shortUrl"], 10)
}

func TestShorten_Idempotent(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "https://idempotent-test.com"})

	resp1, err := http.Post(baseURL()+"/", "application/json", bytes.NewReader(body)) //nolint:noctx
	require.NoError(t, err)

	defer resp1.Body.Close()

	var result1 map[string]string
	require.NoError(t, json.NewDecoder(resp1.Body).Decode(&result1))

	body, _ = json.Marshal(map[string]string{"url": "https://idempotent-test.com"})

	resp2, err := http.Post(baseURL()+"/", "application/json", bytes.NewReader(body)) //nolint:noctx
	require.NoError(t, err)

	defer resp2.Body.Close()

	var result2 map[string]string
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&result2))

	require.Equal(t, result1["shortUrl"], result2["shortUrl"])
}

func TestResolve_Success(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "https://resolve-test.com"})

	shortenResp, err := http.Post(
		baseURL()+"/",
		"application/json",
		bytes.NewReader(body),
	) //nolint:noctx
	require.NoError(t, err)

	defer shortenResp.Body.Close()

	var shortened map[string]string
	require.NoError(t, json.NewDecoder(shortenResp.Body).Decode(&shortened))

	resolveResp, err := http.Get(baseURL() + "/" + shortened["shortUrl"]) //nolint:noctx
	require.NoError(t, err)

	defer resolveResp.Body.Close()

	require.Equal(t, http.StatusOK, resolveResp.StatusCode)

	var resolved map[string]string
	require.NoError(t, json.NewDecoder(resolveResp.Body).Decode(&resolved))
	require.Equal(t, "https://resolve-test.com", resolved["originalUrl"])
}

func TestResolve_NotFound(t *testing.T) {
	t.Parallel()

	resp, err := http.Get(baseURL() + "/notexists1") //nolint:noctx
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestShorten_InvalidURL(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "not-a-url"})

	resp, err := http.Post(baseURL()+"/", "application/json", bytes.NewReader(body)) //nolint:noctx
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShorten_EmptyURL(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": ""})

	resp, err := http.Post(baseURL()+"/", "application/json", bytes.NewReader(body)) //nolint:noctx
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShorten_InvalidScheme(t *testing.T) {
	t.Parallel()

	body, _ := json.Marshal(map[string]string{"url": "ht://example.com"})

	resp, err := http.Post(baseURL()+"/", "application/json", bytes.NewReader(body)) //nolint:noctx
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
