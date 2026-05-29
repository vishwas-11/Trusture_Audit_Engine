package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	pdf "github.com/ledongthuc/pdf"
)

// NoRelevantInformationMessage is shown when OCR/text extraction yields nothing usable.
const NoRelevantInformationMessage = "No relevant information could be extracted from this document."

const maxStoredExtractRunes = 100_000

// DocumentExtractService pulls text from PDFs locally and from images/PDFs via Gemini or OpenAI when keys exist.
type DocumentExtractService struct {
	GeminiAPIKey string
	OpenAIAPIKey string
}

func NewDocumentExtractService(geminiKey, openAIKey string) *DocumentExtractService {
	return &DocumentExtractService{GeminiAPIKey: strings.TrimSpace(geminiKey), OpenAIAPIKey: strings.TrimSpace(openAIKey)}
}

// ExtractFromFile writes extracted plain text and an optional user-facing note when text is empty.
func (s *DocumentExtractService) ExtractFromFile(path string, originalFilename string) (extracted string, note string) {
	ext := strings.ToLower(filepath.Ext(originalFilename))

	data, err := os.ReadFile(path)
	if err != nil {
		return "", "Could not read uploaded file for text extraction."
	}
	if len(data) > 15<<20 {
		data = data[:15<<20]
	}

	switch ext {
	case ".pdf":
		text := extractPDFPlain(path)
		text = strings.TrimSpace(text)
		if text != "" {
			return truncateExtract(text), ""
		}
		if s.GeminiAPIKey != "" {
			t, ok := s.geminiExtract(data, "application/pdf")
			if ok {
				t = strings.TrimSpace(t)
				if isNoRelevantMarker(t) {
					return "", NoRelevantInformationMessage
				}
				if t != "" {
					return truncateExtract(t), ""
				}
			}
		}
		return "", NoRelevantInformationMessage

	case ".jpg", ".jpeg", ".png", ".webp":
		mt := mimeFromImageExt(ext)
		if s.GeminiAPIKey != "" {
			t, ok := s.geminiExtract(data, mt)
			if ok {
				t = strings.TrimSpace(t)
				if isNoRelevantMarker(t) {
					return "", NoRelevantInformationMessage
				}
				if t != "" {
					return truncateExtract(t), ""
				}
			}
		}
		if s.OpenAIAPIKey != "" {
			t, ok := s.openAIVisionExtract(data, mt)
			if ok {
				t = strings.TrimSpace(t)
				if isNoRelevantMarker(t) {
					return "", NoRelevantInformationMessage
				}
				if t != "" {
					return truncateExtract(t), ""
				}
			}
		}
		n := NoRelevantInformationMessage
		if s.GeminiAPIKey == "" && s.OpenAIAPIKey == "" {
			n += " Set GEMINI_API_KEY or OPENAI_API_KEY on the server to enable OCR for images."
		}
		return "", n

	default:
		if len(data) >= 4 && string(data[:4]) == "%PDF" {
			text := strings.TrimSpace(extractPDFPlain(path))
			if text != "" {
				return truncateExtract(text), ""
			}
		}
		return "", NoRelevantInformationMessage
	}
}

func extractPDFPlain(path string) string {
	f, r, err := pdf.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	rd, err := r.GetPlainText()
	if err != nil {
		return ""
	}
	b, err := io.ReadAll(rd)
	if err != nil {
		return ""
	}
	return string(b)
}

func mimeFromImageExt(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func isNoRelevantMarker(s string) bool {
	return strings.TrimSpace(strings.ToUpper(s)) == "NO_RELEVANT_TEXT"
}

func truncateExtract(s string) string {
	if utf8.RuneCountInString(s) <= maxStoredExtractRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxStoredExtractRunes])
}

func (s *DocumentExtractService) geminiExtract(data []byte, mimeType string) (string, bool) {
	payload := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{
						"inline_data": map[string]string{
							"mime_type": mimeType,
							"data":      base64.StdEncoding.EncodeToString(data),
						},
					},
					{
						"text": "Extract all readable text from this document. Output plain text only. " +
							"If there is no meaningful readable content, respond with exactly: NO_RELEVANT_TEXT",
					},
				},
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", false
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s",
		s.GeminiAPIKey,
	)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", false
	}

	var root map[string]any
	if err := json.Unmarshal(respBody, &root); err != nil {
		return "", false
	}
	candidates, _ := root["candidates"].([]any)
	if len(candidates) == 0 {
		return "", true
	}
	first, _ := candidates[0].(map[string]any)
	content, _ := first["content"].(map[string]any)
	parts, _ := content["parts"].([]any)
	if len(parts) == 0 {
		return "", true
	}
	part0, _ := parts[0].(map[string]any)
	text, _ := part0["text"].(string)
	return text, true
}

func (s *DocumentExtractService) openAIVisionExtract(data []byte, mimeType string) (string, bool) {
	b64 := base64.StdEncoding.EncodeToString(data)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, b64)

	payload := map[string]any{
		"model": "gpt-4o-mini",
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{
						"type": "text",
						"text": "Extract all readable text from this image. Output plain text only. " +
							"If there is no meaningful readable content, respond with exactly: NO_RELEVANT_TEXT",
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": dataURL,
						},
					},
				},
			},
		},
		"max_tokens": 4096,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", false
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.OpenAIAPIKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", false
	}

	var root map[string]any
	if err := json.Unmarshal(respBody, &root); err != nil {
		return "", false
	}
	choices, _ := root["choices"].([]any)
	if len(choices) == 0 {
		return "", true
	}
	ch0, _ := choices[0].(map[string]any)
	msg, _ := ch0["message"].(map[string]any)
	text, _ := msg["content"].(string)
	return text, true
}
