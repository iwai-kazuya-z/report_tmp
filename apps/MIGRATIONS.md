# Dorapita DBマイグレーション機能

最終更新: 2025-12-25

---

## 概要

DorapitaプロジェクトではCakePHP Migrations（Phinxベース）を使用したDBスキーマ管理機能が実装されています。

| 項目 | 詳細 |
|------|------|
| **フレームワーク** | CakePHP 4.x Migrations（Phinxベース） |
| **公式ドキュメント** | https://book.cakephp.org/migrations/4/en/index.html |
| **マイグレーションファイル数** | cadm.dorapita.com: 63ファイル<br>dorapita.com: 5ファイル<br>dora-pt.jp: 3ファイル |
| **最新マイグレーション** | 2025-08-25（`CreateApplicationQueue`） |
| **ファイル配置** | `{プロジェクト}/config/Migrations/*.php` |

---

## マイグレーションファイルの構成

### ディレクトリ構造

```
dorapita_code/
├── cadm.dorapita.com/
│   └── config/
│       └── Migrations/
│           ├── 20240314225705_ModifyShokushuItemId.php
│           ├── 20241220094338_CreateApplications.php
│           ├── 20250825085430_CreateApplicationQueue.php
│           └── ... (63ファイル)
├── dorapita.com/
│   └── config/
│       └── Migrations/
│           ├── 20240508054115_UpdateSystemMails.php
│           ├── 20250320043616_AddApplicationIdToSendMails.php
│           └── ... (5ファイル)
└── dora-pt.jp/
    └── config/
        └── Migrations/
            ├── 20240305025309_ViewLogChangeIpField.php
            ├── 20240523092748_CreateRecopAuth.php
            └── 20240826021218_CreatePartnerClients.php
```

### ファイル命名規則

```
{タイムスタンプ}_{マイグレーション名}.php

例:
20250825085430_CreateApplicationQueue.php
  ↓
2025年8月25日 08時54分30秒 - CreateApplicationQueue
```

---

## マイグレーションファイルの例

### 1. テーブル作成（`CreateApplications.php`）

```php
<?php
declare(strict_types=1);

use Migrations\AbstractMigration;

class CreateApplications extends AbstractMigration
{
    /**
     * Change Method.
     *
     * More information on this method is available here:
     * https://book.cakephp.org/phinx/0/en/migrations.html#the-change-method
     * @return void
     */
    public function change(): void
    {
        $table = $this->table('applications', [
            'comment' => 'Table to store information related to applications'
        ]);

        $table->addColumn('user_id', 'integer', [
            'null'    => true,
            'default' => 0,
            'comment' => 'Applicant ID (foreign key to users table)'
        ]);

        $table->addColumn('company_id', 'integer', [
            'null'    => false,
            'default' => 0,
            'comment' => 'Company ID (foreign key to companies table)'
        ]);

        $table->addColumn('recruit_id', 'integer', [
            'null'    => false,
            'default' => 0,
            'comment' => 'Recruitment information ID (foreign key to recruit table)'
        ]);

        $table->addColumn('apply_date', 'date', [
            'null'    => false,
            'default' => null,
            'comment' => 'Application date'
        ]);

        $table->addColumn('current_state', 'tinyinteger', [
            'default' => 1,
            'limit'   => 4,
            'comment' => 'Current application status'
        ]);

        // インデックス追加
        $table->addIndex(['user_id', 'company_id', 'recruit_id'], [
            'name' => 'idx_user_company_recruit_id'
        ]);
        $table->addIndex(['recruit_id'], ['name' => 'idx_recruit_id']);
        $table->addIndex(['user_id'], ['name' => 'idx_user_id']);
        $table->addIndex(['company_id'], ['name' => 'idx_company_id']);

        $table->create();
    }
}
```

### 2. インデックス追加（`AddIndexForSelection.php`）

