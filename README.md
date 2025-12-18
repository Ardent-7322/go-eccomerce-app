# Go E-Commerce Backend

A modular monolithic e-commerce backend built with **Go**, designed for real-world production use.  
The application supports authentication, product management, cart & orders, payments via **Stripe**, and is deployed on **AWS Elastic Beanstalk**.

🔗 **Live Backend API**  
http://go-ecommerce-app.ap-south-1.elasticbeanstalk.com/

🔗 **Live link**  
https://example-frontend-demo-link.com

> A frontend page integrated with backend APIs and used for end-to-end testing.

---

## Features

- User authentication & authorization (JWT-based)
- Product & category management
- Cart and order management
- Secure payment integration using **Stripe**
- Modular monolithic architecture with clear separation of concerns
- RESTful APIs
- Environment-based configuration
- Production deployment on **AWS Elastic Beanstalk**
- Automated CI/CD using **GitHub Actions**

---

## Architecture

This project follows a **modular monolithic architecture**, where each domain is isolated into its own module while running as a single deployable service.

### Core Modules
- Auth / Users
- Products & Categories
- Cart & Orders
- Payments
- Configuration & Middleware

### Design Benefits
- Clear domain boundaries
- Easier maintenance and testing
- Straightforward migration path to microservices if needed
- Simple deployment and scaling

(Architecture diagrams are included in the `/docs` directory.)

---

## Tech Stack

### Backend
- **Go (Golang)**
- **Fiber** (HTTP framework)
- **GORM** (ORM)
- **PostgreSQL**
- **Stripe API** (payments)
- **JWT** for authentication & authorization

### Infrastructure & DevOps
- **AWS Elastic Beanstalk**
- **EC2**
- **RDS (PostgreSQL)**
- **IAM**
- **Docker**
- **GitHub Actions** (CI/CD)
- **GitHub** for version control

---

## CI/CD Pipeline (GitHub Actions)

The project uses **GitHub Actions** to automate validation of code changes.

### Pipeline Overview
- Triggered on pull requests and pushes
- Runs automated tests
- Validates build before deployment
- Ensures code quality and stability before merging

---

## Authentication Flow

- Users authenticate via login/signup
- Backend issues a JWT on successful authentication
- Protected routes are secured using middleware-based token validation

---

## Payment Flow (Stripe)

1. User initiates checkout from the client
2. Backend creates a Stripe payment intent
3. Client completes payment using Stripe Checkout
4. Backend verifies payment via Stripe and updates order status

---

## ⚙️ Environment Variables

```env
APP_ENV=prod
SERVER_PORT=5000
APP_SECRET=*****
DB_HOST=*****
DB_PORT=5432
DB_NAME=*****
DB_USER=*****
DB_PASSWORD=*****
STRIPE_SECRET_KEY=*****
