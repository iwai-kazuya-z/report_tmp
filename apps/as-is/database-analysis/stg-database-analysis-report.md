# STG データベース分析レポート

**作成日**: 2025-12-26
**目的**: Fixture作成戦略の策定

---

## エグゼクティブサマリー

STG環境のPostgreSQLとMySQLデータベースを分析し、以下の結論に達した:

### 主要な発見
1. **データ規模**: PostgreSQL 約8GB、MySQL 約150GB（ログテーブル含む）
2. **PII混在**: 個人情報（氏名、電話番号、メールアドレス等）が複数テーブルに存在
3. **ログテーブルの肥大化**: view_logs, change_logsが非常に大きく、Fixture化は不適切
4. **マスターデータの少なさ**: 都道府県、エリア等のマスターデータは数十行程度

### 推奨戦略
- **マスターデータ**: 全件コピー（加工不要）
- **トランザクションデータ**: 10-50件のサンプル抽出 + PII マスキング
- **ログデータ**: Fixture化しない（テスト時に動的生成）

---

## 1. PostgreSQL (dorapita) 分析

### データベース概要
- **Instance**: pg-120011
- **Database**: dorapita
- **総容量**: 約8GB（ログテーブル含む）

### テーブル分類と推奨事項

#### 1.1 マスターデータ（全件コピー推奨）

| テーブル | 行数 | サイズ | 用途 | Fixture戦略 |
|---------|------|--------|------|-------------|
| prefectures | 47 | 120 KB | 都道府県マスタ | **全件コピー** |
| areas | 9 | 64 KB | エリアマスタ | **全件コピー** |
| kodawari_items | 163 | 64 KB | こだわり条件マスタ | **全件コピー** |
| send_mail_states | 少数 | 32 KB | メール状態マスタ | **全件コピー** |
| jobtype_area_ranks | 少数 | 16 KB | 職種エリアランク | **全件コピー** |
| ng_messages | 少数 | 16 KB | NGメッセージ | **全件コピー** |

**理由**: 行数が少なく、参照整合性維持のため全件必要。PII含まず。

---

#### 1.2 トランザクションデータ（サンプル抽出 + マスキング）

| テーブル | 行数 | サイズ | 主なPII | Fixture戦略 |
|---------|------|--------|---------|-------------|
| **recruits** | 19,487 | 337 MB | 企業名、連絡先、住所 | **20-30件抽出** + 企業情報マスキング |
| **companies** | 2,425 | 4.9 MB | 企業名、担当者、電話、メール | **10-15件抽出** + 担当者情報マスキング |
| **entries** | 47,842 | 73 MB | **氏名、電話、メール、住所、生年月日** | **30-50件抽出** + 全PII マスキング |

##### PIIマスキング方針

**entries テーブル（最も機密性が高い）**:
| カラム | 元データ例 | マスキング例 |
|--------|-----------|-------------|
| name | 山田太郎 | テストユーザー001 |
| kana | ヤマダタロウ | テストユーザー001 |
| tel | 090-1234-5678 | 080-0000-0001 |
| mail | yamada@example.com | test001@example.com |
| zip | 123-4567 | 100-0001 |
| address1 | 東京都渋谷区... | 東京都千代田区1-1-1 |
| address2 | マンション名101 | テストビル101 |
| birthday | 1990-01-15 | 1990-01-01 |

**companies テーブル**:
| カラム | マスキング方針 |
|--------|----------------|
| name | テスト運送株式会社 |
| kana | テストウンソウ |
| catchcopy | そのまま（一般的な文言） |
| pr_txt | そのまま（一般的な文言） |
| zip, address | ダミー住所に置換 |
| site_url | https://example.com |
| c_tanto_mails | test-tanto@example.com |

**recruits テーブル**:
| カラム | マスキング方針 |
|--------|----------------|
| title | そのまま（求人タイトル） |
| c_nm | companies.nameから参照 |
| tel_num | 0120-000-001 |
| saiyo_tanto | 採用担当 太郎 |
| c_tanto_mails | companies参照 |

---

#### 1.3 ログデータ（Fixture化しない）

| テーブル | 行数 | サイズ | 用途 | Fixture戦略 |
|---------|------|--------|------|-------------|
| view_logs | 599 | 748 MB | 閲覧ログ | **Fixture化しない** |
| view_logs_copy_* | - | 3.2 GB × 2 | バックアップ | **除外** |
| send_mails | 270,648 | 776 MB | 送信メールログ | **Fixture化しない** |
| entry_logs | 19,391 | 4.4 MB | 応募ログ | **Fixture化しない** |
| use_numbers | 36,907 | 5.0 MB | 番号使用履歴 | **Fixture化しない** |

**理由**:
- ログデータは単体テストでは不要（E2Eテストで動的生成）
- サイズが大きく、Fixture化すると開発環境の起動が遅くなる
- 古いデータは参照されない

