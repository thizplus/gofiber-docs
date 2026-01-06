package serviceimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/storage"
	"gofiber-template/pkg/oauth"
)

type UserServiceImpl struct {
	userRepo    repositories.UserRepository
	jwtSecret   string
	r2Storage   storage.R2Storage
	r2PublicURL string
}

func NewUserService(userRepo repositories.UserRepository, jwtSecret string, r2Storage storage.R2Storage, r2PublicURL string) services.UserService {
	return &UserServiceImpl{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		r2Storage:   r2Storage,
		r2PublicURL: r2PublicURL,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("อีเมลนี้ถูกใช้งานแล้ว")
	}

	// Auto-generate username from email if not provided
	username := req.Username
	if username == "" {
		username = s.generateUsernameFromEmail(ctx, req.Email)
	} else {
		// Check if provided username already exists
		existingUser, _ = s.userRepo.GetByUsername(ctx, username)
		if existingUser != nil {
			return nil, errors.New("ชื่อผู้ใช้นี้ถูกใช้งานแล้ว")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     username,
		Password:     string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "user",
		IsActive:     true,
		AuthProvider: "email",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req *dto.LoginRequest) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, errors.New("อีเมลหรือรหัสผ่านไม่ถูกต้อง")
	}

	if !user.IsActive {
		return "", nil, errors.New("บัญชีนี้ถูกระงับการใช้งาน")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("อีเมลหรือรหัสผ่านไม่ถูกต้อง")
	}

	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("ไม่พบผู้ใช้งาน")
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("ไม่พบผู้ใช้งาน")
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, userID, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *UserServiceImpl) ListUsers(ctx context.Context, offset, limit int) ([]*models.User, int64, error) {
	users, err := s.userRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *UserServiceImpl) GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserServiceImpl) ValidateJWT(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("รูปแบบ token ไม่ถูกต้อง")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("ข้อมูล token ไม่ถูกต้อง")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.New("รหัสผู้ใช้ใน token ไม่ถูกต้อง")
		}

		user, err := s.userRepo.GetByID(context.Background(), userID)
		if err != nil {
			return nil, errors.New("ไม่พบผู้ใช้งาน")
		}

		return user, nil
	}

	return nil, errors.New("token ไม่ถูกต้อง")
}

// ==================== OAuth Methods ====================

