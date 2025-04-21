# ПВЗ для Приёмки Товаров

Backend-сервис, реализованный на Go,был использован фреймворк GIN,базы данных PortgreSQL и Docker и JWT для авторизации.
Готова все функциональные требования и регистрация с авторизацией по логину,который возвращает токен.
Для интеграционного теста была создана изолированная тестовая база данных с приписками test.
------------------------------------
Вопросы и решения:  
Для развертывания в Doker использовать файл .env: DB_HOST=localhost DB_PORT=5432 DB_USER=admin DB_PASSWORD=password DB_NAME=DataBase  
 Файл .env.test использовать для интеграционного теста,для него была сделана тестовая база данных чтобы провести тест изолированно от реальной бд.  
 .env.test: DB_HOST=localhost DB_PORT=5433 DB_USER=admin DB_PASSWORD=password DB_NAME=TestDataBase  
 Таблицы основной базы данных были очищены перед отправкой.  
 Перед каждый прогоном база данных очищается в интеграционном тесте,для удобства проверки.  
 В юнит тестах для взаимодействия с базой используется мок (заглушка), для изоляции от реальной базы данных.  
 Несоответствия спецификации и описания:  
 В спецификации в user role [employee, moderator]  в описании  (client, moderator), решил оставить clien,employee  
 Использовал hey для проверки нефункциональных требований, но он выводил по запросу Status code distribution:  [401] 20000 responses.  
 В поиске проблемы настроил логирование для поиска ошибок и увидел, что при запросе hey получен заголовок Authorization: '', токен не был передан.  
 Проблема заключалась в том, что hey не поддерживает установку заголовков Authorization с токеном напрямую в простом виде.  
------------------------------------
Подключение к основной базе данных: docker exec -it go_project_db psql -U admin -d DataBase  
Подключение к тестовой базе данных: docker exec -it go_project_test_db psql -U admin -d TestDataBase  
Порт на котором работает сервер: 8080  
http://localhost:8080/  
Эндпоинты(Для проверки использовал Postman):  
/register  
 Регистрация пользователей с ролями `employee` и `client`  
 При регистрации почта не должна совпадать с существующей,пароль минимум из 4 символов  
Post запрос:  
{  
  "email": "user@example.com",  
  "password": "password",  
  "role": "employee"  
}  
ИЛИ В ТЕРМИНАЛЕ  
curl -X POST http://localhost:8080/register \  
  -H "Content-Type: application/json" \  
  -d '{  
    "email": "user@example.com",  
    "password": "password",  
    "role": "employee"  
}'  
------------------------------------ 
Авторизация и получение токена:  
/login  
Post запрос:  
{  
  "email": "user@example.com",  
  "password": "password"  
}  
ИЛИ В ТЕРМИНАЛЕ  
curl -X POST http://localhost:8080/login \  
  -H "Content-Type: application/json" \  
  -d '{  
    "email": "user@example.com",  
    "password": "password"  
}'  
------------------------------------
Авторизация без логина для получения токена:  
роли `employee` и `client`  
/dummyLogin  
Post запрос:  
{  
   "role":"client"  
}  
ИЛИ В ТЕРМИНАЛЕ  
curl -X POST http://localhost:8080/dummyLogin \  
  -H "Content-Type: application/json" \  
  -d '{  
    "role": "client"  
}'  
------------------------------------
------------------------------------
Создание ПВЗ:  
/pvz  
Post запрос:  
//Название города с заглавной буквы.  
{  
   "city":"Москва"   
}  
ИЛИ В ТЕРМИНАЛЕ Вместо <token> вставить валидный токен авторизации  
curl -X POST http://localhost:8080/pvz \  
  -H "Authorization: Bearer <token>" \  
  -H "Content-Type: application/json" \  
  -d '{  
    "city": "Москва"  
}'  
------------------------------------
Создание приемки товара:  
/receptions  
Post запрос:  
{  
   "pvzId":"вставить id pvz"  
}  
ИЛИ В ТЕРМИНАЛЕ Вместо <token> вставить валидный токен авторизации  
curl -X POST http://localhost:8080/receptions \  
  -H "Authorization: Bearer <token>" \  
  -H "Content-Type: application/json" \  
  -d '{  
    "pvzId": "вставить-id-pvz"  
}'  
------------------------------------
Добавление товара в открытую приемку:  
Типы товара с большой буквы: Одежда,Электроника,Обувь  
/products  
Post запрос:  
{  
    "type":"Одежда",  
   "pvzId":"вставить id pvz"  
}  
ИЛИ В ТЕРМИНАЛЕ Вместо <token> вставить валидный токен авторизации  
curl -X POST http://localhost:8080/products \  
  -H "Authorization: Bearer <token>" \  
  -H "Content-Type: application/json" \  
  -d '{  
    "type": "Одежда",  
    "pvzId": "вставить-id-pvz"  
}'  
------------------------------------  
Удаление последнего добавленного товара из открытой приемки:  
Вместо {pvzId} вставить id pvz из которого нужно удалить товар:  
Post запрос:  
/pvz/{pvzId}/delete_last_product  
ИЛИ В ТЕРМИНАЛЕ Вместо <token> вставить валидный токен авторизации  
curl -X POST http://localhost:8080/pvz/<pvzId>/delete_last_product \  
  -H "Authorization: Bearer <token>"  
------------------------------------
Закрытие последней приемки товара:  
Post запрос  
Вместо {pvzId} вставить id pvz в котором нужно закрыть приемку:  
/pvz/{:pvzId}/close_last_reception  
ИЛИ В ТЕРМИНАЛЕ Вместо <token> вставить валидный токен авторизации  
curl -X POST http://localhost:8080/pvz/<pvzId>/close_last_reception \  
  -H "Authorization: Bearer <token>"  
------------------------------------
Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией:  
Get запрос  
где startDate=начальная дата диапозона=2025-04-19T19:30:00Z  
    endDate=конечная дата диапозона=2025-04-19T19:39:00Z  
    page=номер страницы  
    limit=количество запросов на странице  
/pvz?startDate=2025-04-19T19:30:00Z&endDate=2025-04-19T19:39:00Z&page=1&limit=10  
ИЛИ В ТЕРМИНАЛЕ Вместо <token> вставить валидный токен авторизации  
curl -X GET "http://localhost:8080/pvz?startDate=2025-04-19T19:30:00Z&endDate=2025-04-19T19:39:00Z&page=1&limit=10" \  
  -H "Authorization: Bearer <token>"  




