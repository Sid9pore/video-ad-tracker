CREATE TABLE ads (  
    id SERIAL PRIMARY KEY,  
    image_url VARCHAR(255) NOT NULL,  
    target_url VARCHAR(255) NOT NULL  
);


INSERT INTO ads (image_url, target_url) VALUES  
('https://example.com/image1.jpg', 'https://example.com/ad1'),  
('https://example.com/image2.jpg', 'https://example.com/ad2'),  
('https://example.com/image3.jpg', 'https://example.com/ad3'),  
('https://example.com/image4.jpg', 'https://example.com/ad4'),  
('https://example.com/image5.jpg', 'https://example.com/ad5');  


CREATE TABLE clicks (  
    id SERIAL PRIMARY KEY,  
    ad_id INT NOT NULL,  
    timestamp TIMESTAMP NOT NULL,  
    ip VARCHAR(45) NOT NULL,  
    video_playback_time INT NOT NULL,  
    FOREIGN KEY (ad_id) REFERENCES ads(id)  
);  

CREATE INDEX idx_clicks_ad_id_timestamp ON clicks(ad_id, timestamp);

INSERT INTO clicks (ad_id, timestamp, ip, video_playback_time) VALUES  
(1, '2025-04-20 10:00:00', '192.168.1.1', 30),  
(2, '2025-04-21 11:15:00', '192.168.1.2', 45),  
(1, '2025-04-21 12:30:00', '192.168.1.3', 60),  
(3, '2025-04-22 09:45:00', '192.168.1.4', 25),  
(4, '2025-04-22 14:00:00', '192.168.1.5', 40);  
