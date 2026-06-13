# Production hardening: accessibility, bilingual UI, deployment, and CI/CD

Type: AFK
Blocked by: 1. Secure daily journal and mood capture; 2. Detect hidden stress triggers; 3. Generate safe coping guidance; 4. Show mood and stress progress trends
User stories covered: Students can use the app in English or Hindi, access it with assistive technologies, and rely on a deployed production build that is tested before release.

## What to build

Finish the product with the cross-cutting work needed for public use: multilingual interface support, WCAG-minded accessibility, GitHub Actions validation, Docker packaging, secrets handling, and deployment wiring for backend and frontend hosting.

## Acceptance criteria

- [x] Core user-facing screens support English and Hindi content.
- [x] Keyboard navigation, ARIA roles, and screen-reader output are verified on the main flows.
- [x] GitHub Actions runs tests and linting before deployment.
- [x] The backend is containerized and deployment-ready for Render or Koyeb.
- [x] The frontend is deployment-ready for Netlify or Vercel with environment-based backend configuration.
- [x] Sensitive settings are provided through environment variables rather than hard-coded values.
