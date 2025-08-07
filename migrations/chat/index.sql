-- Для быстрого поиска чатов пользователя
CREATE INDEX idx_conversations_client ON chat_conversations(client_id);
CREATE INDEX idx_conversations_contractor ON chat_conversations(contractor_id);
CREATE INDEX idx_conversations_job ON chat_conversations(job_id);

-- Для сообщений
CREATE INDEX idx_messages_conversation ON chat_messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender ON chat_messages(sender_id);
CREATE INDEX idx_messages_unread ON chat_messages(conversation_id, is_read, created_at);