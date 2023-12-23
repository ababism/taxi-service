package generated

import "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"

func ToTripStatusDomain(ts TripStatus) domain.TripStatus {
	converted := domain.TripStatus(ts)
	return converted
}
func GetStatuses() domain.TripStatusCollection {
	return domain.NewTripStatusCollection(
		ToTripStatusDomain(CANCELED),
		ToTripStatusDomain(DRIVERFOUND),
		ToTripStatusDomain(DRIVERSEARCH),
		ToTripStatusDomain(ENDED),
		ToTripStatusDomain(ONPOSITION),
		ToTripStatusDomain(STARTED))
}
