package main

import (
	"encoding/json"
	"os"
)

type FilePublisher struct {
}

func (p *FilePublisher) Publish(s *SkuResponse) error {
	jsonOutpout, _ := json.Marshal(s)

	f, err := os.OpenFile("queue.json", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(jsonOutpout); err != nil {
		return err
	}
	f.WriteString("\n")

	return err
}
