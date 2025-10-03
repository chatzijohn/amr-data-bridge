CREATE TABLE public."rtMeasurements" (
	"time" timestamp NOT NULL,
	"waterMeterDevEUI" varchar(255) NOT NULL,
	"waterSupplyId" int4 NULL,
	measurement int4 NOT NULL,
	alarms varchar(255) NOT NULL
);
