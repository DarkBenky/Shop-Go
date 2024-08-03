import sqlite3
from datetime import datetime, timedelta
from PIL import Image
import random
import io
import base64

# Connect to the SQLite database (or create it if it doesn't exist)
conn = sqlite3.connect('users.db')
cursor = conn.cursor()

# Create tables
cursor.execute('''
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        mail TEXT NOT NULL,
        address TEXT
    );
''')

cursor.execute('''
    CREATE TABLE IF NOT EXISTS items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price INTEGER NOT NULL,
        description TEXT NOT NULL,
        category TEXT,
        images BLOB,
        store_id INTEGER,
        FOREIGN KEY (store_id) REFERENCES stores(id)
    );
''')

cursor.execute('''
    CREATE TABLE IF NOT EXISTS stores (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        address TEXT NOT NULL
    );
''')

cursor.execute('''
    CREATE TABLE IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        store_id INTEGER,
        item_id INTEGER NOT NULL,
        quantity INTEGER NOT NULL,
        time_of_purchase DATE NOT NULL,
        estimated_delivery DATE,
        status TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (store_id) REFERENCES stores(id),
        FOREIGN KEY (item_id) REFERENCES items(id)
    );
''')

cursor.execute('''
    CREATE TABLE IF NOT EXISTS clicks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        item_id INTEGER NOT NULL,
        time_of_click DATE NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (item_id) REFERENCES items(id)
    );
''')

# Insert example data into tables
# Users
cursor.execute('''
    INSERT INTO users (username, password, mail, address)
    VALUES ('john', 'password123', 'johndoe@example.com', '123 Main St')
''')

# Stores
cursor.execute('''
    INSERT INTO stores (name, address)
    VALUES ('Store1', '456 Market St')
''')

def generate_random_image(size=(100, 100)):
    # Create a new random image
    image = Image.new('RGB', size, (
        random.randint(0, 255),
        random.randint(0, 255),
        random.randint(0, 255)
    ))
    
    # Save the image to a BytesIO object
    byte_arr = io.BytesIO()
    image.save(byte_arr, format='PNG')
    
    # Get the binary data from BytesIO
    byte_data = byte_arr.getvalue()
    
    # Encode the binary data to Base64
    base64_encoded = base64.b64encode(byte_data).decode('utf-8')
    
    # Return the Base64 string with the data URI scheme
    return f"data:image/png;base64,{base64_encoded}"

# Items
items = [
    ('Example Item 2', 200, 'This is an example item 2.', 'Pets', generate_random_image((512,480)), 1),
    ('Example Item 3', 300, 'This is an example item 3.', 'Cars', generate_random_image((1024,480)), 1),
    ('Example Item 4', 400, 'This is an example item 4.', 'Cars', generate_random_image((512,2048)), 1),
    ('Example Item 5', 500, 'This is an example item 5.', 'Cars', generate_random_image((1080,480)), 1)
]

cursor.executemany('''
    INSERT INTO items (name, price, description, category , images, store_id)
    VALUES (?, ?, ?, ?, ? , ?)
''', items)

# Orders
time_of_purchase = datetime.now().strftime('%Y-%m-%d')
estimated_delivery = (datetime.now() + timedelta(days=7)).strftime('%Y-%m-%d')
cursor.execute('''
    INSERT INTO orders (user_id, store_id, item_id, quantity, time_of_purchase, estimated_delivery, status)
    VALUES (1, 1, 1, 2, ?, ?, 'Processing')
''', (time_of_purchase, estimated_delivery))

# Clicks
time_of_click = datetime.now().strftime('%Y-%m-%d')
cursor.execute('''
    INSERT INTO clicks (user_id, item_id, time_of_click)
    VALUES (1, 1, ?)
''', (time_of_click,))

# Commit the changes and close the connection
conn.commit()
conn.close()

print("Database setup and example data insertion complete.")
