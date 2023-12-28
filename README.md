# 559

## Запуск

```bash
559.exe https://shonenjumpplus.com/magazine/4856001361564051266 -o output/
```

| Флаг        | Сокращение | Тип данных | Значение поумолчанию | Описание                                   | 
|-------------|------------|------------|----------------------|--------------------------------------------|
| output_path | o          | string     | output/              | Папка для сохранения скачанных изображений |
| clear-files |            | bool       | true                 | Очистить папку перед скачиванием           |

## Сайты

- [manga.fod.fujitv.co.jp](https://manga.fod.fujitv.co.jp/)
- [shonenjumpplus.com](https://shonenjumpplus.com/)
- [comic-walker.com](https://comic-walker.com/) _(без авторизации)_
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
- [storia.takeshobo.co.jp](https://storia.takeshobo.co.jp)

## Настройки

У некоторых парсеров есть настройки, они получаются из `settings.json` в формате ключ/объект. Где ключ — ID/домен парсера, а
объект — объект:

```json
{
  "shonenjumpplus.com": {
    "session": "ISqIN0B2M7zQSf7loxZhxCeC7l23nD2ckV"
  },
  "manga.fod.fujitv.co.jp": {
    "session": "YKt0Ab66gxMxqgvtRXx5takTSuz4np",
    "saveOriginal": false,
    "tryPurchaseBook": false
  }
}
```

Все сайты имеют обязательный параметр `session` для авторизации.

### Дополнительные настройки

| Парсер                    | Параметр        | Тип    | Описание                                                                                                     |
|---------------------------|-----------------|--------|--------------------------------------------------------------------------------------------------------------|
|                           | session         | string | Уникальный идентификатор сеанса или сессии.                                                                  |
| **manga.fod.fujitv.co.jp** | saveOriginal    | bool   | Сохранять ли пазлы вместо декадировоного изображения                                                         |
|                           | tryPurchaseBook | bool   | Проверяет, есть ли книга в списке бесплатных книг, если есть — парсер купит её и загрузит полноценную версию |
