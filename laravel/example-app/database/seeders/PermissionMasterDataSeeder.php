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
