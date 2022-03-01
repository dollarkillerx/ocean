package api

import (
	"encoding/json"
	"log"

	"github.com/dollarkillerx/ocean/pkg/models"
)

func (s *Server) migrate(index string, payload []byte) error {
	var schema models.Schema
	err := json.Unmarshal(payload, &schema)
	if err != nil {
		log.Println(err)
		return err
	}

	return s.storage.CreateIndex(index, schema)
}
