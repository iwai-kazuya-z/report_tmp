# テーブル使用状況分析レポート

**生成日時**: 2025-12-26 21:11:05
**分析対象**: dorapita_code リポジトリ

---

## 分析方法

1. **CakePHP Tableクラス**: 各アプリの src/Model/Table/ 配下のクラスを検出
2. **最終更新日時**: PostgreSQL/MySQLのテーブル更新日時を取得
3. **コード参照頻度**: ソースコード内でのテーブル名出現回数をカウント

---

## 1. CakePHP Table クラス一覧

### dorapita.com (PostgreSQL)

```
_SystemMails
ApplicationHistories
ApplicationQueue
Applications
Areas
Assets
ChargeTypes
Chats
Clients
Companies
CompaniesMs
DeliveryAreas
DesiredWorkShifts
DismissedPopups
EducationLevels
EmploymentDesires
EmploymentStatuses
EmploymentTypes
Entries
EntryLogs
ExperienceTypes
FormTypes
Genders
GraduationCategories
HinmokuItems
Histories
HukuriItems
Hunts
Industries
Informations
JobTypes
KeijouItems
KodawariItems
KyujitsuItems
Locations
MenkyoItems
Messages
MsRecruits
Occupations
OuboTypes
Prefectures
RecruitEmploymentStatuses
RecruitFormTypes
RecruitHinmokuItems
RecruitHukuriItems
RecruitKeijouItems
RecruitKyujitsuItems
Recruits
RecruitSakuruItems
RecruitShokushuItems
RecruitShukintimeItems
RecruitsMs
RecruitTokuyuItems
Redirects
SakuruItems
SalaryTypes
Scouts
SelectionMenkyoItems
Selections
SendMails
SendMailStates
ShokushuItems
ShukintimeItems
SmsSendMessages
SystemMails
TelLogs
TelLogsMs
TokuyuItems
TransportExperiences
TransportItems
TwilioNumbers
UseNumbers
UserBlockCompanies
UserDeliveryAreas
UserEducations
UserExperiences
UserExperiencesDeliveryAreas
UserExperiencesOccupations
UserExperiencesTransports
UserExperiencesVehicleShapes
UserFavorites
UserJobTypesExperiences
UserOccupationsExperiences
UserOccupationsExperiencesTransportItems
UserOtherExperiences
UserProfileMenkyoItems
UserProfiles
UserProviders
UserResumes
Users
UserTokens
UserWorkShifts
VehicleShapes
ViewLogs
VwPrefRanks
WordLogs
YhRecruits
ZipCodes
```

**Total**:       98 tables

### cadm.dorapita.com (MySQL)

```
Ages
ApplicationHistories
Applications
ChargeTypes
Chats
CheckStates
Clients
ClientsHideModal
ClientTokens
Companies
CompanyInformations
CompanyStates
CompanyTokens
DeliveryAreas
DesiredWorkShifts
EducationLevels
EffectReports
EffectReportShokushuItems
EmploymentDesires
EmploymentStatuses
EmploymentTypes
ExperienceTypes
Genders
GraduationCategories
Holidays
Hunts
Industries
Informations
InformationTypes
JobTypes
KakinTypes
LineMessages
Locations
MenkyoItems
Messages
Occupations
OuboTypes
Plans
Prefectures
RecruitAreas
RecruitEmploymentStatuses
RecruitEntries
RecruitEntryMenkyoItems
RecruitFormTypes
Recruits
RecruitShokushuItems
RecruitStates
SakuruItems
SalaryTypes
Salesclerks
Scouts
SelectionMenkyoItems
Selections
SelectionStatuses
SendMails
SendMailStates
ShokushuItems
Shops
ShopUsers
ShopUserTokens
SystemMails
TelLogs
TokuyuItems
TransportExperiences
TransportItems
TwoFactors
UserBlockCompanies
UserDeliveryAreas
UserEducations
UserExperiences
UserExperiencesDeliveryAreas
UserExperiencesOccupations
UserExperiencesTransports
UserExperiencesVehicleShapes
UserOccupationsExperiences
UserOccupationsExperiencesTransportItems
UserOtherExperiences
UserProfileMenkyoItems
UserProfiles
UserResumes
Users
UserVehicleExperienceDetails
UserWorkShifts
VehicleShapes
ViewLogs
VwLocations
ZipCodes
```

