package main

import (
	"errors"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vishwas-11/trusture-backend/internal/api/handlers"
	"github.com/vishwas-11/trusture-backend/internal/api/routes"
	"github.com/vishwas-11/trusture-backend/internal/blockchain"
	"github.com/vishwas-11/trusture-backend/internal/config"
	"github.com/vishwas-11/trusture-backend/internal/repository/mongo"
	"github.com/vishwas-11/trusture-backend/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	// Router
	router := routes.NewRouter(cfg)

	// Mongo
	db := mongo.Connect(cfg.MongoURI)

	// Services
	donationRepo := mongo.NewDonationRepository(db)
	proofRepo := mongo.NewProofRepository(db)
	auditRepo := mongo.NewAuditRepository(db)
	donationService := service.NewDonationService(donationRepo)
	proofService := service.NewProofService(proofRepo)
	vendorService := service.NewVendorService()
	auditService := service.NewAuditService(
		vendorService,
		cfg.AuditPassThreshold,
		cfg.DistanceWarnKM,
		cfg.DistanceFailKM,
		cfg.GSTVerificationEnabled,
	)
	auditPipelineService := service.NewAuditPipelineService(
		donationService,
		proofService,
		auditService,
		auditRepo,
	)
	extractService := service.NewDocumentExtractService(cfg.GeminiAPIKey, cfg.OpenAIAPIKey)
	adminRepo := mongo.NewAdminRepository(db)
	adminAuthService := service.NewAdminAuthService(adminRepo, cfg.JWTSecret)

	handlers.InitNGOHandler(donationService)
	handlers.InitProofHandler(proofService, extractService)
	handlers.InitAdminHandler(adminAuthService, donationService, proofService, auditPipelineService, auditRepo)

	// Blockchain
	client := blockchain.NewClient(cfg.PolygonRPC)

	go blockchain.ListenDonations(
		client,
		cfg.ContractAddress,
		func(vLog types.Log) {

			event, err := blockchain.DecodeDonationEvent(vLog)
			if err != nil {
				log.Println("❌ Decode error:", err)
				return
			}

			amountMatic := new(big.Float).Quo(
				new(big.Float).SetInt(event.Amount),
				big.NewFloat(1e18),
			)
			amountMaticFloat, _ := amountMatic.Float64()

			log.Printf(
				"🧾 Donation event | tx=%s block=%d donor=%s ngo=%s amount=%.6f MATIC purpose=%q eventTs=%s",
				vLog.TxHash.Hex(),
				vLog.BlockNumber,
				event.Donor,
				event.NGO,
				amountMaticFloat,
				event.Purpose,
				event.Timestamp.String(),
			)

			_, err = donationService.RecordDonation(
				vLog.TxHash.Hex(),
				event.Donor,
				event.NGO,
				amountMaticFloat, // wei -> MATIC
				event.Purpose,
			)

			if err != nil {
				log.Printf("❌ Failed to persist donation tx=%s: %v", vLog.TxHash.Hex(), err)
				return
			}

			log.Printf("✅ Donation recorded from blockchain tx=%s", vLog.TxHash.Hex())
		},
	)

	log.Printf("🚀 TRUSTURE running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		// Common local-dev issue: another process is already bound to the same port.
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		if strings.Contains(strings.ToLower(err.Error()), "only one usage") ||
			strings.Contains(strings.ToLower(err.Error()), "address already in use") {
			log.Fatalf("❌ Port %s is already in use. Stop the existing backend process (Windows: netstat -ano | findstr :%s, then taskkill /PID <pid> /F) or change PORT in Backend/.env", cfg.Port, cfg.Port)
		}
		log.Fatalf("❌ HTTP server failed: %v", err)
	}
}
