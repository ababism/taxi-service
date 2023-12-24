package generated

import "gitlab/ArtemFed/mts-final-taxi/projects/driver/internal/domain"

func ToTripStatusDomain(ts TripStatus) domain.TripStatus {
	converted := domain.TripStatus(ts)
	return converted
}
func ScrapeStatusesConstants() {
	domain.InitTripStatusCollection(
		ToTripStatusDomain(CANCELED),
		ToTripStatusDomain(DRIVERFOUND),
		ToTripStatusDomain(DRIVERSEARCH),
		ToTripStatusDomain(ENDED),
		ToTripStatusDomain(ONPOSITION),
		ToTripStatusDomain(STARTED))
}
