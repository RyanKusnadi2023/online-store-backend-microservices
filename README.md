# Online Store Backend - README

## Overview

This repository implements the backend for an online store application, focusing on a microservices architecture with an API Gateway as the single entry point. The backend is responsible for user management, product listing, cart operations, and order checkout, while the frontend (web or mobile) communicates solely through the API Gateway.

## Architecture
![Microservice](https://github.com/user-attachments/assets/86509564-337d-46b2-8c1e-092aab09aac2)

### API Gateway
The **API Gateway** serves as the primary interface for all client requests. Instead of the frontend interacting directly with each service, it communicates exclusively with the gateway. Key responsibilities of the API Gateway include:

- **Single Entry Point:** Frontend clients send all requests to the gateway, greatly simplifying the client-side logic.
- **Request Routing & Load Balancing:** The gateway identifies which microservice should handle the incoming request and forwards it accordingly.
- **Security & Authentication:** Common authentication logic, such as validating JWT tokens, can be applied at the gateway level to ensure that only authorized requests reach the backend services.
- **Cross-Cutting Concerns:** Implement functionalities like rate limiting, request/response logging, and response caching in one place rather than duplicating across services.

By centralizing these responsibilities, the API Gateway helps maintain a clean separation of concerns and reduces complexity in individual services.

### Microservices
The backend is decomposed into distinct microservices, each responsible for a specific domain. This design ensures better scalability, maintainability, and independence in deployment and development.

1. **User Service**  
   - **Functions:** Registration, login, and JWT authentication token handling.
   - **Data:** User credentials, profile information.
   - **Highlights:** Autonomy to evolve user authentication strategies independently.

2. **Product Service**  
   - **Functions:** Product listing, categorization, and search.
   - **Data:** Product catalogs, categories, metadata, and inventory details.
   - **Highlights:** Can leverage appropriate database technologies or indexing strategies to handle large product catalogs efficiently.

3. **Cart Service**  
   - **Functions:** Manages user carts—adding items, listing current cart contents, and removing items.
   - **Data:** Cart items keyed by user ID, product references, and temporary pricing.
   - **Highlights:** Decoupled from product and user logic, focusing solely on cart operations and state.

4. **Order Service**  
   - **Functions:** Handles checkout processes, payment integrations, and order creation.
   - **Data:** Order records, transactions, statuses, and payment confirmations.
   - **Highlights:** Able to focus on robust transactional integrity and payment workflows independently from other services.

### Databases
Each microservice maintains its own dedicated database, embodying the **Database per Service** pattern. This approach ensures loose coupling—services can choose their own data storage technologies and evolve their schemas independently.

- **User DB:** Stores user credentials, profiles, and authentication data.
- **Product DB:** Holds product catalogs, categories, and metadata.
- **Cart DB:** Saves cart items, keyed by user sessions or IDs.
- **Order DB:** Persists order details, statuses, and transaction logs.

This autonomy reduces inter-service dependencies and avoids centralized database bottlenecks.

## Benefits of This Architecture

- **Scalability:** Scale each service independently based on load patterns (e.g., scale the Product service during peak browsing times without affecting the Order service).
- **Resilience & Maintainability:** Isolate faults to a single service, improving system stability and simplifying troubleshooting.
- **Agility:** Update or swap out technologies in one service without impacting others. For example, switch the Cart DB from relational to a key-value store without affecting User or Product services.
- **Security & Observability:** The API Gateway centralizes crucial features like authentication and logging, making it easier to secure and monitor the entire system.

## Getting Started

1. **Clone the Repository:**  
   ```bash
   git clone https://github.com/yourusername/online-store-backend.git
   cd online-store-backend
