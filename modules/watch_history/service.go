package watch_history
import("errors";"time")
var ErrInvalidProgress=errors.New("invalid watch progress")
type Service struct{Repo *Repository}
func NewService(r *Repository)*Service{return &Service{Repo:r}}
func(s *Service)Save(uid,mid uint,p,d int64)error{if uid==0||mid==0||d<=0||p<0||p>d{return ErrInvalidProgress};return s.Repo.CreateOrUpdate(&WatchHistory{UserID:uid,MovieID:mid,ProgressSeconds:p,DurationSeconds:d,Completed:p>=d-30,WatchedAt:time.Now().UTC()})}
func(s *Service)List(uid uint,l,o int)([]WatchHistory,error){return s.Repo.FindByUserID(uid,l,o)}
func(s *Service)Get(uid,mid uint)(*WatchHistory,error){return s.Repo.FindByUserAndMovie(uid,mid)}
func(s *Service)Delete(uid,mid uint)error{return s.Repo.Delete(uid,mid)}
