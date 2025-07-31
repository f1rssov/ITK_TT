# ITK_TT
# WalletT — сервис управления кошельками

Приложение для управления кошельками с балансом.  
Реализовано на Go с использованием Gin и PostgreSQL.  
Поддерживает операции обновления баланса и получение информации о кошельке.

---

## Возможности

- Обновление баланса кошелька (пополнение или списание)
- Получение текущего баланса кошелька по ID
- Swagger UI для удобной работы с API и документацией

---

## Технологии

- Go
- Gin Web Framework
- PostgreSQL (через pgx или pgxpool)
- Swagger (swaggo/swag) для автодокументации

---

## Быстрый запуск

 **Makefile**

```bash
git clone https://github.com/xxxx/xxxx wallet
cd wallett
make
```

## API

POST /wallet — обновить баланс кошелька

GET /wallet/{wallet_id} — получить баланс кошелька по ID

Пример запроса на обновление баланса
```json
{
  "valletId": "123e4567-e89b-12d3-a456-426614174000",
  "operationType": "deposit", // или "withdraw"
  "amount": 150
}
```
Swagger UI содержит подробную документацию
