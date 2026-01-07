#!/bin/bash
# テーブル使用状況分析スクリプト
# 用途: 各アプリケーションで実際に使用されているテーブルを特定

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
DORAPITA_CODE="$REPO_ROOT/dorapita_code"
OUTPUT_FILE="$SCRIPT_DIR/table-usage-report.md"

# 認証情報読み込み
source "$SCRIPT_DIR/.secret"

echo "📊 テーブル使用状況分析を開始します..."
echo ""

# 出力ファイル初期化
cat > "$OUTPUT_FILE" <<EOF
# テーブル使用状況分析レポート

**生成日時**: $(date '+%Y-%m-%d %H:%M:%S')
**分析対象**: dorapita_code リポジトリ

---

## 分析方法

1. **CakePHP Tableクラス**: 各アプリの src/Model/Table/ 配下のクラスを検出
2. **最終更新日時**: PostgreSQL/MySQLのテーブル更新日時を取得
3. **コード参照頻度**: ソースコード内でのテーブル名出現回数をカウント

---

EOF

# ========================================
# 1. CakePHP Tableクラス一覧取得
# ========================================

echo "🔍 Step 1: CakePHP Tableクラスを検索中..."

cat >> "$OUTPUT_FILE" <<EOF
## 1. CakePHP Table クラス一覧

### dorapita.com (PostgreSQL)

EOF

if [ -d "$DORAPITA_CODE/dorapita.com/src/Model/Table" ]; then
    find "$DORAPITA_CODE/dorapita.com/src/Model/Table" -name "*Table.php" | while read file; do
        basename "$file" | sed 's/Table\.php$//' >> "$OUTPUT_FILE.tmp.dorapita"
    done
    if [ -f "$OUTPUT_FILE.tmp.dorapita" ]; then
        echo "\`\`\`" >> "$OUTPUT_FILE"
        cat "$OUTPUT_FILE.tmp.dorapita" | sort >> "$OUTPUT_FILE"
        echo "\`\`\`" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
        echo "**Total**: $(wc -l < "$OUTPUT_FILE.tmp.dorapita") tables" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    fi
fi

cat >> "$OUTPUT_FILE" <<EOF
### cadm.dorapita.com (MySQL)

EOF

if [ -d "$DORAPITA_CODE/cadm.dorapita.com/src/Model/Table" ]; then
    find "$DORAPITA_CODE/cadm.dorapita.com/src/Model/Table" -name "*Table.php" | while read file; do
        basename "$file" | sed 's/Table\.php$//' >> "$OUTPUT_FILE.tmp.cadm"
    done
    if [ -f "$OUTPUT_FILE.tmp.cadm" ]; then
        echo "\`\`\`" >> "$OUTPUT_FILE"
        cat "$OUTPUT_FILE.tmp.cadm" | sort >> "$OUTPUT_FILE"
        echo "\`\`\`" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
        echo "**Total**: $(wc -l < "$OUTPUT_FILE.tmp.cadm") tables" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    fi
fi

cat >> "$OUTPUT_FILE" <<EOF
### kanri.dorapita.com (MySQL)

EOF

if [ -d "$DORAPITA_CODE/kanri.dorapita.com/src/Model/Table" ]; then
    find "$DORAPITA_CODE/kanri.dorapita.com/src/Model/Table" -name "*Table.php" | while read file; do
        basename "$file" | sed 's/Table\.php$//' >> "$OUTPUT_FILE.tmp.kanri"
    done
    if [ -f "$OUTPUT_FILE.tmp.kanri" ]; then
        echo "\`\`\`" >> "$OUTPUT_FILE"
        cat "$OUTPUT_FILE.tmp.kanri" | sort >> "$OUTPUT_FILE"
        echo "\`\`\`" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
        echo "**Total**: $(wc -l < "$OUTPUT_FILE.tmp.kanri") tables" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    fi
fi

cat >> "$OUTPUT_FILE" <<EOF
### dora-pt.jp (MySQL)

EOF

if [ -d "$DORAPITA_CODE/dora-pt.jp/src/Model/Table" ]; then
    find "$DORAPITA_CODE/dora-pt.jp/src/Model/Table" -name "*Table.php" | while read file; do
        basename "$file" | sed 's/Table\.php$//' >> "$OUTPUT_FILE.tmp.dorapt"
    done
    if [ -f "$OUTPUT_FILE.tmp.dorapt" ]; then
        echo "\`\`\`" >> "$OUTPUT_FILE"
        cat "$OUTPUT_FILE.tmp.dorapt" | sort >> "$OUTPUT_FILE"
        echo "\`\`\`" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
        echo "**Total**: $(wc -l < "$OUTPUT_FILE.tmp.dorapt") tables" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    fi
