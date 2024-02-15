# Mary

> [!IMPORTANT]
> The project is made for educational purposes. If you believe your rights are being violated, contact main contributor.

## Usage

```bash
559.exe https://shonenjumpplus.com/magazine/4856001361564051266
```

## Websites

- [manga.fod.fujitv.co.jp](https://manga.fod.fujitv.co.jp/)
- [comic-walker.com](https://comic-walker.com/)
- [www.pixiv.net](https://www.pixiv.net/)
- [comic.webnewtype.com](https://comic.webnewtype.com)
- [manga.bilibili.com](https://manga.bilibili.com)


- GigaViewer group:
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


- SpeedBinb group:
    - [storia.takeshobo.co.jp](https://storia.takeshobo.co.jp)
    - [www.comic-valkyrie.com](https://www.comic-valkyrie.com)
    - [cmoa.jp](https://cmoa.jp)
    - [yanmaga.jp](https://yanmaga.jp)
    - [comic-meteor.jp](https://comic-meteor.jp)

## Config

The application uses [YAML](https://yaml.org/spec/1.2.2/) format for storing settings. All parameters are specified
in `config.yaml`. Example file with all parameters:

```yaml
settings:
  debug:
    enable: true
    url: https://shonenjumpplus.com/episode/14079602755375556618
  outputPath: output/
  clearOutputFolder: true
  threads: 6
  targetMethod: chapter # chapter OR book

sites:
  manga.fod.fujitv.co.jp:
    session: VWdx8id9R0XHjVpvs7s754CxGJpBBl9ZCHCqL1yF
    purchase_free_books: true
```

In case `config.yaml` is absent, the following configuration will be automatically used:

```yaml
settings:
  debug:
    enable: false
  outputPath: output/
  clearOutputFolder: true
  threads: T # Where T is number of logical CPUs usable by the current process
```

Since `settings` are not mandatory, you can leave only the `sites` field.

All websites have a mandatory parameter `session` for authentication.<br>
To add a website to the configuration, add its domain name to `sites`. Then add its `session` afterwards. Details on how
to obtain the `session` for each website will be provided below.

### How to get `session`?

To extract cookies from websites, I use [Cookie-Editor](https://cookie-editor.com).

| Website / Websites group | Type   | Cookie Name / Header Name | Description                                   |
|--------------------------|--------|---------------------------|-----------------------------------------------|
| GigaViewer group         | Cookie | glsc                      |                                               |
| Fod                      | Header | zk-session-key            | Token from Android App / Reverses Engineering |
| Cmoa                     | Header | Cookie                    | Export all cookies as header with extension   |