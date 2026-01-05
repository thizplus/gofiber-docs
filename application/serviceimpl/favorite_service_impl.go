package serviceimpl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type FavoriteServiceImpl struct {
	favoriteRepo repositories.FavoriteRepository
}

func NewFavoriteService(favoriteRepo repositories.FavoriteRepository) services.FavoriteService {
	return &FavoriteServiceImpl{
		favoriteRepo: favoriteRepo,
	}
}

func (s *FavoriteServiceImpl) AddFavorite(ctx context.Context, userID uuid.UUID, req *dto.AddFavoriteRequest) (*dto.FavoriteResponse, error) {
	// Check if already favorited by URL
	exists, err := s.favoriteRepo.ExistsByUserIDAndURL(ctx, userID, req.URL)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("item already in favorites")
	}

	// Check if already favorited by ExternalID (for places)
	if req.ExternalID != "" {
		exists, err = s.favoriteRepo.ExistsByUserIDAndExternalID(ctx, userID, req.ExternalID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("item already in favorites")
		}
	}

	favorite := dto.AddFavoriteRequestToFavorite(req)
	favorite.UserID = userID
	favorite.CreatedAt = time.Now()

	if err := s.favoriteRepo.Create(ctx, favorite); err != nil {
		return nil, err
	}

	return dto.FavoriteToFavoriteResponse(favorite), nil
}

func (s *FavoriteServiceImpl) GetFavorites(ctx context.Context, userID uuid.UUID, req *dto.GetFavoritesRequest) (*dto.FavoriteListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	var favorites []*models.Favorite
	var total int64
	var err error

	if req.Type != "" {
		favorites, err = s.favoriteRepo.GetByUserIDAndType(ctx, userID, req.Type, offset, req.PageSize)
		if err != nil {
			return nil, err
		}
		total, err = s.favoriteRepo.CountByUserIDAndType(ctx, userID, req.Type)
	} else {
		favorites, err = s.favoriteRepo.GetByUserID(ctx, userID, offset, req.PageSize)
		if err != nil {
			return nil, err
		}
		total, err = s.favoriteRepo.CountByUserID(ctx, userID)
	}

	if err != nil {
		return nil, err
	}

	var favoriteResponses []dto.FavoriteResponse
	for _, f := range favorites {
		favoriteResponses = append(favoriteResponses, *dto.FavoriteToFavoriteResponse(f))
	}

	return &dto.FavoriteListResponse{
		Favorites: favoriteResponses,
		Meta: dto.PaginationMeta{
			Total:  total,
			Offset: offset,
			Limit:  req.PageSize,
		},
	}, nil
}

func (s *FavoriteServiceImpl) RemoveFavorite(ctx context.Context, userID uuid.UUID, favoriteID uuid.UUID) error {
	favorite, err := s.favoriteRepo.GetByID(ctx, favoriteID)
	if err != nil {
		return errors.New("favorite not found")
	}

	if favorite.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.favoriteRepo.Delete(ctx, favoriteID)
}

func (s *FavoriteServiceImpl) CheckFavorite(ctx context.Context, userID uuid.UUID, req *dto.CheckFavoriteRequest) (*dto.CheckFavoriteResponse, error) {
	response := &dto.CheckFavoriteResponse{
		IsFavorite: false,
	}

	// Check by ExternalID first (for places)
	if req.ExternalID != "" {
		favorite, err := s.favoriteRepo.GetByUserIDAndExternalID(ctx, userID, req.ExternalID)
		if err == nil && favorite != nil {
			response.IsFavorite = true
			response.FavoriteID = &favorite.ID
			return response, nil
		}
	}

	// Check by URL
	if req.URL != "" {
		favorite, err := s.favoriteRepo.GetByUserIDAndURL(ctx, userID, req.URL)
		if err == nil && favorite != nil {
			response.IsFavorite = true
			response.FavoriteID = &favorite.ID
			return response, nil
		}
	}

	return response, nil
}

func (s *FavoriteServiceImpl) BatchCheckFavorites(ctx context.Context, userID uuid.UUID, externalIDs []string) (*dto.BatchCheckFavoritesResponse, error) {
	favorites, err := s.favoriteRepo.GetByUserIDAndExternalIDs(ctx, userID, externalIDs)
	if err != nil {
		return nil, err
	}

	// Create a map of externalID -> favorite for quick lookup
	favoriteMap := make(map[string]*models.Favorite)
	for _, f := range favorites {
		favoriteMap[f.ExternalID] = f
	}

	// Build response
	response := &dto.BatchCheckFavoritesResponse{
		Items: make(map[string]dto.CheckFavoriteResponse),
	}

	for _, extID := range externalIDs {
		if fav, ok := favoriteMap[extID]; ok {
			response.Items[extID] = dto.CheckFavoriteResponse{
				IsFavorite: true,
				FavoriteID: &fav.ID,
			}
		} else {
			response.Items[extID] = dto.CheckFavoriteResponse{
				IsFavorite: false,
			}
		}
	}

	return response, nil
}

func (s *FavoriteServiceImpl) ToggleFavorite(ctx context.Context, userID uuid.UUID, req *dto.AddFavoriteRequest) (*dto.CheckFavoriteResponse, error) {
	checkReq := &dto.CheckFavoriteRequest{
		Type:       req.Type,
		URL:        req.URL,
		ExternalID: req.ExternalID,
	}

	checkResp, err := s.CheckFavorite(ctx, userID, checkReq)
	if err != nil {
		return nil, err
	}

	if checkResp.IsFavorite {
		// Remove from favorites
		if err := s.favoriteRepo.Delete(ctx, *checkResp.FavoriteID); err != nil {
			return nil, err
		}
		return &dto.CheckFavoriteResponse{
			IsFavorite: false,
		}, nil
	}

	// Add to favorites
	favorite, err := s.AddFavorite(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &dto.CheckFavoriteResponse{
		IsFavorite: true,
		FavoriteID: &favorite.ID,
	}, nil
}
