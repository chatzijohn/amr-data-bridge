CREATE TABLE public.measurements (
	"time" timestamp NOT NULL,
	"waterMeterDevEUI" varchar(255) NOT NULL,
	"waterSupplyId" int4 NULL,
	measurement int4 NOT NULL,
	delta int4 NOT NULL
);
