# Storage Bucket

This module is used to create storage buckets. Buckets are generally configured with the following settings:

1. [Autoclass](#autoclass)
2. [Regional storage location](#regional-storage-location)
3. [Uniform bucket level access](#uniform-bucket-level-access)

## Autoclass

Autoclass is pretty simple way to reduce storage costs unless you have a unique access pattern or specific latency requirements. The pricing model is fair and the management cost is minimal compared to how much cheaper non-standard storage classes are. See the [documentation](https://cloud.google.com/storage/docs/autoclass) for more info.

In the case of this module, we always enable autoclass unless you set `object_ttl_days` to 30 or less. This is because the object will be deleted before autoclass savings can kick in.

## Regional storage location

Multiregion is no longer worth it for the vast majority of use cases. GCP will charge you exorbitant data transfer fees to read your data even if the destination is a sub-region of the multi-region location (e.g. accessing a `us` bucket from a VM in `us-west1`). It seems like multiregion only makes sense if you operate a web service with extremely high availability or low latency requirements.

See the recent [pricing change announcement](https://cloud.google.com/storage/pricing-announce#network) and latest [price tables](https://cloud.google.com/storage/pricing#price-tables) for more info.

## Uniform bucket level access

See [GCP documentation](https://cloud.google.com/storage/docs/uniform-bucket-level-access#should-you-use) for rationale.
