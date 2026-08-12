package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func ExecuteCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Language string `json:"language"`
		Code     string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pistonURL := os.Getenv("PISTON_URL")
	if pistonURL == "" {
		pistonURL = "https://emkc.org/api/v2/piston/execute"
	}

	lang := req.Language
	version := "*"

	if lang == "javascript" || lang == "js" {
		lang = "javascript"
		version = "18.15.0"
	} else if lang == "python" {
		lang = "python"
		version = "3.10.0"
	} else if lang == "go" {
		lang = "go"
		version = "1.16.2"
	}

	pistonReq := map[string]interface{}{
		"language": lang,
		"version":  version,
		"files": []map[string]string{
			{
				"content": req.Code,
			},
		},
	}

	reqBytes, err := json.Marshal(pistonReq)
	if err != nil {
		http.Error(w, "Failed to encode execution request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(
		pistonURL,
		"application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		log.Printf("Error executing code via Piston: %v", err)
		http.Error(w, "Execution service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read execution response", http.StatusInternalServerError)
		return
	}

	var pistonResp struct {
		Message string `json:"message"`
		Run     struct {
			Stdout string `json:"stdout"`
			Stderr string `json:"stderr"`
			Code   int    `json:"code"`
		} `json:"run"`
	}

	if err := json.Unmarshal(bodyBytes, &pistonResp); err != nil {
		log.Printf("Failed to decode Piston response: %s", string(bodyBytes))
		http.Error(w, "Invalid response from execution service", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		errorOutput := fmt.Sprintf("Execution failed: %s", pistonResp.Message)
		http.Error(w, errorOutput, http.StatusBadRequest)
		return
	}

	if pistonResp.Run.Code != 0 || pistonResp.Run.Stderr != "" {
		errorOutput := fmt.Sprintf("Execution failed:\n%s", pistonResp.Run.Stderr)
		http.Error(w, errorOutput, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"output": pistonResp.Run.Stdout,
	})
}
