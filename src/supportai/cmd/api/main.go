package main

import (
	"Moringa_AI/src/supportai/natsclient"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
)

type sendRequest struct {
	Message string `json:"message"`
}

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	log.Printf("Connected to NATS at %s\n", nats.DefaultURL)

	// Start NATS Processor
	processor := natsclient.NewProcessor(nc)
	if err := processor.Start(); err != nil {
		log.Fatalf("failed to start NATS processor: %v", err)
	}

	// Fiber app
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// POST /send - send raw support message which goes into NATS
	app.Post("/send", func(c *fiber.Ctx) error {
		var req sendRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid body",
			})
		}

		if req.Message == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "message is required",
			})
		}

		if err := nc.Publish("support.raw", []byte(req.Message)); err != nil {
			log.Printf("failed to publish to NATS: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to publish message",
			})
		}

		log.Printf("[API] Sent message to support.raw: %s\n", req.Message)

		return c.JSON(fiber.Map{
			"status":  "sent",
			"message": req.Message,
		})
	})

	// Graceful shutdown handling
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Printf("Fiber server stopped: %v\n", err)
		}
	}()
	log.Println("Fiber server running on :8080")

	// Wait for interrupt (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	_ = app.Shutdown()
}