```php
<?php
declare(strict_types=1);

use Migrations\AbstractMigration;

class AddIndexForSelection extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('selections');

        // 複数のインデックスを一度に追加
        $table->addIndex(['company_id'], ['name' => 'idx_company_id'])
            ->addIndex(['recruit_entry_id'], ['name' => 'idx_recruit_entry_id'])
            ->addIndex(['recruit_id'], ['name' => 'idx_recruit_id'])
            ->addIndex(['age_id'], ['name' => 'idx_age_id'])
            ->addIndex(['act'], ['name' => 'idx_act', 'limit' => 10])  // TEXT列には長さ制限が必要
            ->addIndex(['tel'], ['name' => 'idx_tel', 'limit' => 30])  // TEXT列には長さ制限が必要
            ->addIndex(['apply_day'], ['name' => 'idx_apply_day'])
            ->update();
    }
}
```

### 3. カラム追加（`AddColumnsToRecruits.php`）

```php
<?php
declare(strict_types=1);

use Migrations\AbstractMigration;

class AddColumnsToRecruits extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('recruits');

        $table->addColumn('introduction_type', 'string', [
            'default' => 'normal',
            'limit' => 50,
            'null' => false,
            'comment' => 'Introduction type (normal, easy_selection, etc.)'
        ]);

        $table->addColumn('is_bizukomi', 'boolean', [
            'default' => false,
            'null' => false,
            'comment' => 'Flag for bizukomi service'
        ]);

        $table->update();
    }
}
```

---

## マイグレーションコマンド

### 基本コマンド

```bash
# プロジェクトディレクトリに移動
cd /var/www/cadm.dorapita.com

# マイグレーションステータス確認
php bin/cake migrations status

# マイグレーション実行（全ての未実行マイグレーション）
php bin/cake migrations migrate

# 詳細ログ付きで実行
php bin/cake migrations migrate -vvv

# 特定のバージョンまで実行
php bin/cake migrations migrate -t 20250825085430

# ロールバック（1つ前に戻す）
php bin/cake migrations rollback

# 完全にロールバック
php bin/cake migrations rollback -t 0

# 新規マイグレーションファイル作成
php bin/cake bake migration CreateUsers

# シードデータ投入
php bin/cake migrations seed --seed TwilioNumbersSeed
```

### 環境別実行

```bash
# ローカル開発環境
export DATABASE_URL='mysql://root:@127.0.0.1:3306/dorapita1804'
php bin/cake migrations migrate

# ステージング環境
export DATABASE_URL='postgres://dorauser2022:***@staging-db:5432/dorapita'
php bin/cake migrations migrate

# 本番環境
export DATABASE_URL='postgres://dorauser2022:***@pg-120011:5432/dorapita'
php bin/cake migrations migrate -vvv
```

---

## 本番環境での使用実績

### Issue #3294: 本番DBマイグレーション運用の非特権化

**実行日**: 2025-08-26
**目的**: 所有者正規化＋短時間メンテでCREATE INDEX適用

#### 実行手順

```bash
# メンテナンスON（アプリ/ワーカー停止）

# 環境変数設定
export DATABASE_URL='postgres://dorauser2022:***@127.0.0.1:5432/dorapita'

# プロジェクトディレクトリに移動
cd /opt/dorapita/dorapita.com

# マイグレーション前のステータス確認
php bin/cake migrations status

# マイグレーション実行（詳細ログ付き）
time php bin/cake migrations migrate -vvv

# マイグレーション後のステータス確認
php bin/cake migrations status

# メンテナンスOFF
```

#### 実行結果

```
real    0m0.055s
user    0m0.018s
sys     0m0.035s
```

**成果**:
- 所有者をpostgres→dorauser2022に正規化
- 特権不要でマイグレーション実行可能に
- 55ミリ秒で完了（短時間メンテ成功）

---

## シードファイルの使用例

### Seedファイルの配置

```
config/Seeds/
├── TwilioNumbersSeed.php
├── StagingTwilioNumbersSeed.php
└── DeleteOldDataViewLogsSeed.php
```

