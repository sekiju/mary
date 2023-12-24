# 559

## Запуск

```bash
559.exe https://shonenjumpplus.com/magazine/4856001361564051266 -o output/
```

## Сайты

* [manga.fod.fujitv.co.jp](https://manga.fod.fujitv.co.jp/)
* [shonenjumpplus.com](https://shonenjumpplus.com/)
* [comic-walker.com](https://comic-walker.com/) _(без авторизации)_

## Настройки

У некоторых парсеров есть настройки, они получаются из `settings.json` в формате ключ/объект. Где ключ — ID парсера, а
объект — объект:

```json
{
  "shonenjumpplus": {
    "session": "ISqIN0B2M7zQSf7loxZhxCeC7l23nD2ckV"
  },
  "fod": {
    "session": "YKt0Ab66gxMxqgvtRXx5takTSuz4np",
    "saveOriginal": false,
    "tryPurchaseBook": false
  }
}
```

Все сайты имеют обязательный параметр `session` для авторизации.

### Дополнительные настройки

| Парсер                     | Параметр        | Тип    | Описание                                                                                                     |
|----------------------------|-----------------|--------|--------------------------------------------------------------------------------------------------------------|
|                            | session         | string | Уникальный идентификатор сеанса или сессии.                                                                  |
| **manga.fod.fujitv.co.jp** | saveOriginal    | bool   | Сохранять ли пазлы вместо декадировоного изображения                                                         |
|                            | tryPurchaseBook | bool   | Проверяет, есть ли книга в списке бесплатных книг, если есть — парсер купит её и загрузит полноценную версию |
