-- Pastikan sambungan UUID wujud (biasanya sudah ada dari jadual employees)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Jadual meeting_minutes (Jadual Utama)
CREATE TABLE meeting_minutes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    division VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    meeting_date TIMESTAMP NOT NULL,
    meeting_type VARCHAR(100) NOT NULL,
    summary TEXT,
    notes TEXT,
    speaker VARCHAR(255),
    number_of_participants INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 2. Jadual meeting_participants (Relasi Notulen dengan Karyawan)
CREATE TABLE meeting_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    minute_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_minute_participants FOREIGN KEY (minute_id) REFERENCES meeting_minutes (id) ON DELETE CASCADE,
    CONSTRAINT fk_employee_participants FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE
);

-- 3. Jadual meeting_results (Action Items / Tugasan dari Rapat)
CREATE TABLE meeting_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    minute_id UUID NOT NULL,
    employee_id UUID, -- Boleh NULL jika tugasan umum (tiada PIC khusus)
    target_description TEXT NOT NULL,
    target_nominal BIGINT, -- Menggunakan BIGINT untuk selaras dengan Target Tracker
    achievement_status VARCHAR(50) DEFAULT 'To Do',
    target_completion_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_minute_results FOREIGN KEY (minute_id) REFERENCES meeting_minutes (id) ON DELETE CASCADE,
    CONSTRAINT fk_employee_results FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE SET NULL
);

-- 4. Jadual meeting_images (Pautan Gambar / Dokumentasi Rapat)
CREATE TABLE meeting_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    minute_id UUID NOT NULL,
    file_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_minute_images FOREIGN KEY (minute_id) REFERENCES meeting_minutes (id) ON DELETE CASCADE
);

-- (Pilihan) Indexes untuk mempercepatkan carian data (Query Optimization)
CREATE INDEX idx_minute_date ON meeting_minutes(meeting_date);
CREATE INDEX idx_participant_minute ON meeting_participants(minute_id);
CREATE INDEX idx_participant_employee ON meeting_participants(employee_id);
CREATE INDEX idx_result_minute ON meeting_results(minute_id);
CREATE INDEX idx_result_employee ON meeting_results(employee_id);