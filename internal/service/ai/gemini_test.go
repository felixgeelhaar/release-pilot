// Package ai provides AI-powered content generation for ReleasePilot.
package ai

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewGeminiService_NoAPIKey(t *testing.T) {
	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "", // No API key
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("NewGeminiService() unexpected error: %v", err)
	}

	// Should return noop service
	if svc.IsAvailable() {
		t.Error("Service should not be available without API key")
	}
}

func TestNewGeminiService_InvalidAPIKey(t *testing.T) {
	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "invalid-key",
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	_, err := NewGeminiService(cfg)
	if err == nil {
		t.Error("NewGeminiService() should return error for invalid API key")
	}
	if !strings.Contains(err.Error(), "invalid Gemini API key") {
		t.Errorf("Error should mention invalid API key format, got: %v", err)
	}
}

func TestNewGeminiService_ValidConfig(t *testing.T) {
	// Skip this test if no valid API key is available
	// This is an integration test that would require a real API key
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI", // Example format
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("NewGeminiService() error = %v", err)
	}

	if !svc.IsAvailable() {
		t.Error("Service should be available with valid API key format")
	}
}

func TestNewGeminiService_DefaultModel(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
		Model:         "", // No model specified
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("NewGeminiService() error = %v", err)
	}

	geminiSvc, ok := svc.(*geminiService)
	if !ok {
		t.Fatal("Service is not a geminiService")
	}

	if geminiSvc.config.Model != DefaultGeminiModel {
		t.Errorf("Default model = %v, want %v", geminiSvc.config.Model, DefaultGeminiModel)
	}
}

func TestNewGeminiService_CustomModel(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
		Model:         "gemini-1.5-pro",
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("NewGeminiService() error = %v", err)
	}

	geminiSvc, ok := svc.(*geminiService)
	if !ok {
		t.Fatal("Service is not a geminiService")
	}

	if geminiSvc.config.Model != "gemini-1.5-pro" {
		t.Errorf("Custom model = %v, want gemini-1.5-pro", geminiSvc.config.Model)
	}
}

func TestNewGeminiService_CustomPrompts(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		CustomPrompts: CustomPrompts{
			ChangelogSystem: "Custom system prompt",
			ChangelogUser:   "Custom user prompt",
		},
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("NewGeminiService() error = %v", err)
	}

	geminiSvc, ok := svc.(*geminiService)
	if !ok {
		t.Fatal("Service is not a geminiService")
	}

	if geminiSvc.prompts.changelogSystem != "Custom system prompt" {
		t.Errorf("Changelog system prompt = %v, want Custom system prompt", geminiSvc.prompts.changelogSystem)
	}
	if geminiSvc.prompts.changelogUser != "Custom user prompt" {
		t.Errorf("Changelog user prompt = %v, want Custom user prompt", geminiSvc.prompts.changelogUser)
	}
}

func TestGeminiService_GenerateChangelog_EmptyChanges(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	result, err := svc.GenerateChangelog(context.Background(), nil, DefaultGenerateOptions())
	if err != nil {
		t.Errorf("GenerateChangelog() error = %v", err)
	}
	if result != "" {
		t.Errorf("GenerateChangelog() = %v, want empty string", result)
	}
}

func TestGeminiService_GenerateReleaseNotes_EmptyChangelog(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	result, err := svc.GenerateReleaseNotes(context.Background(), "", DefaultGenerateOptions())
	if err != nil {
		t.Errorf("GenerateReleaseNotes() error = %v", err)
	}
	if result != "" {
		t.Errorf("GenerateReleaseNotes() = %v, want empty string", result)
	}
}

func TestGeminiService_GenerateMarketingBlurb_EmptyNotes(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	result, err := svc.GenerateMarketingBlurb(context.Background(), "", DefaultGenerateOptions())
	if err != nil {
		t.Errorf("GenerateMarketingBlurb() error = %v", err)
	}
	if result != "" {
		t.Errorf("GenerateMarketingBlurb() = %v, want empty string", result)
	}
}

func TestGeminiService_SummarizeChanges_EmptyChanges(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	result, err := svc.SummarizeChanges(context.Background(), nil, DefaultGenerateOptions())
	if err != nil {
		t.Errorf("SummarizeChanges() error = %v", err)
	}
	if result != "" {
		t.Errorf("SummarizeChanges() = %v, want empty string", result)
	}
}

func TestNewService_GeminiProvider(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	svc, err := NewService(
		WithProvider("gemini"),
		WithAPIKey("AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI"),
		WithModel("gemini-1.5-pro"),
		WithTimeout(30*time.Second),
	)
	if err != nil {
		t.Fatalf("NewService(gemini) error = %v", err)
	}

	if svc == nil {
		t.Error("NewService(gemini) returned nil")
	}

	// Verify it's a geminiService
	_, ok := svc.(*geminiService)
	if !ok {
		t.Error("NewService(gemini) did not return a geminiService")
	}
}

func TestGeminiDefaultConstants(t *testing.T) {
	if DefaultGeminiModel != "gemini-2.0-flash-exp" {
		t.Errorf("DefaultGeminiModel = %v, want gemini-2.0-flash-exp", DefaultGeminiModel)
	}
}

