update events set language_code = null where language_code = '';
update events set chat_id = null where chat_id = 0;
update events set chat_type = null where chat_type = '';