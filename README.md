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

- [About](#about)
- [Getting Started](#getting_started)
- [Deployment](#deployment)
- [Usage](#usage)
- [Built Using](#built_using)
- [Authors](#authors)


## 🧐 About <a name = "about"></a>

This project is used to backup files, sql and nosql databases in s3 and keep a version history.

- Mysql or MariaDb
- Postgresql
- Mongodb
- Folder
- File


## 🏁 Getting Started <a name = "getting_started"></a>

### Prerequisites

```
go 1.20
Account S3
```

### Installing

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
Set your environments vars

## 🎈 Usage <a name="usage"></a>

### Run from source

```bash
go run . 
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

- [Go](https://go.dev/) - Programming Language.
- [Fx](https://github.com/uber-go/fx) - Fx is a dependency injection system for Go..
- [S3](https://pkg.go.dev/github.com/aws/aws-sdk-go/service/s3) - Package s3 provides the client and types for making API requests to Amazon Simple Storage Service. 


## ✍️ Authors <a name = "authors"></a>

- [@sistemasnegros](https://github.com/sistemasnegros) - Idea & Initial work


