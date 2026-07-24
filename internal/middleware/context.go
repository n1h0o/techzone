package middleware

type ContextKey string

// хранит ключ с данными пользователя после авторизации
const UserKey ContextKey = "user"
