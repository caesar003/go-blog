# GO-BLOG

Yeah, you read it right. It's a blog written in go.


Go-Blog is a simple server application written in Go that allows users to write articles and comment on them. It provides basic CRUD functionalities for users, posts, and comments via RESTful APIs.

## Pre-requisite

1. Ensure you have Go installed on your system. You can download it from the [official website](https://go.dev).
2. An SQL Server.

## Installation

2. Clone the repository:
```shell
git clone https://github.com/caesar003/go-blog.git
cd goblog
```

3. Install dependencies:
```shell
go mod tidy
```
4. Set up MySQL database by running the script provided:
```shell
mysql -u username -p < db.sql
```
## Configuration

**Copy `credentials-copy.json` to `credentials.json`:**

```shell
cp credentials-copy.json credentials.json
```
    
and modify `credentials.json` with you database connection details.

## Usage

Run the server:
```shell
go run main.go
```

The server will start running at `http://localhost/8090`

## Endpoints

- **User Endpoints**:
    - **GET /api/user**: Get all users.
    - **POST /api/user**: Create a new user.
    - **GET /api/user/{id}**: Get a user by ID.
    - **PUT /api/user/{id}**: Update a user by ID.
    - **DELETE /api/user/{id}**: Delete a user by ID.
- **Post Endpoints**:
    - **GET /api/post**: Get all posts.
    - **POST /api/post**: Create a new post.
    - **GET /api/post/{id}**: Get a post by ID.
    - **PUT /api/post/{id}**: Update a post by ID.
    - **DELETE /api/post/{id}**: Delete a post by ID.
    - **GET /api/post/{id}/comment**: Get all comments for a post.
- **Comment Endpoints**:
    - **POST /api/comment**: Create a new comment.
    - **PUT /api/comment/{id}**: Update a comment by ID.
    - **DELETE /api/comment/{id}**: Delete a comment by ID.

## Author
that's me: [caesar003](https://github.com/caesar003)


## Acknowledgements
- Salem Olorundare for the insightful blog post ["How to Create a CRUD Application with GoLang and MySQL"](https://www.honeybadger.io/blog/how-to-create-crud-application-with-golang-and-mysql/), which served as the foundation for this project
- Thanks to [Go](https://go.dev) for the awesome programming language
