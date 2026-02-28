# A6X Test Fixtures

Device: **Supernote Nomad (gen1)** (`APPLY_EQUIPMENT=A6X`)
Auto-detected as: A6X2 decoder (default fallback)
Source: [philips/supernote-typescript](https://github.com/philips/supernote-typescript/tree/main/tests/input) (firmware 3.15.27)

Naming convention: `{content}_{orient}_{pages}p[_rtr][_{lang}].note`

| File | Orient | Pages | RTR | Lang | Style | Notes |
|---|---|---|---|---|---|---|
| `blank_v1000_2p.note` | vertical/1000 | 2 | no | — | white/blank | 2-page, writing tools demo; page 2 has link back to page 1 |
| `shapes_v1000_1p_rtr.note` | vertical/1000 | 1 | yes | en_US | white/blank | shapes, patterns, headings, keyword highlights |
| `blank_v1000_1p_rtr_tr.note` | vertical/1000 | 1 | yes | tr_TR | white/blank | Turkish RTR |
