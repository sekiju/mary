application {
  check_updates = true
  max_parallel_downloads = 4
  max_parallel_chapters = 1
}

output {
  directory      = "downloads"
  clean_on_start = false
  file_format    = "auto" # auto | jpeg | png | avif | webp
}

site {
  "shonenjumpplus.com" {
    cookie = ""
  }

  "comic-zenon.com" {
    cookie = ""
  }

  "pocket.shonenmagazine.com" {
    cookie = ""
  }

  "comic-gardo.com" {
    cookie = ""
  }

  "magcomi.com" {
    cookie = ""
  }

  "comic-action.com" {
    cookie = ""
  }

  "comic-days.com" {
    cookie = ""
  }

  "kuragebunch.com" {
    cookie = ""
  }

  "viewer.heros-web.com" {
    cookie = ""
  }

  "www.cmoa.jp" {
    cookie = ""
  }

  "www.corocoro.jp" {
    cookie = ""
  }
}
