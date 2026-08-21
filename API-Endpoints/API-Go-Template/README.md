# Technology Stack

## Server / Operating System

- **Ubuntu Server 24.04 LTS**
- **Linux**
- **systemd** for automatic service startup and process management

## Backend

- **Go (Golang)**
- Standard Go HTTP server using `net/http`
- The API is compiled into a standalone binary and managed through a `systemd` service.

The backend provides a REST-style API with:

- Public `GET` endpoints for reading data
- Authenticated `POST`, `PUT`, and `DELETE` endpoints for modifying data
- JSON request and response handling

## Database

- **MySQL**
- Database: `olympics`

The following tables are used:

- `country`
- `athletes`
- `medals`
- `sports`
- `medals_athletes_sports`
- `authorized_users`
- `api_tokens`

The `medals_athletes_sports` table connects countries, athletes, medal types, and sports.

## Authentication

User authentication is handled through the database.

Users are stored in:

```text
olympics.authorized_users
```

Passwords are stored as hashes rather than plaintext.

After successful authentication, the API issues a token. Protected requests must include the token in the HTTP `Authorization` header:

```text
Authorization: User <token>
```

The API uses the following access model:

| HTTP Method | Authentication |
|---|---|
| `GET` | Not required |
| `POST` | Required |
| `PUT` | Required |
| `DELETE` | Required |

## API Endpoints

The API currently provides CRUD functionality for:

```text
/api/country
/api/athletes
/api/medals
/api/sports
/api/medals-athletes-sports
```

Each endpoint supports:

```text
GET
POST
PUT
DELETE
```

`GET` requests are publicly accessible, while all modifying operations require authentication.

## Service Management

The Go application is compiled into a standalone Linux binary:

```text
/var/www/olympics-api/olympics-api
```

It is started and managed by `systemd` using:

```text
olympics-api.service
```

The service is configured to:

- Start automatically during system boot
- Run the compiled Go API
- Restart automatically if the application crashes

## Overall Architecture

```text
Client / Website
       │
       │ HTTP + JSON
       ▼
┌─────────────────┐
│     Go API      │
│                 │
│  Public GET     │
│  Authenticated  │
│  POST/PUT/DELETE│
└────────┬────────┘
         │
         │ MySQL
         ▼
┌────────────────────────────┐
│     Olympics Database      │
│                            │
│ • country                  │
│ • athletes                 │
│ • medals                   │
│ • sports                   │
│ • medals_athletes_sports   │
│ • authorized_users         │
│ • api_tokens               │
└────────────────────────────┘
```

This stack was chosen to keep the project relatively simple.
The Go backend runs as a standalone compiled application, so no separate PHP runtime or traditional web server is required for the API itself.

# Database Creation

```sql
CREATE DATABASE olympics;

USE olympics;


CREATE TABLE country (
    land_id SMALLINT UNSIGNED NOT NULL,
    land_code VARCHAR(2) NOT NULL,
    land_name VARCHAR(255) NOT NULL,

    PRIMARY KEY (land_id),
    UNIQUE (land_name)
);


CREATE TABLE athletes (
    athleten_id INT UNSIGNED NOT NULL,
    athleten_vorname VARCHAR(255) NOT NULL,
    athleten_name VARCHAR(255) NOT NULL,

    PRIMARY KEY (athleten_id)
);


CREATE TABLE medals (
    medaillen_id SMALLINT UNSIGNED NOT NULL,
    art_medaille VARCHAR(255) NOT NULL,

    PRIMARY KEY (medaillen_id)
);


CREATE TABLE sports (
    sportart_id INT UNSIGNED NOT NULL,
    sportart VARCHAR(255) NOT NULL,

    PRIMARY KEY (sportart_id),
    UNIQUE (sportart)
);


CREATE TABLE medals_athletes_sports (
    land_id SMALLINT UNSIGNED NOT NULL,
    athleten_id INT UNSIGNED NOT NULL,
    medaillen_id SMALLINT UNSIGNED NOT NULL,
    sportart_id INT UNSIGNED NOT NULL,
    anzahl_medaillen INT UNSIGNED NOT NULL DEFAULT 0,

    PRIMARY KEY (
        land_id,
        athleten_id,
        medaillen_id,
        sportart_id
    ),

    FOREIGN KEY (land_id)
        REFERENCES country(land_id),

    FOREIGN KEY (athleten_id)
        REFERENCES athletes(athleten_id),

    FOREIGN KEY (medaillen_id)
        REFERENCES medals(medaillen_id),

    FOREIGN KEY (sportart_id)
        REFERENCES sports(sportart_id)
);
```
