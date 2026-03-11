# Go E-Commerce Backend

## Overview

This is a personal learning project where I built a complete e-commerce backend in Go to understand how real-world backend systems are structured, deployed, and operated in production-like environments.

The goal of this project was not just to expose APIs, but to work through the full lifecycle of a backend application — authentication, domain separation, payments, configuration management, deployment, and CI/CD — while keeping the architecture understandable and maintainable.

Although the system is deployed and functional, the primary focus was learning solid backend engineering practices rather than building a commercial product.

🔗 Live Backend API  
http://go-ecommerce-app.ap-south-1.elasticbeanstalk.com/

🔗 Sample Frontend (used for testing APIs end-to-end)  
https://example-frontend-demo-link.com

---

## What This Backend Does

- Handles user signup, login, and authentication using JWT
- Manages products and categories
- Supports cart creation and order placement
- Integrates Stripe for secure online payments
- Exposes RESTful APIs for frontend or client consumption
- Uses environment-based configuration for different setups
- Is deployed on AWS with automated CI/CD

---

## How It Works

High-level flow:

Client Request  
→ API Gateway (Fiber HTTP server)  
→ Domain Module (Auth / Products / Orders / Payments)  
→ Database or External Service (PostgreSQL / Stripe)  
→ JSON Response

Key ideas:
- Each domain is isolated into its own module
- Business logic is kept separate from HTTP handlers
- Middleware is used for authentication and request validation
- Configuration is injected via environment variables
- The application runs as a single deployable service

---

## Architecture

This project follows a **modular monolithic architecture**.

Instead of splitting everything into microservices, I chose a modular monolith to:
- Keep deployment and debugging simple
- Maintain clear domain boundaries
- Avoid unnecessary operational complexity

### Core Modules

- Auth / Users  
- Products & Categories  
- Cart & Orders  
- Payments  
- Configuration & Middleware  

Each module owns its routes, services, and data access logic, while sharing common infrastructure like logging, configuration, and middleware.

This structure also makes it easier to break modules into separate services in the future if required.

---

## Design Choices

- I chose a modular monolith to focus on clean domain separation without microservice overhead
- Fiber was selected for its simplicity and performance
- GORM was used to speed up development and schema management
- JWT-based authentication was implemented to understand stateless auth flows
- Stripe was integrated to learn real payment workflows and webhooks
- AWS Elastic Beanstalk was used to simplify deployment while still working with EC2 and RDS

---

## Tech Stack

### Backend
- Go (Golang)
- Fiber (HTTP framework)
- GORM (ORM)
- PostgreSQL
- JWT for authentication
- Stripe API for payments

### Infrastructure & DevOps
- AWS Elastic Beanstalk
- EC2
- RDS (PostgreSQL)
- IAM
- Docker
- GitHub Actions for CI/CD
- GitHub for version control

---

## Authentication Flow

- Users register or log in via API
- Backend validates credentials
- A JWT is issued on successful authentication
- Protected routes use middleware to validate tokens
- Requests without valid tokens are rejected

This helped me understand how stateless authentication works in real systems.

---

## Payment Flow (Stripe)

1. Client initiates checkout
2. Backend creates a Stripe payment intent
3. Client completes payment via Stripe Checkout
4. Backend verifies payment status
5. Order status is updated accordingly

This part of the project was especially useful for understanding third-party integrations and handling asynchronous payment flows.

---

## CI/CD Pipeline

The project uses GitHub Actions to automate basic checks.

Pipeline steps:
- Triggered on push and pull requests
- Runs build and basic validation
- Ensures the application can be built successfully before deployment

This helped me understand how automated pipelines fit into real development workflows.

---

## Environment Configuration

The application is configured entirely through environment variables.

```env
APP_ENV=prod
SERVER_PORT=9000
APP_SECRET=*****
DB_HOST=*****
DB_PORT=5432
DB_NAME=*****
DB_USER=*****
DB_PASSWORD=*****
STRIPE_SECRET_KEY=*****
