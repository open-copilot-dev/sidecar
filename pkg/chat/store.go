package chat

import (
	"encoding/json"
	"io/fs"
	"open-copilot.dev/sidecar/pkg/chat/domain"
	"open-copilot.dev/sidecar/pkg/common"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Store interface {
	ListChats(curPage, pageSize int) (int, []*domain.Chat, error)
	SaveChat(chat *domain.Chat) error
	GetChat(chatID string) (*domain.Chat, error)
	DeleteChatMessage(chatID string, messageID string) error
	DeleteChat(chatID string) error
	DeleteAllChats() error
}

//----------------------------------------------------------------

type LocalStore struct {
	dir string
}

func (l *LocalStore) ListChats(curPage, pageSize int) (int, []*domain.Chat, error) {
	// read file under dir
	files := make([]fs.FileInfo, 0)
	err := filepath.Walk(l.dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			files = append(files, info)
		}
		return nil
	})
	if err != nil {
		return 0, nil, common.NewErrWithCause(common.ErrCodeIo, "walk dir failed", err)
	}
	if len(files) == 0 {
		return 0, []*domain.Chat{}, nil
	}

	// sort by time
	sort.Slice(files, func(i, j int) bool {
		iTime := files[i].ModTime()
		jTime := files[j].ModTime()
		return iTime.After(jTime)
	})

	// 分页过滤
	if curPage < 1 {
		curPage = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if curPage*pageSize > len(files) {
		curPage = len(files) / pageSize
		if len(files)%pageSize > 0 {
			curPage++
		}
	}
	chats := make([]*domain.Chat, 0, pageSize)
	for i := (curPage - 1) * pageSize; i < curPage*pageSize; i++ {
		if i >= len(files) {
			break
		}
		entry := files[i]
		chat, err := l.GetChat(strings.TrimSuffix(entry.Name(), ".json"))
		if err != nil {
			return 0, nil, err
		}
		chats = append(chats, chat)
	}
	return len(files), chats, nil
}

func (l *LocalStore) SaveChat(chat *domain.Chat) error {
	path := filepath.Join(l.dir, chat.ChatID+".json")
	content, err := json.Marshal(chat)
	if err != nil {
		return common.NewErrWithCause(common.ErrCodeMarshal, "marshal chat failed", err)
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		return common.NewErrWithCause(common.ErrCodeIo, "write file failed", err)
	}
	return nil
}

func (l *LocalStore) GetChat(chatID string) (*domain.Chat, error) {
	path := filepath.Join(l.dir, chatID+".json")
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil, common.ErrNotFound
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, common.NewErrWithCause(common.ErrCodeIo, "read file failed", err)
	}
	var chat = domain.Chat{}
	if err := json.Unmarshal(content, &chat); err != nil {
		return nil, common.NewErrWithCause(common.ErrCodeMarshal, "unmarshal file failed", err)
	}
	return &chat, nil
}

func (l *LocalStore) DeleteChatMessage(chatID string, messageID string) error {
	chat, err := l.GetChat(chatID)
	if err != nil {
		return err
	}
	newMessages := make([]*domain.ChatMessage, 0)
	for _, message := range chat.Messages {
		if message.MessageID != messageID {
			newMessages = append(newMessages, message)
		}
	}
	chat.Messages = newMessages
	if len(chat.Messages) == 0 {
		return l.DeleteChat(chatID)
	}
	return l.SaveChat(chat)
}

func (l *LocalStore) DeleteChat(chatID string) error {
	path := filepath.Join(l.dir, chatID+".json")

	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(path); err != nil {
		return common.NewErrWithCause(common.ErrCodeIo, "delete file failed", err)
	}
	return nil
}

func (l *LocalStore) DeleteAllChats() error {
	err := os.RemoveAll(l.dir)
	if err != nil {
		return common.NewErrWithCause(common.ErrCodeIo, "delete dir failed", err)
	}
	_ = os.MkdirAll(l.dir, 0o755)
	return err
}

func NewLocalStore(dir string) *LocalStore {
	_ = os.MkdirAll(dir, 0o755)
	return &LocalStore{dir: dir}
}