// GetGoogleAuthURL returns Google OAuth authorization URL
func (s *UserServiceImpl) GetGoogleAuthURL(state string) string {
	return oauth.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// HandleGoogleCallback handles Google OAuth callback
func (s *UserServiceImpl) HandleGoogleCallback(ctx context.Context, code string, stateData dto.OAuthStateData) (string, *models.User, bool, error) {
	// Exchange code for token
	token, err := oauth.GoogleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return "", nil, false, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from Google
	client := oauth.GoogleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", nil, false, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser dto.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", nil, false, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Find or create user
	user, isNewUser, err := s.findOrCreateGoogleUser(ctx, &googleUser)
	if err != nil {
		return "", nil, false, err
	}

	// Generate JWT
	jwtToken, err := s.GenerateJWT(user)
	if err != nil {
		return "", nil, false, fmt.Errorf("failed to generate token: %w", err)
	}

	return jwtToken, user, isNewUser, nil
}

// findOrCreateGoogleUser finds existing user or creates new one from Google OAuth
func (s *UserServiceImpl) findOrCreateGoogleUser(ctx context.Context, googleUser *dto.GoogleUserInfo) (*models.User, bool, error) {
	// Try to find by Google ID
	existingUser, _ := s.userRepo.GetByGoogleID(ctx, googleUser.ID)
	if existingUser != nil {
		return existingUser, false, nil
	}

	// Try to find by email
	existingUser, _ = s.userRepo.GetByEmail(ctx, googleUser.Email)
	if existingUser != nil {
		// Link Google account to existing user
		existingUser.GoogleID = &googleUser.ID
		existingUser.AuthProvider = "google"
		if existingUser.Avatar == "" {
			existingUser.Avatar = googleUser.Picture
		}
		if err := s.userRepo.Update(ctx, existingUser.ID, existingUser); err != nil {
			return nil, false, fmt.Errorf("failed to link Google account: %w", err)
		}
		return existingUser, false, nil
	}

	// Create new user
	username := s.generateUsernameFromEmail(ctx, googleUser.Email)
	newUser := &models.User{
		ID:           uuid.New(),
		Email:        googleUser.Email,
		Username:     username,
		Password:     "", // No password for OAuth users
		FirstName:    googleUser.GivenName,
		LastName:     googleUser.FamilyName,
		Avatar:       googleUser.Picture,
		Role:         "user",
		IsActive:     true,
		GoogleID:     &googleUser.ID,
		AuthProvider: "google",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, false, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, true, nil
}

// GetLineAuthURL returns LINE OAuth authorization URL
func (s *UserServiceImpl) GetLineAuthURL(state string) string {
	return oauth.LineOAuthConfig.AuthCodeURL(state)
}

// HandleLineCallback handles LINE OAuth callback
func (s *UserServiceImpl) HandleLineCallback(ctx context.Context, code string, stateData dto.OAuthStateData) (string, *models.User, bool, error) {
	// Exchange code for token
	tokenResp, err := oauth.LineOAuthConfig.ExchangeCodeForToken(code)
	if err != nil {
		return "", nil, false, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from LINE
	lineUser, err := oauth.LineOAuthConfig.GetUserInfo(tokenResp)
	if err != nil {
		return "", nil, false, fmt.Errorf("failed to get user info: %w", err)
	}

	// Find or create user
	user, isNewUser, err := s.findOrCreateLineUser(ctx, lineUser)
	if err != nil {
		return "", nil, false, err
	}

	// Generate JWT
	jwtToken, err := s.GenerateJWT(user)
	if err != nil {
		return "", nil, false, fmt.Errorf("failed to generate token: %w", err)
	}

	return jwtToken, user, isNewUser, nil
}

// findOrCreateLineUser finds existing user or creates new one from LINE OAuth
func (s *UserServiceImpl) findOrCreateLineUser(ctx context.Context, lineUser *dto.LineUserInfo) (*models.User, bool, error) {
	// Try to find by LINE ID
	existingUser, _ := s.userRepo.GetByLineID(ctx, lineUser.ID)
	if existingUser != nil {
		return existingUser, false, nil
	}

	// Try to find by email (if available)
	if lineUser.Email != "" {
		existingUser, _ = s.userRepo.GetByEmail(ctx, lineUser.Email)
		if existingUser != nil {
			// Link LINE account to existing user
			existingUser.LineID = &lineUser.ID
			existingUser.AuthProvider = "line"
			if existingUser.Avatar == "" {
				existingUser.Avatar = lineUser.PictureURL
			}
			if err := s.userRepo.Update(ctx, existingUser.ID, existingUser); err != nil {
				return nil, false, fmt.Errorf("failed to link LINE account: %w", err)
			}
			return existingUser, false, nil
		}
	}

	// Create new user
	var username string
	if lineUser.DisplayName != "" {
		username = s.generateUsernameFromDisplayName(ctx, lineUser.DisplayName)
	} else {
		username = "line_" + lineUser.ID[:8]
	}

	// Generate email placeholder if no email from LINE
	email := lineUser.Email
	if email == "" {
		email = fmt.Sprintf("line_%s@stou.placeholder", lineUser.ID[:8])
	}

	// Parse display name into first and last name
	firstName, lastName := s.parseDisplayName(lineUser.DisplayName)

	newUser := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		Password:     "", // No password for OAuth users
		FirstName:    firstName,
		LastName:     lastName,
		Avatar:       lineUser.PictureURL,
		Role:         "user",
		IsActive:     true,
		LineID:       &lineUser.ID,
		AuthProvider: "line",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, false, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, true, nil
}

// generateUsernameFromEmail generates a unique username from email
func (s *UserServiceImpl) generateUsernameFromEmail(ctx context.Context, email string) string {
	parts := strings.Split(email, "@")
	baseUsername := strings.ToLower(parts[0])

	// Remove special characters
	var cleanUsername strings.Builder
	for _, r := range baseUsername {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			cleanUsername.WriteRune(r)
		}
	}
	baseUsername = cleanUsername.String()

	if len(baseUsername) > 20 {
		baseUsername = baseUsername[:20]
	}

	if baseUsername == "" {
		baseUsername = "user"
	}

	// Check for uniqueness
	username := baseUsername
	counter := 1
	for {
		existing, _ := s.userRepo.GetByUsername(ctx, username)
		if existing == nil {
			break
		}
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
		if counter > 100 {
			username = fmt.Sprintf("%s_%d", baseUsername, time.Now().UnixNano()%10000)
			break
		}
	}

	return username
}

// generateUsernameFromDisplayName generates a unique username from display name
func (s *UserServiceImpl) generateUsernameFromDisplayName(ctx context.Context, displayName string) string {
	var cleanUsername strings.Builder
	for _, r := range strings.ToLower(displayName) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			cleanUsername.WriteRune(r)
		}
	}
	baseUsername := cleanUsername.String()

	if len(baseUsername) < 3 {
		baseUsername = "line_user"
	}

	if len(baseUsername) > 20 {
		baseUsername = baseUsername[:20]
	}

	// Check for uniqueness
	username := baseUsername
	counter := 1
	for {
		existing, _ := s.userRepo.GetByUsername(ctx, username)
		if existing == nil {
			break
		}
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++
		if counter > 100 {
			username = fmt.Sprintf("%s_%d", baseUsername, time.Now().UnixNano()%10000)
			break
		}
	}

	return username
}

