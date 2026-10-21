package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"raffle-saas/internal/config"
	"raffle-saas/internal/db"
	"raffle-saas/internal/handler"
	"raffle-saas/internal/middleware"
	"raffle-saas/internal/repository"
	"raffle-saas/internal/service"
	"raffle-saas/internal/utils"

	"github.com/go-chi/chi/v5"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err := utils.ConnectMongoDB(ctx, cfg.MongoURI)
	if err != nil {
		slog.Error("failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}

	mongoDB := mongoClient.Database(cfg.MongoDB)

	if err := db.RunMigrations(ctx, mongoDB); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)

	userRepo := repository.NewUserRepository(mongoDB)
	raffleRepo := repository.NewRaffleRepository(mongoDB)
	ticketRepo := repository.NewTicketRepository(mongoDB)
	auditRepo := repository.NewAuditRepository(mongoDB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(mongoDB)

	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, jwtManager)
	raffleSvc := service.NewRaffleService(raffleRepo, ticketRepo, auditRepo)
	ticketSvc := service.NewTicketService(ticketRepo, raffleRepo, auditRepo, mongoClient)

	authH := handler.NewAuthHandler(authSvc)
	raffleH := handler.NewRaffleHandler(raffleSvc)
	ticketH := handler.NewTicketHandler(ticketSvc)
	webhookH := handler.NewWebhookHandler()
	adminH := handler.NewAdminHandler(userRepo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.Post("/oauth/twitch", authH.OturnhTwitch)
			r.Post("/refresh", authH.Refresh)
			r.Group(func(r chi.Router) {
				r.Use(middleware.Auth(jwtManager))
				r.Post("/logout", authH.Logout)
				r.Get("/me", authH.Me)
			})
		})

		r.Route("/raffles", func(r chi.Router) {
			r.Get("/", raffleH.List)
			r.Get("/{id}", raffleH.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(middleware.Auth(jwtManager))
				r.Post("/", raffleH.Create)
				r.Get("/my", raffleH.ListMy)
				r.Patch("/{id}", raffleH.Update)
				r.Post("/{id}/publish", raffleH.Publish)
				r.Post("/{id}/cancel", raffleH.Cancel)
				r.Post("/{id}/sort", raffleH.Sort)
				r.Post("/{id}/resort", raffleH.ReSort)

				r.Post("/{id}/tickets", ticketH.Purchase)
				r.Get("/{id}/tickets", ticketH.ListByRaffle)
			})
		})

		r.Route("/tickets", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.Auth(jwtManager))
				r.Get("/my", ticketH.ListMy)
				r.Post("/{id}/upload-proof", ticketH.UploadProof)
				r.Post("/{id}/verify", ticketH.Verify)
				r.Post("/{id}/reject", ticketH.Reject)
			})
		})

		r.Route("/webhooks", func(r chi.Router) {
			r.Post("/stripe", webhookH.Stripe)
		})

		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.Auth(jwtManager))
				r.Get("/", adminH.ListUsers)
				r.Patch("/{id}/role", adminH.UpdateRole)
				r.Patch("/{id}/status", adminH.UpdateStatus)
			})
		})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	if err := mongoClient.Disconnect(shutdownCtx); err != nil {
		slog.Error("failed to disconnect MongoDB", "error", err)
	}

	slog.Info("server stopped")
}
