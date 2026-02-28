# A6X2 Test Fixtures

Device: **Supernote Nomad X2** (`APPLY_EQUIPMENT=N6`)
Resolution: 1404×1872 px

Sources:
- Files marked [orig]: created for this project
- Files marked [snlib]: [Walnut356/snlib](https://github.com/Walnut356/snlib/tree/main/test_notes)
- Files marked [pts]: [philips/supernote-typescript](https://github.com/philips/supernote-typescript/tree/main/tests/input)

Naming convention: `{content}_{orient}_{pages}p[_rtr][_{lang}].note`
Orient: `v{n}` = vertical (orientation code n), `h{n}` = horizontal

| File | Orient | Pages | RTR | Lang | Source | Notes |
|---|---|---|---|---|---|---|
| `blank_v1180_1p.note` | vertical/1180 | 1 | no | — | [orig] | non-default vertical orientation code |
| `blank_h1270_1p.note` | horizontal/1270 | 1 | no | — | [orig] | horizontal orientation |
| `task_v1000_1p_rtr_de.note` | vertical/1000 | 1 | yes | de_DE | [pts] | task list template, German RTR |
| `multilayer_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | 3 layers: LAYER1 + MAINLAYER + BGLAYER |
| `4edges_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | strokes near all four page edges |
| `bmp_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | bitmap background layer |
| `markers_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | marker tool color variations |
| `parallel_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | dense parallel line patterns |
| `pressure_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | pen pressure variation |
| `tilt_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | pen tilt full circle |
| `erase_pre_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | state before erase operation |
| `erase_post_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | state after erase operation |
| `move_pre_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | state before move operation |
| `move_post_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | state after move operation |
| `rotate90_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | content rotated 90° |
| `rotate180_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | content rotated 180° |
| `stroke_v1000_1p.note` | vertical/1000 | 1 | no | — | [snlib] | stroke copy operation |
