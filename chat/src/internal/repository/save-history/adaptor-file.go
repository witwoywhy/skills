package savehistory

import (
	"agent-chat/infrastructure"
	"agent-chat/internal/domain/history"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type adaptorFile struct {
}

func NewAdaptorFile() Port {
	return &adaptorFile{}
}

func (a *adaptorFile) Execute(request *Request) (*Response, error) {
	history := history.History{
		Datetime: request.Datetime.UnixMilli(),
		From:     request.From,
		To:       request.To,
		Message:  request.Message,
	}

	b, err := json.Marshal(history)
	if err != nil {
		return nil, fmt.Errorf("failed when json.Marshal history: %v", err)
	}
	b = append(b, 10)

	historyFromDir := infrastructure.Path.HistoryDir + "/" + request.From
	historyFromFilePath := historyFromDir + "/" + request.Datetime.Format(time.DateOnly) + ".json"

	historyToDir := infrastructure.Path.HistoryDir + "/" + request.To
	historyToFilePath := historyToDir + "/" + request.Datetime.Format(time.DateOnly) + ".json"

	_, err = os.Stat(historyFromDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(historyFromDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed when create folder [%s]: %v", historyFromDir, err)
		}

		err = os.WriteFile(historyFromFilePath, b, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed when write file [%s]: %v", historyFromFilePath, err)
		}
	} else {
		f, err := os.OpenFile(historyFromFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed when open file [%s]: %v", historyFromFilePath, err)
		}
		defer f.Close()

		f.Write(b)
	}

	_, err = os.Stat(historyToDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(historyToDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed when create folder [%s]: %v", historyToDir, err)
		}

		err = os.WriteFile(historyToFilePath, b, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed when write file [%s]: %v", historyToFilePath, err)
		}
	} else {
		f, err := os.OpenFile(historyToFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed when open file [%s]: %v", historyToFilePath, err)
		}
		defer f.Close()

		f.Write(b)
	}

	return &Response{}, nil
}
