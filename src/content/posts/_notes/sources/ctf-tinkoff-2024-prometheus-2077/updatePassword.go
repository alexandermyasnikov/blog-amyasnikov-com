func UpdatePassword(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	session := sessions.Default(c)
	username := session.Get("user_id")
	// username := c.PostForm("username")
	password := c.PostForm("password")

	password = strings.ReplaceAll(password, "🔥", "💨")

	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	go func(db *sql.DB, password, username interface{}) {

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		query := fmt.Sprintf("UPDATE users SET password = '%s' \n", password) // hashPassword(password))
		query = query + "WHERE username = $1"

		_, err := db.ExecContext(ctx, query, username)

		if err != nil {
			log.Printf("Failed to insert message: %v", err)
		}

	}(db, password, username)

	Logout(c)
}