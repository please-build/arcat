Arcat
=====

Arcat (formerly jarcat) is a tool for creating, extracting and concatenating archives
(specifically `.zip`, `.tar` and .`ar` archives). It also attempts to be as
deterministic as possible in its output (for example, all timestamps are fixated
in order to avoid unnecessarily perturbing outputs).

Its main original use was for efficiently concatenating zip archives (in the
form of `.jar` files). Since then it has grown into a general archiving tool,
but still retains that critical original function.

Its only main consumer is [Please](https://please.build), and at present the
CLI is fairly bespoke to a bunch of historical needs there. It may change
significantly in future if we decide to make it more generally useful.
