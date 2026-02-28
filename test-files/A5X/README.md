# A5X Test Fixtures

Device: **Supernote Manta (gen1)** (`APPLY_EQUIPMENT=A5X`)
Auto-detected as: A6X2 decoder (default fallback)
Source: [philips/supernote-typescript](https://github.com/philips/supernote-typescript/tree/main/tests/input)

Naming convention: `{content}_{orient}_{pages}p[_rtr][_{lang}].note`

| File | Orient | Pages | RTR | Lang | Style | Notes |
|---|---|---|---|---|---|---|
| `ruled_v1000_10p_rtr.note` | vertical/1000 | 10 | yes | en_US | 10mm ruled | handwritten 1–10, writing tools demo |
| `ruled8mm_v1000_2p.note` | vertical/1000 | 2 | no | — | 8mm ruled | old format (no FILE_ID) |
| `lined_v1000_1p_rtr.note` | vertical/1000 | 1 | yes | en_US | user custom lined | firmware 2.14.28 |
