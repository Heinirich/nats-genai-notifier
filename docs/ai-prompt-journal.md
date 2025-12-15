# AI Prompt Journal
_Real-Time Event Notification Service Using Go, NATS, Fiber, and Local Gemini (Ollama)_

This journal documents how AI assistance was used to design, debug, and complete the project while learning Go and event‑driven architecture.

---

## 1. Prompt: “Explain NATS pub/sub like I'm new to backend development.”
### AI Response Summary
- Compared NATS to a high‑speed postal system.
- Explained subjects, publishers, subscribers.
- Clarified why NATS is lightweight and ideal for microservices.

### Reflection
Helped me understand NATS quickly and gave confidence to design the raw and enriched subjects.

---

## 2. Prompt: “Generate Go code for a basic NATS subscriber.”
### AI Response Summary
- Provided a runnable subscriber example.
- Showed message callback patterns.
- Demonstrated auto‑reconnect features.

### Reflection
Formed the foundation of my event‑processing pipeline.

---

## 3. Prompt: “How do I call a local Ollama model from Go?”
### AI Response Summary
- Showed POST structure for `/api/generate`.
- Provided JSON payload example.
- Clarified streaming vs non‑streaming options.

### Reflection
Removed confusion about the Ollama API and enabled fast LLM integration.

---

## 4. Prompt: “Write a prompt that returns JSON with title, body, priority, and action.”
### AI Response Summary
AI generated a structured response template:

```json
{
 "title": "",
 "body": "",
 "priority": "",
 "action": ""
}
```

### Reflection
This became the core prompt in my `ollama.go` file and ensured consistent enrichments.

---

## 5. Prompt: “Fix this Go error: undefined GenerateSupportResponse.”
### AI Response Summary
- Identified missing imports.
- Suggested correcting module paths.
- Explained how Go resolves internal packages.

### Reflection
This fixed my project structure and improved understanding of Go modularization.

---

## 6. Prompt: “Improve my Fiber endpoint for sending messages to NATS.”
### AI Response Summary
- Improved JSON parsing.
- Suggested cleaner error responses.
- Improved API structure.

### Reflection
My `/send` endpoint became cleaner and more reliable.

---

## 7. Prompt: “AI response is not valid JSON. Fix the prompt.”
### AI Response Summary
- Recommended strict rule: **Return ONLY JSON.**
- Suggested removing explanations or extra text.

### Reflection
After applying this, my entire NATS → AI → NATS processing chain became stable.

---

## 8. Prompt: “Explain priority classification logic for support systems.”
### AI Response Summary
- High → severe frustration, outages
- Medium → repeated issues
- Low → general inquiries

### Reflection
This was integrated directly into my LLM prompt.

---

## 9. Prompt: “Help me design project architecture.”
### AI Response Summary
Suggested the exact flow:

```
Fiber API → NATS 'raw' → AI Processor → NATS 'enriched' → Consumers
```

### Reflection
This became the final project architecture.

---

## 10. Prompt: “Summarize everything into a final system description I can use in docs.”
### AI Response Summary
Provided clear descriptions for:
- Architecture
- Workflow
- Data transformation
- AI usage

### Reflection
Used this in the final Toolkit document and README.

---

# End of AI Prompt Journal