### Seedファイルの例（`TwilioNumbersSeed.php`）

```php
<?php
declare(strict_types=1);

use Migrations\AbstractSeed;

class TwilioNumbersSeed extends AbstractSeed
{
    public function run(): void
    {
        $data = [
            [
                'phone_number' => '+81312345678',
                'country_code' => 'JP',
                'is_active'    => true,
                'created'      => date('Y-m-d H:i:s'),
            ],
            // ... 他のデータ
        ];

        $table = $this->table('twilio_numbers');
        $table->insert($data)->save();
    }
}
```

### Seed実行コマンド

```bash
# 本番環境でのSeed実行例（release-note.mdより）
cd /var/www/dorapita.com
./bin/cake migrations seed --seed TwilioNumbersSeed

# ステージング環境でのSeed実行例
cd /var/www/dorapita.com
./bin/cake migrations seed --seed StagingTwilioNumbersSeed
```

---

## 問題点と課題

### ⚠️ マイグレーション機能が完全には活用されていない

| 問題 | 証拠 | 影響 |
|------|------|------|
| **手動SQL更新が頻発** | issue #3088, #3927, #3864など**11件** | ・マイグレーション履歴が不完全<br>・監査証跡が残らない<br>・ロールバックが困難 |
| **本番とステージングで手順が異なる** | 本番: 手動SQL、ステージング: マイグレーション | ・デプロイ手順の不統一<br>・環境差異の発生リスク |
| **ドキュメント不足** | README.mdにマイグレーション実行手順なし | ・属人化リスク<br>・新規メンバーのオンボーディング困難 |
| **CI/CD未統合** | GitHub Actionsでマイグレーション自動実行なし | ・手動実行ミスのリスク<br>・デプロイ自動化が不完全 |

### 手動SQL更新の例

#### Issue #3927: 本番SQL更新「勤務開始時間」の手動追加

```sql
-- 手動SQL更新（マイグレーションファイルなし）
INSERT INTO user_profile_work_start_times (name, created, modified)
VALUES ('00:00〜', NOW(), NOW());
```

**問題点**:
- マイグレーションファイルが作成されていない
- phinxlogテーブルに履歴が残らない
- ステージング環境と本番環境で実行タイミングが異なる

#### Issue #3864: 本番SQL更新（複数テーブル）

```sql
-- 手動SQL更新（マイグレーションファイルなし）
UPDATE user_profile_menkyo_items SET ... WHERE ...;
UPDATE user_resumes SET ... WHERE ...;
UPDATE user_profiles SET ... WHERE ...;
UPDATE user_tokens SET ... WHERE ...;
UPDATE users SET ... WHERE ...;
```

**問題点**:
- 複数テーブルの更新が一括実行
- トランザクション管理が不明確
- ロールバック手順が未定義

---

## Rails `db:migrate` との比較

| 機能 | Rails `db:migrate` | CakePHP Migrations | 備考 |
|------|-------------------|-------------------|------|
| **マイグレーションファイル** | ✅ `db/migrate/*.rb` | ✅ `config/Migrations/*.php` | |
| **タイムスタンプ形式** | ✅ `20250101120000_create_users.rb` | ✅ `20250101120000_CreateUsers.php` | |
| **アップ/ダウン** | ✅ `up` / `down` | ✅ `change` (自動リバース) | CakePHPは`change`メソッドのみでOK |
| **ステータス確認** | ✅ `rails db:migrate:status` | ✅ `cake migrations status` | |
| **ロールバック** | ✅ `rails db:rollback` | ✅ `cake migrations rollback` | |
| **シードデータ** | ✅ `rails db:seed` | ✅ `cake migrations seed` | |
| **環境別管理** | ✅ `RAILS_ENV=production` | ✅ `DATABASE_URL=...` | |
| **トランザクション** | ✅ 自動 | ✅ 自動 | |
| **外部キー制約** | ✅ `add_foreign_key` | ✅ `addForeignKey` | |
| **マイグレーションバージョン** | ✅ `schema_migrations` テーブル | ✅ `phinxlog` テーブル | |
| **本番運用の標準化** | ✅ **一般的** | ⚠️ **Dorapitaでは部分的** | Dorapitaは手動SQLも併用 |

