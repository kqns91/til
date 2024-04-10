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

/config/app.php に追加 と公式ドキュメントにあるが、これをおこなってからプロバイダーを登録するとエラーが発生する。バージョン的に不要？

```php
'providers' => [
    // ...省略
    Spatie\Permission\PermissionServiceProvider::class,
],
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

権限を追加するモデルにトレイトを追加。

```php
use Spatie\Permission\Traits\HasRoles; // 追加

class User extends Authenticatable
{

    use HasRoles; // 追加

    // ...省略
}
```

権限のマスターデータ用のシーダーを作成

```sh
<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use Spatie\Permission\Models\Role; // 追加
use Spatie\Permission\Models\Permission; // 追加
use App\Models\User; // 追加

class PermissionMasterDataSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        // 管理者ユーザー作成
        $user = User::create([
            'name' => '管理者ユーザー',
            'email' => 'admin@sample.com',
            'password' => bcrypt('password'),
        ]);

        // 管理者ロール作成
        $adminRole = Role::create(['name' => 'admin']);

        $full_access_permission = Permission::create(['name' => 'full_access']);

        // 管理者ロールに permission を付与
        $adminRole->givePermissionTo($full_access_permission);

        // 管理者ユーザーに管理者ロールを付与
        $user->assignRole($adminRole);

        // 一般ユーザー作成
        User::create([
            'name'=> '一般ユーザー',
            'email'=> 'staff@sample.com',
            'password'=> bcrypt('password'),
        ]);
    }
}
```

シーダーを実行

```sh
php artisan db:seed --class=PermissionMasterDataSeeder
```

いい感じにデータが登録されている

```sql
postgres=# select * from users;
 id |      name      |      email       | email_verified_at |                           password                           | remember_token |     created_at      |     updated_at      
----+----------------+------------------+-------------------+--------------------------------------------------------------+----------------+---------------------+---------------------
  1 | 管理者ユーザー | admin@sample.com |                   | $2y$12$AQYxyRKnTh1bCTZiS..ZP.7WzRqWuRMDCz5e7Iw7ZRI/QlnsxSxSS |                | 2024-04-10 15:00:41 | 2024-04-10 15:00:41
  2 | 一般ユーザー   | staff@sample.com |                   | $2y$12$jLFtipIiIadFMID8aUrg7ePOapB7eJ5Vt9uLwunIZTI330WkESCBC |                | 2024-04-10 15:00:41 | 2024-04-10 15:00:41
(2 rows)

postgres=# 
postgres=# select * from roles;
 id | name  | guard_name |     created_at      |     updated_at      
----+-------+------------+---------------------+---------------------
  1 | admin | web        | 2024-04-10 15:00:41 | 2024-04-10 15:00:41
(1 row)

postgres=# select * from permissions;
 id |    name     | guard_name |     created_at      |     updated_at      
----+-------------+------------+---------------------+---------------------
  1 | full_access | web        | 2024-04-10 15:00:41 | 2024-04-10 15:00:41
(1 row)

postgres=# select * from model_has_roles ;
 role_id |   model_type    | model_id 
---------+-----------------+----------
       1 | App\Models\User |        1
(1 row)

postgres=# select * from model_has_permissions ;
 permission_id | model_type | model_id 
---------------+------------+----------
(0 rows)

postgres=# select * from permissions_id_seq ;
 last_value | log_cnt | is_called 
------------+---------+-----------
          1 |      32 | t
(1 row)

postgres=# select * from roles_id_seq ;
 last_value | log_cnt | is_called 
------------+---------+-----------
          1 |      32 | t
(1 row)

postgres=# select * from role_has_permissions ;
 permission_id | role_id 
---------------+---------
             1 |       1
(1 row)
```

