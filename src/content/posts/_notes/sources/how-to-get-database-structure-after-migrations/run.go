func run() error {
	params := []migrationParam{
		{schema: taskrepository.MigrationsSchema, migrationsFs: taskrepository.MigrationsFS},
		{schema: userrepository.MigrationsSchema, migrationsFs: userrepository.MigrationsFS},
		// ...
	}

	for _, param := range params {
		if err := saveSchema(param); err != nil {
			return err
		}
	}

	return nil
}
