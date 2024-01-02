<p align="center">
  <a href="" rel="noopener">
 <img height=200px src="https://miro.medium.com/max/900/1*5JXt0wiQjX_FDwYvrxPN9Q.png" alt="Project logo"></a>
</p>

<h3 align="center">Golang rotate backup s3</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/kylelobo/The-Documentation-Compendium.svg)](https://github.com/kylelobo/The-Documentation-Compendium/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

---

<p align="center"> This project is used to backup files in S3.
    <br> 
</p>

## 📝 Table of Contents

-   [About](#about)
-   [Getting Started](#getting_started)
-   [Deployment](#deployment)
-   [Usage](#usage)
-   [Built Using](#built_using)
-   [Authors](#authors)

## 🧐 About <a name = "about"></a>

This project is used to backup files, sql and nosql databases in s3 and keep a version history.

-   Mysql or MariaDb
-   Postgresql
-   Mongodb
-   Folder
-   File

## 🏁 Getting Started <a name = "getting_started"></a>

### Prerequisites

```
Account S3
```

### Installing from binary

```bash
mkdir ~/go-rotate-backup-s3

curl -L "https://github.com/sistemasnegros/go-rotate-backup-s3/releases/download/v1.0.1/backups3" -o ~/go-rotate-backup-s3/backups3

curl -L  "https://raw.githubusercontent.com/sistemasnegros/go-rotate-backup-s3/master/.env.default" -o ~/go-rotate-backup-s3/.env

ln -sf ~/go-rotate-backup-s3/backups3 /usr/local/bin/backups3
```

### Installing from source

Prerequisites

```
go 1.20
```

```
git clone https://github.com/sistemasnegros/go-rotate-backup-s3
cd go-rotate-backup-s3/src
```

Once you've cloned the project, install dependencies with

```
go mod download
```

copy .env.default to .env

```
cp .env.default .env
```

### Configuration .env

Set your environments vars

```bash
BACKUP_SRC=/tmp/nosql.db
BACKUP_COMMAND=mongodump --uri="mongodb://root:toor@localhost/?authSource=admin" --db=certimail --archive="${BACKUP_SRC}"
BACKUP_DST=prefixs3/mongodb
BACKUP_KEEP=5
BACKUP_PREFiX_NAME=nosql.db
BACKUP_COMMAND_TIMEOUT=120

# email service
SMTP_ENABLED=true
SMTP_HOST=smtp.server.com
SMTP_PORT=25
SMTP_USER=myUser
SMTP_PASS=myPassword
SMTP_FROM=no-reply@mydomain.com
SMTP_TO=user1@mydomain.com,user2@mydomain.com
SMTP_TEMPLATE=commons/infra/html/emails/

# Files Service
AWS_REGION=us-east-2
AWS_ACCESS_KEY_ID=XXXXXXXXXXXXXX
AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxx
AWS_BUCKET=mybucket
AWS_ENDPOINT=https://s3.us-east-2.amazonaws.com
AWS_URL_PUBLIC=https://s3.us-east-2.amazonaws.com/mybucket

```

### Build production

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o backups3 .
```

## 🎈 Usage <a name="usage"></a>

### Run from binary

```bash
backups3 -config .env
```

### Run from source

```bash
go run .
```

output

```
INFO[2023-04-10T18:44:07-05:00] connection S3 successful
INFO[2023-04-10T18:44:22-05:00] executed command successful
INFO[2023-04-10T18:47:34-05:00] create backup successful: prefixs3/mongodb/v0/2023-04-10_18:44:22_nosql.db
```

## 🚀 Deployment <a name = "deployment"></a>

To create a production version of your app:

```bash
go build -o backups3
chmod +x  backups3
```

### Run from binary

```bash
./backups3 -config .env
```

## ⛏️ Built Using <a name = "built_using"></a>

-   [Go](https://go.dev/) - Programming Language.
-   [Fx](https://github.com/uber-go/fx) - Fx is a dependency injection system for Go..
-   [S3](https://pkg.go.dev/github.com/aws/aws-sdk-go/service/s3) - Package s3 provides the client and types for making API requests to Amazon Simple Storage Service.

## ✍️ Authors <a name = "authors"></a>

-   [@sistemasnegros](https://github.com/sistemasnegros) - Idea & Initial work
