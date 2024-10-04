package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/config"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/routes"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/services"
	"github.com/chloejepson16/Go-API-Tech-Challenge/internal/swagger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"

	_ "github.com/lib/pq"
)

func main(){
	ctx:= context.Background()
	if err:= run(ctx); err != nil{
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run(ctx context.Context) error{
	//setting up configuration for local environment
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("[in run]: %w", err)
	}

	logger := httplog.NewLogger("courses-microservice", httplog.Options{
		LogLevel:        cfg.LogLevel,
		JSON:            false,
		Concise:         true,
		ResponseHeaders: false,
	})

	//Connecting to DB
	dbUser:= os.Getenv("DATABASE_USER")
	dbPassword:= os.Getenv("DATABASE_PASSWORD")
	dbHost:= os.Getenv("DATABASE_HOST")
	dbPort:= os.Getenv("DATABASE_PORT")
	dbName:= "coursesDB"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err:= sql.Open("postgres", dsn)
	if err != nil{
		log.Fatalf("Failed to open a DB connection: %v", err)
	}
	defer db.Close()
	err= db.Ping()
	if err != nil{
		log.Fatalf("Failed to open a DB connection on ping: %v", err)
	}else{
		fmt.Print("Success, connected to DB")
	}

	//creating a router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	svs:= services.NewPersonService(db)
	routes.RegisterRoutes(r, logger, svs, routes.WithRegisterHealthRoute(true))

	if cfg.HTTPUseSwagger {
		swagger.RunSwagger(r, logger, cfg.HTTPDomain+cfg.HTTPPort)
	}

	serverInstance := &http.Server{
		Addr:              cfg.HTTPDomain + cfg.HTTPPort,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 500 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
		WriteTimeout:      500 * time.Millisecond,
		Handler:           r,
	}

	// Graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		fmt.Println()
		logger.Info("Shutdown signal received")

		shutdownCtx, err := context.WithTimeout(
			serverCtx, time.Duration(cfg.HTTPShutdownDuration)*time.Second,
		)
		if err != nil {
			log.Fatalf("Error creating context.WithTimeout. err: %v", err)
		}

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := serverInstance.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Error shutting down server. err: %v", err)
		}
		serverStopCtx()
	}()

	err = http.ListenAndServe("localhost:3000", r)
	if err!= nil{
		return err
	}
	return nil
}