---

### 1.4 バックアップテーブル（除外）

| テーブル | サイズ | 戦略 |
|---------|--------|------|
| recruits_backup_20251215 | 284 MB | **除外** |
| recruits_backup_20251208 | 284 MB | **除外** |

---

## 2. MySQL (dorapita1804_db) 分析

### データベース概要
- **Instance**: db-120011
- **Database**: dorapita1804_db
- **総容量**: 約150GB（ログテーブル含む）

### テーブル分類と推奨事項

#### 2.1 超大規模ログテーブル（完全除外）

| テーブル | 行数 | サイズ | 用途 | Fixture戦略 |
|---------|------|--------|------|-------------|
| **view_logs** | 29,907,847 | **58 GB** | 閲覧ログ | **完全除外** |
| **change_logs** | 2,190,329 | **40 GB** | 変更ログ | **完全除外** |
| **effect_report_shokushu_items** | 68,284,806 | **23 GB** | レポートデータ | **完全除外** |
| **oubo_analyzes** | 10,335,112 | **5 GB** | 応募分析 | **完全除外** |
| **text_search_logs** | 7,224,515 | **1.6 GB** | 検索ログ | **完全除外** |
| **effect_reports** | 1,620,443 | **1 GB** | 効果レポート | **完全除外** |

**重要**: これらのテーブルはFixture化すると開発環境が破綻する。テストでは空テーブルまたは少数のモックデータで対応。

---

#### 2.2 中規模関連テーブル（ごく少数のサンプルのみ）

| テーブル | 行数 | サイズ | Fixture戦略 |
|---------|------|--------|-------------|
| recruit_hukuri_items | 556,816 | 142 MB | **5-10件** |
| cnt_dy_items | 2,480,624 | 98 MB | **Fixture化しない** |
| recruit_tokuyu_items | 233,283 | 59 MB | **5-10件** |
| recruit_sakuru_items | 143,823 | 53 MB | **5-10件** |
| recruit_shukintime_items | 124,671 | 40 MB | **5-10件** |
| recruit_hinmoku_items | 138,341 | 32 MB | **5-10件** |
| recruit_form_types | 75,193 | 18 MB | **5-10件** |
| recruit_keijou_items | 59,465 | 15 MB | **5-10件** |
| recruit_entry_menkyo_items | 173,301 | 13 MB | **5-10件** |
| recruit_shokushu_items | 40,200 | 9 MB | **5-10件** |
| recruit_employment_statuses | 38,462 | 7 MB | **5-10件** |
| recruit_areas | 33,565 | 10 MB | **5-10件** |

**方針**: これらは主テーブル（recruits, selections）に紐づく多対多の関連データ。主テーブルのFixtureに対応する最小限のデータのみ抽出。

---

#### 2.3 主要トランザクションデータ（サンプル抽出 + マスキング）

| テーブル | 行数 | サイズ | 主なPII | Fixture戦略 |
|---------|------|--------|---------|-------------|
| **recruits** | 21,583 | 590 MB | 企業名、連絡先 | **20-30件抽出** + マスキング |
| **selections** | 58,764 | 46 MB | **氏名、電話、メール、住所、生年月日** | **30-50件抽出** + 全PII マスキング |
| **recruit_entries** | 75,020 | 17 MB | **氏名、電話、メール、住所、生年月日** | **30-50件抽出** + 全PII マスキング |
| **companies** | 2,082 | 3.5 MB | 企業名、担当者、**パスワード** | **10-15件抽出** + 担当者情報・パスワード削除 |
| **send_mails** | 255,130 | 398 MB | メールアドレス | **Fixture化しない** |

##### 特記事項: companies.pw フィールド

**⚠️ 重大な発見**: `companies` テーブルに `pw` カラム（パスワード）が平文またはハッシュで保存されている可能性がある。

**対策**:
1. Fixture生成時に `pw` カラムを **完全に削除** または固定値（例: `bcrypt('password')`）に置換
2. または、`pw` カラムをFixtureから除外し、テストコード内で動的生成

---

#### 2.4 マスターデータ（全件またはサンプル）

| テーブル | 行数 | サイズ | Fixture戦略 |
|---------|------|--------|-------------|
| zip_codes | 120,360 | 11 MB | **全件コピー** または郵便番号API利用 |
| information | 少数 | - | **全件コピー**（お知らせマスタ） |
| login_logs | 14,808 | 3.5 MB | **Fixture化しない** |
| assets | 46,282 | 5.5 MB | **主要画像のみ10-20件** |
| contacts | 2,618 | 4.5 MB | **5-10件** |

---

## 3. Fixture作成の優先順位

### Phase 1: マスターデータ（最優先）

