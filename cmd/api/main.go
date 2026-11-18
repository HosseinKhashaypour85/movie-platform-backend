package main

import("context";"errors";"log/slog";"net/http";"os";"os/signal";"syscall";"time";"movie-platform/internal/config";"movie-platform/internal/database";"movie-platform/internal/httpx";"movie-platform/internal/middleware";"movie-platform/modules/auth";"movie-platform/modules/users";"movie-platform/modules/watch_history";"github.com/go-chi/chi/v5")
func main(){
 logger:=slog.New(slog.NewJSONHandler(os.Stdout,nil));cfg,e:=config.Load();if e!=nil{logger.Error("configuration failed","error",e);os.Exit(1)}
 db,e:=database.Open(cfg);if e!=nil{logger.Error("database failed","error",e);os.Exit(1)};sqlDB,e:=db.DB();if e!=nil{logger.Error("database pool failed","error",e);os.Exit(1)};defer sqlDB.Close()
 if cfg.Environment!="production"{if e=db.AutoMigrate(&users.User{},&watch_history.WatchHistory{});e!=nil{logger.Error("migration failed","error",e);os.Exit(1)}}
 userRepo:=users.NewRepository(db);ah:=auth.NewHandler(auth.NewService(userRepo,cfg));wh:=watch_history.NewHandler(watch_history.NewService(watch_history.NewRepository(db)))
 router:=chi.NewRouter();router.Use(middleware.RequestID,middleware.Recover(logger),middleware.SecurityHeaders,middleware.BodyLimit(1<<20),middleware.Timeout(15*time.Second),middleware.RateLimit(120,time.Minute),middleware.CORS(cfg.AllowedOrigins))
 router.Get("/",func(w http.ResponseWriter,r *http.Request){httpx.JSON(w,200,map[string]string{"service":"movie-platform-api","status":"running"})})
 router.Get("/health/live",func(w http.ResponseWriter,r *http.Request){httpx.JSON(w,200,map[string]string{"status":"alive"})})
 router.Get("/health/ready",func(w http.ResponseWriter,r *http.Request){ctx,cancel:=context.WithTimeout(r.Context(),2*time.Second);defer cancel();if e:=sqlDB.PingContext(ctx);e!=nil{httpx.Error(w,503,"not_ready","dependencies unavailable");return};httpx.JSON(w,200,map[string]string{"status":"ready"})})
 router.Route("/api/v1",func(r chi.Router){auth.RegisterRoutes(r,ah);r.Group(func(p chi.Router){p.Use(middleware.Auth(cfg.JWTSecret,cfg.JWTIssuer));p.Get("/profile",func(w http.ResponseWriter,r *http.Request){id,ok:=middleware.UserIDFromContext(r.Context());if !ok{httpx.Error(w,401,"unauthorized","authentication required");return};u,e:=userRepo.FindByID(id);if e!=nil||!u.IsActive{httpx.Error(w,401,"unauthorized","authentication required");return};httpx.JSON(w,200,u)});watch_history.RegisterRoutes(p,wh)})})
 srv:=&http.Server{Addr:":"+cfg.Port,Handler:router,ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:30*time.Second,IdleTimeout:60*time.Second,MaxHeaderBytes:1<<20}
 go func(){logger.Info("server starting","addr",srv.Addr);if e:=srv.ListenAndServe();e!=nil&&!errors.Is(e,http.ErrServerClosed){logger.Error("server failed","error",e);os.Exit(1)}}()
 stop,cancel:=signal.NotifyContext(context.Background(),syscall.SIGINT,syscall.SIGTERM);defer cancel();<-stop.Done();ctx,done:=context.WithTimeout(context.Background(),10*time.Second);defer done();if e:=srv.Shutdown(ctx);e!=nil{logger.Error("shutdown failed","error",e)}
}
