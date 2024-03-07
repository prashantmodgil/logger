
# logger# Go Log Search API

  

This is a simple REST API built in Go using the Gorilla Mux router to search log files stored remotely based on criteria such as date range and search string.

  

## Features

  

- Search log files by date range and search string

- Flexible and extensible architecture

- Easy to deploy and use

  

## Prerequisites

  

- Docker installed on your machine

- Remote storage accessible with log files stored in the specified structure

- Basic understanding of Go programming language

  

## Getting Started

  

1. Clone the repository:

  

```bash
git  clone  https://github.com/your-username/go-log-search-api.git
```

2.  Run Service
```bash

docker-compose up
```
Thats All You need

## Curl Requests
1. Seed A log  
```
curl --location --request GET '127.0.01:8080/seed'
```
2. Get Api for logs
```
curl --location --request GET '127.0.01:8080/logs?searchKeyword=hello&from=2024-01-02T15:04:05Z&to=2024-03-07T15:04:05Z'
```
