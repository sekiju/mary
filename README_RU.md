# Mary

> [!IMPORTANT]
> The project is made for educational purposes. If you believe your rights are being violated, contact main contributor.

## Использование

Скачайте последнюю версию из [Github Releases](https://github.com/sekiju/mary/releases/latest). Далее откройте приложение, оно попросит ввести вас ссылку на главу.<br>
Также всё можно запустить через терминал с указанием ссылки:

```bash
mary.exe https://shonenjumpplus.com/magazine/4856001361564051266
```

> [!IMPORTANT]
> Платные главы можно скачать только после того, как вы её купили и указали в [настройках Mary](#настройки)!

## Веб-сайты

- [manga.fod.fujitv.co.jp](https://manga.fod.fujitv.co.jp/)
- [comic-walker.com](https://comic-walker.com/)
- [www.pixiv.net](https://www.pixiv.net/)
- [comic.webnewtype.com](https://comic.webnewtype.com)
- [manga.bilibili.com](https://manga.bilibili.com)


- GigaViewer группа:
    - [shonenjumpplus.com](https://shonenjumpplus.com/)
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
    - [comicbushi-web.com](https://comicbushi-web.com)


- SpeedBinb группа:
    - [storia.takeshobo.co.jp](https://storia.takeshobo.co.jp)
    - [www.comic-valkyrie.com](https://www.comic-valkyrie.com)
    - [cmoa.jp](https://cmoa.jp)
    - [yanmaga.jp](https://yanmaga.jp)
    - [comic-meteor.jp](https://comic-meteor.jp)

## Настройки

Приложение использует [YAML](https://yaml.org/spec/1.2.2/) формат для хранения настроек. Все настройки прописываются
в `config.yaml`. Пример:

```yaml
settings:
  debug:
    enable: true
    url: https://shonenjumpplus.com/episode/14079602755375556618
  outputPath: output/
  clearOutputFolder: true
  threads: 6

sites:
  manga.fod.fujitv.co.jp:
    session: VWdx8id9R0XHjVpvs7s754CxGJpBBl9ZCHCqL1yF
    purchase_free_books: true
```

У приложения есть настройки по умолчанию, которые работают без создания конфига `config.yaml`:

```yaml
settings:
  debug:
    enable: false
  outputPath: output/
  clearOutputFolder: true
  threads: T # T - количество логических ядер процессора
```

И так как `settings` не обязательный параметр, то вы можете прописать только `sites`.

Каждый веб-сайт из `sites` имеет обязательный параметр `session` для авторизации.<br>
Чтобы добавить новый сайт в `sites`, сделайте отступ и
введите [доменное имя](https://blog.skillfactory.ru/wp-content/uploads/2023/02/domen-4-3253604.png) сайта. Ещё раз
сделайте отступ и добавьте ключ `session`.
Подробнее о том как получить `session` для каждого веб-сайт можно посмотреть [здесь](#как-получить-session).

### Как получить `session`?

Большинство сайтов хранит сессии (вашу авторизацию) в Cookie. Для того чтобы их получить, я
пользуюсь [Cookie-Editor](https://cookie-editor.com).<br>
Некоторые веб-сайты требуют несколько Cookie файлов для полноценной работы авторизации, поэтому иногда придется
экспортировать все Cookie в ввиде Header (заголовка).

| Веб-сайт / Группа веб-сайтов | Тип    | Cookie Name / Header Name | Description                                                        |
|------------------------------|--------|---------------------------|--------------------------------------------------------------------|
| GigaViewer группа            | Cookie | glsc                      |                                                                    |
| Fod                          | Header | zk-session-key            | Токен из Android приложения (можно получить только отреверсив его) |
| Cmoa                         | Header | Cookie                    | Экспортируйте все куки как Header String                           |