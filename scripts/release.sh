#!/usr/bin/env bash
# 推送规则：获取最新版本 -> 增量 bump -> 创建根+各协议模块 tags -> 有 gh 则建 release。
# 不做打包/构建（打包由发布方另行处理）。
set -euo pipefail

DRY_RUN="${1:-}"
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_DIR"

# 先同步远端标签，保证“最新版本”取的是远端已发布的最新版
git fetch --tags origin 2>/dev/null || true

latest_tag() {
  git tag -l 'v[0-9]*' --merged HEAD | sort -V | tail -1
}

bump_patch() {
  local v="$1" maj minor patch
  IFS='.' read -r maj minor patch <<< "${v#v}"
  printf 'v%d.%d.%d' "$maj" "$minor" "$((patch + 1))"
}

run() {
  if [ "$DRY_RUN" = "--dry-run" ]; then echo "[dry-run] $*"; else "$@"; fi
}

LATEST="$(latest_tag)"
[ -n "$LATEST" ] || { echo "错误：未找到已有版本标签"; exit 1; }
NEW="$(bump_patch "$LATEST")"
echo "最新版本: $LATEST -> 新版本: $NEW"

run git tag -a "$NEW" -m "release $NEW"

# 各协议模块增量打标签（与现有 protocols/<path>/vX.Y.Z 模式一致）
while IFS= read -r mod; do
  t="protocols/$mod/$NEW"
  run git tag -a "$t" -m "release $NEW"
done < <(find protocols -name go.mod -printf '%h\n' | sed 's|^protocols/||' | sort)

if [ "$DRY_RUN" = "--dry-run" ]; then
  echo "[dry-run] 标签将推送至 origin: git push origin --tags"
  exit 0
fi

run git push origin --tags

if command -v gh >/dev/null 2>&1; then
  run gh release create "$NEW" --repo erikwang2013/industrial-protocols-go --title "v$NEW" --generate-notes
else
  echo "提示：未安装 gh CLI，跳过 GitHub Release 创建（tags 已推送）。"
fi
