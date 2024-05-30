# Integrand

<div align="center">
<p align="center">
    
![LucidMQ](https://integrand.io/static/images/logos/Integrand-logo.svg)

**Simple-Ops Webhook Streaming. Build your real time integrations without the headache of infrastructure.**

<a href="https://integrand.io">Website</a> •
<a href="https://integrand.io/docs/">Documentation</a> 
    
![CI](https://github.com/integrandio/integrand-app/actions/workflows/.github/workflows/integrand-app.yml/badge.svg)
![Last Commit](https://img.shields.io/github/last-commit/integrandio/integrand-app)
![Github Stars](https://img.shields.io/github/stars/integrandio/integrand-app)
![Github Issues](https://img.shields.io/github/issues/integrandio/integrand-app)

</p>
> :warning: **This project is in Alpha Stage**: Expect breaking changes
</div>

## What is Integrand?

Integrand is an infrastructure application that focuses on handling webhooks and providing easy to use API's that provide a streaming interface. It enables the creation of stream or queue based applications by providing a rock solid foundation and simple API's.

### Repo Structure

Integrand is a monolith application that contains both backend and frontend making up the application.

    ├── commitlog      # The base library containing code for the commitlog
    ├── persistence    # Contains code for interacting with persistant state
    ├── services       # Wrappers around persistence and other services
    ├── Web            # Api and web client code
    ├── utils          # Utils used across the application
    └── data           # Where our persistent data is stored, including SQL definitions

### Goals of Integrand
- Allow integration developers to focus on building robust integrations, not deal with web hook and underlying server infrastructure.

- Developer first. Platforms like Zapier and Robot are great, except when complexity hits. Sometimes you're going to need to create custom code that doesn't need restrictions of another platform.

- Easy to operate as Integrand is self contained. Don't worry about viewing deployment instructions and having to install multiple databases just to get it running. 

## ⚡️ Quick start

### 🐳 Docker
```bash
docker run -d \
      --name integrand \
      -v $PWD/data:/data \
      -p 8000:8000 \
      -e ROOT_EMAIL="test@example.com" \
      -e ROOT_PASSWORD="MyPassword" \
      registry.vineglue.com/integrand-app:latest
```

### 🐙 Docker Compose:
```yaml
services:
  openobserve:
    image: registry.vineglue.com/integrand-app:latest
    restart: unless-stopped
    environment:
      ROOT_EMAIL: "root@example.com"
      ROOT_PASSWORD: "MyPassword"
    ports:
      - "8000:8000"
    volumes:
      - data:/data
volumes:
  data:
```