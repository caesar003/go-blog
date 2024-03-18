package handlers

import (
	"database/sql"
	"fmt"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ReplyTo   *int      `json:"reply_to"`
}

func CreateComment(db *sql.DB, postID int, userID int, content string, replyTo *int) error {
	query := "INSERT INTO comments (post_id, user_id, content, reply_to) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, postID, userID, content, replyTo)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func GetPostComments(db *sql.DB, postID int) ([]Comment, error) {
	query := "SELECT * FROM comments where post_id = ?"
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.ReplyTo); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func DeleteComment(db *sql.DB, id int) error {
	query := "DELETE FROM comments WHERE id = ?"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateComment(db *sql.DB, id int, content string) error {
	query := "UPDATE comments SET content = ? WHERE id = ?"
	_, err := db.Exec(query, content, id)
	if err != nil {
		return err
	}

	return nil
}
