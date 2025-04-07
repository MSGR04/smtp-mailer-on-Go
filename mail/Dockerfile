FROM ubuntu:22.04

# Обновляем систему и устанавливаем Postfix и opendkim
RUN apt-get update && apt-get install -y postfix opendkim opendkim-tools && rm -rf /var/lib/apt/lists/*

# Копируем конфигурационные файлы в контейнер
COPY postfix/main.cf /etc/postfix/main.cf
COPY opendkim/opendkim.conf /etc/opendkim.conf
COPY opendkim/KeyTable /etc/opendkim/KeyTable
COPY opendkim/SigningTable /etc/opendkim/SigningTable
COPY opendkim/keys /etc/opendkim/keys

# Настраиваем opendkim для работы через inet:8891@localhost (или другой, согласно вашему конфигурационному файлу)
# Убедитесь, что в /etc/opendkim.conf есть строка: Socket inet:8891@localhost

# Открываем порты для SMTP (25) и submission (587)
EXPOSE 25 587

# Копируем скрипт запуска (опционально можно использовать CMD для запуска Postfix и opendkim)
COPY start.sh /start.sh
RUN chmod +x /start.sh

CMD ["/start.sh"]
