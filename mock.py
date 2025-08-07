import psycopg2
import random
from datetime import datetime, timedelta

# Настройки подключения
DB_SETTINGS = {
    "dbname": "pepsi",
    "user": "pepsi",
    "password": "pepsi123",
    "host": "localhost",  # или "database", если внутри docker
    "port": 5432
}

YOUR_USER_ID = 1
JOB_ID = 1

# Простые примеры сообщений
CLIENT_MESSAGES = [
    "Здравствуйте!", "Когда приступите?", "Есть ли вопросы по заданию?",
    "Ожидаю выполнения.", "Спасибо за обратную связь.", "Все понятно?",
    "Напишите, когда начнете.", "Хорошего дня!"
]
CONTRACTOR_MESSAGES = [
    "Здравствуйте!", "Приступаю сейчас.", "Есть вопрос по адресу.",
    "Все понятно, начинаю.", "Когда крайний срок?", "Спасибо!", "Уточню и вернусь.",
    "Сделал, проверьте, пожалуйста."
]

def generate_message_text(is_client):
    if is_client:
        return random.choice(CLIENT_MESSAGES)
    else:
        return random.choice(CONTRACTOR_MESSAGES)

def main():
    conn = psycopg2.connect(**DB_SETTINGS)
    cur = conn.cursor()

    # Получаем всех юзеров кроме Sabalaq
    cur.execute("SELECT id FROM users WHERE id != %s", (YOUR_USER_ID,))
    contractor_ids = [row[0] for row in cur.fetchall()]

    for contractor_id in contractor_ids:
        # Создаём переписку, если её нет
        cur.execute("""
            INSERT INTO chat_conversations (job_id, client_id, contractor_id)
            VALUES (%s, %s, %s)
            ON CONFLICT (job_id, client_id, contractor_id) DO NOTHING
            RETURNING id
        """, (JOB_ID, YOUR_USER_ID, contractor_id))
        row = cur.fetchone()

        if row:
            conversation_id = row[0]
        else:
            # Если уже существует, получаем её ID
            cur.execute("""
                SELECT id FROM chat_conversations
                WHERE job_id = %s AND client_id = %s AND contractor_id = %s
            """, (JOB_ID, YOUR_USER_ID, contractor_id))
            conversation_id = cur.fetchone()[0]

        # Генерация 50–80 сообщений с чередованием отправителя
        message_count = random.randint(50, 80)
        sender_flag = True  # True = клиент, False = контрактор
        created_time = datetime.now() - timedelta(days=3)

        for _ in range(message_count):
            sender_id = YOUR_USER_ID if sender_flag else contractor_id
            message_text = generate_message_text(sender_flag)
            sender_flag = not sender_flag
            created_time += timedelta(seconds=random.randint(30, 300))

            cur.execute("""
                INSERT INTO chat_messages (conversation_id, sender_id, message_text, created_at, updated_at)
                VALUES (%s, %s, %s, %s, %s)
            """, (conversation_id, sender_id, message_text, created_time, created_time))

        print(f"Добавлено {message_count} сообщений в чат с пользователем ID {contractor_id}")
        conn.commit()

    cur.close()
    conn.close()

if __name__ == "__main__":
    main()
