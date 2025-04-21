CREATE TABLE clicks (  
    id SERIAL PRIMARY KEY,  
    ad_id INT NOT NULL,  
    timestamp TIMESTAMP NOT NULL,  
    ip VARCHAR(45) NOT NULL,  
    video_playback_time INT NOT NULL,  
    FOREIGN KEY (ad_id) REFERENCES ads(id)  
);  

CREATE INDEX idx_clicks_ad_id_timestamp ON clicks(ad_id, timestamp);