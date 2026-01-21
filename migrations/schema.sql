-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 2. Dramas Table
CREATE TABLE IF NOT EXISTS dramas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    synopsis TEXT,
    poster_url TEXT,
    year INTEGER,
    rating DECIMAL(3, 1) DEFAULT 0,
    total_seasons INTEGER DEFAULT 1,
    status VARCHAR(20) CHECK (status IN ('ongoing', 'completed')),
    view_count INTEGER DEFAULT 0,
    source_url TEXT,
    added_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dramas_year ON dramas(year);
CREATE INDEX IF NOT EXISTS idx_dramas_status ON dramas(status);
CREATE INDEX IF NOT EXISTS idx_dramas_view_count ON dramas(view_count DESC);
CREATE INDEX IF NOT EXISTS idx_dramas_title ON dramas USING gin(to_tsvector('english', title));

-- 3. Genres Table
CREATE TABLE IF NOT EXISTS genres (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_genres_slug ON genres(slug);

-- 4. Drama Genres Table (Many-to-Many)
CREATE TABLE IF NOT EXISTS drama_genres (
    drama_id UUID REFERENCES dramas(id) ON DELETE CASCADE,
    genre_id UUID REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (drama_id, genre_id)
);

-- 5. Actors Table
CREATE TABLE IF NOT EXISTS actors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    photo_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 6. Drama Actors Table (Many-to-Many)
CREATE TABLE IF NOT EXISTS drama_actors (
    drama_id UUID REFERENCES dramas(id) ON DELETE CASCADE,
    actor_id UUID REFERENCES actors(id) ON DELETE CASCADE,
    role VARCHAR(20) CHECK (role IN ('main', 'support')),
    PRIMARY KEY (drama_id, actor_id)
);

-- 7. Seasons Table
CREATE TABLE IF NOT EXISTS seasons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    drama_id UUID REFERENCES dramas(id) ON DELETE CASCADE,
    season_number INTEGER NOT NULL,
    title VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(drama_id, season_number)
);

-- 8. Episodes Table
CREATE TABLE IF NOT EXISTS episodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    season_id UUID REFERENCES seasons(id) ON DELETE CASCADE,
    episode_number INTEGER NOT NULL,
    title VARCHAR(255),
    video_url TEXT NOT NULL,
    duration INTEGER, -- in seconds
    thumbnail_url TEXT,
    view_count INTEGER DEFAULT 0,
    source_url TEXT,
    added_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(season_id, episode_number)
);

CREATE INDEX IF NOT EXISTS idx_episodes_view_count ON episodes(view_count DESC);

-- 9. Watch History Table
CREATE TABLE IF NOT EXISTS watch_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    episode_id UUID REFERENCES episodes(id) ON DELETE CASCADE,
    progress_seconds INTEGER DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    last_watched_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, episode_id)
);

CREATE INDEX IF NOT EXISTS idx_watch_history_user ON watch_history(user_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_last_watched ON watch_history(last_watched_at DESC);

-- 10. Watchlist Table
CREATE TABLE IF NOT EXISTS watchlist (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    drama_id UUID REFERENCES dramas(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, drama_id)
);

-- 11. Reviews Table
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    drama_id UUID REFERENCES dramas(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 10),
    review_text TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, drama_id)
);

CREATE INDEX IF NOT EXISTS idx_reviews_drama ON reviews(drama_id);

-- 12. Comments Table
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    episode_id UUID REFERENCES episodes(id) ON DELETE CASCADE,
    comment_text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_comments_episode ON comments(episode_id, created_at DESC);
