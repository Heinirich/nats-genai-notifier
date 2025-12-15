# Real-Time Event Notification Service — Toolkit Document

## 1) Title & Objective
**Tech:** Go + Fiber + NATS + Local Gemini (Ollama)  
**Objective:** Build a real-time event enrichment service that receives raw support messages, processes them using a local Gemini model via Ollama, rewrites the message, assigns a priority level (High/Medium/Low), and generates a recommended action. Output is published back to NATS as structured JSON.

---

## 2) Quick Summary of the Technology
**Go (Golang):** A fast, compiled, concurrent language ideal for backend and distributed systems.  
**Fiber:** A lightweight, Express-like web framework for Go.  
**NATS:** A high‑performance messaging system used for real-time event streaming in microservices.  
**Ollama + Gemini:** Runs LLMs locally. Gemini is used to summarize text, rewrite messages, classify priority, and generate support actions.

**Real-world usage:**  
NATS powers microservice communication systems in fintech, telecom, IoT, and monitoring platforms where real‑time message flow and low latency are essential.

---

## 3) System Requirements
- OS: Windows/macOS/Linux
- Go: v1.21+
- Editor: VS Code or GoLand
- Dependencies: Fiber, nats.go
- Tools: git, curl
- Ollama installed locally
- Gemini model:
  ```bash
  ollama pull gemini
  ```

---

## 4) Installation & Setup
### Install Go
```bash
sudo apt install golang-go
go version
```

### Create Project
```bash
mkdir support-ai
cd support-ai
go mod init support-ai
```

### Install Dependencies
```bash
go get github.com/gofiber/fiber/v2
go get github.com/nats-io/nats.go
```

### Run NATS Server
```bash
docker run -p 4222:4222 nats
```

### Install and Start Ollama
```bash
ollama pull gemini
ollama serve
```

### Recommended Structure
```
support-ai/
 ├── cmd/api/main.go
 ├── internal/
 │     ├── ai/ollama.go
 │     ├── models/ticket.go
 │     └── natsclient/subscriber.go
 └── go.mod
```

---

## 5) Minimal Working Example

### **Model — ticket.go**
```go
type Ticket struct {
    Title    string `json:"title"`
    Body     string `json:"body"`
    Priority string `json:"priority"`
    Action   string `json:"action"`
}
```

---

### **AI Processor — ollama.go**
```go
prompt := `
Rewrite the user's message professionally, summarize the issue,
categorize priority (High, Medium, Low), and generate a recommended action.

Return ONLY valid JSON:
{
 "title": "",
 "body": "",
 "priority": "",
 "action": ""
}

User message:
` + raw
```

Ollama call:
```go
resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(bodyBytes))
```

---

### **NATS Subscriber — subscriber.go**
```go
nc.Subscribe("support.raw", func(msg *nats.Msg) {
    enriched, _ := ai.GenerateSupportResponse(string(msg.Data))
    nc.Publish("support.enriched", []byte(enriched))
})
```

---

### **Fiber API Endpoint — main.go**
```go
app.Post("/send", func(c *fiber.Ctx) error {
    var req Req
    c.BodyParser(&req)
    nc.Publish("support.raw", []byte(req.Message))
    return c.JSON(fiber.Map{"status": "sent"})
})
```

---

### Expected Output (AI Enriched JSON)
```json
{
  "title": "Marxwel 0746273465 - Frustrated by Internet Disruptions",
  "body": "The user Maxwell is experiencing slow internet and frequent disconnections. He is frustrated and requires support assistance.",
  "priority": "High",
  "action": "Escalate to network diagnostics and contact the user immediately."
}
```

---

## 6) AI Prompt Journal (with refined prompts)
See `ai-prompt-journal.md` for full details.

Snapshot:

| Step | Prompt Theme | Key Outcome |
|-----:|--------------|-------------|
| 1 | Understanding NATS | Learned pub/sub concepts |
| 2 | Build Go subscriber | Enabled event pipeline |
| 3 | Ollama integration | Correct API endpoint & format |
| 4 | JSON-only AI output | Ensured stable parsing |
| 5 | Fiber endpoint design | Improved structure & clarity |

---

## 7) Common Issues & Fixes
| Issue | Cause | Fix |
|------|-------|------|
| NATS connection refused | Server not running | Start NATS or change URL |
| Ollama gives invalid JSON | Model adds extra text | Add strict JSON-only command |
| Fiber cannot parse requests | Missing struct tags | Ensure proper JSON tags |
| Missing modules | Dependency mismatch | Run `go mod tidy` |

---

## 8) References
- Go Documentation
- Fiber Web Framework Docs
- NATS Official Documentation
- Ollama API Docs
- Gemini Model Guide
- Moringa School Shared Resources

---

# End of Toolkit Document
