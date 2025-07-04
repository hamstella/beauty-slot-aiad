# API設計 未解決事項 Q&A

予約管理APIの実装を進めるにあたり、決定が必要な項目をQ&A形式で整理しました。

## 1. 営業時間・定休日設定

### Q1: 営業時間の管理方法はどうしますか？
**A:** デフォルト営業時間をシステム設定として保持し、スタッフ個別のシフト登録で上書き可能とする。定休日は営業カレンダーテーブルで管理する。デフォルト営業時間はDBでもつ。APIから変更可能とする。

### Q2: 祝日対応は必要ですか？
**A:** 現段階では不要。手動で定休日設定可能な仕組みで対応する。将来的に祝日API連携を検討。

### Q3: 特別営業日（年末年始等）の設定は必要ですか？
**A:** 現段階では不要。営業カレンダーでの手動設定で対応する。

## 2. 予約フロー詳細仕様

### Q4: 予約可能期間の設定はどうしますか？
- 最短：当日予約可能 / 翌日から
- 最長：30日後 / 60日後 / 90日後

**A:** 最短：翌日から（当日予約は混乱回避のため不可）、最長：90日後まで。最長に関してはシステム設定で変更可能とする。

### Q5: 仮予約の保持時間は何分が適切ですか？
- 15分 / 30分 / 60分

**A:** 仮予約機能は不要です

### Q6: キャンセル・変更の制限時間はどうしますか？
- キャンセル：1時間前 / 2時間前 / 24時間前
- 変更：30分前 / 1時間前 / 2時間前

**A:** キャンセル：前日まで。

### Q7: 同一顧客の連続予約を制限しますか？
- 制限なし / 1日1回まで / 1週間に3回まで

**A:** 制限なし。

## 3. 空き時間検索アルゴリズム

### Q8: 予約間隔は何分に設定しますか？
- 0分（連続予約可能）/ 15分 / 30分

**A:** 15分。施術後の清掃や準備時間を確保し、スタッフの休憩時間も考慮。

### Q9: 複数スタッフが対応可能な場合の表示順序は？
- 空き時間が多い順 / ランダム / 人気順

**A:** 空き時間が多い順。スケジュールの効率化と予約の取りやすさを重視。

### Q10: 人気スタッフの優先表示は必要ですか？
**A:** 現段階では不要。公平性を保ち、全スタッフに予約機会を提供する。将来的に指名料システムと合わせて検討。

## 4. 通知システム仕様

### Q11: リマインダー通知のタイミングはいつですか？
- 24時間前のみ / 24時間前と2時間前 / 24時間前、2時間前、30分前

**A:** 24時間前と2時間前の2回。過度な通知を避けつつ、忘れ防止効果を最大化。

### Q12: 通知チャネルは何を優先しますか？
- メールのみ / メール+SMS / メール+SMS+プッシュ通知

**A:** メールのみ。実装コストを抑え、確実な配信を重視。将来的にSMS追加を検討。

### Q13: 通知失敗時のリトライ回数と間隔は？
- 3回（5分間隔）/ 5回（10分間隔）/ 3回（15分間隔）

**A:** 3回（15分間隔）。サーバー負荷を抑えつつ、一時的な障害に対応。

## 5. 認証・認可方式

### Q14: 顧客認証は電話番号ベースで良いですか？
- 電話番号+SMS認証 / メールアドレス+パスワード / 両方対応

**A:** メールアドレス+パスワード。実装が簡単で、SMS認証のコストを避けられる。

### Q15: SMS認証の実装コストは許容範囲ですか？
**A:** 現段階では許容範囲外。MVPではメール認証に絞り、将来的に検討。

### Q16: 管理者の役割分担は必要ですか？
- admin（全権限）のみ / admin + staff（制限あり）/ 細かい権限設定

**A:** admin + staff（制限あり）。adminは全操作可能、staffは自分の予約管理のみ。

## 6. エラーハンドリング

### Q17: エラーメッセージの日本語化は必要ですか？
**A:** 必要。ユーザー向けメッセージは日本語、開発者向けログは英語とする。

### Q18: より詳細なエラーコードが必要ですか？
現在の案：VALIDATION_ERROR, NOT_FOUND, CONFLICT, BUSINESS_RULE_ERROR, INTERNAL_ERROR

**A:** 現在の案で十分。シンプルな構成を維持し、詳細はメッセージで補完。

## 7. パフォーマンス・運用

### Q19: 想定同時接続数はどの程度ですか？
- 10人以下 / 50人以下 / 100人以上

**A:** 50人以下。中規模美容室を想定し、余裕を持った設計とする。

### Q20: データ量増加時のパフォーマンス対策は必要ですか？
- 特に不要 / キャッシュ導入 / パーティション分割

**A:** 特に不要

## 8. その他

### Q21: 本番環境での監視項目は何が必要ですか？
- 基本的なヘルスチェックのみ / レスポンス時間監視 / 詳細なメトリクス監視

**A:** レスポンス時間監視。API応答時間、データベース接続状況、エラー率を監視対象とする。 

---

## 回答記入後の次ステップ

1. 上記の質問にお答えいただく
2. 回答内容をもとにAPI仕様を詳細化
3. バックエンドハンドラーの実装を開始

回答をいただければ、それに基づいて具体的な実装を進めます。