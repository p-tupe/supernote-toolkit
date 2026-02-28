# A5X2 Test Fixtures

Device: **Supernote Manta X2** (`APPLY_EQUIPMENT=N5`)
Resolution: 1920×2560 px
Source: created for this project

Naming convention: `{content}_{orient}_{pages}p[_rtr][_{lang}].note`
Orient: `v{n}` = vertical (orientation code n), `h{n}` = horizontal

| File | Orient | Pages | RTR | Lang | Notes |
|---|---|---|---|---|---|
| `blank_v1000_1p.note` | vertical/1000 | 1 | no | — | baseline conversion test |
| `blank_v1000_1p_artifacts.note` | vertical/1000 | 1 | no | — | contains layer rendering artifacts |
| `wip_v1000_1p_rtr.note` | vertical/1000 | 1 | in-progress | en_US | real-time text recognition |
| `blank_h1090_1p_rtr.note` | horizontal/1090 | 1 | yes | en_US | horizontal orientation |
| `text_v1000_1p_rtr.note` | vertical/1000 | 1 | yes | — | real handwritten text content |
