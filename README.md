# Go Ecommerce Backend

I built this to understand how a production backend actually comes together — not just writing APIs, but going through the full lifecycle: architecture decisions, authentication, payments, deployment, and CI/CD.

The backend is live on AWS. It handles real requests, connects to a managed PostgreSQL instance on RDS, and processes payments through Stripe.

**Live API:** http://go-ecommerce-app.ap-south-1.elasticbeanstalk.com/


## What it does

- User registration, login, and JWT-based authentication
- Product and category management
- Cart creation and order placement
- Stripe payment processing
- Role-based access control on protected routes
- Fully containerized and deployed on AWS with a CI pipeline


## Architecture

Modular monolith — not microservices. Each domain (auth, products, orders, payments) owns its routes, services, and data access logic, but they all run as one deployable unit.

I picked this over microservices intentionally. For a project at this scale, splitting into services would add operational overhead without real benefit. The module boundaries are clean enough that breaking things apart later would be straightforward if needed.

```
Client Request
    → Fiber HTTP server
    → Middleware (auth, validation)
    → Domain handler (Auth / Products / Orders / Payments)
    → Service layer (business logic)
    → Repository layer (database access)
    → PostgreSQL / Stripe
    → JSON Response
```


## Tech stack

**Backend:** Go, Fiber, GORM, PostgreSQL, JWT, Stripe API

**Infrastructure:** AWS Elastic Beanstalk, EC2, RDS, IAM, Docker, GitHub Actions


## Authentication

Standard JWT flow — user logs in, gets a signed token, sends it on protected requests. Middleware validates the token and rejects anything invalid before it hits the handler. Passwords are hashed with bcrypt.

Roles are checked at the middleware level so business logic stays clean.


## Payments (Stripe)

1. Client initiates checkout
2. Backend creates a Stripe payment intent
3. Client completes payment on the Stripe side
4. Backend verifies payment status
5. Order is confirmed or marked as failed

The async nature of payment verification was the most interesting part to get right.


## CI/CD

GitHub Actions runs on every push — builds the binary, runs checks, and blocks merges if the build fails. Deployment to Elastic Beanstalk is handled separately after passing CI.


## Running locally

```bash
git clone https://github.com/Ardent-7322/go-eccomerce-app
cd go-eccomerce-app
cp .env.example .env   # fill in your values
docker-compose up --build
```

The server starts on port 9000 by default.


## Environment variables

```env
APP_ENV=prod
SERVER_PORT=9000
APP_SECRET=
DB_HOST=
DB_PORT=5432
DB_NAME=
DB_USER=
DB_PASSWORD=
STRIPE_SECRET_KEY=
```
