Запуск 
  
    make build
  
    make run

Миграция бд

    make migrate

Запуск тестов

    make test

Часть сервиса аутентификации

    Три REST маршрута:
    
      /signUp - регистрация пользователя в системе
          body example
          {
            "email": "example@gmail.com",
            "password": "qwerty12345"
          }
          
      /sign-in - авторизация
          body example
          {
            "email": "example@gmail.com",
            "password": "qwerty12345"
          }
          в теле ответа будет access token
          
      /refresh - обновляет и возвращает access token; создает новый refresh токен, сохранаяет в бд и в cookie

      рефреш токен сохраняется в cookie
      ip адрес пользователя сервис получает из заголовка "X-Forwarded-For"

