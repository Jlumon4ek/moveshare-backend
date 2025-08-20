-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Notification metadata
    type VARCHAR(50) NOT NULL, -- 'job_application', 'job_update', 'payment', 'document_upload', 'new_job', 'review', 'system'
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    
    -- Related entities (optional)
    job_id BIGINT REFERENCES jobs(id) ON DELETE SET NULL,
    chat_id BIGINT REFERENCES chat_conversations(id) ON DELETE SET NULL,
    related_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL, -- who triggered the notification
    
    -- Status and priority
    is_read BOOLEAN DEFAULT FALSE,
    priority VARCHAR(20) DEFAULT 'normal', -- 'low', 'normal', 'high', 'urgent'
    
    -- Actions (JSON array of available actions)
    actions JSONB DEFAULT '[]',
    
    -- Metadata (additional data specific to notification type)
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    read_at TIMESTAMP WITH TIME ZONE NULL,
    expires_at TIMESTAMP WITH TIME ZONE NULL -- for temporary notifications
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_job_id ON notifications(job_id);

-- Composite index for common queries
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread_created ON notifications(user_id, is_read, created_at DESC);