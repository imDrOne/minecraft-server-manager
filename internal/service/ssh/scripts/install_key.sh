#!/bin/sh

KEY="$1"

# Если ключ не передан как аргумент, читаем из stdin
if [ -z "$KEY" ]; then
  read -r KEY
fi

mkdir -p ~/.ssh
chmod 700 ~/.ssh

# Добавление ключа, если его ещё нет
echo "$KEY" >> ~/.ssh/authorized_keys

# Удаление дубликатов
sort -u ~/.ssh/authorized_keys -o ~/.ssh/authorized_keys

# Установка прав
chmod 600 ~/.ssh/authorized_keys
