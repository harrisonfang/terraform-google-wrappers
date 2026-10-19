# Storage Bucket

This module is used to create storage buckets. Buckets are always configured with the following settings:

1. [Autoclass](#autoclass)
2. [Regional storage location](#regional-storage-location)
3. [Uniform bucket level access](#uniform-bucket-level-access)

## Autoclass

Autoclass leads to some pretty good reductions in storage costs unless you have a unique access pattern or specific latency requirements. The pricing model is fair because you only incur additional costs when accessing objects that haven't been used in a long time. This is worth it given how much cheaper non-standard storage classes are.

See [GCP docuemntation](https://cloud.google.com/storage/docs/autoclass) for more info.

## Regional storage location

Multiregion is no longer worth it for the vast majority of use cases. GCP will charge you exorbitant data transfer fees just to read your data - even if the destination is a sub-region of the multi-region location (e.g. accessing a `us` bucket from a VM in `us-west1`). 

See the recent [pricing change announcement](https://cloud.google.com/storage/pricing-announce#network) and latest [price tables](https://cloud.google.com/storage/pricing#price-tables) for more info.

## Uniform bucket level access
See [GCP documentation](https://cloud.google.com/storage/docs/uniform-bucket-level-access#should-you-use) for rationale.