**Total**:       87 tables

### kanri.dorapita.com (MySQL)

```
Accounts
Ages
ApplicationHistories
Applications
Areas
Assets
ChangeLogs
ChargeTypes
CheckStates
Companies
CompanyStates
Counts
EffectReports
EffectReportShokushuItems
EmploymentStatuses
FormTypes
Genders
HinmokuItems
HolidayTypes
HukuriItems
Informations
KakinTypes
KeijouItems
Locations
MenkyoItems
OuboTypes
PickupTypes
Plans
Prefectures
RecruitAreas
RecruitAssets
RecruitDrafts
RecruitEditRequests
RecruitEmploymentStatuses
RecruitEntries
RecruitFormTypes
RecruitHinmokuItems
RecruitHukuriItems
RecruitKeijouItems
Recruits
RecruitSakuruItems
RecruitShokushuItems
RecruitShukintimeItems
RecruitStates
RecruitStatus
RecruitTokuyuItems
RecruitTypes
SakuruItems
SalaryTypes
Salesclerks
Scouts
Selections
SelectionStatuses
SendMails
SendMailStates
ShokushuItems
Shops
ShopTokens
ShopUsers
ShopUserTokens
ShukintimeItems
SystemMails
TokuyuItems
TwoFactors
ViewLogs
VwCompanyLastDays
VwLocations
```

**Total**:       67 tables

### dora-pt.jp (MySQL)

```
Ages
ApplicationQueue
Clients
Companies
CompanyInformations
EffectReports
EmploymentStatuses
FormTypes
Genders
HinmokuItems
HolidayTypes
HukuriItems
InformationTypes
IntroductionTypes
KeijouItems
KyujitsuItems
Locations
MenkyoItems
OuboTypes
PartnerClients
Plans
Prefectures
RecopAuths
RecruitEmploymentStatuses
RecruitEntries
RecruitEntryMenkyoItems
RecruitFormTypes
RecruitHinmokuItems
RecruitHukuriItems
RecruitKeijouItems
RecruitKyujitsuItems
Recruits
RecruitSakuruItems
RecruitShokushuItems
RecruitShukintimeItems
RecruitTokuyuItems
SakuruItems
SalaryTypes
Scouts
SelectionMenkyoItems
Selections
SelectionStatuses
ShokushuItems
ShopUsers
ShukintimeItems
TelLogs
TokuyuItems
ViewLogs
ZipCodes
```

**Total**:       49 tables

---

## 2. PostgreSQL テーブル最終更新日時

| テーブル名 | 最終更新日時 | 行数 | 判定 |
|-----------|-------------|------|------|
| areas | NULL | 9 | 判定不可 |
| companies | 2025-12-26 05:03:11.092803 | 2425 | ✅ 使用中 |
| entries | 2025-12-17 14:00:03.330559 | 47842 | ✅ 使用中 |
| entry_logs | 2025-12-17 13:59:37.453689 | 19391 | ✅ 使用中 |
| kodawari_items | NULL | 163 | 判定不可 |
| pg_stat_statements | N/A | 3465 | 判定不可 |
| phinxlog | N/A | 6 | 判定不可 |
| prefectures | 2022-06-28 07:16:08.39367 | 47 | ⚠️ 長期未更新 |
| recruits | 2025-12-26 05:03:03.39389 | 19487 | ✅ 使用中 |
| recruits_backup_20251208 | 2025-12-04 13:02:40.928881 | 36044 | ✅ 使用中 |
| recruits_backup_20251215 | 2025-12-04 13:02:40.928881 | 36044 | ✅ 使用中 |
| redirects | 2023-05-25 13:46:13.162156+09 | 1 | ⚠️ 長期未更新 |
| send_mails | 2025-12-17 14:02:03.470011 | 270648 | ✅ 使用中 |
| send_mail_states | 2022-06-23 15:26:07.482087 | 4 | ⚠️ 長期未更新 |
| system_mails | 2025-10-01 15:04:07.48379 | 10 | ✅ 使用中 |
| twilio_numbers | 2025-10-15 14:25:29.108484 | 4 | ✅ 使用中 |
| use_numbers | 2025-10-15 14:25:29.113843 | 36907 | ✅ 使用中 |
| view_logs | 2025-12-24 17:44:02.21549 | 599 | ✅ 使用中 |
| view_logs_copy_202510061628 | 2025-10-06 12:09:05.305205 | 10725476 | ✅ 使用中 |
| view_logs_copy_20251217 | 2025-10-06 12:09:05.305205 | 10725476 | ✅ 使用中 |
| vw_pref_ranks | N/A | 10 | 判定不可 |
|  | N/A |  | 判定不可 |

