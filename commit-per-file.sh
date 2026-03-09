#!/bin/bash
# Adiciona e faz commit de cada arquivo alterado separadamente (conventional commits).
# O tipo é inferido automaticamente pelo caminho do arquivo.
# Uso: ./commit-per-file.sh

set -e

# Infere o tipo do conventional commit pelo path do arquivo
infer_type() {
  local file="$1"
  case "$file" in
    README*|*.md|doc/*|docs/*) echo "docs" ;;
    *_test.go|*_test.ts|test/*|tests/*) echo "test" ;;
    go.mod|go.sum|.gitignore|.env*|Makefile|Dockerfile*|*.yml|*.yaml|*.sh) echo "chore" ;;
    internal/*|cmd/*|pkg/*|*.go) echo "feat" ;;
    *) echo "chore" ;;
  esac
}

# Lista arquivos modificados e não rastreados ($NF pega o path mesmo em renames)
# Usa while read em vez de mapfile para compatibilidade com Bash 3 (macOS)
FILES=()
while IFS= read -r line; do FILES+=("$line"); done < <(git status --short | awk '{ if (NF >= 2) print $NF }')

if [ ${#FILES[@]} -eq 0 ]; then
  echo "Nenhum arquivo para commitar."
  exit 0
fi

echo "Encontrados ${#FILES[@]} arquivo(s). Tipo inferido por arquivo."
echo ""

for file in "${FILES[@]}"; do
  # Pula se for diretório (git status pode listar dir com /)
  [ -d "$file" ] && continue

  TYPE=$(infer_type "$file")
  status=$(git status --short "$file" | awk '{print $1}')
  if [[ "$status" == "??" ]] || [[ "$status" == "A" ]]; then
    desc="add $file"
  else
    desc="update $file"
  fi

  echo ">>> git add $file"
  git add "$file"
  echo ">>> git commit -m \"$TYPE: $desc\""
  git commit -m "$TYPE: $desc"
  echo ""
done

echo "Pronto. ${#FILES[@]} commit(s) criado(s)."
