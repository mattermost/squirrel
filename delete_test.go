package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteBuilderToSql(t *testing.T) {
	b := Delete("").
		Prefix("WITH prefix AS ?", 0).
		From("a").
		What("a.*, b.*").
		Where("b = ?", 1).
		Using("other").
		Join("other2 ON (other.id = other2.id)").
		LeftJoin("other3 ON (other.id = other3.id)").
		RightJoin("other4 ON (other.id = other4.id)").
		CrossJoin("other5 ON (other.id = other5.id)").
		InnerJoin("other6 ON (other.id = other6.id)").
		JoinClause(Select("*").From("other7").Prefix("JOIN (").Suffix(")")).
		OrderBy("c").
		Limit(2).
		Offset(3).
		Suffix("RETURNING ?", 4)

	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql :=
		"WITH prefix AS ? " +
			"DELETE a.*, b.* FROM a USING other JOIN other2 ON (other.id = other2.id) LEFT JOIN other3 ON (other.id = other3.id) RIGHT JOIN other4 ON (other.id = other4.id) CROSS JOIN other5 ON (other.id = other5.id) INNER JOIN other6 ON (other.id = other6.id) JOIN ( SELECT * FROM other7 ) WHERE b = ? ORDER BY c LIMIT 2 OFFSET 3 " +
			"RETURNING ?"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{0, 1, 4}
	assert.Equal(t, expectedArgs, args)
}

func TestDeleteBuilderToSqlErr(t *testing.T) {
	_, _, err := Delete("").ToSql()
	assert.Error(t, err)
}

func TestDeleteBuilderMustSql(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("TestDeleteBuilderMustSql should have panicked!")
		}
	}()
	Delete("").MustSql()
}

func TestDeleteBuilderPlaceholders(t *testing.T) {
	b := Delete("test").Where("x = ? AND y = ?", 1, 2)

	sql, _, _ := b.PlaceholderFormat(Question).ToSql()
	assert.Equal(t, "DELETE FROM test WHERE x = ? AND y = ?", sql)

	sql, _, _ = b.PlaceholderFormat(Dollar).ToSql()
	assert.Equal(t, "DELETE FROM test WHERE x = $1 AND y = $2", sql)
}

func TestDeleteBuilderRunners(t *testing.T) {
	db := &DBStub{}
	b := Delete("test").Where("x = ?", 1).RunWith(db)

	expectedSql := "DELETE FROM test WHERE x = ?"

	b.Exec()
	assert.Equal(t, expectedSql, db.LastExecSql)
}

func TestDeleteBuilderNoRunner(t *testing.T) {
	b := Delete("test")

	_, err := b.Exec()
	assert.Equal(t, RunnerNotSet, err)
}

func TestDeleteWithQuery(t *testing.T) {
	db := &DBStub{}
	b := Delete("test").Where("id=55").Suffix("RETURNING path").RunWith(db)

	expectedSql := "DELETE FROM test WHERE id=55 RETURNING path"
	b.Query()

	assert.Equal(t, expectedSql, db.LastQuerySql)
}
