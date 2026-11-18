package catalog

import "errors"
var ErrInvalidID = errors.New("invalid id or pagination")
type Service struct { Repo *Repository }
func NewService(r *Repository) *Service { return &Service{Repo:r} }
func (s *Service) Movies(q string,g uint,y uint16,l,o int)([]Movie,int64,error){ if l>100||o<0{return nil,0,ErrInvalidID}; return s.Repo.ListMovies(q,g,y,l,o) }
func (s *Service) Movie(id uint)(*Movie,error){ if id==0{return nil,ErrInvalidID}; return s.Repo.MovieByID(id) }
func (s *Service) SeriesList(q string,l,o int)([]Series,int64,error){ if l>100||o<0{return nil,0,ErrInvalidID}; return s.Repo.ListSeries(q,l,o) }
func (s *Service) SeriesItem(id uint)(*Series,error){ if id==0{return nil,ErrInvalidID}; return s.Repo.SeriesByID(id) }
func (s *Service) Seasons(id uint)([]Season,error){ if id==0{return nil,ErrInvalidID}; return s.Repo.Seasons(id) }
func (s *Service) Episodes(id uint)([]Episode,error){ if id==0{return nil,ErrInvalidID}; return s.Repo.Episodes(id) }
