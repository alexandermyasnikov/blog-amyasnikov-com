func Migrate(ctx context.Context, logger *slog.Logger, dbHandler *sql.DB, schema string, migrationsFS fs.FS) error {
	_, err := dbHandler.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS "+schema)
	if err != nil {
		return err
	}

	store, err := database.NewStore(goose.DialectPostgres, schema+".goose_db_version")
	if err != nil {
		return err
	}

	migrationsFS, err = fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	provider, err := goose.NewProvider("", dbHandler, migrationsFS,
		goose.WithStore(store), goose.WithVerbose(true), goose.WithLogger(sharedlogginggoose.NewGooseLogger(logger)))
	if err != nil {
		return err
	}

	migrations, err := provider.Up(ctx)
	if err != nil {
		return err
	}

	logger.InfoContext(ctx, "migrations applied", "migrations", migrations)

	return nil
}
