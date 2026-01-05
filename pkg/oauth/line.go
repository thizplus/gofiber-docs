package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gofiber-template/domain/dto"
)

var LineOAuthConfig *LineConfig

// LineConfig holds LINE OAuth configuration
type LineConfig struct {
	ChannelID     string
	ChannelSecret string
	RedirectURL   string
}

// InitLineOAuth initializes the LINE OAuth configuration
func InitLineOAuth() {
	LineOAuthConfig = &LineConfig{
		ChannelID:     os.Getenv("LINE_CHANNEL_ID"),
		ChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		RedirectURL:   os.Getenv("LINE_REDIRECT_URL"),
	}
}

// AuthCodeURL generates the LINE authorization URL
func (c *LineConfig) AuthCodeURL(state string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", c.ChannelID)
	params.Set("redirect_uri", c.RedirectURL)
	params.Set("state", state)
	params.Set("scope", "profile openid email")

	return "https://access.line.me/oauth2/v2.1/authorize?" + params.Encode()
}

// ExchangeCodeForToken exchanges the authorization code for tokens
func (c *LineConfig) ExchangeCodeForToken(code string) (*dto.LineTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURL)
	data.Set("client_id", c.ChannelID)
	data.Set("client_secret", c.ChannelSecret)

	resp, err := http.Post(
		"https://api.line.me/oauth2/v2.1/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LINE token exchange failed: %s", string(body))
	}

	var tokenResp dto.LineTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserProfile fetches user profile from LINE API
func (c *LineConfig) GetUserProfile(accessToken string) (*dto.LineProfile, error) {
	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LINE profile request failed: %s", string(body))
	}

	var profile dto.LineProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("failed to decode profile: %w", err)
	}

	return &profile, nil
}

// VerifyIDToken verifies and decodes LINE ID token
func (c *LineConfig) VerifyIDToken(idToken string) (*dto.LineIDTokenPayload, error) {
	data := url.Values{}
	data.Set("id_token", idToken)
	data.Set("client_id", c.ChannelID)

	resp, err := http.Post(
		"https://api.line.me/oauth2/v2.1/verify",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LINE ID token verification failed: %s", string(body))
	}

	var payload dto.LineIDTokenPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode ID token payload: %w", err)
	}

	// Validate issuer
	if payload.Iss != "https://access.line.me" {
		return nil, errors.New("invalid ID token issuer")
	}

	// Validate audience
	if payload.Aud != c.ChannelID {
		return nil, errors.New("invalid ID token audience")
	}

	return &payload, nil
}

// GetUserInfo gets complete user info from LINE
func (c *LineConfig) GetUserInfo(tokenResp *dto.LineTokenResponse) (*dto.LineUserInfo, error) {
	// Get profile from API
	profile, err := c.GetUserProfile(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	userInfo := &dto.LineUserInfo{
		ID:          profile.UserID,
		DisplayName: profile.DisplayName,
		PictureURL:  profile.PictureURL,
	}

	// Try to get email from ID token
	if tokenResp.IDToken != "" {
		payload, err := c.VerifyIDToken(tokenResp.IDToken)
		if err == nil && payload.Email != "" {
			userInfo.Email = payload.Email
		}
	}

	return userInfo, nil
}
