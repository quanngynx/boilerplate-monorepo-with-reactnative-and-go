# boilerplate-monorepo-with-reactnative-and-go

This project was created with [Better-T-Stack](https://github.com/AmanVarshney01/create-better-t-stack), a modern TypeScript stack that combines React, TanStack Router, and more.

## Features

- **TypeScript** - For type safety and improved developer experience
- **TanStack Router** - File-based routing with full type safety
- **React Native** - Build mobile apps using React
- **Expo** - Tools for React Native development
- **TailwindCSS** - Utility-first CSS for rapid UI development
- **Shared UI package** - shadcn/ui primitives live in `packages/ui`
- **Turborepo** - Optimized monorepo build system
- **Go Gin API** - Backend server with GORM, Postgres, Redis, and Docker support

## Getting Started

First, install the dependencies:

```bash
pnpm install
```

### Backend (Go Gin)

Start Postgres and Redis:

```bash
pnpm --filter server dev:deps
```

Copy the server environment file and start the API:

```bash
cp apps/server/.env.example apps/server/.env
pnpm dev:server
```

Or run the full Docker stack (Postgres, Redis, and API):

```bash
pnpm dev:server:docker
```

The API listens on [http://localhost:8080](http://localhost:8080).

### Frontend

Then, run the development servers:

```bash
pnpm run dev
```

Open [http://localhost:3001](http://localhost:3001) in your browser to see the web application. Use the Expo Go app to run the mobile application.

For native development, copy the Expo environment file:

```bash
cp apps/native/.env.example apps/native/.env
```

## UI Customization

React web apps in this stack share shadcn/ui primitives through `packages/ui`.

- Change design tokens and global styles in `packages/ui/src/styles/globals.css`
- Update shared primitives in `packages/ui/src/components/*`
- Adjust shadcn aliases or style config in `packages/ui/components.json` and `apps/web/components.json`

### Add more shared components

Run this from the project root to add more primitives to the shared UI package:

```bash
npx shadcn@latest add accordion dialog popover sheet table -c packages/ui
```

Import shared components like this:

```tsx
import { Button } from "@boilerplate-monorepo-with-reactnative-and-go/ui/components/button";
```

### Add app-specific blocks

If you want to add app-specific blocks instead of shared primitives, run the shadcn CLI from `apps/web`.

## Deployment (Cloudflare via Alchemy)

- Target: web
- Dev: pnpm run dev
- Deploy: pnpm run deploy
- Destroy: pnpm run destroy

For more details, see the guide on [Deploying to Cloudflare with Alchemy](https://www.better-t-stack.dev/docs/guides/cloudflare-alchemy).

## Project Structure

```
boilerplate-monorepo-with-reactnative-and-go/
├── apps/
│   ├── web/         # Frontend application (React + TanStack Router)
│   ├── native/      # Mobile application (React Native, Expo)
│   └── server/      # Backend API (Go Gin + GORM + Redis)
├── packages/
│   ├── ui/          # Shared shadcn/ui components and styles
│   └── env/         # Shared environment validation
```

## Server Environment Variables

Configure `apps/server/.env` using `apps/server/.env.example`:

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8080` | API bind port |
| `GIN_MODE` | `debug` | Gin mode (`release` in production) |
| `ALLOWED_ORIGINS` | `http://localhost:3001` | Comma-separated CORS origins |
| `DATABASE_URL` | — | Postgres connection string |
| `REDIS_URL` | — | Redis connection string |
| `LOG_LEVEL` | `info` | Structured log level |
| `APP_NAME` | `server` | Application name in logs |

## Available Scripts

- `pnpm run dev`: Start all applications in development mode
- `pnpm run build`: Build all applications
- `pnpm run dev:web`: Start only the web application
- `pnpm run dev:server`: Start only the Go API (requires Postgres + Redis)
- `pnpm run dev:server:docker`: Start Postgres, Redis, and API via Docker Compose
- `pnpm run check-types`: Check TypeScript types across all apps
- `pnpm run dev:native`: Start the React Native/Expo development server
