package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Soroka-EDMS/svc/sessions/pkgs/config"
)

func TestSave_Exists(t *testing.T) {
	db, err := Connection(config.GetLogger().Logger, "stub")
	assert.NoError(t, err)

	var testData = struct {
		sub   string
		token string
	}{
		"user@example.com",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZ2xhZHlzLmNoYW1wbEBlZG1zLmNvbSIsImh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiXSwiZXhwIjoxNTY4Mzc4NTg0LCJpYXQiOjE1NjgyOTIxODQsImlzcyI6Imh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiLCJtYXNrIjozMjc2NywibmJmIjoxNTY4MjkyMTg0LCJzdWIiOiJnbGFkeXMuY2hhbXBsQGVkbXMuY29tIn0.xUKkiNOClwhhGlgXWj_J9u0t_ImJKsW-mbK9xuTiF5o",
	}

	db.Save(testData.sub, testData.token)

	assert.True(t, db.Exist(testData.sub, testData.token))
}

func TestGet_Exists(t *testing.T) {
	db, err := Connection(config.GetLogger().Logger, "stub")
	assert.NoError(t, err)

	var testData = struct {
		sub   string
		token string
	}{
		"user@example.com",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZ2xhZHlzLmNoYW1wbEBlZG1zLmNvbSIsImh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiXSwiZXhwIjoxNTY4Mzc4NTg0LCJpYXQiOjE1NjgyOTIxODQsImlzcyI6Imh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiLCJtYXNrIjozMjc2NywibmJmIjoxNTY4MjkyMTg0LCJzdWIiOiJnbGFkeXMuY2hhbXBsQGVkbXMuY29tIn0.xUKkiNOClwhhGlgXWj_J9u0t_ImJKsW-mbK9xuTiF5o",
	}

	db.Save(testData.sub, testData.token)

	token, err := db.Get(testData.sub)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGet_DoesNotExist(t *testing.T) {
	db, err := Connection(config.GetLogger().Logger, "stub")
	assert.NoError(t, err)

	var testData = struct {
		sub   string
		token string
	}{
		"user@example.com",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZ2xhZHlzLmNoYW1wbEBlZG1zLmNvbSIsImh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiXSwiZXhwIjoxNTY4Mzc4NTg0LCJpYXQiOjE1NjgyOTIxODQsImlzcyI6Imh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiLCJtYXNrIjozMjc2NywibmJmIjoxNTY4MjkyMTg0LCJzdWIiOiJnbGFkeXMuY2hhbXBsQGVkbXMuY29tIn0.xUKkiNOClwhhGlgXWj_J9u0t_ImJKsW-mbK9xuTiF5o",
	}

	db.Save(testData.sub, testData.token)

	token, err := db.Get("user@email.com")
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestDelete(t *testing.T) {
	db, err := Connection(config.GetLogger().Logger, "stub")
	assert.NoError(t, err)

	var testData = struct {
		sub   string
		token string
	}{
		"user@example.com",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZ2xhZHlzLmNoYW1wbEBlZG1zLmNvbSIsImh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiXSwiZXhwIjoxNTY4Mzc4NTg0LCJpYXQiOjE1NjgyOTIxODQsImlzcyI6Imh0dHBzOi8vZWRtcy5jb20vc2Vzc2lvbnMiLCJtYXNrIjozMjc2NywibmJmIjoxNTY4MjkyMTg0LCJzdWIiOiJnbGFkeXMuY2hhbXBsQGVkbXMuY29tIn0.xUKkiNOClwhhGlgXWj_J9u0t_ImJKsW-mbK9xuTiF5o",
	}

	err = db.Save(testData.sub, testData.token)
	assert.NoError(t, err)

	err = db.Delete(testData.sub, testData.token)
	assert.NoError(t, err)

	assert.False(t, db.Exist(testData.sub, testData.token))
}
