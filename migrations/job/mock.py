import psycopg2
from psycopg2.extras import execute_batch
from random import randint, choice, uniform
from datetime import datetime, timedelta

# Параметры подключения к БД
conn_params = {
    'dbname': 'pepsi',
    'user': 'pepsi',
    'password': 'pepsi123',
    'host': 'localhost',
    'port': 5432,
}

job_statuses = ['pending', 'in_progress', 'completed', 'cancelled', 'rejected', 'active']
job_types = ['residential', 'commercial', 'storage']
truck_sizes = ['small', 'medium', 'large']
building_types = ['stairs', 'elevator', 'stairs_and_elevator']
packing_boxes_options = ['t', 'f']
bulky_items_options = ['t', 'f']
inventory_list_options = ['t', 'f']
hoisting_options = ['t', 'f']

def random_date(start_days=0, end_days=30):
    """Возвращает случайную дату в будущем в диапазоне от start_days до end_days от сегодня."""
    start = datetime.now() + timedelta(days=start_days)
    end = datetime.now() + timedelta(days=end_days)
    return start + (end - start) * uniform(0, 1)

def generate_jobs_for_user(user_id, min_jobs=10, max_jobs=20):
    jobs = []
    n = randint(min_jobs, max_jobs)
    for _ in range(n):
        pickup_date = random_date()
        delivery_date = pickup_date + timedelta(days=randint(1, 7))
        
        job = {
            'contractor_id': user_id,
            'job_type': choice(job_types),
            'number_of_bedrooms': randint(1, 5),
            'packing_boxes': choice(packing_boxes_options),
            'bulky_items': choice(bulky_items_options),
            'inventory_list': choice(inventory_list_options),
            'hoisting': choice(hoisting_options),
            'additional_services_description': 'Additional service ' + str(randint(1, 1000)),
            'estimated_crew_assistants': randint(1, 5),
            'truck_size': choice(truck_sizes),
            'pickup_address': f'{randint(1, 100)} Main St',
            'pickup_floor': randint(0, 10),
            'pickup_building_type': choice(building_types),
            'pickup_walk_distance': choice(['0-50 ft', '50-100 ft', '100-200 ft']),
            'delivery_address': f'{randint(101, 200)} Market St',
            'delivery_floor': randint(0, 10),
            'delivery_building_type': choice(building_types),
            'delivery_walk_distance': choice(['0-50 ft', '50-100 ft', '100-200 ft']),
            'distance_miles': round(uniform(1, 100), 2),
            'job_status': choice(job_statuses),
            'pickup_date': pickup_date.date(),
            'pickup_time_from': (pickup_date + timedelta(hours=randint(6, 12))).time(),
            'pickup_time_to': (pickup_date + timedelta(hours=randint(13, 18))).time(),
            'delivery_date': delivery_date.date(),
            'delivery_time_from': (delivery_date + timedelta(hours=randint(6, 12))).time(),
            'delivery_time_to': (delivery_date + timedelta(hours=randint(13, 18))).time(),
            'cut_amount': round(uniform(50, 500), 2),
            'payment_amount': round(uniform(500, 5000), 2),
            'weight_lbs': randint(500, 5000),
            'volume_cu_ft': randint(100, 10000),
            'created_at': datetime.now(),
            'updated_at': datetime.now(),
        }
        jobs.append(job)
    return jobs

def main():
    try:
        conn = psycopg2.connect(**conn_params)
        cursor = conn.cursor()
        
        # Получаем список id пользователей, которым надо добавить задания
        cursor.execute("SELECT id FROM users WHERE id BETWEEN 1 AND 12")
        user_ids = [row[0] for row in cursor.fetchall()]
        
        all_jobs = []
        for user_id in user_ids:
            user_jobs = generate_jobs_for_user(user_id)
            all_jobs.extend(user_jobs)
        
        insert_query = """
            INSERT INTO jobs (
                contractor_id, job_type, number_of_bedrooms, packing_boxes, bulky_items, inventory_list, hoisting,
                additional_services_description, estimated_crew_assistants, truck_size, pickup_address, pickup_floor,
                pickup_building_type, pickup_walk_distance, delivery_address, delivery_floor, delivery_building_type,
                delivery_walk_distance, distance_miles, job_status, pickup_date, pickup_time_from, pickup_time_to,
                delivery_date, delivery_time_from, delivery_time_to, cut_amount, payment_amount, weight_lbs,
                volume_cu_ft, created_at, updated_at
            ) VALUES (
                %(contractor_id)s, %(job_type)s, %(number_of_bedrooms)s, %(packing_boxes)s, %(bulky_items)s,
                %(inventory_list)s, %(hoisting)s, %(additional_services_description)s, %(estimated_crew_assistants)s,
                %(truck_size)s, %(pickup_address)s, %(pickup_floor)s, %(pickup_building_type)s, %(pickup_walk_distance)s,
                %(delivery_address)s, %(delivery_floor)s, %(delivery_building_type)s, %(delivery_walk_distance)s,
                %(distance_miles)s, %(job_status)s, %(pickup_date)s, %(pickup_time_from)s, %(pickup_time_to)s,
                %(delivery_date)s, %(delivery_time_from)s, %(delivery_time_to)s, %(cut_amount)s, %(payment_amount)s,
                %(weight_lbs)s, %(volume_cu_ft)s, %(created_at)s, %(updated_at)s
            )
        """
        
        execute_batch(cursor, insert_query, all_jobs, page_size=100)
        conn.commit()
        print(f'Inserted {len(all_jobs)} jobs successfully.')
        
    except Exception as e:
        print('Error:', e)
    finally:
        cursor.close()
        conn.close()

if __name__ == '__main__':
    main()
