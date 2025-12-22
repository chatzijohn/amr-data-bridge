CREATE TABLE public."waterSupplies" (
    id SERIAL PRIMARY KEY,
    "supplyNumber" VARCHAR(255) NOT NULL UNIQUE,
    geometry public.geometry(Point) NOT NULL,
    "waterMeterDevEUI" VARCHAR(255) UNIQUE,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "isActive" BOOLEAN DEFAULT true
);