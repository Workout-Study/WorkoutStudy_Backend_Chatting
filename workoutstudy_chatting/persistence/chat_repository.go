package persistence

import (
	"database/sql"
	"fmt"
	"time"
	"workoutstudy_chatting/model"
)

// ChatRepository 인터페이스 정의
type ChatRepository interface {
	RetrieveMessages(fitGroupID int, since time.Time) ([]model.ChatMessage, error)
	SaveMessage(msg model.ChatMessage) error
}

type ChatRepositoryImpl struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepository {
	return &ChatRepositoryImpl{DB: db}
}

func (repo *ChatRepositoryImpl) RetrieveMessages(fitGroupID int, since time.Time) ([]model.ChatMessage, error) {
	// SQL 쿼리를 작성하여, 주어진 fit_group_id에 대한 메시지를 시간 기준 내림차순으로 조회합니다.
	query := `
	SELECT id, fit_group_id, fit_mate_id, message, message_time, message_type
	FROM message
	WHERE fit_group_id = $1 AND message_time > $2
	ORDER BY message_time DESC
	`
	rows, err := repo.DB.Query(query, fitGroupID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.ChatMessage
	for rows.Next() {
		var msg model.ChatMessage
		if err := rows.Scan(&msg.ID, &msg.FitGroupID, &msg.FitMateID, &msg.Message, &msg.MessageTime, &msg.MessageType); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (repo *ChatRepositoryImpl) SaveMessage(msg model.ChatMessage) error {
	query := `
    INSERT INTO message (message_id, fit_group_id, fit_mate_id, message, message_time, message_type, created_at, created_by, updated_at, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, NOW(), $7)
    `
	_, err := repo.DB.Exec(query, msg.ID, msg.FitGroupID, msg.FitMateID, msg.Message, msg.MessageTime, msg.MessageType, fmt.Sprintf("%d", msg.FitMateID))
	return err
}
