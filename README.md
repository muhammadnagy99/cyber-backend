# 🚀 Role-Based Access Control (RBAC) System in Golang

## 📌 Overview
This project is a **Role-Based Access Control (RBAC) system** built using **Golang**, **PebbleDB (Key-Value Store)**, and **Clean Architecture**. It provides secure and efficient **role, permission, and role-permission management**, following **SOLID principles**.

The system ensures **fine-grained access control** by allowing roles to be assigned permissions dynamically. It is designed for **scalability, security, and ease of integration** into any application requiring authentication and authorization.

---

## 🎯 Goals
- ✅ **Provide a flexible and scalable RBAC system** for API authentication and authorization.
- ✅ **Ensure high performance** with fast **key-value storage (PebbleDB)**.
- ✅ **Maintain code modularity** using **Clean Architecture** and **SOLID principles**.
- ✅ **Enable easy integration** into existing applications with a well-documented REST API.
- ✅ **Follow best security practices** with **Bearer Token authentication**.
- ✅ **Improve maintainability** with structured **logging** and meaningful error handling.

---

## 🔥 Key Features
- **🔑 Role Management** → Create, update, delete, and fetch roles.
- **🛡️ Permission Management** → Assign, update, delete, and fetch permissions.
- **🔗 Role-Permission Linking** → Assign and verify permissions for roles.
- **🔍 Check Access** → API to check if a role has a specific permission.
- **🔐 Secure Authentication** → Uses **Bearer Token authentication** to protect API endpoints.
- **🚀 High Performance** → Utilizes **PebbleDB**, a fast key-value storage for instant lookups.
- **🛠️ Clean Architecture** → Decoupled layers for maintainability and scalability.
- **📜 Structured Logging** → Detailed logs for debugging and security tracking.

---

## ⚡ Installation & Running the Project

### **Using Docker (Recommended)**
This project is **containerized** using **Docker and Docker Compose**, making it easy to deploy.

#### **1️⃣ Build and Run the Containers**
```sh
docker-compose up --build -d
```

#### **2️⃣ Check Running Containers**
```sh
docker ps
```

#### **3️⃣ View Application Logs**
```sh
docker-compose logs -f
```

#### **4️⃣ Stop the Containers**
```sh
docker-compose down
```

---

## 🔥 Running Locally (Without Docker)

If you prefer to run the application manually:

### **1️⃣ Clone the Repository**
```sh
git clone https://github.com/muhammadnagy99/gRBAC.git
cd cyber-rbac
```

### **2️⃣ Set Up Environment Variables**
Create a `.env` file and define:
```sh
SERVER_PORT=8080
DATABASE_PATH=rbac.db
BEARER_TOKEN=your-secret-token
```

### **3️⃣ Install Dependencies**
```sh
go mod tidy
```

### **4️⃣ Run the API Server**
```sh
go run cmd/main.go
```

---

🚀 **This RBAC system is a powerful demonstration of real-world Golang backend development, security, and performance optimization.**  