| 優先度 | テーブル | DB | 行数目安 | 理由 |
|-------|---------|-----|---------|------|
| ★★★ | prefectures | PG | 47（全件） | 都道府県マスタ、外部キー参照多数 |
| ★★★ | areas | PG | 9（全件） | エリアマスタ |
| ★★★ | kodawari_items | PG | 163（全件） | こだわり条件 |
| ★★ | information | MySQL | 全件 | お知らせマスタ |

### Phase 2: 企業・求人データ（コア機能）

| 優先度 | テーブル | DB | 行数目安 | マスキング要否 |
|-------|---------|-----|---------|---------------|
| ★★★ | companies | PG | 10-15件 | ✅ 担当者情報 |
| ★★★ | companies | MySQL | 10-15件 | ✅ 担当者・パスワード |
| ★★★ | recruits | PG | 20-30件 | ✅ 連絡先 |
| ★★★ | recruits | MySQL | 20-30件 | ✅ 連絡先 |

### Phase 3: 応募データ（テスト重要度高）

| 優先度 | テーブル | DB | 行数目安 | マスキング要否 |
|-------|---------|-----|---------|---------------|
| ★★★ | entries | PG | 30-50件 | ✅ **全PII** |
| ★★★ | selections | MySQL | 30-50件 | ✅ **全PII** |
| ★★ | recruit_entries | MySQL | 30-50件 | ✅ **全PII** |

### Phase 4: 関連データ（最小限）

| 優先度 | テーブル | DB | 行数目安 |
|-------|---------|-----|---------|
| ★ | recruit_*_items（各種） | MySQL | 各5-10件 |
| ★ | assets | MySQL | 10-20件 |

---

## 4. PIIマスキング実装方針

### 4.1 マスキング対象カラム一覧

| テーブル | カラム | タイプ | マスキング方法 |
|---------|--------|--------|----------------|
| entries, selections, recruit_entries | name | 氏名 | テストユーザー{連番} |
| 同上 | kana | フリガナ | テストユーザー{連番} |
| 同上 | tel | 電話番号 | 080-0000-{連番4桁} |
| 同上 | mail | メール | test{連番}@example.com |
| 同上 | zip_code | 郵便番号 | 100-0001 |
| 同上 | address, address1 | 住所 | 東京都千代田区1-1-1 |
| 同上 | address2, detail_address | 住所詳細 | テストビル{連番}号室 |
| 同上 | birthday | 生年月日 | 1990-01-01 |
| companies | pw | パスワード | **削除またはbcrypt('password')** |
| companies | tanto, tanto_kana | 担当者 | 担当 太郎 |
| companies | email | メール | tanto{連番}@example.com |
| companies | telnum | 電話番号 | 03-0000-{連番4桁} |

### 4.2 マスキングスクリプト実装例（擬似コード）

```sql
-- entries テーブルのマスキング
UPDATE entries
SET
    name = 'テストユーザー' || LPAD(ROW_NUMBER() OVER (ORDER BY id), 3, '0'),
    kana = 'テストユーザー' || LPAD(ROW_NUMBER() OVER (ORDER BY id), 3, '0'),
    tel = '080-0000-' || LPAD((ROW_NUMBER() OVER (ORDER BY id))::text, 4, '0'),
    mail = 'test' || LPAD(ROW_NUMBER() OVER (ORDER BY id), 3, '0') || '@example.com',
    zip = '100-0001',
    address1 = '東京都千代田区1-1-1',
    address2 = 'テストビル' || ROW_NUMBER() OVER (ORDER BY id) || '号室',
    birthday = '1990-01-01';
```

---

## 5. Fixture生成ワークフロー

### Step 1: STG DBからサンプルデータ抽出

```bash
# PostgreSQL
PGPASSWORD=<※.secret参照> pg_dump \
  -h 127.0.0.1 -p 35432 \
  -U dorauser2022 -d dorapita \
  -t prefectures -t areas -t kodawari_items \
  --data-only --column-inserts \
  > /tmp/pgsql_master.sql

# MySQL
mysqldump \
  -h 127.0.0.1 -P 33306 \
  -u root -p<※.secret参照> \
  dorapita1804_db \
  information \
  --no-create-info --skip-extended-insert \
  > /tmp/mysql_master.sql
```

### Step 2: トランザクションデータの抽出（LIMIT付き）

```sql
-- PostgreSQL: recruits 30件抽出
COPY (
  SELECT * FROM recruits
  WHERE closed = false
  ORDER BY created DESC
  LIMIT 30
) TO '/tmp/recruits_sample.csv' CSV HEADER;
```

### Step 3: Claude CodeによるPIIマスキング

```bash
# CSVを読み込み、PIIをマスキングしてCakePHP Fixture形式に変換
# → Claude Codeがこの処理を実装
```

### Step 4: CakePHP Fixtureファイル生成

