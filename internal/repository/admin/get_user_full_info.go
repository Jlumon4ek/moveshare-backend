package admin

import (
	"context"
	"moveshare/internal/models"
)

func (r *repository) GetUserFullInfo(ctx context.Context, userID int64) (*models.UserFullInfo, error) {
	userInfo := &models.UserFullInfo{}

	userQuery := `
		SELECT id, username, email, role, status, profile_photo_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := r.db.QueryRow(ctx, userQuery, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.Status,
		&user.ProfilePhotoID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	userInfo.User = user

	companyQuery := `
		SELECT id, user_id, company_name, email_address, address, state, mc_license_number,
		       company_description, contact_person, phone_number, city, zip_code, dot_number,
		       created_at, updated_at
		FROM companies
		WHERE user_id = $1
	`
	var company models.Company
	err = r.db.QueryRow(ctx, companyQuery, userID).Scan(
		&company.ID,
		&company.UserID,
		&company.CompanyName,
		&company.EmailAddress,
		&company.Address,
		&company.State,
		&company.MCLicenseNumber,
		&company.CompanyDescription,
		&company.ContactPerson,
		&company.PhoneNumber,
		&company.City,
		&company.ZipCode,
		&company.DotNumber,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err == nil {
		userInfo.Company = &company
	}

	trucksQuery := `
		SELECT id, user_id, truck_name, license_plate, make, model, year, color,
		       length, width, height, max_weight, truck_type, climate_control,
		       liftgate, pallet_jack, security_system, refrigerated, furniture_pads,
		       created_at, updated_at
		FROM trucks
		WHERE user_id = $1
	`
	trucksRows, err := r.db.Query(ctx, trucksQuery, userID)
	if err == nil {
		defer trucksRows.Close()
		for trucksRows.Next() {
			var truck models.TruckSwagger
			err := trucksRows.Scan(
				&truck.ID,
				&truck.UserID,
				&truck.TruckName,
				&truck.LicensePlate,
				&truck.Make,
				&truck.Model,
				&truck.Year,
				&truck.Color,
				&truck.Length,
				&truck.Width,
				&truck.Height,
				&truck.MaxWeight,
				&truck.TruckType,
				&truck.ClimateControl,
				&truck.Liftgate,
				&truck.PalletJack,
				&truck.SecuritySystem,
				&truck.Refrigerated,
				&truck.FurniturePads,
				&truck.CreatedAt,
				&truck.UpdatedAt,
			)
			if err == nil {
				userInfo.Trucks = append(userInfo.Trucks, truck)
			}
		}
	}

	jobsQuery := `
		SELECT id, contractor_id, job_type, number_of_bedrooms, packing_boxes,
		       bulky_items, inventory_list, hoisting, additional_services_description,
		       estimated_crew_assistants, truck_size, pickup_address, pickup_floor,
		       pickup_building_type, pickup_walk_distance, delivery_address,
		       delivery_floor, delivery_building_type, delivery_walk_distance,
		       distance_miles, job_status, pickup_date, pickup_time_from,
		       pickup_time_to, delivery_date, delivery_time_from, delivery_time_to,
		       cut_amount, payment_amount, weight_lbs, volume_cu_ft, created_at, updated_at
		FROM jobs
		WHERE contractor_id = $1
	`
	jobsRows, err := r.db.Query(ctx, jobsQuery, userID)
	if err == nil {
		defer jobsRows.Close()
		for jobsRows.Next() {
			var job models.Job
			err := jobsRows.Scan(
				&job.ID,
				&job.ContractorID,
				&job.JobType,
				&job.NumberOfBedrooms,
				&job.PackingBoxes,
				&job.BulkyItems,
				&job.InventoryList,
				&job.Hoisting,
				&job.AdditionalServicesDescription,
				&job.EstimatedCrewAssistants,
				&job.TruckSize,
				&job.PickupAddress,
				&job.PickupFloor,
				&job.PickupBuildingType,
				&job.PickupWalkDistance,
				&job.DeliveryAddress,
				&job.DeliveryFloor,
				&job.DeliveryBuildingType,
				&job.DeliveryWalkDistance,
				&job.DistanceMiles,
				&job.JobStatus,
				&job.PickupDate,
				&job.PickupTimeFrom,
				&job.PickupTimeTo,
				&job.DeliveryDate,
				&job.DeliveryTimeFrom,
				&job.DeliveryTimeTo,
				&job.CutAmount,
				&job.PaymentAmount,
				&job.WeightLbs,
				&job.VolumeCuFt,
				&job.CreatedAt,
				&job.UpdatedAt,
			)
			if err == nil {
				userInfo.Jobs = append(userInfo.Jobs, job)
			}
		}
	}

	reviewsQuery := `
		SELECT id, job_id, reviewer_id, reviewee_id, rating, comment, created_at, updated_at
		FROM reviews
		WHERE reviewer_id = $1 OR reviewee_id = $1
	`
	reviewsRows, err := r.db.Query(ctx, reviewsQuery, userID)
	if err == nil {
		defer reviewsRows.Close()
		for reviewsRows.Next() {
			var review models.Review
			err := reviewsRows.Scan(
				&review.ID,
				&review.JobID,
				&review.ReviewerID,
				&review.RevieweeID,
				&review.Rating,
				&review.Comment,
				&review.CreatedAt,
				&review.UpdatedAt,
			)
			if err == nil {
				userInfo.Reviews = append(userInfo.Reviews, review)
			}
		}
	}

	paymentsQuery := `
		SELECT id, user_id, job_id, stripe_payment_intent_id, stripe_payment_method_id,
		       stripe_customer_id, amount_cents, currency, status, description,
		       failure_reason, created_at, updated_at
		FROM payments
		WHERE user_id = $1
	`
	paymentsRows, err := r.db.Query(ctx, paymentsQuery, userID)
	if err == nil {
		defer paymentsRows.Close()
		for paymentsRows.Next() {
			var payment models.Payment
			err := paymentsRows.Scan(
				&payment.ID,
				&payment.UserID,
				&payment.JobID,
				&payment.StripePaymentIntentID,
				&payment.StripePaymentMethodID,
				&payment.StripeCustomerID,
				&payment.AmountCents,
				&payment.Currency,
				&payment.Status,
				&payment.Description,
				&payment.FailureReason,
				&payment.CreatedAt,
				&payment.UpdatedAt,
			)
			if err == nil {
				userInfo.Payments = append(userInfo.Payments, payment)
			}
		}
	}

	verificationQuery := `
		SELECT id, user_id, object_name, file_type, status, created_at
		FROM verification_file
		WHERE user_id = $1
	`
	verificationRows, err := r.db.Query(ctx, verificationQuery, userID)
	if err == nil {
		defer verificationRows.Close()
		for verificationRows.Next() {
			var verification models.VerificationFile
			var id int64
			var userIDCheck int64
			var createdAt interface{}
			err := verificationRows.Scan(
				&id,
				&userIDCheck,
				&verification.ObjectName,
				&verification.FileType,
				&verification.Status,
				&createdAt,
			)
			if err == nil {
				userInfo.Verification = append(userInfo.Verification, verification)
			}
		}
	}

	return userInfo, nil
}