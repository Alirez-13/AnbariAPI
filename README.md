# AnbariAPI 📦

> A robust, strictly-architected Modular Monolith backend for a Warehouse Management System, built in Go.

AnbariAPI is a high-performance inventory management API designed to handle complex warehouse operations. Engineered with a strong emphasis on **Domain-Driven Design (DDD)** and **Clean Architecture**, the codebase is divided into isolated bounded contexts. This ensures that the system remains scalable, maintainable, and highly resilient to schema and business logic changes over time.

## ✨ Core Domains & Features

The system is structurally divided into three independent modules:

* **📦 Inventory Management (`internal/inventory/`)**
    * **Immutable Transactions:** Processes inbound and outbound stock movements via an immutable transaction ledger (`Transaction`, `TransactionLine`).
    * **Batch Tracking:** Precise tracking of stock through `InventoryBatch` and `BatchAllocation` models.
    * **Atomic Operations:** Safe data handling using a dedicated GORM transaction runner.
* **🗂️ Catalog Management (`internal/catalog/`)**
    * Handles hierarchical product categorization and product metadata.
    * Full CRUD operations with strict DTO validations at the application layer.
* **🔐 Authentication & Sessions (`internal/auth/`)**
    * Secure user registration and login using bcrypt password hashing.
    * Stateful session management and validation mapped to database-backed session models.

## 🛠️ Tech Stack

* **Language:** Go (Golang)
* **Web Framework:** Gin (Handling API routing and CORS)
* **ORM:** GORM (With custom repository implementations per domain)
* **Database:** SQLite (Configured with WAL mode for high concurrency)
* **Architecture Pattern:** Domain-Driven Design (DDD), Modular Monolith, Hexagonal Architecture.

## 📂 Architecture & Project Structure

The project strictly follows DDD principles, avoiding technical grouping in favor of domain grouping. Each bounded context under `internal/` maintains its own layers:

```text
AnbariAPI/
├── api/             
│   └── routes/                 # Global HTTP routing aggregation
├── internal/                   # Core Bounded Contexts
│   ├── auth/                   # Identity and Session Management
│   ├── catalog/                # Product and Category Management
│   └── inventory/              # Core Warehouse Logic
│       ├── application/        # Use cases (Process transactions, validate lines)
│       ├── domain/             # Core business models and interface ports
│       ├── infrastructure/     # GORM repositories and database adapters
│       └── interfaces/         # Delivery mechanisms (HTTP handlers)
├── shared/                     # Cross-cutting concerns
│   ├── config/                 # Global configurations (CORS, etc.)
│   ├── database/               # DB connection pooling
│   └── migration/              # Auto-migration utilities
├── data/                       # Local SQLite database files
├── main.go                     # Application bootstrap
└── go.mod / go.sum             # Go module dependencies



## 🚀 Getting Started

### Prerequisites
* [Go](https://golang.org/doc/install) (1.20+ recommended)
* SQLite3

### Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Alirez-13/AnbariAPI.git
   cd AnbariAPI
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the application:**
   Because AnbariAPI uses SQLite locally, there is no need to set up an external database server to test the API. The auto-migration engine will automatically prepare the schema on startup.
   ```bash
   go run main.go
   ```

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

## 📄 License

This project is licensed under the MIT License.