---

## 推奨アクション

### Phase 1: 緊急対応（1週間以内）

#### 1.1 本番SQL更新のマイグレーションファイル化

**対象**: 過去11件の手動SQL更新issue

| Issue# | タイトル | 対応 |
|--------|---------|------|
| #3927 | 本番SQL更新「勤務開始時間」の手動追加 | マイグレーションファイル作成 |
| #3864 | 本番SQL更新（複数テーブル） | マイグレーションファイル作成 |
| #3839 | 本番SQL更新（特定原稿削除） | マイグレーションファイル作成 |
| ... | ... | ... |

**手順**:
```bash
# 1. マイグレーションファイル作成
cd /var/www/cadm.dorapita.com
php bin/cake bake migration AddWorkStartTimesToUserProfileWorkStartTimes

# 2. マイグレーションファイル編集
vim config/Migrations/20250101000000_AddWorkStartTimesToUserProfileWorkStartTimes.php

# 3. ステージング環境で検証
php bin/cake migrations migrate

# 4. 本番環境で実行
php bin/cake migrations migrate -vvv
```

#### 1.2 デプロイ手順書の作成

**作成すべきドキュメント**:
- `docs/MIGRATION.md` - マイグレーション実行手順
- `docs/DEPLOYMENT.md` - デプロイフロー全体
- `.github/PULL_REQUEST_TEMPLATE.md` - マイグレーションチェックリスト

**テンプレート例**（`.github/PULL_REQUEST_TEMPLATE.md`）:

```markdown
## チェックリスト

### マイグレーション
- [ ] マイグレーションファイルを作成した
- [ ] ローカル環境でマイグレーションを実行し、動作確認した
- [ ] ロールバック手順を確認した
- [ ] マイグレーション実行時間を計測した（本番メンテナンス時間の見積もり）

### スキーマ変更
- [ ] インデックス追加の場合、CONCURRENT指定を検討した
- [ ] 外部キー制約の追加の場合、既存データの整合性を確認した
- [ ] カラム追加の場合、デフォルト値を設定した

### 本番デプロイ
- [ ] ステージング環境でマイグレーションを実行し、問題なく完了した
- [ ] 本番メンテナンス時間を確保した
- [ ] ロールバック手順を準備した
```

---

### Phase 2: 短期対応（1ヶ月以内）

#### 2.1 CI/CDパイプラインへの組み込み

**GitHub Actions例**:

```yaml
# .github/workflows/deploy-staging.yml
name: Deploy to Staging

on:
  push:
    branches:
      - staging

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup PHP
        uses: shivammathur/setup-php@v2
        with:
          php-version: '8.1'

      - name: Install dependencies
        run: composer install --no-dev --optimize-autoloader

      - name: Run migrations
        env:
          DATABASE_URL: ${{ secrets.STAGING_DATABASE_URL }}
        run: |
          cd cadm.dorapita.com
          php bin/cake migrations migrate -vvv

      - name: Deploy to Cloud Run
        run: |
          gcloud run deploy cadm-dorapita-com \
            --image gcr.io/${{ secrets.GCP_PROJECT }}/cadm-dorapita-com:${{ github.sha }} \
            --region asia-northeast1
```

#### 2.2 マイグレーション実行ログの一元管理

**Cloud Logging統合**:

```php
// config/bootstrap.php
use Cake\Log\Log;

// マイグレーション実行時のログをCloud Loggingに送信
Log::setConfig('migrations', [
    'className' => 'GoogleCloudLogging',
    'path' => LOGS,
    'levels' => ['info', 'warning', 'error'],
    'scopes' => ['migrations'],
]);
```

