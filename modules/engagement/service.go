package engagement

import "errors"

var ErrInvalidInput=errors.New("invalid engagement input")
type Service struct { repo *Repository }
func NewService(repo *Repository)*Service{return &Service{repo:repo}}
func(s *Service)AddFavorite(uid,mid uint)error{if uid==0||mid==0{return ErrInvalidInput};return s.repo.AddFavorite(uid,mid)}
func(s *Service)RemoveFavorite(uid,mid uint)error{if uid==0||mid==0{return ErrInvalidInput};return s.repo.RemoveFavorite(uid,mid)}
func(s *Service)Favorites(uid uint,limit,offset int)([]Favorite,error){if uid==0{return nil,ErrInvalidInput};if limit<1||limit>100{limit=20};if offset<0{offset=0};return s.repo.Favorites(uid,limit,offset)}
func(s *Service)Rate(uid,mid uint,score uint8)error{if uid==0||mid==0||score<1||score>10{return ErrInvalidInput};return s.repo.Rate(uid,mid,score)}
func(s *Service)RatingSummary(mid uint)(map[string]any,error){if mid==0{return nil,ErrInvalidInput};return s.repo.RatingSummary(mid)}
