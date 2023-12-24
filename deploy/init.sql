\c image_maesto;

CREATE TABLE
    image (
        id VARCHAR(27) NOT NULL PRIMARY KEY,
        url TEXT NOT NULL,
        content_type VARCHAR(255) NOT NULL,
        exif JSON,
        status VARCHAR(20) NOT NULL,
        checksum VARCHAR(255) NOT NULL,
        width BIGINT NOT NULL,
        height BIGINT NOT NULL,
        byte_size BIGINT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    )