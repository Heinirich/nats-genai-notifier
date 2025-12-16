package natsclient

import (
	"Moringa_AI/src/supportai/ai"
	"Moringa_AI/src/supportai/models"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

const (
	rawSubject      = "support.raw"
	enrichedSubject = "support.enriched"
)

// Processor holds dependencies for NATS processing.
type Processor struct {
	nc *nats.Conn
}

// NewProcessor creates a new Processor.
func NewProcessor(nc *nats.Conn) *Processor {
	return &Processor{nc: nc}
}

// Start sets up the subscriber and processes incoming messages.
func (p *Processor) Start() error {
	_, err := p.nc.Subscribe(rawSubject, func(msg *nats.Msg) {
		rawMsg := string(msg.Data)
		log.Printf("[NATS] Received raw message: %s\n", rawMsg)

		ticket, err := ai.GenerateSupportResponse(rawMsg)
		if err != nil {
			log.Printf("[AI] Error generating support response: %v\n", err)
			return
		}

		payload, err := json.Marshal(ticket)
		if err != nil {
			log.Printf("[JSON] Error marshalling ticket: %v\n", err)
			return
		}

		if err := p.nc.Publish(enrichedSubject, payload); err != nil {
			log.Printf("[NATS] Error publishing enriched message: %v\n", err)
			return
		}

		log.Printf("[NATS] Published enriched message to %s\n", enrichedSubject)
	})

	if err != nil {
		return err
	}

	log.Printf("[NATS] Subscribed to %s\n", rawSubject)
	return nil
}

// Helper for consumers who want to subscribe to enriched messages in code.
func SubscribeToEnriched(nc *nats.Conn, handler func(*models.Ticket)) error {
	_, err := nc.Subscribe(enrichedSubject, func(msg *nats.Msg) {
		var t models.Ticket
		if err := json.Unmarshal(msg.Data, &t); err != nil {
			log.Printf("[JSON] Unable to unmarshal enriched ticket: %v\n", err)
			return
		}
		handler(&t)
	})

	if err != nil {
		return err
	}

	log.Printf("[NATS] Subscribed to %s\n", enrichedSubject)
	return nil
}