fi

echo "✅ Step 1 完了"
echo ""

# ========================================
# 2. PostgreSQL テーブル最終更新日時
# ========================================

echo "🔍 Step 2: PostgreSQL テーブルの最終更新日時を取得中..."

cat >> "$OUTPUT_FILE" <<EOF
---

## 2. PostgreSQL テーブル最終更新日時

EOF

# cloud-sql-proxy起動
pkill -f "cloud-sql-proxy.*pg-120011" 2>/dev/null || true
cloud-sql-proxy "dorapita-core-dev:asia-northeast1:pg-120011" --port=35432 --gcloud-auth &
PG_PROXY_PID=$!
sleep 5

export PATH="/opt/homebrew/opt/postgresql@16/bin:$PATH"

# modified列があるテーブルの最終更新日時を取得
PGPASSWORD="$PGSQL_PASSWORD" psql -h 127.0.0.1 -p "$PGSQL_PORT" -U "$PGSQL_USER" -d "$PGSQL_DATABASE" -t -c "
SELECT
    table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;
" > "$OUTPUT_FILE.tmp.pgtables"

cat >> "$OUTPUT_FILE" <<EOF
| テーブル名 | 最終更新日時 | 行数 | 判定 |
|-----------|-------------|------|------|
EOF

while read table_name; do
    table_name=$(echo "$table_name" | xargs) # trim whitespace

    # modified列の有無を確認
    has_modified=$(PGPASSWORD="$PGSQL_PASSWORD" psql -h 127.0.0.1 -p "$PGSQL_PORT" -U "$PGSQL_USER" -d "$PGSQL_DATABASE" -t -c "
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = '$table_name' AND column_name = 'modified';
    " | xargs)

    if [ "$has_modified" -eq 1 ]; then
        # modified列がある場合
        last_modified=$(PGPASSWORD="$PGSQL_PASSWORD" psql -h 127.0.0.1 -p "$PGSQL_PORT" -U "$PGSQL_USER" -d "$PGSQL_DATABASE" -t -c "
        SELECT COALESCE(MAX(modified)::text, 'NULL') FROM $table_name;
        " 2>/dev/null | xargs || echo "ERROR")

        row_count=$(PGPASSWORD="$PGSQL_PASSWORD" psql -h 127.0.0.1 -p "$PGSQL_PORT" -U "$PGSQL_USER" -d "$PGSQL_DATABASE" -t -c "
        SELECT COUNT(*) FROM $table_name;
        " 2>/dev/null | xargs || echo "0")
    else
        last_modified="N/A"
        row_count=$(PGPASSWORD="$PGSQL_PASSWORD" psql -h 127.0.0.1 -p "$PGSQL_PORT" -U "$PGSQL_USER" -d "$PGSQL_DATABASE" -t -c "
        SELECT COUNT(*) FROM $table_name;
        " 2>/dev/null | xargs || echo "0")
    fi

    # 判定ロジック
    judgment="不明"
    if [[ "$last_modified" == "ERROR" ]]; then
        judgment="権限エラー"
    elif [[ "$last_modified" == "NULL" || "$last_modified" == "N/A" ]]; then
        judgment="判定不可"
    elif [[ "$last_modified" > "2024-06-01" ]]; then
        judgment="✅ 使用中"
    else
        judgment="⚠️ 長期未更新"
    fi

    echo "| $table_name | $last_modified | $row_count | $judgment |" >> "$OUTPUT_FILE"
done < "$OUTPUT_FILE.tmp.pgtables"

echo "" >> "$OUTPUT_FILE"

kill $PG_PROXY_PID 2>/dev/null || true

echo "✅ Step 2 完了"
echo ""

# ========================================
# 3. MySQL テーブル最終更新日時
# ========================================

echo "🔍 Step 3: MySQL テーブルの最終更新日時を取得中..."

cat >> "$OUTPUT_FILE" <<EOF
---

## 3. MySQL テーブル最終更新日時

EOF

# cloud-sql-proxy起動
pkill -f "cloud-sql-proxy.*db-120011" 2>/dev/null || true
cloud-sql-proxy "dorapita-core-dev:asia-northeast1:db-120011" --port=33306 --gcloud-auth &
MYSQL_PROXY_PID=$!
sleep 5

mysql -h 127.0.0.1 -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" -N -e "
SELECT
    table_name,
    COALESCE(update_time, 'N/A'),
    table_rows,
    CASE
        WHEN update_time IS NULL THEN '判定不可'
        WHEN update_time > '2024-06-01' THEN '✅ 使用中'
        ELSE '⚠️ 長期未更新'
    END AS judgment
FROM information_schema.TABLES
WHERE table_schema = '$MYSQL_DATABASE'
ORDER BY
    CASE WHEN update_time IS NULL THEN 1 ELSE 0 END,
    update_time DESC
LIMIT 100;
" 2>&1 | grep -v "Using a password" | awk -F'\t' '{printf "| %s | %s | %s | %s |\n", $1, $2, $3, $4}' > "$OUTPUT_FILE.tmp.mysql"

cat >> "$OUTPUT_FILE" <<EOF
| テーブル名 | 最終更新日時 | 行数 | 判定 |
|-----------|-------------|------|------|
EOF

cat "$OUTPUT_FILE.tmp.mysql" >> "$OUTPUT_FILE"

echo "" >> "$OUTPUT_FILE"

# MySQL proxyはStep 4でも使うのでまだ終了しない

echo "✅ Step 3 完了"
echo ""

# ========================================
# 4. コード内参照頻度カウント
# ========================================

echo "🔍 Step 4: コード内でのテーブル参照頻度をカウント中..."

cat >> "$OUTPUT_FILE" <<EOF
---

## 4. コード内参照頻度

### PostgreSQL テーブル (dorapita.com)

| テーブル名 | 参照回数 |
|-----------|---------|
EOF

if [ -f "$OUTPUT_FILE.tmp.pgtables" ]; then
    while read table_name; do
        table_name=$(echo "$table_name" | xargs)
        if [ -n "$table_name" ]; then
            count=$(grep -r "$table_name" "$DORAPITA_CODE/dorapita.com/src" --include="*.php" 2>/dev/null | wc -l | xargs)
            echo "| $table_name | $count |" >> "$OUTPUT_FILE"
        fi
    done < "$OUTPUT_FILE.tmp.pgtables"
fi

echo "" >> "$OUTPUT_FILE"

cat >> "$OUTPUT_FILE" <<EOF
### MySQL テーブル (cadm/kanri/dora-pt)

| テーブル名 | cadm | kanri | dora-pt |
|-----------|------|-------|---------|
EOF

# MySQL主要テーブルのみ（上位30件）
mysql -h 127.0.0.1 -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" -N -e "
SELECT table_name
FROM information_schema.TABLES
WHERE table_schema = '$MYSQL_DATABASE'
ORDER BY table_rows DESC
LIMIT 30;
" 2>&1 | grep -v "Using a password" | while read table_name; do
    cadm_count=$(grep -r "$table_name" "$DORAPITA_CODE/cadm.dorapita.com/src" --include="*.php" 2>/dev/null | wc -l | xargs)
    kanri_count=$(grep -r "$table_name" "$DORAPITA_CODE/kanri.dorapita.com/src" --include="*.php" 2>/dev/null | wc -l | xargs)
    dorapt_count=$(grep -r "$table_name" "$DORAPITA_CODE/dora-pt.jp/src" --include="*.php" 2>/dev/null | wc -l | xargs)
    echo "| $table_name | $cadm_count | $kanri_count | $dorapt_count |" >> "$OUTPUT_FILE"
done

echo "" >> "$OUTPUT_FILE"

# MySQL proxy終了
kill $MYSQL_PROXY_PID 2>/dev/null || true

echo "✅ Step 4 完了"
echo ""

# ========================================
# 5. 統合分析・推奨事項
# ========================================

echo "📝 Step 5: 統合分析と推奨事項を生成中..."

cat >> "$OUTPUT_FILE" <<EOF
---

## 5. 統合分析と推奨事項

### Fixture化優先度

#### 🟢 優先度: 最高（必須）

- Tableクラスが存在する
- 最終更新日時が6ヶ月以内
- コード参照回数が多い（50回以上）

#### 🟡 優先度: 中（推奨）

- Tableクラスは存在しないが、コード参照あり
- または、最終更新日時が6ヶ月以内

#### 🔴 優先度: 低（除外検討）

- Tableクラスなし
- コード参照なし
- 最終更新日時が6ヶ月以上前

---

## 6. 次のアクション

1. **高優先度テーブル**: 上記🟢テーブルのFixture生成を最優先
2. **中優先度テーブル**: 必要に応じてFixture化
3. **低優先度テーブル**: Fixture化しない（ストレージ削減）

EOF

# 一時ファイル削除
rm -f "$OUTPUT_FILE.tmp."*

echo ""
echo "✅ 分析完了！"
echo ""
echo "📄 レポート: $OUTPUT_FILE"
echo ""
