# 559

## Запуск

```bash
559.exe https://shonenjumpplus.com/magazine/4856001361564051266 -o output/
```

## Сайты

- [manga.fod.fujitv.co.jp](https://manga.fod.fujitv.co.jp/)
- [shonenjumpplus.com](https://shonenjumpplus.com/)
- [comic-walker.com](https://comic-walker.com/) _(регистрация в коде не поддерживается)_
- [www.pixiv.net](https://www.pixiv.net/)
- [pocket.shonenmagazine.com](https://pocket.shonenmagazine.com)
- [comic-action.com](https://comic-action.com)
- [comic-days.com](https://comic-days.com)
- [comic-growl.com](https://comic-growl.com)
- [comic-earthstar.com](https://comic-earthstar.com)
- [comic-gardo.com](https://comic-gardo.com)
- [comic-trail.com](https://comic-trail.com)
- [comic-zenon.com](https://comic-zenon.com)
- [comicborder.com](https://comicborder.com)
- [kuragebunch.com](https://kuragebunch.com)
- [magcomi.com](https://magcomi.com)
- [tonarinoyj.jp](https://tonarinoyj.jp)
- [viewer.heros-web.com](https://viewer.heros-web.com)
- [www.sunday-webry.com](https://www.sunday-webry.com)
- [storia.takeshobo.co.jp](https://storia.takeshobo.co.jp) (без регистрации)
- [comicbushi-web.com](https://comicbushi-web.com)
- [comic.webnewtype.com](https://comic.webnewtype.com) (без регистрации)
- [www.comic-valkyrie.com](https://www.comic-valkyrie.com)
- [cmoa.jp](https://cmoa.jp)

## Настройки

Все параметры прописываются в `config.yaml`. Пример файла:

```yaml
settings:
  debug: 
    enable: true
    url: https://shonenjumpplus.com/episode/14079602755375556618
  output_path: output/
  clear_output_folder: true
  threads: 6

sites:
  manga.fod.fujitv.co.jp:
    session: VWdx8id9R0XHjVpvs7s754CxGJpBBl9ZCHCqL1yF
    purchase_free_books: true
```

Все сайты имеют обязательный параметр `session` для авторизации.
Но некоторые имеют свои уникальные настройки, как `purchase_free_books` у `manga.fod.fujitv.co.jp`.

### Дополнительные настройки

| Парсер                     | Параметр        | Тип    | Описание                                                                                                     |
|----------------------------|-----------------|--------|--------------------------------------------------------------------------------------------------------------|
|                            | session         | string | Уникальный идентификатор сеанса или сессии.                                                                  |
| **manga.fod.fujitv.co.jp** | tryPurchaseBook | bool   | Проверяет, есть ли книга в списке бесплатных книг, если есть — парсер купит её и загрузит полноценную версию |
