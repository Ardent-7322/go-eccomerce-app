# Go E-Commerce Backend

A modular monolithic e-commerce backend built with **Go**, designed for real-world production use.  
The application supports authentication, product management, cart & orders, payments via Stripe, and is deployed on **AWS Elastic Beanstalk**.

🔗 Live Backend URL:  
http://go-ecommerce-app.ap-south-1.elasticbeanstalk.com/


🔗 Live URL:  
https://example-frontend-demo-link.com



---

## Features

- User authentication & authorization (JWT based)
- Product & category management
- Cart and order management
- Secure payment integration using **Stripe**
- Modular monolithic architecture (clean separation of concerns)
- RESTful APIs
- Environment-based configuration
- Production deployment on AWS Elastic Beanstalk

---

## Architecture

This project follows a **modular monolithic architecture**, where each domain is isolated into its own module while running as a single deployable unit.

**Core Modules**
- Auth / Users
- Products & Categories
- Cart & Orders
- Payments
- Configuration & Middleware

The architecture allows:
- Clear domain boundaries
- Easy migration to microservices in the future
- Simple deployment & scaling

_(Architecture diagrams included in `/docs` folder)_

---

## Tech Stack

### Backend
- **Go (Golang)**
- **Fiber** (web framework)
- **GORM** (ORM)
- **PostgreSQL**
- **Stripe API** (payments)
- **JWT** for authentication

### Infrastructure & DevOps
- **AWS Elastic Beanstalk**
- **EC2**
- **RDS (PostgreSQL)**
- **IAM**
- **Environment variables for config**
- **GitHub for version control**

---

## Authentication Flow

- Users authenticate via login/signup
- JWT token is issued on successful login
- Token is validated on protected routes using middleware

---

## Payment Flow (Stripe)

1. User creates an order
2. Backend creates a Stripe payment intent
3. Client completes payment using Stripe Checkout
4. Backend verifies payment status and updates order state

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
