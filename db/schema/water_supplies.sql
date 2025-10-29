CREATE TABLE public."waterSupplies" (
    id SERIAL PRIMARY KEY,
    "supplyNumber" VARCHAR(255) NOT NULL,
    geometry public.geometry(point) NOT NULL,
    "waterMeterSerialNumber" VARCHAR(255),
    "currentImage" VARCHAR(255),
    "previousImage" VARCHAR(255),
    "createdAt" TIMESTAMP NOT NULL,
    "updatedAt" TIMESTAMP NOT NULL
);
