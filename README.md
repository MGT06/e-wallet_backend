# [Backend] - E-Walett API

<p>
  <a href="https://golang.org/"><img style="border-radius:4px" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"></a>
  <a href="https://github.com/gin-gonic/gin"><img style="border-radius:4px" src="https://img.shields.io/badge/Gin-008B8B?style=for-the-badge&logo=go&logoColor=white" alt="Gin"></a>
  <a href="https://github.com/joho/godotenv"><img style="border-radius:4px" src="https://img.shields.io/badge/Godotenv-ECD53F?style=for-the-badge&logo=dotenv&logoColor=black" alt="Godotenv"></a>
  <a href="https://www.postgresql.org/"><img style="border-radius:4px" src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"></a>
  <a href="https://github.com/jackc/pgx"><img style="border-radius:4px" src="https://img.shields.io/badge/pgx-339933?style=for-the-badge&logo=postgresql&logoColor=white" alt="pgx"></a>
  <a href="https://redis.io/"><img style="border-radius:4px" src="https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis"></a>
</p>

A robust backend RESTful API built with Go and the Gin Framework, developed as part of the collaburation team regu mawar. This project serves as a comprehensive backend service for managing transactions, authentication and managing data.

## Tech Stack

- [![Go](https://img.shields.io/badge/Go-v1.27.1-blue?logo=go&logoColor=white)](https://go.dev/learn/) <br>Digunakan sebagai bahasa pemrograman utama untuk membangun seluruh logika backend aplikasi _Gin Exercise Koda Batch 9_ ini.

- [![Gin-Gonic](https://img.shields.io/badge/Gin_Gonic-v1.12.0-green?logo=Gin&logoColor=white)](https://gin-gonic.com/en/)<br>Dipakai sebagai _HTTP router_ dan web framework untuk mengatur _endpoint_ RESTful API fitur autentikasi, manajemen _user_, komunitas, dan _events_.
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v18.6-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/) <br>Menjadi basis data relasional utama untuk menyimpan data aplikasi secara permanen seperti akun pengguna, profil, daftar komunitas, detail _events_, serta relasi _join/saved events_.
- [![pgx](https://img.shields.io/badge/Pgx-v5.11.0-blue?)](https://github.com/jackc/pgx) <br>Digunakan sebagai _driver_ dan _toolkit_ penghubung antara kode Go dan PostgreSQL untuk mengeksekusi _query_ database secara cepat dan efisien pada fitur-fitur aplikasi.
- [![Redis](https://img.shields.io/badge/Redis-v8.10.2-blue?logo=go&logoColor=white)](https://redis.io/) <br>Diterapkan untuk mengelola sesi pengguna, seperti penyimpanan token autentikasi atau _blacklist_ token saat _logout_ dari aplikasi.
- [![GodotEnv](https://img.shields.io/badge/GodotEnv-v1.5.1-blue?)](https://github.com/jackc/pgx) <br>Berfungsi untuk membaca file konfigurasi `.env` guna mengelola variabel rahasia dan dinamis (seperti koneksi database PostgreSQL, port server, dan kunci rahasia token) di aplikasi.

## Features

- Authentication (Register, Login, Logout, Forgot Password, Reset passsword)
- User (change pin, get favorite user transfer, check balance, history, etc.)
- Transactions (transfer, search people, topup, history)

## Usage Instruction

### Clone this repository

```bash
$ git clone https://github.com/MGT06/e-wallet_backend
```

### Environment Setup

Create your environment on the root directory named **.env** example use [.env.example](https://github.com/MGT06/e-wallet_backend/blob/main/.env.example)

### Redis Config Setup

Setup your configuration for your use cache on the root directory named **redis.conf** example [redis.conf](https://github.com/MGT06/e-wallet_backend/blob/main/redis.conf.example)

### Docker Pull Image from [**DokerHub**](https://hub.docker.com/)

1. PosgresSQL [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v18.6-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/)

```bash
$ docker pull postgres:18.6-alpine3.24
```

2. Redis [![Redis](https://img.shields.io/badge/Redis-v8.10.2-blue?logo=go&logoColor=white)](https://redis.io/)

```bash
$ docker pull redis:8.10.2-alpine3.23
```

### Running Database

1. PosgresSQL [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-v18.6-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/)

```bash
$ docker run --name {NAME_CONTAINER} -d --env-file .env -p {PORT_FORWARD} -v {YOUR_DOCKER_VOLUME}:/var/lib/postgresql/18/docker postgres:18.6-alpine3.24
```

2. Redis [![Redis](https://img.shields.io/badge/Redis-v8.10.2-blue?logo=go&logoColor=white)](https://redis.io/)

```bash
$ docker run -d --name {NAME_CONTAINER} -p {PORT-FORWARDING} -v {PATH_CONFIG}:/usr/local/etc/redis/redis.conf" -v "{DOCKER_VOLUME}:/data" --restart unless-stopped redis:8.8.3-alpine3.23 redis-server /usr/local/etc/redis/redis.conf
```

3. Create Schema Table and used data dummy

```zsh
$ make migrate-up
$ make db-seed
```

4. Install dependency

```bash
$ go mod download
```

3. etc.

## Routes

| Endpoint        | Method | Description     |
| :-------------- | :----- | :-------------- |
| `/auth/`        | `POST` | Authentictions  |
| `/user/`        | `GET`  | Management user |
| `/transaction/` | `GET`  | Transactions    |

### Documentation

For complete documentation, visit [docs swagger endpooint **`/docs/v1/index.html`** ](http://localhost/docs/index.html)

## Changelog

| Version | Desctiption |
| ------- | ----------- |
| 1.0.0   | Initial app |

- 1.1.1
  - Feat history .... by [almarufh](https://github.com/almarufh/)
  - Feat transfer .... by [given](https://github.com/MGT06/)
  - Feat topup .... by [given](https://github.com/MGT06/)
  - Feat forgot password .... by [given](https://github.com/almarufh/)
  - Feat cors .... by [almarufh](https://github.com/almarufh/)
  - Feat hash password .... by [almarufh](https://github.com/almarufh)
  - Feat Autentication .... by [almarufh](https://github.com/almarufh)

## How to Contribute

- Fork this repository
- Create your changes
- Pull Request

## License [![License: MIT](https://img.shields.io/badge/License-MIT-red)](https://opensource.org/license/mit)

This project is licensed under the MIT License

## Related Project

[![WhatsApp](https://img.shields.io/badge/FRONTEND--blue?logo=googlechrome&logoColor=blue)](https://github.com/MGT06/e-wallet_frontend)
