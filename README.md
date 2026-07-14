# Social Feed Backend

A scalable RESTful backend for a social media application built with **Golang**, **PostgreSQL**, and **Clean Architecture**. The API provides secure authentication, post management, comments, replies, likes, image uploads, and public/private post visibility.

## ✨ Features

- 🔐 JWT Authentication & Authorization
- 👤 User Registration & Login
- 📝 Create, Update & Delete Posts
- 🖼️ Image Upload Support
- 🌍 Public & Private Posts
- ❤️ Like & Unlike Posts
- 💬 Comment System
- ↩️ Nested Replies
- 👍 Like & Unlike Comments/Replies
- 📋 View Users Who Liked Posts & Comments
- ✅ Input Validation
- 🛡️ Secure Password Hashing (bcrypt)
- 🏗️ Clean Architecture
- 🚀 PostgreSQL Integration

---

## 🛠️ Tech Stack

- **Language:** Golang
- **Router:** Chi
- **Database:** PostgreSQL
- **Authentication:** JWT
- **Password Hashing:** bcrypt
- **Validation:** go-playground/validator
- **Architecture:** Clean Architecture

---

## 📁 Project Structure

```text
.
├── cmd/
│   └── api/
│       ├── main.go
│       ├── routes.go
│       └── ...
├── internal/
│   ├── auth/
│   ├── db/
│   ├── mail/
│   ├── middleware/
│   ├── store/
│   ├── validator/
│   └── ...
├── uploads/
├── migrations/
├── .env
├── go.mod
└── README.md
```

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/anikmahidul9/social-backend.git
```

```bash
cd social-backend
```

---

### 2. Install Dependencies

```bash
go mod tidy
```

---

### 3. Configure Environment Variables

Create a `.env` file in the project root.

```env
ADDR=:8080

DB_ADDR=postgres://postgres:password@localhost:5432/social?sslmode=disable

JWT_SECRET=my-super-secret-key-123456
```

---

### 4. Create PostgreSQL Database

```sql
CREATE DATABASE social;
```

Run your database migrations (if applicable) before starting the server.

---

### 5. Run the Server

```bash
go run cmd/api/*.go
```

The API will start at:

```
http://localhost:8080
```

---

## 📌 API Features

### Authentication

- Register
- Login
- JWT Authentication
- Protected Routes

### Posts

- Create Post
- Update Post
- Delete Post
- Public & Private Visibility
- Upload Images
- Feed Pagination

### Comments

- Add Comment
- Edit Comment
- Delete Comment
- Nested Replies

### Likes

- Like/Unlike Posts
- Like/Unlike Comments
- View Users Who Liked

---

## 🗄️ Database

PostgreSQL is used as the primary database with a normalized relational schema.

### Main Tables

- users
- posts
- post_images
- comments
- post_likes
- comment_likes

---

## 🔒 Security

- JWT Authentication
- Password Hashing with bcrypt
- Request Validation
- Authorization Middleware
- SQL Injection Protection (Parameterized Queries)
- Protected Endpoints

---

## 🏗️ Architecture

The project follows **Clean Architecture**, separating business logic from infrastructure and delivery layers.

```text
HTTP Request
      │
      ▼
Router (Chi)
      │
      ▼
Middleware
      │
      ▼
HTTP Handlers
      │
      ▼
Business Logic
      │
      ▼
Repository Layer
      │
      ▼
PostgreSQL
```

---

## 📮 Example API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/auth/register` | Register a new user |
| POST | `/v1/auth/login` | User login |
| GET | `/v1/posts` | Get feed |
| POST | `/v1/posts` | Create post |
| PATCH | `/v1/posts/{id}` | Update post |
| DELETE | `/v1/posts/{id}` | Delete post |
| POST | `/v1/posts/{id}/likes` | Like a post |
| DELETE | `/v1/posts/{id}/likes` | Unlike a post |
| POST | `/v1/posts/{id}/comments` | Add a comment |
| POST | `/v1/comments/{id}/replies` | Add a reply |
| POST | `/v1/comments/{id}/likes` | Like a comment |

---

## 🧪 Testing

Run all tests:

```bash
go test ./...
```

---

## 🤝 Contributing

1. Fork the repository.
2. Create a new branch.

```bash
git checkout -b feature/new-feature
```

3. Commit your changes.

```bash
git commit -m "Add new feature"
```

4. Push to your branch.

```bash
git push origin feature/new-feature
```

5. Open a Pull Request.

---

## 📄 License

This project is licensed under the MIT License.

---

## 👨‍💻 Author

**Mahidul Anik**

- GitHub: https://github.com/anikmahidul9

---

⭐ If you find this project useful, consider giving it a star on GitHub!
