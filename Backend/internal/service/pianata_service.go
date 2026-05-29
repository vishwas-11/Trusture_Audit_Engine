package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type PinataService struct {
	client *resty.Client
}

func NewPinataService() *PinataService {
	client := resty.New()
	client.SetBaseURL("https://api.pinata.cloud")
	client.SetTimeout(20 * time.Second)

	return &PinataService{client: client}
}

func (p *PinataService) UploadFile(filePath string) (string, error) {
	apiKey := strings.TrimSpace(os.Getenv("PINATA_API_KEY"))
	secret := strings.TrimSpace(os.Getenv("PINATA_SECRET_API_KEY"))
	if secret == "" {
		// Backward compatibility with existing .env naming used in this project.
		secret = strings.TrimSpace(os.Getenv("PINATA_API_SECRET_KEY"))
	}
	if apiKey == "" || secret == "" {
		return "", errors.New("pinata credentials missing: set PINATA_API_KEY and PINATA_SECRET_API_KEY")
	}

	resp, err := p.client.R().
		SetHeader("pinata_api_key", apiKey).
		SetHeader("pinata_secret_api_key", secret).
		SetFile("file", filePath).
		Post("/pinning/pinFileToIPFS")

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("pinata upload failed: status=%d body=%s", resp.StatusCode(), strings.TrimSpace(resp.String()))
	}

	var result struct {
		IPFSHash string `json:"IpfsHash"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("pinata response parse failed: %w", err)
	}
	if strings.TrimSpace(result.IPFSHash) == "" {
		return "", fmt.Errorf("pinata upload failed: missing IpfsHash in response: %s", strings.TrimSpace(resp.String()))
	}

	return result.IPFSHash, nil
}
