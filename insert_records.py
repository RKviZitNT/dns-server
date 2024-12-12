import sqlite3

conn = sqlite3.connect('data/domains.db')
cursor = conn.cursor()

cursor.execute('''
CREATE TABLE IF NOT EXISTS domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    domain TEXT NOT NULL,
    "address" TEXT,
    redirect_to TEXT
)
''')

domains = [
    ("example.com", "93.184.216.34", None),
    ("redirect.com", None, "target.com"),
    ("google.com", "142.250.190.78", None),
    ("yahoo.com", None, "altavista.com"),
    ("test.com", "192.168.0.1", None),
]

cursor.executemany('''
INSERT INTO domains (domain, "address", redirect_to) VALUES (?, ?, ?)
''', domains)

conn.commit()

conn.close()

print("Таблица успешно заполнена данными.")