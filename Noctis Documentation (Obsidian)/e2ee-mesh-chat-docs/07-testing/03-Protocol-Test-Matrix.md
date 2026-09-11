# Protocol Test Matrix

  Test                              Expected result
  --------------------------------- -----------------------
  Valid packet                      Accepted
  Unsupported version               Rejected
  Missing required field            Rejected
  Invalid field type                Rejected
  Oversized field                   Rejected
  Invalid packet type               Rejected
  Modified authenticated metadata   Rejected
  Duplicate packet ID               Rejected/deduplicated
  Replayed message                  Rejected
  Invalid nonce encoding            Rejected
  Invalid cryptographic envelope    Rejected
  Unexpected protocol state         Rejected
  Unknown destination               Controlled failure
