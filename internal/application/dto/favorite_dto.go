package dto

import (
	"time"

	"github.com/mehdi1443/cinema-booking-system/internal/domain/entity"
)

type AddFavoriteMovieRequest struct {
	MovieID int64 `json:"movie_id" binding:"required"`
}

type AddFavoriteCinemaRequest struct {
	CinemaID int64 `json:"cinema_id" binding:"required"`
}

type FavoriteFilterRequest struct {
	Target string `form:"target"`
	Limit  int    `form:"limit,default=10"`
	Offset int    `form:"offset,default=0"`
}

type FavoriteResponse struct {
	ID        int64                 `json:"id"`
	UserID    int64                 `json:"user_id"`
	MovieID   *int64                `json:"movie_id,omitempty"`
	CinemaID  *int64                `json:"cinema_id,omitempty"`
	Movie     *MovieListResponse    `json:"movie,omitempty"`
	Cinema    *CinemaDetailResponse `json:"cinema,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
}

type FavoriteListWithCountResponse struct {
	Favorites []*FavoriteResponse `json:"favorites"`
	Total     int64               `json:"total"`
}

func ToFavoriteResponse(favorite *entity.Favorite) *FavoriteResponse {
	if favorite == nil {
		return nil
	}

	resp := &FavoriteResponse{
		ID:        favorite.ID,
		UserID:    favorite.UserID,
		MovieID:   favorite.MovieID,
		CinemaID:  favorite.CinemaID,
		CreatedAt: favorite.CreatedAt,
	}

	if favorite.Movie != nil {
		resp.Movie = ToMovieListResponse(favorite.Movie)
	}

	if favorite.Cinema != nil {
		resp.Cinema = ToCinemaDetailResponse(favorite.Cinema)
	}

	return resp
}

func ToFavoriteResponseSlice(favorites []entity.Favorite) []*FavoriteResponse {
	if len(favorites) == 0 {
		return []*FavoriteResponse{}
	}

	responses := make([]*FavoriteResponse, 0, len(favorites))
	for ind := range favorites {
		if resp := ToFavoriteResponse(&favorites[ind]); resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses
}
