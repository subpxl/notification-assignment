-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    recipient VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    message_id VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMP
);

-- Insert test data
INSERT INTO messages (content, recipient, status) VALUES
('Test 1', '+903333333331', 'pending'),
('Test 2', '+903333333332', 'pending'),
('Test 3', '+903333333333', 'pending'),
('Test 4', '+903333333334', 'pending'),
('Test 5', '+903333333335', 'pending'),
('Test 6', '+903333333336', 'pending'),
('Test 7', '+903333333337', 'pending'),
('Test 8', '+903333333338', 'pending'),
('Test 9', '+903333333339', 'pending'),
('Test 10','+940333333341', 'pending'),
('Test 11','+900333333342', 'pending'),
('Test 12','+900333333343', 'pending'),
('Test 13','+900333333344', 'pending');
