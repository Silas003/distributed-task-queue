# Distributed Task Queue (DTQ)

A lightweight distributed task queue system written in Go, designed for reliability and horizontal scaling.
It is aimed at sending of emails.

## Features

- ✅ Distributed task processing  
- 🔄 Automatic retries for failed tasks  
- 🔒 Redis-backed persistence  
 

## Installation

```bash
git clone https://github.com/yourusername/distributed-task-queue.git
cd distributed-task-queue
go mod download
add .env with the email configurations to google'/s smtp
go run .
```