```php
<?php
// dorapita.com/tests/Fixture/RecruitsFixture.php
declare(strict_types=1);

namespace App\Test\Fixture;

use Cake\TestSuite\Fixture\TestFixture;

class RecruitsFixture extends TestFixture
{
    public function init(): void
    {
        $this->records = [
            [
                'id' => 1,
                'title' => 'ドライバー募集',
                'c_nm' => 'テスト運送株式会社',
                'tel_num' => '0120-000-001',
                // ... マスキング済みデータ
            ],
            // ...
        ];
        parent::init();
    }
}
```

---

## 6. 推奨Fixtureサイズと開発環境への影響

### Fixture総データ量見積もり

| カテゴリ | テーブル数 | 総行数 | 推定サイズ | 初期化時間 |
|---------|-----------|-------|-----------|-----------|
| マスターデータ | 6 | 300 | 100 KB | < 1秒 |
| 企業・求人 | 4 | 80 | 500 KB | 1-2秒 |
| 応募データ | 3 | 120 | 300 KB | 1-2秒 |
| 関連データ | 10 | 100 | 200 KB | 1-2秒 |
| **合計** | **23** | **600** | **~1 MB** | **< 10秒** |

**結論**: 適切にサンプリングすれば、Fixtureデータは1MB以下に収まり、PHPUnit実行時の初期化時間は10秒未満。

---

## 7. 除外すべきデータ（重要）

### 完全除外リスト

| テーブル | 理由 |
|---------|------|
| view_logs（PG/MySQL） | 巨大（58GB）、テストに不要 |
| view_logs_copy_* | バックアップ、不要 |
| change_logs | 巨大（40GB）、テストに不要 |
| effect_report_shokushu_items | 巨大（23GB）、レポート機能専用 |
| oubo_analyzes | 巨大（5GB）、分析専用 |
| send_mails | 巨大（776MB + 398MB）、テストで動的生成 |
| text_search_logs | ログデータ、テストに不要 |
| entry_logs | ログデータ、テストに不要 |
| use_numbers | セッション管理、テストで動的生成 |
| cnt_dy_items | 集計データ、テストに不要 |

---

## 8. セキュリティ考慮事項

### 8.1 Fixture配布時の注意

1. **Fixtureファイルは外部公開しない**
   - GitHubにpushする場合、privateリポジトリ必須
   - または、Fixtureを `.gitignore` に追加し、チーム内で別途共有

2. **PII完全マスキングの確認**
   - 本番データが残っていないか、目視確認必須
   - 自動化ツールでの検証（例: `grep -r "090-" tests/Fixture/`）

3. **パスワードの取り扱い**
   - `companies.pw` は絶対にFixtureに含めない
   - テストコード内で `$user->setPassword('test_password')` を使用

### 8.2 STG DB接続の制限

- STG DBへの接続は最小限に（Fixture生成時のみ）
- 接続ログを記録し、不正アクセスを監視
- `.secret` ファイルの厳格な管理（git管理外、権限600）

---

## 9. 次のアクション

### Immediate (今すぐ)
- [ ] マスターデータFixture生成スクリプト作成
- [ ] PIIマスキングロジックの実装
- [ ] Fixtureファイル自動生成ツールの開発

### Short-term (1週間以内)
- [ ] Phase 1-2のFixture生成（マスター + 企業・求人）
- [ ] PHPUnitでのFixture投入テスト
- [ ] Fixture配布方法の確定

### Mid-term (1ヶ月以内)
- [ ] Phase 3-4のFixture生成（応募 + 関連）
- [ ] CI/CDへのFixture統合
- [ ] Fixtureメンテナンス方針の策定

---

## 付録A: テーブル詳細情報

### PostgreSQL (dorapita)

```
全24テーブル
├── マスター: 6テーブル（300行）
├── トランザクション: 3テーブル（69K行）
├── ログ: 5テーブル（327K行） ← 除外
└── バックアップ: 2テーブル ← 除外
```

### MySQL (dorapita1804_db)

```
全159テーブル
├── マスター: 少数
├── トランザクション: 約10テーブル（150K行）
├── 関連: 約50テーブル（多対多）
└── ログ: 約6テーブル（100M行、150GB） ← 完全除外
```

---

## 付録B: 参照資料

- [CakePHP 4.x Fixtures](https://book.cakephp.org/4/ja/development/testing.html#fixture-factories)
- [PostgreSQL pg_dump](https://www.postgresql.org/docs/current/app-pgdump.html)
- [MySQL mysqldump](https://dev.mysql.com/doc/refman/8.0/en/mysqldump.html)
- [GDPR Personal Data](https://gdpr.eu/eu-gdpr-personal-data/)

---

**Document version**: 1.0
**Last updated**: 2025-12-26
**Author**: Claude Code