---

### Phase 3: 中長期対応（3ヶ月以内）

#### 3.1 マイグレーション履歴の可視化

**Slack通知の実装**:

```bash
# bin/cake migrations migrate 実行後にSlack通知
MIGRATION_STATUS=$(php bin/cake migrations status)
curl -X POST $SLACK_WEBHOOK_URL \
  -H 'Content-Type: application/json' \
  -d "{\"text\":\"本番環境マイグレーション完了\n\`\`\`\n$MIGRATION_STATUS\n\`\`\`\"}"
```

#### 3.2 マイグレーションレビューフローの確立

**レビュー基準**:
1. **パフォーマンス影響評価**
   - インデックス追加: CONCURRENT指定の有無
   - テーブルロック時間の見積もり
   - 本番メンテナンス時間の妥当性

2. **データ整合性チェック**
   - 外部キー制約の追加: 既存データの整合性
   - NOT NULL制約の追加: デフォルト値の設定
   - カラム削除: 参照箇所の確認

3. **ロールバック戦略**
   - ロールバック手順の明記
   - データ損失リスクの評価
   - 緊急時の対応手順

---

## ベストプラクティス

### 1. マイグレーションファイルの命名規則

```
✅ Good:
20250825085430_CreateApplicationQueue.php
20250825090000_AddIndexToSelections.php
20250825091500_RemoveUnusedColumnsFromRecruits.php

❌ Bad:
20250825085430_Migration1.php
20250825090000_Fix.php
20250825091500_Update.php
```

### 2. 1ファイル1操作の原則

```php
// ✅ Good: 1つのマイグレーションで1つの操作
class CreateApplications extends AbstractMigration
{
    public function change(): void
    {
        // applicationsテーブルの作成のみ
        $table = $this->table('applications');
        // ...
        $table->create();
    }
}

// ❌ Bad: 1つのマイグレーションで複数の操作
class UpdateDatabase extends AbstractMigration
{
    public function change(): void
    {
        // applicationsテーブル作成
        $table1 = $this->table('applications');
        $table1->create();

        // selectionsテーブル更新
        $table2 = $this->table('selections');
        $table2->addColumn('new_column', 'string')->update();

        // recruitsテーブルにインデックス追加
        $table3 = $this->table('recruits');
        $table3->addIndex(['status'])->update();
    }
}
```

### 3. インデックス追加時のCONCURRENT指定

```php
// PostgreSQL + 大きなテーブルの場合
class AddIndexToRecruits extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('recruits');

        // CONCURRENT指定（ロックなしでインデックス作成）
        $table->addIndex(['status', 'created_at'], [
            'name' => 'idx_recruits_status_created',
            // ⚠️ CakePHP Migrationsは直接CONCURRENTをサポートしていない
            // 手動でexecute()を使用する必要がある
        ])->update();
    }

    // PostgreSQLのCREATE INDEX CONCURRENTLYを使用する場合
    public function up(): void
    {
        $this->execute('CREATE INDEX CONCURRENTLY idx_recruits_status_created ON recruits (status, created_at)');
    }

    public function down(): void
    {
        $this->execute('DROP INDEX CONCURRENTLY IF EXISTS idx_recruits_status_created');
    }
}
```

### 4. NOT NULL制約追加時のデフォルト値設定

```php
// ✅ Good: デフォルト値を設定
class AddStatusToApplications extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('applications');
        $table->addColumn('status', 'string', [
            'default' => 'pending',  // デフォルト値を設定
            'limit' => 50,
            'null' => false,
            'after' => 'apply_date',
        ])->update();
    }
}

// ❌ Bad: デフォルト値なしでNOT NULL
class AddStatusToApplications extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('applications');
        $table->addColumn('status', 'string', [
            'limit' => 50,
            'null' => false,  // 既存レコードがある場合エラー
        ])->update();
    }
}
```

