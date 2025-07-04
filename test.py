from datetime import date, timedelta

num_users = 10
jobs_per_user = 15
start_pickup_date = date(2025, 7, 10)

def bool_str(val):
    return 'TRUE' if val else 'FALSE'

sql_values = []

for user_id in range(1, num_users + 1):
    for job_num in range(1, jobs_per_user + 1):
        job_title = f"Delivery Job {user_id}-{job_num}"
        description = "Auto-generated job description"
        cargo_type = "General" if job_num % 2 == 0 else "Electronics"
        urgency = "Urgent" if job_num % 3 == 0 else "Normal"
        truck_size = "Large" if job_num % 4 == 0 else "Medium"
        loading_assistance = bool_str(job_num % 2 == 0)
        pickup_date = start_pickup_date + timedelta(days=job_num)
        delivery_date = pickup_date + timedelta(days=1)
        pickup_time_window = "08:00-10:00"
        delivery_time_window = "14:00-16:00"
        pickup_location = f"Warehouse {chr(65 + (job_num % 5))}"
        delivery_location = f"Destination {chr(70 + (job_num % 5))}"
        payout_amount = 100 + job_num * 5
        early_delivery_bonus = 10.0
        payment_terms = "Net 15" if job_num % 2 == 0 else "COD"
        weight_lb = 500 + job_num * 10
        volume_cu_ft = 100 + job_num * 5

        liftgate = bool_str(job_num % 2 == 0)
        fragile_items = bool_str(job_num % 5 == 0)
        climate_control = bool_str(job_num % 4 == 0)
        assembly_required = bool_str(job_num % 3 == 0)
        extra_insurance = bool_str(job_num % 6 == 0)
        additional_packing = bool_str(job_num % 7 == 0)

        values = f"({user_id}, '{job_title}', '{description}', '{cargo_type}', '{urgency}', '{truck_size}', {loading_assistance}, '{pickup_date}', '{pickup_time_window}', '{delivery_date}', '{delivery_time_window}', '{pickup_location}', '{delivery_location}', {payout_amount:.2f}, {early_delivery_bonus:.2f}, '{payment_terms}', {weight_lb:.2f}, {volume_cu_ft:.2f}, {liftgate}, {fragile_items}, {climate_control}, {assembly_required}, {extra_insurance}, {additional_packing})"
        sql_values.append(values)

# Combine into a single SQL VALUES block
print("VALUES\n" + ",\n".join(sql_values) + ";")
