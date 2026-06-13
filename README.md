Here’s the **expanded end‑to‑end technical statement** for the *Mental Wellness Tracker* challenge, with a recommended tech stack and evaluation parameters fully integrated:

---

## 📝 Technical Problem Statement: *Mental Wellness Tracker*

### **Objective**
Develop a Generative AI‑powered solution that helps students monitor and improve their mental well‑being during high‑stakes exams (NEET, JEE, CUET, CAT, GATE, UPSC). The app should analyze journaling and mood logs to uncover hidden stress triggers and provide contextual, empathetic wellness support.

---

### **Inputs**
- **Daily Journaling Entries**: Free‑form text reflecting thoughts and emotions  
- **Mood Logs**: Self‑reported stress levels, emotions, energy states  
- **Optional Metadata**: Study hours, sleep patterns, exam countdown  

---

### **Outputs**
1. **Stress Analysis** – Detect hidden triggers and emotional patterns  
2. **Conversational Support** – Personalized coping strategies, mindfulness exercises, motivational encouragement  
3. **Progress Tracking** – Visual mood/stress trends over time  
4. **Safe Companion Role** – Empathetic, always‑available digital support  

---

### **Tech Stack**

#### 🔧 Backend
- **Language:** Golang (Fiber/Gin)  
- **AI Layer:**  
  - LLM integration (Ollama, LM Studio, or API like OpenAI/Claude) for journaling analysis  
  - Rule‑based + ML hybrid for mood classification  
- **Database:** PostgreSQL (encrypted journaling storage)  
- **Deployment:** Render or Koyeb free tier (Dockerized microservice)  
- **Containerization:** Dockerfile (`FROM golang:1.22-alpine`  

#### 🎨 Frontend
- **Framework:** React (or Vue)  
- **Hosting:** Netlify or Vercel free tier  
- **UI Features:**  
  - Journaling input box + mood sliders  
  - Stress trend charts (Chart.js/D3.js)  
  - Coping strategy cards (breathing, mindfulness, motivational prompts)  
- **Accessibility:** WCAG 2.1 compliance, ARIA roles, multilingual support (English + Hindi), mobile‑friendly design  

#### ⚙️ DevOps
- **CI/CD:** GitHub Actions → run tests + linting before deploy  
- **Monitoring:** Render/Koyeb logs + optional Prometheus/Grafana  
- **Secrets Management:** Environment variables for API keys and DB credentials  

---

### **Constraints (Evaluation Parameters)**

#### 🔴 High Impact
- **Code Quality**  
  - Clean, modular, well‑structured code  
  - Clear separation of API, AI logic, and data layers  
- **Problem Statement Alignment**  
  - Must deliver stress detection, coping strategies, and empathetic support  

#### 🟠 Medium Impact
- **Security**  
  - Encrypt journaling data (AES at rest, TLS in transit)  
  - Input sanitization to prevent injection attacks  
  - Privacy by design: no external sharing of sensitive logs  
- **Efficiency**  
  - Goroutines (Go) or for parallel AI calls  
  - Caching of repeated journaling analysis  
  - Lightweight Docker image (< 200 MB)  

#### 🟢 Low Impact
- **Testing**  
  - Unit tests for journaling parser, mood classifier, coping strategy generator  
  - Integration tests for API endpoints  
- **Accessibility**  
  - WCAG 2.1 compliance  
  - Keyboard navigation, screen reader support  
  - Mobile optimization for low‑end devices  

---

### **Success Criteria**
- Students can log entries daily in < 30 seconds.  
- AI detects stress triggers with > 80% accuracy compared to manual review.  
- Coping strategies are contextually relevant and safe.  
- App passes security audits (encrypted storage, sanitized inputs).  
- Accessible UI usable on mobile devices with low bandwidth.  
- Fully deployed on Render (backend) + Netlify/Vercel (frontend) with public URLs.  

---

### **Deployment Flow**
1. Push backend code to GitHub → Render auto‑deploys via connected repo.  
2. Push frontend code to GitHub → Netlify/Vercel auto‑deploys.  
3. Backend URL consumed by frontend via `.env` config.  
4. CI/CD ensures tests pass before deployment.  
5. Students access via public frontend URL, journaling data flows securely to backend, AI generates insights, frontend displays coping strategies.  

---

✅ This is now a **competition‑ready technical statement**: clear objectives, defined tech stack, constraints tied to evaluation parameters, and deployment flow.  

Would you like me to **map this into a scoring rubric** (points per parameter, weightage) so you can evaluate submissions or guide your own build priorities?

---

## Backend implementation

The workspace now includes a runnable Go backend for the first issue slice: secure daily journal and mood capture.

It now supports PostgreSQL persistence, Docker packaging, and Render deployment.

### Run locally

1. Generate a base64-encoded 32-byte encryption key.
2. Set `JOURNAL_ENCRYPTION_KEY` to that value.
3. Set `DATABASE_URL` to a PostgreSQL connection string if you want Postgres, or leave it empty to use the local encrypted file store.
4. Start the API with `go run ./cmd/api`.

Example environment variables:

```bash
JOURNAL_ENCRYPTION_KEY=...base64-32-bytes...
JOURNAL_DATA_FILE=data/journals.json
ADDR=:8080
```

### Docker

Build and run the container with your environment variables:

```bash
docker build -t mental-wellness-tracker-api .
docker run -p 8080:8080 -e JOURNAL_ENCRYPTION_KEY=... -e DATABASE_URL=... mental-wellness-tracker-api
```

### Render

The `render.yaml` blueprint provisions a web service and managed PostgreSQL database.
Set `JOURNAL_ENCRYPTION_KEY` as a secret in the Render dashboard, then deploy the repository with the blueprint.

### Frontend chart

A React/Vite frontend lives in `frontend/` and reads the trend endpoint at `GET /v1/trends?user_id=...`.

Run it locally with:

```bash
cd frontend
npm install
npm run dev
```

Set `VITE_API_BASE_URL` if the frontend is not talking to `http://localhost:8080`.

### API

- `GET /healthz` returns a simple readiness check.
- `POST /v1/entries` creates a journal entry.
- `GET /v1/entries?user_id=student-1` lists saved entries for a student.
- `GET /v1/trends?user_id=student-1` returns chart-ready daily mood and stress averages.

### Payload

```json
{
  "user_id": "student-1",
  "entry_text": "Today felt heavy, but I finished my revision plan.",
  "mood_level": 4,
  "energy_level": 5,
  "sleep_hours": 6.5,
  "study_hours": 8,
  "exam_countdown_days": 12
}
```