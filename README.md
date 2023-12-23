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
Все сессии и токены авторизации прописываются в `settings.json` в формате ключ/объект:
```json
{
  "shonenjumpplus": {
    "Session": "ISqIN0B2M7zQSf7loxZhxCeC7l23nD2ckV"
  },
  "fod": {
    "Session": "YKt0Ab66gxMxqgvtRXx5takTSuz4np"
  }
}
```

Все сайты имеют обязательный параметр `Session`.