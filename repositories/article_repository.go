package repositories

import (
	"article-api/db"
	"article-api/models"
	"log"
)

func CreateArticle(article models.Article) (int64, error) {
	query := `
		INSERT INTO posts (Title, Content, Category, Status)
		VALUES (?, ?, ?, ?)
	`
	result, err := db.DB.Exec(query, article.Title, article.Content, article.Category, article.Status)
	if err != nil {
		log.Println("Gagal menyimpan artikel:", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Gagal mendapatkan ID artikel:", err)
		return 0, err
	}

	return id, nil
}

// GetArticles mengambil daftar artikel dengan pagination
func GetArticles(limit int, offset int) ([]models.Article, error) {
	query := `
		SELECT Id, Title, Content, Category, Status
		FROM posts
		LIMIT ? OFFSET ?
	`
	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Gagal mengambil artikel:", err)
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Category, &article.Status)
		if err != nil {
			log.Println("Gagal membaca baris:", err)
			return nil, err
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error setelah iterasi:", err)
		return nil, err
	}

	return articles, nil
}

// GetArticleByID mengambil detail artikel berdasarkan ID
func GetArticleByID(id int) (*models.Article, error) {
	query := `
		SELECT Id, Title, Content, Category, Status
		FROM posts
		WHERE Id = ?
	`
	row := db.DB.QueryRow(query, id)

	var article models.Article
	err := row.Scan(&article.ID, &article.Title, &article.Content, &article.Category,  &article.Status)
	if err != nil {
		log.Println("Gagal mengambil artikel:", err)
		return nil, err
	}

	return &article, nil
}

// UpdateArticle mengubah data artikel berdasarkan ID
func UpdateArticle(id int, article models.Article) error {
	query := `
		UPDATE posts
		SET Title = ?, Content = ?, Category = ?, Status = ?
		WHERE Id = ?
	`
	_, err := db.DB.Exec(query, article.Title, article.Content, article.Category, article.Status, id)
	if err != nil {
		log.Println("Gagal mengubah artikel:", err)
		return err
	}

	return nil
}

// DeleteArticle menghapus data artikel berdasarkan ID
func DeleteArticle(id int) error {
	query := `
		DELETE FROM posts
		WHERE Id = ?
	`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Println("Gagal menghapus artikel:", err)
		return err
	}

	return nil
}
