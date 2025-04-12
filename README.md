# llm_coding_challenge
Репозиторий по решению кейса неструктурированных сообщений в хакатоне LLM Coding Challenge 2025

### Запуск контейнера
```docker
docker build -t llm-structure-service .
docker run -d --env-file .env -p 8000:8000 --name llm-structure-service llm-structure-service
```
