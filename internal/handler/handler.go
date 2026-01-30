package handler

import (
	"chats-api/internal/model"
	"chats-api/internal/repository"
	"chats-api/internal/services"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type Handler struct {
	chats    services.ChatsService
	messages services.MessagesService
	logger   *slog.Logger
}

func NewHandler(chats services.ChatsService, messages services.MessagesService, logger *slog.Logger) *Handler {
	return &Handler{
		chats:    chats,
		messages: messages,
		logger:   logger,
	}
}

func (h *Handler) HandleChatsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handling create chat")

		type CreateChatReq struct {
			Title string `json:"title"`
		}

		var req CreateChatReq
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json body: "+err.Error())
			h.logger.Error("got invalid json body")
			return
		}

		title, err := h.chats.ValidateChatCreate(req.Title)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			h.logger.Error("chat request is invalid")
			return
		}

		chat, err := h.chats.CreateChat(r.Context(), title)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			h.logger.Error(fmt.Sprintf("failed to create chat %v", err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(chat)
		h.logger.Info(fmt.Sprintf("successfully created chat with title: %s", title))
	}
}

func (h *Handler) HandleMessagesCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handling create chat")

		chatIdStr := r.PathValue("id")
		chatId, err := strconv.Atoi(chatIdStr)
		if err != nil || chatId == 0 {
			writeError(w, http.StatusBadRequest, "invalid chat_id")
			h.logger.Error("chat id is invalid")
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
			h.logger.Error("got invalid json body")
			return
		}

		if err := h.messages.ValidateMessageCreate(req.Text); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			h.logger.Error("message request is invalid")
			return
		}

		message, err := h.messages.CreateMessage(r.Context(), req.Text, chatId)
		if errors.Is(err, repository.ErrChatNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			h.logger.Error(fmt.Sprintf("chat with id %d not found", chatId))
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			h.logger.Error(fmt.Sprintf("failed to create message %v", err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
		h.logger.Info(fmt.Sprintf("successfully created message in chat %d with id: %d", chatId, message.Id))
	}
}

func (h *Handler) HandleMessagesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handling get chat")
		chatIdStr := r.PathValue("id")
		chatId, err := strconv.Atoi(chatIdStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid chat_id")
			h.logger.Error("chat id is invalid")
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 20

		if limitStr != "" {
			l, err := strconv.Atoi(limitStr)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid limit")
				h.logger.Error("limit is invalid")
				return
			}
			if l > 100 {
				l = 100
				h.logger.Warn("limit is too large, setting to 100")
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
			h.logger.Error(fmt.Sprintf("chat with id %d not found", chatId))
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			h.logger.Error(fmt.Sprintf("failed to get chat with id %d: %v", chatId, err))
			return
		}

		messages, err := h.messages.GetAllMessagesFromChat(chatId, limit)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			h.logger.Error(fmt.Sprintf("failed to get messages from chat with id %d: %v", chatId, err))
			return
		}
		resp := Response{
			Chat:     chat,
			Messages: messages,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&resp)
		h.logger.Info(fmt.Sprintf("successfully fetched chat with id %d and limit %d", chatId, limit))
	}
}

func (h *Handler) HandleChatsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("handling delete chat")

		chatIdStr := r.PathValue("id")
		chatId, err := strconv.Atoi(chatIdStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid chat_id")
			h.logger.Error("chat id is invalid")
			return
		}

		if err := h.chats.DeleteChat(chatId); errors.Is(err, repository.ErrChatNotFound) {
			writeError(w, http.StatusNoContent, err.Error())
			h.logger.Error(fmt.Sprintf("chat with id %d not found", chatId))
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			h.logger.Error(fmt.Sprintf("failed to delete chat with id %d: %v", chatId, err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		h.logger.Info(fmt.Sprintf("successfully deleted chat with id %d", chatId))
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
