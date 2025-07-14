templ:
	@templ generate --watch
tailwind:
	@tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
air:
	@air
migrate_up:
	sqlite3 internal/db/app.db < internal/db/migrations/schema.up.sql
migrate_down:
	sqlite3 internal/db/app.db < internal/db/migrations/schema.down.sql