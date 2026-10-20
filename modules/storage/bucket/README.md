# Storage Bucket

This module is used to create storage buckets. Buckets use the following settings:

1. [Autoclass](#autoclass)
2. [Regional storage location](#regional-storage-location)
3. [Uniform bucket level access](#uniform-bucket-level-access)

## Autoclass

Autoclass can reduce storage costs by moving infrequently accessed objects to colder storage classes.

This module enables Autoclass when `object_ttl_days` is unset, `null`, or greater than 30, and disables it for TTLs of 30 days or less. Objects of at least 128 KiB become eligible for Nearline after 30 days without data access. This is a simple cutoff, not a cost-savings guarantee: savings also depend on object sizes, access patterns, and Autoclass management fees.

See [GCP docuemntation](https://cloud.google.com/storage/docs/autoclass) for more info.

## Regional storage location

Multiregion is no longer worth it for the vast majority of use cases. GCP will charge you exorbitant data transfer fees just to read your data - even if the destination is a sub-region of the multi-region location (e.g. accessing a `us` bucket from a VM in `us-west1`). 

See the recent [pricing change announcement](https://cloud.google.com/storage/pricing-announce#network) and latest [price tables](https://cloud.google.com/storage/pricing#price-tables) for more info.

## Uniform bucket level access
See [GCP documentation](https://cloud.google.com/storage/docs/uniform-bucket-level-access#should-you-use) for rationale.
