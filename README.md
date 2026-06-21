# TEUS (t's)

:contentReference[oaicite:0]{index=0}

TEUS (pronounced **"t's"**) is a lightweight **Go foundation module** that provides shared tools and utilities for backend development.

It is designed as a reusable **config and infrastructure toolkit** used across multiple projects such as ERP, MES, and microservices.

---

## 📦 Purpose

TEUS provides commonly used building blocks for backend systems:

- Configuration loader
- Database helper (SQLX ready)
- Logger utilities
- HTTP response helpers
- Error definitions
- General utilities

Instead of rewriting these components in every project, TEUS centralizes them into a single reusable module.

---

## 🧠 Philosophy

TEUS is NOT a framework and does NOT contain business logic.

It is a **foundation toolkit** that sits below all applications.

> "Build once, reuse everywhere."

---

## 📁 Packages

- `config` → environment & configuration loader  
- `database` → PostgreSQL / SQLX helper  
- `logger` → simple logging utilities  
- `http` → HTTP response helpers  
- `utils` → helper functions  
- `errors` → shared error definitions  

---

## 🚀 Installation

```bash
go get github.com/eonebyte/go-teus