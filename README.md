# Integrand

<div align="center">
<p align="center">
    
![LucidMQ](https://user-images.githubusercontent.com/25624274/218341069-514ac1ec-0a06-4260-a229-c047dd531ac2.png)

**Simple-Ops Event Streaming. Build your real time applications without the headache of ops overhead.**

<a href="https://integrand.io">Website</a>
<a href="https://integrand.io/docs/">Documentation</a> •
    
![CI](https://github.com/lucidmq/lucidmq/actions/workflows/.github/workflows/integrand-app.yml/badge.svg)
![Last Commit](https://img.shields.io/github/last-commit/integrand/integrand-app)
![Github Stars](https://img.shields.io/github/stars/integrand/integrand-app)
![Github Issues](https://img.shields.io/github/issues/integrand/integrand-app)

</p>
</div>

**Simple-Ops Webhook Streaming. Build your real time integrations without the headache of infrastructure.**


> :warning: **This project is in Alpha Stage**: Expect breaking changes

## What is Integrand?

Integrand is an application infrastructure that focuses on handling webhooks and providing easy to use API's that provide a streaming interface. It enables the creation of stream or queue based applications by providing a rock solid foundation and simple API's.

### Repo Structure

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
      registry.vineglue.com/integrand/integrand:latest
```

### 🐙 Docker Compose:

```yaml
services:
  openobserve:
    image: registry.vineglue.com/integrand/integrand:latest
    restart: unless-stopped
    environment:
      ZO_ROOT_USER_EMAIL: "root@example.com"
      ZO_ROOT_USER_PASSWORD: "Complexpass#123"
    ports:
      - "5080:5080"
    volumes:
      - data:/data
volumes:
  data:
```