// parseDisplayName parses display name into first and last name
func (s *UserServiceImpl) parseDisplayName(displayName string) (string, string) {
	parts := strings.Fields(displayName)
	if len(parts) == 0 {
		return "User", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

// ==================== Profile Methods ====================

// UpdateProfileInfo updates user profile information
func (s *UserServiceImpl) UpdateProfileInfo(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("ไม่พบผู้ใช้งาน")
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.StudentID != "" {
		// Check if student ID already exists
		existing, _ := s.userRepo.GetByStudentID(ctx, req.StudentID)
		if existing != nil && existing.ID != userID {
			return nil, errors.New("รหัสนักศึกษานี้ถูกใช้งานแล้ว")
		}
		user.StudentID = &req.StudentID
	}
	if req.Language != "" {
		user.Language = req.Language
	}
	if req.Theme != "" {
		user.Theme = req.Theme
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, userID, user); err != nil {
		return nil, fmt.Errorf("อัปเดตข้อมูลไม่สำเร็จ: %w", err)
	}

	return user, nil
}

// UpdateAvatar updates user avatar
func (s *UserServiceImpl) UpdateAvatar(ctx context.Context, userID uuid.UUID, fileData []byte, contentType string) (*dto.UpdateAvatarResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("ไม่พบผู้ใช้งาน")
	}

	// Validate file type
	if !isValidImageType(contentType) {
		return nil, errors.New("ไฟล์ต้องเป็นรูปภาพ (JPEG, PNG, GIF, WEBP)")
	}

	// Validate file size (max 5MB)
	if len(fileData) > 5*1024*1024 {
		return nil, errors.New("ไฟล์ต้องมีขนาดไม่เกิน 5MB")
	}

	// Delete old avatar if exists and is from R2
	if user.Avatar != "" && strings.Contains(user.Avatar, s.r2PublicURL) {
		oldKey := strings.TrimPrefix(user.Avatar, s.r2PublicURL+"/")
		_ = s.r2Storage.DeleteFile(oldKey)
	}

	// Generate unique filename
	ext := getExtensionFromContentType(contentType)
	filename := fmt.Sprintf("avatars/%s_%d%s", userID.String(), time.Now().UnixNano(), ext)

	// Upload new avatar
	avatarURL, err := s.r2Storage.UploadFile(bytes.NewReader(fileData), filename, contentType)
	if err != nil {
		return nil, fmt.Errorf("อัปโหลดรูปไม่สำเร็จ: %w", err)
	}

	// Update user
	user.Avatar = avatarURL
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, userID, user); err != nil {
		return nil, fmt.Errorf("อัปเดตข้อมูลไม่สำเร็จ: %w", err)
	}

	return &dto.UpdateAvatarResponse{
		AvatarURL: avatarURL,
	}, nil
}

// DeleteAvatar removes user avatar
func (s *UserServiceImpl) DeleteAvatar(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("ไม่พบผู้ใช้งาน")
	}

	// Delete from R2 if exists
	if user.Avatar != "" && strings.Contains(user.Avatar, s.r2PublicURL) {
		oldKey := strings.TrimPrefix(user.Avatar, s.r2PublicURL+"/")
		_ = s.r2Storage.DeleteFile(oldKey)
	}

	user.Avatar = ""
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, userID, user)
}

// isValidImageType checks if content type is a valid image
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
	}
	for _, t := range validTypes {
		if t == contentType {
			return true
		}
	}
	return false
}

// getExtensionFromContentType returns file extension from content type
func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}
