# Laravel

## 環境構築

```sh
// インストール
brew install php
brew install composer
```

```
// プロジェクト作成
composer create-project laravel/laravel:^11.0 example-app
```

## permission

権限管理機能を簡単に実装できるパッケージ

```sh
composer require spatie/laravel-permission
```

プロバイダーを登録

```sh
php artisan vendor:publish --provider="Spatie\Permission\PermissionServiceProvider"
// 2024_04_10_141815_create_permission_tables.php が作成された
```

マイグレーション

```sh
php artisan config:clear
php artisan migrate
```

```sql
// マイグレーション前
postgres=# \d
                  List of relations
 Schema |         Name          |   Type   |  Owner   
--------+-----------------------+----------+----------
 public | cache                 | table    | postgres
 public | cache_locks           | table    | postgres
 public | failed_jobs           | table    | postgres
 public | failed_jobs_id_seq    | sequence | postgres
 public | job_batches           | table    | postgres
 public | jobs                  | table    | postgres
 public | jobs_id_seq           | sequence | postgres
 public | migrations            | table    | postgres
 public | migrations_id_seq     | sequence | postgres
 public | password_reset_tokens | table    | postgres
 public | sessions              | table    | postgres
 public | users                 | table    | postgres
 public | users_id_seq          | sequence | postgres
(13 rows)

// マイグレーション後
postgres=# \d
                  List of relations
 Schema |         Name          |   Type   |  Owner   
--------+-----------------------+----------+----------
 public | cache                 | table    | postgres
 public | cache_locks           | table    | postgres
 public | failed_jobs           | table    | postgres
 public | failed_jobs_id_seq    | sequence | postgres
 public | job_batches           | table    | postgres
 public | jobs                  | table    | postgres
 public | jobs_id_seq           | sequence | postgres
 public | migrations            | table    | postgres
 public | migrations_id_seq     | sequence | postgres
 public | model_has_permissions | table    | postgres new!
 public | model_has_roles       | table    | postgres new!
 public | password_reset_tokens | table    | postgres
 public | permissions           | table    | postgres new!
 public | permissions_id_seq    | sequence | postgres new!
 public | role_has_permissions  | table    | postgres new!
 public | roles                 | table    | postgres new!
 public | roles_id_seq          | sequence | postgres new!
 public | sessions              | table    | postgres
 public | users                 | table    | postgres
 public | users_id_seq          | sequence | postgres
(20 rows)
```

