func saveSchema(param migrationParam) error {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{ //nolint:exhaustruct
		Image:        "postgres:18-alpine",
		Env:          map[string]string{"POSTGRES_PASSWORD": "pass", "POSTGRES_DB": "testdb"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ //nolint:exhaustruct
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}

	defer func() { _ = postgresC.Terminate(ctx) }()

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://postgres:pass@%s/testdb?sslmode=disable", net.JoinHostPort(host, port.Port()))

	dbHandler, err := shareddatabase.New(ctx, dsn)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := sharedsqlmigrator.Migrate(ctx, logger, dbHandler, param.schema, param.migrationsFs); err != nil {
		return err
	}

	cmd := []string{
		"pg_dump",
		"--schema-only",
		"--format=plain",
		"--no-owner",
		"--username=postgres",
		"--dbname=testdb",
		"--restrict-key=000000000000000000000000000000000000000000000000000000000000000",
	}

	_, r, err := postgresC.Exec(ctx, cmd, exec.Multiplexed())
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./docs/sqlschema/%s.sql", param.schema))
	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}
