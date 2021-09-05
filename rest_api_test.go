package main

import (
	"testing"

	"github.com/Mechwarrior1/PGL_backend/word2vec"
	"github.com/stretchr/testify/assert"
)

func TestStartServer(t *testing.T) {
	var embed *word2vec.Embeddings
	s, _, _, err := StartServer(embed)
	// if err = s.ListenAndServeTLS("secure//cert.pem", "secure//key.pem"); err != nil && err != http.ErrServerClosed {
	// 	e.Logger.Fatal(err)
	// }
	if assert.NoError(t, err) {
		s.Close()
	}
}
