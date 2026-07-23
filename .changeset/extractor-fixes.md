---
"mdl": patch
---

Rewrote the corocoro extractor to parse the site's new Next.js RSC data format.

Fixed a nil-pointer dereference in the cmoa extractor when no cookie is configured, fixed paid chapter detection for the ganma extractor, and replaced `FindChapters` panics with graceful unsupported-chapter-listing errors. Added `manga.ChapterListingFeature` for extractor capability detection and proper sentinel errors for paid chapters / unsupported listing.
