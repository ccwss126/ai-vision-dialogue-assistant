package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const defaultPort = "8080"
const defaultModel = "gemini-2.5-flash"

type analyzeRequest struct {
	ImageBase64 string `json:"imageBase64"`
	Question    string `json:"question"`
}

type analyzeResponse struct {
	Success bool   `json:"success"`
	Mode    string `json:"mode"`
	Answer  string `json:"answer"`
}

type inlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type geminiPart struct {
	Text       string      `json:"text,omitempty"`
	InlineData *inlineData `json:"inline_data,omitempty"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func main() {
	loadEnvFiles(".env", "../.env", "backend/.env")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	client := newHTTPClient()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/analyze", withCORS(func(w http.ResponseWriter, r *http.Request) {
		handleAnalyze(w, r, client)
	}))

	log.Printf("Go backend running at http://localhost:%s", port)
	log.Printf("Gemini API Key configured: %s", yesNo(os.Getenv("GEMINI_API_KEY") != ""))
	log.Printf("Proxy enabled: %s", yesNo(proxyURL() != ""))

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func handleAnalyze(w http.ResponseWriter, r *http.Request, client *http.Client) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, analyzeResponse{
			Success: false,
			Mode:    "demo",
			Answer:  "仅支持 POST /api/analyze。",
		})
		return
	}

	var req analyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, analyzeResponse{
			Success: false,
			Mode:    "demo",
			Answer:  "请求格式错误。",
		})
		return
	}

	req.Question = strings.TrimSpace(req.Question)
	if strings.TrimSpace(req.ImageBase64) == "" || req.Question == "" {
		writeJSON(w, http.StatusBadRequest, analyzeResponse{
			Success: false,
			Mode:    "demo",
			Answer:  "请提供画面截图和语音问题。",
		})
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		writeJSON(w, http.StatusOK, demoAnswer())
		return
	}

	answer, err := callGemini(r.Context(), client, apiKey, req)
	if err != nil {
		log.Printf("Gemini request failed, fallback to demo mode: %v", err)
		writeJSON(w, http.StatusOK, demoAnswer())
		return
	}

	writeJSON(w, http.StatusOK, analyzeResponse{
		Success: true,
		Mode:    "gemini",
		Answer:  answer,
	})
}

func callGemini(ctx context.Context, client *http.Client, apiKey string, req analyzeRequest) (string, error) {
	mimeType, imageData := parseImage(req.ImageBase64)
	prompt := strings.Join([]string{
		"你是 AI 视觉对话助手。",
		"你正在通过摄像头观察用户当前画面，并通过麦克风听到用户说的话。",
		"请结合图片内容和用户语音问题回答，回答要自然、具体、简洁。",
		"如果图片信息不足或无法确定，请明确说明不确定点。",
		"用户说的话：" + req.Question,
	}, "\n")

	body := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{InlineData: &inlineData{MimeType: mimeType, Data: imageData}},
					{Text: prompt},
				},
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = defaultModel
	}

	endpoint := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent",
		url.PathEscape(model),
	)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", apiKey)

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("gemini status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed geminiResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", err
	}
	for _, candidate := range parsed.Candidates {
		for _, part := range candidate.Content.Parts {
			if strings.TrimSpace(part.Text) != "" {
				return strings.TrimSpace(part.Text), nil
			}
		}
	}
	return "", errors.New("empty gemini response")
}

func parseImage(imageBase64 string) (string, string) {
	if strings.HasPrefix(imageBase64, "data:") {
		parts := strings.SplitN(imageBase64, ",", 2)
		if len(parts) == 2 {
			mime := strings.TrimPrefix(parts[0], "data:")
			mime = strings.TrimSuffix(mime, ";base64")
			return mime, parts[1]
		}
	}
	return "image/jpeg", imageBase64
}

func demoAnswer() analyzeResponse {
	return analyzeResponse{
		Success: true,
		Mode:    "demo",
		Answer:  "这是 Demo Mode 回复：我会基于当前摄像头画面和你的语音问题进行回答。配置 GEMINI_API_KEY 后即可启用真实 Gemini Vision 分析。",
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, data analyzeResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func newHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	return &http.Client{
		Timeout:   35 * time.Second,
		Transport: transport,
	}
}

func proxyURL() string {
	if v := os.Getenv("HTTPS_PROXY"); v != "" {
		return v
	}
	return os.Getenv("HTTP_PROXY")
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}

func loadEnvFiles(paths ...string) {
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.Trim(strings.TrimSpace(value), `"'`)
			if key != "" && os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}
	}
}
