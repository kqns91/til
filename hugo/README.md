# Hugo

## インストール

```cmd
$ brew install hugo
```

## セットアップ

とりあえず、テーマ通りのサイトをローカルに立ち上げる。

```cmd
$ hugo new site mypage
$ cd mypage
```

作成された`mypage`に移動する。

[公式テーマ](https://themes.gohugo.io/)から気に入ったものを選ぶ。

```cmd
$ git clone https://github.com/hugo-sid/hugo-blog-awesome.git themes/hugo-blog-awesome
```

`./themes/{テーマ名}/exampleSite`配下のディレクトリやファイルをmypage直下に上書きする。

```cmd
$ hugo server -D
```

起動し、`http://localhost:1313`にアクセスする。

