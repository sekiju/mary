#

В структуре класса всегда создаём `Storage`:
```go
type ShonenJumpPlus struct {
    Storage readers.ReaderStorage
}
```

#

Обязательно делаем указатели памяти `*` в классах:

```go
func (s *ShonenJumpPlus) SetSession(str string) {
    s.Storage.Session = &str
}
```

Иначе все изменения данных просто потеряются.