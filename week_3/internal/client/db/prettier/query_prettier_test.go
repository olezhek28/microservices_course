package prettier

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestQueryPrettier(t *testing.T) {
	t.Parallel()

	var (
		id       = gofakeit.Int64()
		name     = gofakeit.FirstName()
		safeName = fmt.Sprintf("\"%s\"", name)
		age      = gofakeit.Int64()
	)

	t.Run("without squirrel", func(t *testing.T) {
		t.Parallel()
		t.Run("query with no params", func(t *testing.T) {
			t.Parallel()

			q := "SELECT id, name, age FROM users"
			pretty := Pretty(q, PlaceholderDollar)
			require.Equal(t, q, pretty)
		})

		t.Run("query with spaces", func(t *testing.T) {
			t.Parallel()

			q := "   SELECT id, name, age FROM users   "
			pretty := Pretty(q, PlaceholderDollar)
			require.Equal(t, "SELECT id, name, age FROM users", pretty)
		})

		t.Run("select by one param", func(t *testing.T) {
			t.Parallel()

			q := "SELECT id, name, age FROM users WHERE id = $1"
			pretty := Pretty(q, PlaceholderDollar, id)
			require.Equal(t, fmt.Sprintf("SELECT id, name, age FROM users WHERE id = %v", id), pretty)
		})

		t.Run("select with IN clause", func(t *testing.T) {
			t.Parallel()

			q := "SELECT id, name, age FROM users WHERE id IN ($1,$2,$3)"
			pretty := Pretty(q, PlaceholderDollar, id, 123, 321)
			require.Equal(t, fmt.Sprintf("SELECT id, name, age FROM users WHERE id IN (%v,%v,%v)", id, 123, 321), pretty)
		})

		t.Run("insert", func(t *testing.T) {
			t.Parallel()

			q := "INSERT INTO users (name, age) VALUES ($1,$2)"
			pretty := Pretty(q, PlaceholderDollar, name, age)
			require.Equal(t, fmt.Sprintf("INSERT INTO users (name, age) VALUES (%v,%v)", safeName, age), pretty)
		})

		t.Run("update", func(t *testing.T) {
			t.Parallel()

			q := "UPDATE users SET name = $1, age = $2 WHERE name = $1"
			pretty := Pretty(q, PlaceholderDollar, []byte(name), age)
			require.Equal(t, fmt.Sprintf("UPDATE users SET name = %v, age = %v WHERE name = %v", safeName, age, safeName), pretty)
		})
	})

	t.Run("with squirrel", func(t *testing.T) {
		t.Parallel()
		t.Run("query with no params", func(t *testing.T) {
			t.Parallel()

			builder := sq.Select("id, name, age").
				PlaceholderFormat(sq.Dollar).
				From("users")

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, "SELECT id, name, age FROM users", pretty)
		})

		t.Run("query with spaces", func(t *testing.T) {
			t.Parallel()

			builder := sq.Select("id, name, age").
				PlaceholderFormat(sq.Dollar).
				From("users    ")

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)

			require.Equal(t, "SELECT id, name, age FROM users", pretty)
		})

		t.Run("SELECT by one param", func(t *testing.T) {
			t.Parallel()

			builder := sq.Select("id, name, age").
				PlaceholderFormat(sq.Dollar).
				From("users").
				Where(sq.Eq{"id": id})

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, fmt.Sprintf("SELECT id, name, age FROM users WHERE id = %v", id), pretty)
		})

		t.Run("select with IN clause", func(t *testing.T) {
			t.Parallel()

			builder := sq.Select("id, name, age").
				PlaceholderFormat(sq.Dollar).
				From("users").
				Where(sq.Eq{"id": []int64{id, 123, 321}})

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, fmt.Sprintf("SELECT id, name, age FROM users WHERE id IN (%v,%v,%v)", id, 123, 321), pretty)
		})

		t.Run("insert", func(t *testing.T) {
			t.Parallel()

			builder := sq.Insert("users").
				PlaceholderFormat(sq.Dollar).
				Columns("name, age").
				Values(name, age)

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, fmt.Sprintf("INSERT INTO users (name, age) VALUES (%v,%v)", safeName, age), pretty)
		})

		t.Run("update", func(t *testing.T) {
			t.Parallel()

			builder := sq.Update("users").
				PlaceholderFormat(sq.Dollar).
				Set("name", name).
				Set("age", age).
				Where(sq.Eq{"name": name})

			query, v, err := builder.ToSql()
			require.NoError(t, err)

			pretty := Pretty(query, PlaceholderDollar, v...)
			require.Equal(t, fmt.Sprintf("UPDATE users SET name = %v, age = %v WHERE name = %v", safeName, age, safeName), pretty)
		})
	})
}
