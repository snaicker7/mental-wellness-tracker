# Secure User Authentication & Login Page

Type: Feature
Blocked by: 6. Mobile-first polished mood dashboard
User stories covered: Students can sign up and log in securely to their personal mental wellness tracker using their credentials, protecting their confidential journal entries.

## What to build

Design and implement a clean, premium, theme-cohesive Login/Signup screen that serves as the entry gate for the application. Upon successful authentication, the user session is preserved (locally/token-based) and the student greeting dynamically displays their real name instead of falling back to a hardcoded string.

## Acceptance criteria

- [ ] A dedicated Login/Signup view is displayed to unauthenticated users.
- [ ] Users can toggle between Login and Sign Up forms.
- [ ] Sign Up allows creating a new account (username and password) via the Go backend.
- [ ] Login verifies the credentials via the Go backend and issues an authentication state.
- [ ] The authenticated username is set as the active `userId` and saved in local storage to keep the session alive across page refreshes.
- [ ] The settings panel inputs display clean focus state outlines, and the backend URL is properly configured.
- [ ] Login screen follows the playful wellness design system (soft background, rounded inputs, bold header typography, and micro-interactions).
