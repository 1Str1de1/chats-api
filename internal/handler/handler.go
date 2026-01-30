package handler

import (
	"chats-api/internal/model"
	"chats-api/internal/repository"
	"chats-api/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Handler struct {
	chats    *services.ChatsService
	messages *services.MessagesService
}

func NewHandler(db *gorm.DB) *Handler {
	chatsRepo := repository.NewChatsRepo(db)
	messagesRepo := repository.NewMessagesRepo(db)

	chats := services.NewChatsRepository(chatsRepo)
	messages := services.NewMessagesRepository(messagesRepo)

	return &Handler{
		chats:    chats,
		messages: messages,
	}
}

func (h *Handler) HandleChatsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type CreateChatReq struct {
			Title string `json:"title"`
		}

		var req CreateChatReq
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json body: "+err.Error())
			return
		}

		title, err := h.chats.ValidateChatCreate(req.Title)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		chat, err := h.chats.CreateChat(r.Context(), title)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(chat)
	}
}

func (h *Handler) HandleMessagesCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIdStr := r.PathValue("id")
		chatId, err := strconv.Atoi(chatIdStr)
		if err != nil || chatId == 0 {
			writeError(w, http.StatusBadRequest, "invalid chat_id")
			return
		}

		type CreateMessageReq struct {
			Text string `json:"text"`
		}

		var req CreateMessageReq
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json body: "+err.Error())
			return
		}

		if err := h.messages.ValidateMessageCreate(req.Text); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		message, err := h.messages.CreateMessage(r.Context(), req.Text, chatId)
		if errors.Is(err, repository.ErrChatNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}
}

func (h *Handler) HandleMessagesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIdStr := r.PathValue("id")
		chatId, err := strconv.Atoi(chatIdStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid chat_id")
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 20

		if limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid limit")
				return
			}
			if l > 100 {
				l = 100
			}
			limit = l
		}

		type Response struct {
			Chat     *model.Chat      `json:"chat"`
			Messages []*model.Message `json:"messages"`
		}

		chat, err := h.chats.GetChat(chatId)
		if errors.Is(err, repository.ErrChatNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		messages, err := h.messages.GetAllMessagesFromChat(chatId, limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		resp := Response{
			Chat:     chat,
			Messages: messages,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)
	}
}

func (h *Handler) HandleChatsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIdStr := r.PathValue("id")
		chatId, err := strconv.Atoi(chatIdStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid chat_id")
			return
		}

		if err := h.chats.DeleteChat(chatId); errors.Is(err, repository.ErrChatNotFound) {
			writeError(w, http.StatusNoContent, err.Error())
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
