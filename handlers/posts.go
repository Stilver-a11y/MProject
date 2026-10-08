package handlers

import (
	"encoding/json"
	models "mproject/Models"
	"mproject/database"
	"mproject/middleware"
	"net/http"
	"strings"
)

type PostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Media       string `json:"media"`
}

// CreatePost -- POST /posts
func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	var post models.Post
	err := database.Db.QueryRow(
		r.Context(),
		"INSERT INTO posts (title, description, media, user_id) VALUES ($1, $2, $3, $4) RETURNING id, title, description, media, user_id, created_at",
		req.Title, req.Description, req.Media, userID,
	).Scan(&post.ID, &post.Title, &post.Description, &post.Media, &post.UserID, &post.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// GetPost - GET /posts
func GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.Db.Query(
		r.Context(),
		"SELECT id, title, description, media, user_id, created_at FROM posts ORDER BY created_at DESC",
	)
	if err != nil {
		http.Error(w, "Failed to get posts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Media, &p.UserID, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// UpdatePost - PUT /posts/:id
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(uint)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID := parts[2]

	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Checking that the post belongs to this user
	var ownerID uint
	err := database.Db.QueryRow(r.Context(), "SELECT user_id FROM posts WHERE id = $1", postID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "You can only edit your own post", http.StatusForbidden)
		return
	}

	_, err = database.Db.Exec(
		r.Context(),
		"UPDATE posts SET title = $1, description = $2, media = $3 WHERE id = $4",
		req.Title, req.Description, req.Media, postID,
	)

	if err != nil {
		http.Error(w, "Failed to update post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})
}

// DeletePost - DELETE /posts/:id
func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	postID := parts[2]

	var ownerID uint
	err := database.Db.QueryRow(r.Context(), "SELECT user_id FROM posts WHERE id = $1", postID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

	_, err = database.Db.Exec(r.Context(), "DELETE FROM posts WHERE id = $1", postID)
	if err != nil {
		http.Error(w, "Failed to delete post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}
