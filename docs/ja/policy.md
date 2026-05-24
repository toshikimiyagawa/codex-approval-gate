# Minimal Policy Precheck

policy precheck は、local LLM が decision を出す前に、明らかな command risk を要約する。

この結果だけで request を自動 allow / deny しない。結果は judge prompt の context として追加する。provider failure や曖昧な結果は引き続き `ask` を返す。

## Verdict

- `safe`: command が狭い read-only pattern に一致する。
- `risky`: command が破壊的または権限変更系の pattern に一致する。
- `unknown`: command が既知 pattern に一致しない。

## 初期 Pattern

read-only の例:

- `git status`
- `git diff`
- `ls`
- `pwd`
- `cat`

risky の例:

- `rm`
- `sudo`
- `chmod`
- `chown`
- `>` を使う shell redirection

policy は意図的に conservative にする。unknown command は provider に委ね、provider failure 時は従来通り `ask` に fallback する。