---

## 3. MySQL テーブル最終更新日時

| テーブル名 | 最終更新日時 | 行数 | 判定 |
|-----------|-------------|------|------|
| effect_reports | 2025-12-26 21:10:04 | 1620443 | ✅ 使用中 |
| effect_report_shokushu_items | 2025-12-26 21:10:04 | 68284852 | ✅ 使用中 |
| recruits | 2025-12-26 06:05:08 | 21583 | ✅ 使用中 |
| change_logs | 2025-12-26 06:05:08 | 2190329 | ✅ 使用中 |
| companies | 2025-12-26 01:03:04 | 2082 | ✅ 使用中 |
| favor_items | 2025-12-26 01:02:04 | 10 | ✅ 使用中 |
| application_histories | 2025-12-25 18:44:14 | 2120 | ✅ 使用中 |
| applications | 2025-12-25 18:44:11 | 908 | ✅ 使用中 |
| chats | 2025-12-25 18:43:50 | 461 | ✅ 使用中 |
| view_logs | 2025-12-24 17:44:02 | 29907847 | ✅ 使用中 |
| users | 2025-12-24 10:52:17 | 174 | ✅ 使用中 |
| send_mails | 2025-12-24 04:04:01 | 255130 | ✅ 使用中 |
| accounts | 2025-12-23 18:56:22 | 75 | ✅ 使用中 |
| login_logs | 2025-12-23 18:53:19 | 14808 | ✅ 使用中 |
| user_resumes | 2025-12-23 12:51:30 | 76 | ✅ 使用中 |
| user_profiles | 2025-12-23 12:51:30 | 130 | ✅ 使用中 |
| user_profile_menkyo_items | 2025-12-23 12:49:55 | 134 | ✅ 使用中 |
| user_other_experiences | 2025-12-23 12:49:55 | 30 | ✅ 使用中 |
| user_educations | 2025-12-23 12:49:37 | 41 | ✅ 使用中 |
| user_experiences | 2025-12-23 12:49:33 | 46 | ✅ 使用中 |
| two_factors | 2025-12-23 09:54:02 | 2732 | ✅ 使用中 |
| selections | 2025-12-17 14:05:02 | 58764 | ✅ 使用中 |
| line_messages | 2025-12-17 14:05:02 | 535 | ✅ 使用中 |
| clients | 2025-12-17 14:02:27 | 971 | ✅ 使用中 |
| user_tokens | 2025-12-17 14:00:18 | 306 | ✅ 使用中 |
| company_informations | 2025-12-17 14:00:03 | 63356 | ✅ 使用中 |
| recruit_entries | 2025-12-17 14:00:03 | 75020 | ✅ 使用中 |
| recruit_entry_menkyo_items | 2025-12-17 13:54:01 | 173301 | ✅ 使用中 |
| user_occupations_experiences | 2025-12-17 13:51:19 | 16 | ✅ 使用中 |
| user_work_shifts | 2025-12-17 13:51:19 | 19 | ✅ 使用中 |
| user_occupations_experiences_transport_items | 2025-12-17 13:51:19 | 53 | ✅ 使用中 |
| messages | 2025-12-17 12:10:24 | 549 | ✅ 使用中 |
| user_experiences_transports | 2025-12-17 12:06:34 | 38 | ✅ 使用中 |
| user_experiences_occupations | 2025-12-17 12:06:34 | 39 | ✅ 使用中 |
| user_experiences_delivery_areas | 2025-12-17 12:06:34 | 49 | ✅ 使用中 |
| user_experiences_vehicle_shapes | 2025-12-17 12:06:34 | 26 | ✅ 使用中 |
| selection_menkyo_items | 2025-12-17 09:55:02 | 9699 | ✅ 使用中 |
| recruit_keijou_items | 2025-12-17 09:47:53 | 59465 | ✅ 使用中 |
| recruit_menkyo_items | 2025-12-17 09:47:53 | 563 | ✅ 使用中 |
| recruit_hukuri_items | 2025-12-17 09:47:53 | 556816 | ✅ 使用中 |
| recruit_hinmoku_items | 2025-12-17 09:47:53 | 138341 | ✅ 使用中 |
| recruit_form_types | 2025-12-17 09:47:53 | 75193 | ✅ 使用中 |
| recruit_tokuyu_items | 2025-12-17 09:47:53 | 233283 | ✅ 使用中 |
| recruit_employment_statuses | 2025-12-17 09:47:53 | 38462 | ✅ 使用中 |
| recruit_shukintime_items | 2025-12-17 09:47:53 | 124671 | ✅ 使用中 |
| recruit_areas | 2025-12-17 09:47:53 | 33565 | ✅ 使用中 |
| recruit_shokushu_items | 2025-12-17 09:47:53 | 40200 | ✅ 使用中 |
| recruit_ages | 2025-12-17 09:47:53 | 957 | ✅ 使用中 |
| recruit_sakuru_items | 2025-12-17 09:47:53 | 143823 | ✅ 使用中 |
| dismissed_popups | 2025-12-16 14:57:29 | 15 | ✅ 使用中 |
| user_delivery_areas | 2025-12-15 16:31:17 | 20 | ✅ 使用中 |
| application_queue | 2025-12-10 10:18:05 | 54 | ✅ 使用中 |
| tel_logs | 2025-12-08 19:40:21 | 32281 | ✅ 使用中 |
| scouts | 2025-12-05 19:32:14 | 90 | ✅ 使用中 |
| hunts | 2025-12-05 16:57:48 | 36 | ✅ 使用中 |
| shop_users | 2025-12-05 16:38:49 | 125 | ✅ 使用中 |
| locations | 2025-12-04 10:39:44 | 3609 | ✅ 使用中 |
| shops | 2025-11-19 14:05:30 | 30 | ✅ 使用中 |
| clients_hide_modal | 2025-11-18 12:31:09 | 4 | ✅ 使用中 |
| client_tokens | 2025-11-18 12:30:31 | 1271 | ✅ 使用中 |
| shop_user_tokens | 2025-11-18 12:13:26 | 238 | ✅ 使用中 |
| information | 2025-11-11 19:53:08 | 18 | ✅ 使用中 |
| assets | 2025-11-11 16:34:09 | 46282 | ✅ 使用中 |
| recruit_assets | 2025-11-11 16:34:09 | 41860 | ✅ 使用中 |
| salesclerks | 2025-11-11 15:55:19 | 64 | ✅ 使用中 |
| user_favorites | 2025-11-07 17:20:18 | 17 | ✅ 使用中 |
| vehicle_shapes | N/A | 26 | 判定不可 |
| text_search_logs | N/A | 7224515 | 判定不可 |
| book_allows | N/A | 391 | 判定不可 |
| employment_statuses | N/A | 7 | 判定不可 |
| promotions | N/A | 6 | 判定不可 |
| hukuri_items | N/A | 38 | 判定不可 |
| cnt_pv_items | N/A | 21893 | 判定不可 |
| transport_experiences | N/A | 6 | 判定不可 |
| recop_auth | N/A | 6 | 判定不可 |
| experience_types | N/A | 7 | 判定不可 |
| book_states | N/A | 4 | 判定不可 |
| vw_locations | N/A | NULL | 判定不可 |
| utm_logs | N/A | 21 | 判定不可 |
| banners | N/A | 72 | 判定不可 |
| employment_desires | N/A | 6 | 判定不可 |
| prefectures | N/A | 47 | 判定不可 |
| holidays | N/A | 26 | 判定不可 |
| cnt_fm_items | N/A | 3343 | 判定不可 |
| rankings | N/A | 344 | 判定不可 |
| menkyo_items | N/A | 18 | 判定不可 |
| entry_types | N/A | 2 | 判定不可 |
| book_edit_logs | N/A | 2473 | 判定不可 |
| industries | N/A | 13 | 判定不可 |
| recruit_kyujitsu_items | N/A | 0 | 判定不可 |
| selection_statuses | N/A | 5 | 判定不可 |
| account_types | N/A | 2 | 判定不可 |
| vw_company_last_days | N/A | NULL | 判定不可 |
| tokuyu_items | N/A | 33 | 判定不可 |
| tags | N/A | 6 | 判定不可 |
| banner_types | N/A | 8 | 判定不可 |
| plans | N/A | 7 | 判定不可 |
| holiday_types | N/A | 5 | 判定不可 |
| cnt_dy_items | N/A | 2480624 | 判定不可 |
| salary_values | N/A | 38 | 判定不可 |