func TestGeminiKeyPattern(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		valid bool
	}{
		{"valid key", "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI", true},
		{"valid key with underscore", "AIzaSyDdI0hCZtE6vySjMm_WEfRq3CPzqKqqsHI", true},
		{"valid key with hyphen", "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI", true},
		{"valid key alphanumeric", "AIzaSyDdI0hCZtE6vySjMmWEfRq3CPzqKqqsHI1", true},
		{"invalid prefix AIzb", "AIzbSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI", false},
		{"invalid prefix Aiza", "AizaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI", false},
		{"missing prefix", "SyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI", false},
		{"too short", "AIza-short", false},
		{"empty", "", false},
		{"openai format", "sk-1234567890abcdef1234567890abcdef", false},
		{"anthropic format", "sk-ant-api03-validkeyformat12345678901234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := geminiKeyPattern.MatchString(tt.key)
			if result != tt.valid {
				t.Errorf("geminiKeyPattern.MatchString(%q) = %v, want %v", tt.key, result, tt.valid)
			}
		})
	}
}

func TestGeminiService_IsAvailable(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		wantAvail bool
		wantErr   bool
	}{
		{
			name:      "with valid API key format",
			apiKey:    "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
			wantAvail: false, // Will be false because client creation will fail without real API key
			wantErr:   false, // But no error during construction with valid format
		},
		{
			name:      "without API key",
			apiKey:    "",
			wantAvail: false,
			wantErr:   false,
		},
		{
			name:      "with invalid API key format",
			apiKey:    "invalid-key",
			wantAvail: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip tests that would make real API calls
			if tt.apiKey != "" && !tt.wantErr {
				t.Skip("Skipping test that would require valid Gemini API key")
			}

			cfg := ServiceConfig{
				Provider:      "gemini",
				APIKey:        tt.apiKey,
				MaxTokens:     2048,
				Temperature:   0.7,
				Timeout:       30 * time.Second,
				RetryAttempts: 3,
			}

			svc, err := NewGeminiService(cfg)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewGeminiService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if got := svc.IsAvailable(); got != tt.wantAvail {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.wantAvail)
			}
		})
	}
}

func TestGeminiService_ResilienceConfig(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	cfg := ServiceConfig{
		Provider:      "gemini",
		APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
		MaxTokens:     2048,
		Temperature:   0.7,
		Timeout:       60 * time.Second,
		RetryAttempts: 5,
		RateLimitRPM:  100,
	}

	svc, err := NewGeminiService(cfg)
	if err != nil {
		t.Fatalf("NewGeminiService() error = %v", err)
	}

	geminiSvc, ok := svc.(*geminiService)
	if !ok {
		t.Fatal("Service is not a geminiService")
	}

	if geminiSvc.resilience == nil {
		t.Error("Resilience should be configured")
	}

	// Verify resilience config was applied
	if geminiSvc.config.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want 60s", geminiSvc.config.Timeout)
	}
	if geminiSvc.config.RetryAttempts != 5 {
		t.Errorf("RetryAttempts = %v, want 5", geminiSvc.config.RetryAttempts)
	}
	if geminiSvc.config.RateLimitRPM != 100 {
		t.Errorf("RateLimitRPM = %v, want 100", geminiSvc.config.RateLimitRPM)
	}
}

func TestGeminiService_TemperatureConversion(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	tests := []struct {
		name        string
		temperature float64
	}{
		{"zero", 0.0},
		{"low", 0.3},
		{"medium", 0.7},
		{"high", 1.0},
		{"very high", 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ServiceConfig{
				Provider:      "gemini",
				APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
				MaxTokens:     2048,
				Temperature:   tt.temperature,
				Timeout:       30 * time.Second,
				RetryAttempts: 3,
			}

			svc, err := NewGeminiService(cfg)
			if err != nil {
				t.Fatalf("NewGeminiService() error = %v", err)
			}

			geminiSvc, ok := svc.(*geminiService)
			if !ok {
				t.Fatal("Service is not a geminiService")
			}

			if geminiSvc.config.Temperature != tt.temperature {
				t.Errorf("Temperature = %v, want %v", geminiSvc.config.Temperature, tt.temperature)
			}
		})
	}
}

func TestGeminiService_MaxTokensConfig(t *testing.T) {
	// Skip this test if no valid API key is available
	t.Skip("Skipping test that requires valid Gemini API key")

	tests := []struct {
		name      string
		maxTokens int
	}{
		{"small", 512},
		{"medium", 2048},
		{"large", 4096},
		{"very large", 8192},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := ServiceConfig{
				Provider:      "gemini",
				APIKey:        "AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI",
				MaxTokens:     tt.maxTokens,
				Temperature:   0.7,
				Timeout:       30 * time.Second,
				RetryAttempts: 3,
			}

			svc, err := NewGeminiService(cfg)
			if err != nil {
				t.Fatalf("NewGeminiService() error = %v", err)
			}

			geminiSvc, ok := svc.(*geminiService)
			if !ok {
				t.Fatal("Service is not a geminiService")
			}

			if geminiSvc.config.MaxTokens != tt.maxTokens {
				t.Errorf("MaxTokens = %v, want %v", geminiSvc.config.MaxTokens, tt.maxTokens)
			}
		})
	}
}
