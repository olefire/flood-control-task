# Flood Control

Для решения проблемы ограничения скорости запросов я выбрал алгоритм Token Bucket, поскольку он никода не превышает заданный лимит (в других алгоритмах при всплеске трафика около границы окна лимит может быть превышен).
В этом алгоритме каждому пользователю выделяется бакет, в который  помещается определенное количество токенов.
При обработке запроса к токенам пользователя прибавляется **floor(elapsed*rate/window)**, где **elapsed** - разница между текущим временем и временем последнего запроса,
**rate** - количество добавляющихся в секунду токенов, **window** - N секунд. Затем токен из бакета удаляется, а если токенов не достаточно, то запрос отбрасывается.<br>

В качестве общего хранилища данных я использовал Redis, так как он является резидентным хранилищем типа key-value, т.е. ключи и значения хранятся в оперативной памяти, поэтому обращения к данным выполняются очень быстро. В key хранится userId, а в value - бакет.

При одновременной работе нескольких инстансов два параллельных запросах на одного юзера могут привести к грязному чтению.
Чтобы избежать этой проблемы, все операции с базой данных необходимо проводить в атомарной транзакции.

Текущее время получаю из редиса, чтобы решение было устойчивым к отсутствию у инстансов точной синхронизации с NTP.