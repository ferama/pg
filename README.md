# pg

`pg` is an alternative CLI tool for PostgreSQL

It aims to be simple to use with less commands to remember

## Features

* very easy to use
* support for multiple database connections
* handy subcommands for database structure exploration
* commands autocomplete ( we all love `tab -> tab` right? :) )
* embedded sql editor for more advanced exploration
* query history navigation

## Handy subcommands

![pg subcommands](https://raw.githubusercontent.com/ferama/pg/main/media/commands.png)


## Sql Editor

![pg sql editor](https://raw.githubusercontent.com/ferama/pg/main/media/editor.png)

## Config file example

Put this in your $HOME/.pg dir

```yaml
connections:
  - name: db1
    url: postgres://user1:pass1@host1:5432/db1
  - name: db2
    url: postgres://user2:pass2@host2:5432/db2
```