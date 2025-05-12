-- Создание таблицы клиентов
CREATE TABLE clients (
                         client_id VARCHAR(36) PRIMARY KEY,
                         login VARCHAR(255) NOT NULL,
                         age INTEGER,
                         location VARCHAR(255),
                         gender VARCHAR(10),
                         CHECK (gender IN ('MALE', 'FEMALE'))
);

-- Создание таблицы рекламодателей
CREATE TABLE advertisers (
                             advertiser_id VARCHAR(36) PRIMARY KEY,
                             name VARCHAR(255) NOT NULL
);

-- Создание таблицы ML скоринга
CREATE TABLE ml_scores (
                           client_id VARCHAR(36) NOT NULL,
                           advertiser_id VARCHAR(36) NOT NULL,
                           score INTEGER,
                           PRIMARY KEY (client_id, advertiser_id),
                           FOREIGN KEY (client_id) REFERENCES clients(client_id) ON DELETE CASCADE,
                           FOREIGN KEY (advertiser_id) REFERENCES advertisers(advertiser_id) ON DELETE CASCADE
);

-- Создание таблицы рекламных кампаний
CREATE TABLE campaigns (
                           campaign_id VARCHAR(36) PRIMARY KEY,
                           advertiser_id VARCHAR(36) NOT NULL,
                           impressions_limit INTEGER,
                           clicks_limit INTEGER,
                           cost_per_impression FLOAT,
                           cost_per_click FLOAT,
                           ad_title VARCHAR(255),
                           ad_text TEXT,
                           start_date INTEGER,
                           end_date INTEGER,
                           image_url VARCHAR,
                           targeting_gender VARCHAR(10),
                           targeting_age_from INTEGER,
                           targeting_age_to INTEGER,
                           targeting_location VARCHAR(255),
                           impressions_count INTEGER DEFAULT 0,
                           clicks_count INTEGER DEFAULT 0,
                           FOREIGN KEY (advertiser_id) REFERENCES advertisers(advertiser_id) ON DELETE CASCADE,
                           CHECK (targeting_gender IS NULL OR targeting_gender IN ('MALE', 'FEMALE', 'ALL'))
);

-- Создание таблицы статистики кампаний по дням
CREATE TABLE campaign_daily_stats (
                                      campaign_id VARCHAR(36) NOT NULL,
                                      date INTEGER NOT NULL,
                                      impressions_count INTEGER DEFAULT 0,
                                      clicks_count INTEGER DEFAULT 0,
                                      conversion FLOAT DEFAULT 0,
                                      spent_impressions FLOAT DEFAULT 0,
                                      spent_clicks FLOAT DEFAULT 0,
                                      spent_total FLOAT DEFAULT 0,
                                      PRIMARY KEY (campaign_id, date),
                                      FOREIGN KEY (campaign_id) REFERENCES campaigns(campaign_id) ON DELETE CASCADE
);

-- Создание таблицы ежедневной статистики рекламодателей
CREATE TABLE advertiser_daily_stats (
                                        advertiser_id VARCHAR(36) NOT NULL,
                                        date INTEGER NOT NULL,
                                        impressions_count INTEGER DEFAULT 0,
                                        clicks_count INTEGER DEFAULT 0,
                                        conversion FLOAT DEFAULT 0,
                                        spent_impressions FLOAT DEFAULT 0,
                                        spent_clicks FLOAT DEFAULT 0,
                                        spent_total FLOAT DEFAULT 0,
                                        PRIMARY KEY (advertiser_id, date),
                                        FOREIGN KEY (advertiser_id) REFERENCES advertisers(advertiser_id) ON DELETE CASCADE
);

-- Создание таблицы объявлений (вариант с отдельной таблица объявлений)
CREATE TABLE ads (
                     ad_id VARCHAR(36) PRIMARY KEY, -- Совпадает с campaign_id
                     advertiser_id VARCHAR(36) NOT NULL,
                     ad_title VARCHAR(255),
                     ad_text TEXT,
                     FOREIGN KEY (ad_id) REFERENCES campaigns(campaign_id) ON DELETE CASCADE,
                     FOREIGN KEY (advertiser_id) REFERENCES advertisers(advertiser_id) ON DELETE CASCADE
);

-- Создание таблицы для учета кликов по объявлениям
CREATE TABLE ad_clicks (
                           ad_id VARCHAR(36) NOT NULL,
                           client_id VARCHAR(36) NOT NULL,
                           click_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           FOREIGN KEY (ad_id) REFERENCES campaigns(campaign_id) ON DELETE CASCADE,
                           FOREIGN KEY (client_id) REFERENCES clients(client_id) ON DELETE CASCADE
);

-- Создание таблицы для учета показов объявлений
CREATE TABLE ad_impressions (
                                ad_id VARCHAR(36) NOT NULL,
                                client_id VARCHAR(36) NOT NULL,
                                impression_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                FOREIGN KEY (ad_id) REFERENCES campaigns(campaign_id) ON DELETE CASCADE,
                                FOREIGN KEY (client_id) REFERENCES clients(client_id) ON DELETE CASCADE
);

-- Создание таблицы для управления временем
CREATE TABLE current_dates (
                              id SMALLINT PRIMARY KEY CHECK (id = 1),
                              date INTEGER NOT NULL
);

-- Инициализация текущей даты
INSERT INTO current_dates (id, date) VALUES (1, 0);
