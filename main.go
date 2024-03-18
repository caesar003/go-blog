package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/caesar003/go-blog/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Config struct {
	DBUser string `json:"dbUser"`
	DBPass string `json:"dbPass"`
	DBName string `json:"dbName"`
	DBHost string `json:"dbHost"`
	DBPort string `json:"dbPort"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

const dbDriver = "mysql"

func openDB() (*sql.DB, error) {
	config, _ := loadConfig("./credentials.json")
	dbUser := config.DBUser
	dbPass := config.DBPass
	dbName := config.DBName
	dbHost := config.DBHost
	dbPort := config.DBPort

	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?parseTime=true`, dbUser, dbPass, dbHost, dbPort, dbName)
	// db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?parseTime=true")
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		panic(err.Error())
	}

	return db, nil
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/user", getAllUsersHandler).Methods("GET")
	r.HandleFunc("/api/user", createUserHandler).Methods("POST")
	r.HandleFunc("/api/user/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/api/user/{id}", updateUserHandler).Methods("PUT")
	r.HandleFunc("/api/user/{id}", deleteUserHandler).Methods("DELETE")

	r.HandleFunc("/api/post", getAllPostHandler).Methods("GET")
	r.HandleFunc("/api/post", createPostHandler).Methods("POST")
	r.HandleFunc("/api/post/{id}", getPostHandler).Methods("GET")
	r.HandleFunc("/api/post/{id}", updatePostHandler).Methods("PUT")
	r.HandleFunc("/api/post/{id}", deletePostHandler).Methods("DELETE")
	r.HandleFunc("/api/post/{id}/comment", getPostCommentHandler).Methods("GET") // list all comments

	r.HandleFunc("/api/comment", createCommentHandler).Methods("POST")
	r.HandleFunc("/api/comment/{id}", updateCommentHandler).Methods("PUT")
	r.HandleFunc("/api/comment/{id}", deleteCommentHandler).Methods("DELETE")

	log.Println("Server listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()
	var user handlers.User
	json.NewDecoder(r.Body).Decode(&user)

	err = handlers.CreateUser(db, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User created successfully")
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	users, err := handlers.GetAllUsers(db)
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()
	vars := mux.Vars(r)
	idStr := vars["id"]
	userID, err := strconv.Atoi(idStr)
	user, err := handlers.GetUser(db, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)

	var user handlers.User
	err = json.NewDecoder(r.Body).Decode(&user)

	handlers.UpdateUser(db, userID, user.Name, user.Email)
	if err != nil {
		http.Error(w, "User not found!", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "User updated successfully")
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()
	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	err = handlers.DeleteUser(db, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "User deleted successfully")

	// Convert the user object to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	var post handlers.Post
	json.NewDecoder(r.Body).Decode(&post)

	err = handlers.CreatePost(db, post.Title, post.Content, post.UserID)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Post created successfully")
}

func getAllPostHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	posts, err := handlers.GetAllPosts(db)
	if err != nil {
		http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := openDB()
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	postID, err := strconv.Atoi(idStr)

	post, err := handlers.GetPost(db, postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	postID, err := strconv.Atoi(idStr)

	var post handlers.Post
	err = json.NewDecoder(r.Body).Decode(&post)

	err = handlers.UpdatePost(db, postID, post.Title, post.Content)
	if err != nil {
		http.Error(w, "Post not found!", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Post updated successfully")
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	err = handlers.DeletePost(db, postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Post deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}

func createCommentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()
	// parse JSON data from the request body
	var comment handlers.Comment
	json.NewDecoder(r.Body).Decode(&comment)

	err = handlers.CreateComment(db, comment.PostID, comment.UserID, comment.Content, comment.ReplyTo)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "comment created successfully")
}

func getPostCommentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]
	postId, err := strconv.Atoi(idStr)
	users, err := handlers.GetPostComments(db, postId)
	if err != nil {
		http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	commentID, err := strconv.Atoi(idStr)

	var comment handlers.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)

	err = handlers.UpdateComment(db, commentID, comment.Content)
	if err != nil {
		http.Error(w, "Comment not found!", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Comment updated successfully")
}

func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	err = handlers.DeleteComment(db, commentID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Comment deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