---

## 4. コード内参照頻度

### PostgreSQL テーブル (dorapita.com)

| テーブル名 | 参照回数 |
|-----------|---------|
| areas | 47 |
| companies | 53 |
| entries | 75 |
| entry_logs | 1 |
| kodawari_items | 1 |
| pg_stat_statements | 0 |
| phinxlog | 0 |
| prefectures | 40 |
| recruits | 194 |
| recruits_backup_20251208 | 0 |
| recruits_backup_20251215 | 0 |
| redirects | 6 |
| send_mails | 3 |
| send_mail_states | 1 |
| system_mails | 2 |
| twilio_numbers | 1 |
| use_numbers | 1 |
| view_logs | 12 |
| view_logs_copy_202510061628 | 0 |
| view_logs_copy_20251217 | 0 |
| vw_pref_ranks | 1 |

### MySQL テーブル (cadm/kanri/dora-pt)

| テーブル名 | cadm | kanri | dora-pt |
|-----------|------|-------|---------|
| effect_report_shokushu_items | 3 | 3 | 0 |
| view_logs | 3 | 2 | 11 |
| oubo_analyzes | 0 | 0 | 0 |
| text_search_logs | 0 | 0 | 0 |
| cnt_dy_items | 0 | 0 | 0 |
| change_logs | 0 | 3 | 0 |
| effect_reports | 47 | 52 | 8 |
| recruit_hukuri_items | 2 | 5 | 8 |
| send_mails | 5 | 6 | 0 |
| recruit_tokuyu_items | 4 | 5 | 7 |
| recruit_entry_menkyo_items | 7 | 3 | 7 |
| recruit_sakuru_items | 4 | 5 | 7 |
| recruit_hinmoku_items | 2 | 5 | 7 |
| recruit_shukintime_items | 2 | 5 | 8 |
| zip_codes | 1 | 0 | 1 |
| recruit_form_types | 3 | 5 | 6 |
| recruit_entries | 20 | 12 | 11 |
| company_informations | 50 | 2 | 5 |
| recruit_keijou_items | 2 | 5 | 7 |
| textsearch_words | 2 | 2 | 2 |
| selections | 77 | 114 | 34 |
| assets | 2 | 11 | 5 |
| recruit_assets | 0 | 1 | 0 |
| recruit_shokushu_items | 18 | 12 | 15 |
| recruit_employment_statuses | 5 | 7 | 17 |
| recruit_areas | 3 | 5 | 2 |
| tel_logs | 3 | 2 | 3 |
| cnt_pv_items | 0 | 0 | 0 |
| recruits | 62 | 49 | 49 |
| login_logs | 0 | 0 | 0 |

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

