# Cryptographic Test Matrix

  Test                    Expected result            Security property
  ----------------------- -------------------------- --------------------------
  Generate identity       Success                    identity availability
  Reload identity         Same public identity       persistence
  Valid signature         Accepted                   authentication
  Modified message        Rejected                   integrity
  Invalid signature       Rejected                   authentication
  X25519 agreement        Same shared secret         key agreement
  Identity substitution   Rejected                   MITM resistance
  HKDF same inputs        Same output                deterministic derivation
  HKDF changed context    Different output           domain separation
  Encrypt/decrypt         Success                    confidentiality
  Modified ciphertext     Rejected                   integrity
  Modified AAD            Rejected                   metadata integrity
  Wrong key               Rejected                   key separation
  Nonce reuse scenario    Prevented/rejected         AEAD safety
  Empty message           Correct behavior           edge case
  Large message           Correct bounded behavior   robustness
  Unicode message         Correct behavior           application correctness
