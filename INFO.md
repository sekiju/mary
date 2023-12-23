# ChoGazoViewerBrowser
## [manga.fod.fujitv.co.jp](https://manga.fod.fujitv.co.jp/)

Отличие полноценной книги от демо прописано в префиксе айди эпизода.

BT000137301100200201900209_000

* _000 - demo
* _001 - full

Есть два типа API:
* Web
* Для Android приложения

Отличий немного, у android к примеру 2 роутера для получения ключей и все они возвращают чуток детальнее информацию чем обычный веб:
* POST https://manga.fod.fujitv.co.jp/api/books/licenceKey
* POST https://manga.fod.fujitv.co.jp/api/books/licenceKeyForBrowser

Android требует для запросов 2 хедера:
* zk-app-version: 1.1.24
* zk-os-type: 1

Web Api:
* GET https://manga.fod.fujitv.co.jp/web/books/licenceKey?book_id=1373011&episode_id=BT000137301100200201

Web требует один:
* Zk-Web-version: 1.3.3

Если приложение обновили, то и версии в запросах тоже надо обновлять, потому что есть проверка чтоб актуально было всё.

Сессия живёт 48 часов?
### Web

Дополнительные ключи получаются если подтверждена покупка.

_Возможно ключи - хэш изображения_

### Epub
.ebg формат (используется только в Японии?)

Мы можем скачать Epub файл в RIFF контейнере. Ключи контейнера:
* META
* data
* KEY
* KTST

data - drm epub archive. demo and full has **_different keys_**