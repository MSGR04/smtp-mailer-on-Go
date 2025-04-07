#!/bin/bash
# Запускаем opendkim
service opendkim start
# Запускаем Postfix
service postfix start
# Держим контейнер живым
tail -f /var/log/mail.log
