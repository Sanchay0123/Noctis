# M10 Final Evidence Inventory

## Validation Evidence Location

`~/Documents/Noctis/CrytpProject/m10_final_validation/`

The following artifacts were reported present after final validation. SHA-256 values provide an integrity snapshot of the evidence directory at the time of inventory.

| File | Size | SHA-256 |
|---|---:|---|
| `FINAL-SYSTEM-VALIDATION-REPORT.md` | 3.1K | `aeb7d41c942d72833c917e43c6b6004e7d367fd778174d08a647dc7734f72269` |
| `build-test-results.txt` | 5.8M | `e68c04e6025e232ab2a6e67a0fe1fc82ac58249f4191a1ce2ae09fdf65a7eb7a` |
| `crypto-results.txt` | 1.3K | `3ad2c4469eca57de61091928141c1de05da4aa535cb8c27a325af4238e0c659d` |
| `direct-messaging-results.txt` | 932B | `3ccc56acb3ff056ab05286123dd50a03dfedf751f6f21901649735f93f8734b0` |
| `environment.txt` | 3.5K | `d3a4a0437c925b449e9514061ad4475abc93d9accf1b104e438020e35b26dfae` |
| `final-summary.md` | 3.1K | `aeb7d41c942d72833c917e43c6b6004e7d367fd778174d08a647dc7734f72269` |
| `gui-results.txt` | 212B | `8e0f30e5bceb1f0ffbfcbdc601895862a90370ab1b8022fa022c4eb303ac19e8` |
| `multihop-results.txt` | 959B | `4d499958149dd9ac63e8735dbc71a706a937261e57735b6488b09dd624b469fa` |
| `protocol-results.txt` | 1.6K | `0a601feb15840e3d1ff149ec96cd0beb0b76cb8c837d69c0292b2f6ba60b9483` |
| `resource-results.txt` | 866B | `0867c8f71877fa8ee761465fbe277245bc3f54d1589b5552e78728e0efd5da71` |
| `routing-results.txt` | 1.4K | `9af978bb60c6fd272af274bd541f525165dcf8c4bf21755057075b2c1ffbc184` |
| `security-results.txt` | 2.1K | `bd27e575609f84f1b3ac4c2d97d0e8f5353041bc94ffce0fd23e7874737e7cdf` |

## Repository-State Note

Final validation reported no modifications to production source under `internal/`, `cmd/` or `protocol/`. The repository may contain evaluation-only evidence directories and the untracked M10 benchmark harness; these are not production implementation changes.

## Evidence Interpretation

- Final validation evidence is distinct from historical M9 and M10.2 evidence.
- GUI final execution was NOT MEASURED in the headless validation environment.
- M10.2 measurements remain controlled local loopback results and are not Internet deployment benchmarks.
