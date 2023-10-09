# Id-iota

This is simple unique identifier that is not overcomplicated. It does not scale because we are not building another twitter here.
Based on timestamps as uint32 so its usability is till 2106 year.

64 bits

- 32 bits of timestamp
- 32 bits of randomness

Default encoding is in base36 (its popular in crypto and we are trendy as fuck) and that gives us 13 characters.
