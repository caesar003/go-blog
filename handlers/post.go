package handlers

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PostWithAuthor struct {
	Post
	Author
	Comments int `json:"comments"`
}

func CreatePost(db *sql.DB, title string, content string, userId int) error {
	query := "INSERT INTO posts (title, content, user_id) VALUES(?, ?, ?)"
	_, err := db.Exec(query, title, content, userId)
	if err != nil {
		return err
	}
	return nil
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
	query := "SELECT * FROM posts"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPost(db *sql.DB, id int) (*PostWithAuthor, error) {
	query := `
	SELECT 
		posts.id,
		posts.title,
		posts.content,
		posts.user_id,
		posts.created_at,
		posts.updated_at,
		users.name,
		users.email,
		(
			SELECT COUNT(*)
			FROM comments
			WHERE post_id = posts.id
	)
	FROM 
	posts 
	INNER JOIN
	users
	ON
	posts.user_id = users.id
	WHERE 
	posts.id = ?`
	row := db.QueryRow(query, id)

	post := &PostWithAuthor{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.Name, &post.Email, &post.Comments)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return post, nil
}

func UpdatePost(db *sql.DB, id int, title, content string) error {
	query := "UPDATE posts SET title = ?, content = ? WHERE id = ?"
	_, err := db.Exec(query, title, content, id)
	if err != nil {
		return err
	}
	return nil
}

func DeletePost(db *sql.DB, id int) error {
	query := "DELETE FROM posts WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