---

## トラブルシューティング

### 問題1: マイグレーション実行時に "must be owner" エラー

**エラーメッセージ**:
```
ERROR 42501: must be owner of table phinxlog
```

**原因**: phinxlogテーブルの所有者がpostgresで、dorauser2022に権限がない

**解決方法**（issue #3294の対応）:
```sql
-- postgres権限で実行
GRANT SELECT, INSERT, UPDATE, DELETE ON public.phinxlog TO dorauser2022;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE public.phinxlog_id_seq TO dorauser2022;
```

---

### 問題2: マイグレーション実行後にロールバックできない

**原因**: `change()`メソッドで非可逆的な操作を実行

**解決方法**: `up()`と`down()`メソッドを明示的に定義

```php
class DropUnusedTable extends AbstractMigration
{
    public function up(): void
    {
        $this->table('old_table')->drop()->save();
    }

    public function down(): void
    {
        // テーブル削除は非可逆的なので、ロールバック時に再作成
        $table = $this->table('old_table');
        $table->addColumn('id', 'integer', ['identity' => true])
            ->addColumn('name', 'string', ['limit' => 255])
            ->create();
    }
}
```

---

### 問題3: マイグレーション実行時間が長すぎる

**症状**: 本番メンテナンス時間内に完了しない

**原因**:
- 大きなテーブルに対するインデックス追加
- NOT NULL制約の追加でテーブルフルスキャン

**解決方法**:
1. **CREATE INDEX CONCURRENTLY を使用**（PostgreSQL）
2. **バッチ処理で段階的に実行**
3. **メンテナンス時間外に実行**（読み取り専用モード）

```php
// バッチ処理の例
class AddNotNullToLargeTable extends AbstractMigration
{
    public function up(): void
    {
        // 1. デフォルト値を設定（既存レコードに適用）
        $this->execute("UPDATE recruits SET status = 'active' WHERE status IS NULL");

        // 2. NOT NULL制約を追加
        $this->table('recruits')
            ->changeColumn('status', 'string', ['null' => false])
            ->update();
    }
}
```

---

## まとめ

### 現状

| 項目 | 状態 |
|------|------|
| **マイグレーション機能の存在** | ✅ CakePHP Migrations実装済み |
| **開発環境での利用** | ✅ 活発に使用（63ファイル） |
| **本番環境での利用** | ⚠️ 部分的（手動SQL更新が併用） |
| **ドキュメント** | ❌ 不足（README.mdに記載なし） |
| **CI/CD統合** | ❌ 未実装 |

### 推奨アクション（優先度順）

1. **P0（1週間）**: 全ての本番SQL更新をマイグレーションファイル化
2. **P0（1週間）**: デプロイ手順書の作成
3. **P1（1ヶ月）**: CI/CDパイプラインへの組み込み
4. **P1（1ヶ月）**: マイグレーション実行ログの一元管理
5. **P2（3ヶ月）**: マイグレーション履歴の可視化
6. **P2（3ヶ月）**: マイグレーションレビューフローの確立

### 期待される効果

| 効果 | 詳細 |
|------|------|
| **監査証跡の完全性** | 全てのスキーマ変更がphinxlogテーブルに記録 |
| **ロールバック可能性** | 問題発生時に迅速にロールバック可能 |
| **環境差異の削減** | ステージング・本番で同一手順を実行 |
| **デプロイ自動化** | CI/CDパイプラインで自動実行 |
| **属人化の解消** | ドキュメント化により誰でも実行可能 |

---

**参考資料**:
- CakePHP Migrations公式ドキュメント: https://book.cakephp.org/migrations/4/en/index.html
- Phinx公式ドキュメント: https://book.cakephp.org/phinx/0/en/index.html
- Issue #3294: 本番DBマイグレーション運用の非特権化
- Issue #3088, #3927, #3864: 本番SQL更新issue

**最終更新**: 2025-12-25
