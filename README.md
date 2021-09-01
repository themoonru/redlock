redlock
=======

Распределенная блокировка заданий с помощью redis.

Поддерживает fancing при потере связи с redis, а также удержание мастера при новых выборах.

Быстрый старт
-------------
1.Склонируйте репозиторий
``` 
git clone https://github.com/themoonru/redlock.git
```
2.Запустите redis
``` 
cd redlock/redis
docker-compose up -d
```
3.Соберите приложение и запустите нужное количество экземпляров, указав в качестве параметра номер экземпляра
``` 
go build
NUM=1 ./redlock
NUM=2 ./redlock
```

4.Попробуйте уронить мастер или redis
```
cd redlock/redis
docker-compose down
```