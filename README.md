# koha-i18n

A little program that modifies Koha templates to wrap translatable text inside a
t() call, to make the text translatable by Koha::I18n

## Installation

```
go get github.com/jajm/koha-i18n
```

## Usage

```
koha-i18n [--in-place] FILES...

Options:
    --in-place      Modify files in place instead of printing to STDOUT
```

## License

MIT
